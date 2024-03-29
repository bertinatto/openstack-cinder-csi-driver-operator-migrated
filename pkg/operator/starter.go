package operator

import (
	"context"
	"fmt"
	"time"

	dynamicclient "k8s.io/client-go/dynamic"
	kubeclient "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/klog"

	"github.com/openshift/library-go/pkg/controller/controllercmd"
	csicontrollerset "github.com/openshift/library-go/pkg/operator/csi/controllerset"
	goc "github.com/openshift/library-go/pkg/operator/genericoperatorclient"
	"github.com/openshift/library-go/pkg/operator/v1helpers"

	"github.com/openshift/openstack-cinder-csi-driver-operator/pkg/apis/operator/v1alpha1"
	"github.com/openshift/openstack-cinder-csi-driver-operator/pkg/generated"
)

const (
	operandName       = "openstack-cinder-csi-driver"
	operandNamespace  = "openshift-openstack-cinder-csi-driver"
	operatorNamespace = "openshift-openstack-cinder-csi-driver-operator"

	resync = 20 * time.Minute
)

func RunOperator(ctx context.Context, controllerConfig *controllercmd.ControllerContext) error {
	// Create clientsets and informers
	dynamicClient := dynamicclient.NewForConfigOrDie(rest.AddUserAgent(controllerConfig.KubeConfig, "dynamic-client"))
	kubeClient := kubeclient.NewForConfigOrDie(rest.AddUserAgent(controllerConfig.KubeConfig, "kube-client"))
	kubeInformersForNamespaces := v1helpers.NewKubeInformersForNamespaces(kubeClient, "", operandNamespace, operatorNamespace)

	// Create GenericOperatorclient. This is used by controllers created down below
	gvr := v1alpha1.SchemeGroupVersion.WithResource("openstackcinderdrivers")
	operatorClient, dynamicInformers, err := goc.NewClusterScopedOperatorClient(controllerConfig.KubeConfig, gvr)
	if err != nil {
		return err
	}

	csiControllerSet := csicontrollerset.New(
		operatorClient,
		controllerConfig.EventRecorder,
	).WithLogLevelController().WithManagementStateController(
		operandName,
		false,
	).WithStaticResourcesController(
		"OpenStackCinderDriverStaticResources",
		kubeClient,
		kubeInformersForNamespaces,
		generated.Asset,
		[]string{
			"namespace.yaml",
			"csidriver.yaml",
			"storageclass.yaml",
			"controller_sa.yaml",
			"node_sa.yaml",
			"rbac/attacher_binding.yaml",
			"rbac/attacher_role.yaml",
			"rbac/controller_privileged_binding.yaml",
			"rbac/node_privileged_binding.yaml",
			"rbac/privileged_role.yaml",
			"rbac/provisioner_binding.yaml",
			"rbac/provisioner_role.yaml",
			"rbac/resizer_binding.yaml",
			"rbac/resizer_role.yaml",
			"rbac/snapshotter_binding.yaml",
			"rbac/snapshotter_role.yaml",
		},
	).WithCSIDriverController(
		"OpenStackCinderDriverController",
		operandName,
		operandNamespace,
		generated.MustAsset,
		kubeClient,
		kubeInformersForNamespaces.InformersFor(operandNamespace),
		csicontrollerset.WithCloudCredentials(dynamicClient, "credentials.yaml"),
		csicontrollerset.WithControllerService("controller.yaml"),
		csicontrollerset.WithNodeService("node.yaml"),
	)

	if err != nil {
		return err
	}

	klog.Info("Starting the informers")
	go kubeInformersForNamespaces.Start(ctx.Done())
	go dynamicInformers.Start(ctx.Done())

	klog.Info("Starting controllerset")
	go csiControllerSet.Run(ctx, 1)

	<-ctx.Done()

	return fmt.Errorf("stopped")
}
