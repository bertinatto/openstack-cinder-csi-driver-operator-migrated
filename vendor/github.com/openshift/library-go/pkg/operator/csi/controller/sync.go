package controller

import (
	"fmt"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"

	operatorv1 "github.com/openshift/api/operator/v1"
	"github.com/openshift/library-go/pkg/operator/resource/resourceapply"
	"github.com/openshift/library-go/pkg/operator/resource/resourcemerge"
	"github.com/openshift/library-go/pkg/operator/resource/resourceread"
	"github.com/openshift/library-go/pkg/operator/v1helpers"
)

const (
	// Index of a container in assets/controller.yaml and assets/node.yaml
	csiDriverContainerIndex           = 0 // Both Deployment and DaemonSet
	provisionerContainerIndex         = 1
	attacherContainerIndex            = 2
	resizerContainerIndex             = 3
	snapshottterContainerIndex        = 4
	nodeDriverRegistrarContainerIndex = 1
	livenessProbeContainerIndex       = 2 // Only in DaemonSet
)

func (c *Controller) syncCredentialsRequest(status *operatorv1.OperatorStatus) (*unstructured.Unstructured, error) {
	cr := readCredentialRequestsOrDie(c.credentialsManifest)

	// Set spec.secretRef.namespace
	err := unstructured.SetNestedField(cr.Object, c.csiDriverNamespace, "spec", "secretRef", "namespace")
	if err != nil {
		return nil, err
	}

	var expectedGeneration int64 = -1
	generation := resourcemerge.GenerationFor(
		status.Generations,
		schema.GroupResource{Group: credentialsRequestGroup, Resource: credentialsRequestResource},
		cr.GetNamespace(),
		cr.GetName())
	if generation != nil {
		expectedGeneration = generation.LastGeneration
	}

	cr, _, err = applyCredentialsRequest(c.dynamicClient, c.eventRecorder, cr, expectedGeneration)
	return cr, err
}

func (c *Controller) syncDeployment(spec *operatorv1.OperatorSpec, status *operatorv1.OperatorStatus) (*appsv1.Deployment, error) {
	deploy := c.getExpectedDeployment(spec)

	deploy, _, err := resourceapply.ApplyDeployment(
		c.kubeClient.AppsV1(),
		c.eventRecorder,
		deploy,
		resourcemerge.ExpectedDeploymentGeneration(deploy, status.Generations))
	if err != nil {
		return nil, err
	}

	return deploy, nil
}

func (c *Controller) syncDaemonSet(spec *operatorv1.OperatorSpec, status *operatorv1.OperatorStatus) (*appsv1.DaemonSet, error) {
	daemonSet := c.getExpectedDaemonSet(spec)

	daemonSet, _, err := resourceapply.ApplyDaemonSet(
		c.kubeClient.AppsV1(),
		c.eventRecorder,
		daemonSet,
		resourcemerge.ExpectedDaemonSetGeneration(daemonSet, status.Generations))
	if err != nil {
		return nil, err
	}

	return daemonSet, nil
}

func (c *Controller) syncStatus(
	meta *metav1.ObjectMeta,
	status *operatorv1.OperatorStatus,
	deployment *appsv1.Deployment,
	daemonSet *appsv1.DaemonSet,
	credentialsRequest *unstructured.Unstructured) error {

	// Set the last generation change we dealt with
	status.ObservedGeneration = meta.Generation

	// Node Service is mandatory, so always set set  generation
	resourcemerge.SetDaemonSetGeneration(&status.Generations, daemonSet)

	// Set number of replicas
	if daemonSet.Status.NumberUnavailable == 0 {
		status.ReadyReplicas = daemonSet.Status.UpdatedNumberScheduled
	}

	// Credentials are optional, so maybe set generation
	if credentialsRequest != nil {
		resourcemerge.SetGeneration(&status.Generations, operatorv1.GenerationStatus{
			Group:          credentialsRequestGroup,
			Resource:       credentialsRequestResource,
			Namespace:      credentialsRequest.GetNamespace(),
			Name:           credentialsRequest.GetName(),
			LastGeneration: credentialsRequest.GetGeneration(),
		})
	}

	// Controller Service is not mandatory as well
	if c.controllerManifest != nil {
		// CSI Controller Service was deployed, set deployment generation
		resourcemerge.SetDeploymentGeneration(&status.Generations, deployment)

		// Add number of CSI controllers to the number of replicas ready
		if deployment != nil {
			if deployment.Status.UnavailableReplicas == 0 && daemonSet.Status.NumberUnavailable == 0 {
				status.ReadyReplicas += deployment.Status.UpdatedReplicas
			}
		}
	}

	// Finally, set the conditions

	// The operator does not have any prerequisites (at least now)
	v1helpers.SetOperatorCondition(&status.Conditions,
		operatorv1.OperatorCondition{
			Type:   operatorv1.OperatorStatusTypePrereqsSatisfied,
			Status: operatorv1.ConditionTrue,
		})

	// The operator is always upgradeable (at least now)
	v1helpers.SetOperatorCondition(&status.Conditions,
		operatorv1.OperatorCondition{
			Type:   operatorv1.OperatorStatusTypeUpgradeable,
			Status: operatorv1.ConditionTrue,
		})

	// The operator is avaiable for now
	v1helpers.SetOperatorCondition(&status.Conditions,
		operatorv1.OperatorCondition{
			Type:   operatorv1.OperatorStatusTypeAvailable,
			Status: operatorv1.ConditionTrue,
		})

	// Make it not available if daemonSet hasn't deployed the pods
	if !isDaemonSetAvailable(daemonSet) {
		v1helpers.SetOperatorCondition(&status.Conditions,
			operatorv1.OperatorCondition{
				Type:    operatorv1.OperatorStatusTypeAvailable,
				Status:  operatorv1.ConditionFalse,
				Message: "Waiting for the DaemonSet to deploy the CSI Node Service",
				Reason:  "AsExpected",
			})
	}

	// Make it not available if deployment hasn't deployed the pods
	if c.controllerManifest != nil {
		if !isDeploymentAvailable(deployment) {
			v1helpers.SetOperatorCondition(&status.Conditions,
				operatorv1.OperatorCondition{
					Type:    operatorv1.OperatorStatusTypeAvailable,
					Status:  operatorv1.ConditionFalse,
					Message: "Waiting for Deployment to deploy the CSI Controller Service",
					Reason:  "AsExpected",
				})
		}
	}

	// The operator is not progressing for now
	v1helpers.SetOperatorCondition(&status.Conditions,
		operatorv1.OperatorCondition{
			Type:   operatorv1.OperatorStatusTypeProgressing,
			Status: operatorv1.ConditionFalse,
			Reason: "AsExpected",
		})

	isProgressing, msg := c.getDaemonSetProgress(status, daemonSet)
	if isProgressing {
		v1helpers.SetOperatorCondition(&status.Conditions,
			operatorv1.OperatorCondition{
				Type:    operatorv1.OperatorStatusTypeProgressing,
				Status:  operatorv1.ConditionTrue,
				Message: msg,
				Reason:  "AsExpected",
			})
	}

	if c.controllerManifest != nil {
		// CSI Controller deployed, let's check its progressing state
		isProgressing, msg := c.getDeploymentProgress(status, deployment)
		if isProgressing {
			v1helpers.SetOperatorCondition(&status.Conditions,
				operatorv1.OperatorCondition{
					Type:    operatorv1.OperatorStatusTypeProgressing,
					Status:  operatorv1.ConditionTrue,
					Message: msg,
					Reason:  "AsExpected",
				})
		}
	}

	return nil
}

func (c *Controller) getExpectedDeployment(spec *operatorv1.OperatorSpec) *appsv1.Deployment {
	deployment := resourceread.ReadDeploymentV1OrDie(c.controllerManifest)

	if c.images.csiDriver != "" {
		deployment.Spec.Template.Spec.Containers[0].Image = c.images.csiDriver
	}
	if c.images.provisioner != "" {
		deployment.Spec.Template.Spec.Containers[provisionerContainerIndex].Image = c.images.provisioner
	}
	if c.images.attacher != "" {
		deployment.Spec.Template.Spec.Containers[attacherContainerIndex].Image = c.images.attacher
	}
	if c.images.resizer != "" {
		deployment.Spec.Template.Spec.Containers[resizerContainerIndex].Image = c.images.resizer
	}
	if c.images.snapshotter != "" {
		deployment.Spec.Template.Spec.Containers[snapshottterContainerIndex].Image = c.images.snapshotter
	}

	// TODO: add LivenessProbe when

	logLevel := getLogLevel(spec.LogLevel)
	for i, container := range deployment.Spec.Template.Spec.Containers {
		for j, arg := range container.Args {
			if strings.HasPrefix(arg, "--v=") {
				deployment.Spec.Template.Spec.Containers[i].Args[j] = fmt.Sprintf("--v=%d", logLevel)
			}
		}
	}

	return deployment
}

func (c *Controller) getExpectedDaemonSet(spec *operatorv1.OperatorSpec) *appsv1.DaemonSet {
	daemonSet := resourceread.ReadDaemonSetV1OrDie(c.nodeManifest)

	if c.images.csiDriver != "" {
		daemonSet.Spec.Template.Spec.Containers[csiDriverContainerIndex].Image = c.images.csiDriver
	}
	if c.images.nodeDriverRegistrar != "" {
		daemonSet.Spec.Template.Spec.Containers[nodeDriverRegistrarContainerIndex].Image = c.images.nodeDriverRegistrar
	}
	if c.images.livenessProbe != "" {
		daemonSet.Spec.Template.Spec.Containers[livenessProbeContainerIndex].Image = c.images.livenessProbe
	}

	logLevel := getLogLevel(spec.LogLevel)
	for i, container := range daemonSet.Spec.Template.Spec.Containers {
		for j, arg := range container.Args {
			if strings.HasPrefix(arg, "--v=") {
				daemonSet.Spec.Template.Spec.Containers[i].Args[j] = fmt.Sprintf("--v=%d", logLevel)
			}
		}
	}

	return daemonSet
}

func (c *Controller) getDaemonSetProgress(status *operatorv1.OperatorStatus, daemonSet *appsv1.DaemonSet) (bool, string) {
	switch {
	case daemonSet == nil:
		return true, "Waiting for DaemonSet to be created"
	case daemonSet.Generation != daemonSet.Status.ObservedGeneration:
		return true, "Waiting for DaemonSet to act on changes"
	case daemonSet.Status.NumberUnavailable > 0:
		return true, "Waiting for DaemonSet to deploy node pods"
	}
	return false, ""
}

func (c *Controller) getDeploymentProgress(status *operatorv1.OperatorStatus, deployment *appsv1.Deployment) (bool, string) {

	var deploymentExpectedReplicas int32
	if deployment != nil && deployment.Spec.Replicas != nil {
		deploymentExpectedReplicas = *deployment.Spec.Replicas
	}

	switch {
	case deployment == nil:
		return true, "Waiting for Deployment to be created"
	case deployment.Generation != deployment.Status.ObservedGeneration:
		return true, "Waiting for Deployment to act on changes"
	case deployment.Status.UnavailableReplicas > 0:
		return true, "Waiting for Deployment to deploy controller pods"
	case deployment.Status.UpdatedReplicas < deploymentExpectedReplicas:
		return true, "Waiting for Deployment to update pods"
	case deployment.Status.AvailableReplicas < deploymentExpectedReplicas:
		return true, "Waiting for Deployment to deploy pods"
	}

	return false, ""
}

func isDaemonSetAvailable(d *appsv1.DaemonSet) bool {
	return d != nil && d.Status.NumberAvailable > 0
}

func isDeploymentAvailable(d *appsv1.Deployment) bool {
	return d != nil && d.Status.AvailableReplicas > 0
}

func getLogLevel(logLevel operatorv1.LogLevel) int {
	switch logLevel {
	case operatorv1.Normal, "":
		return 2
	case operatorv1.Debug:
		return 4
	case operatorv1.Trace:
		return 6
	case operatorv1.TraceAll:
		return 100
	default:
		return 2
	}
}
