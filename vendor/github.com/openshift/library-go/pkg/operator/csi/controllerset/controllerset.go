package controllerset

import (
	"context"

	"github.com/openshift/library-go/pkg/controller/factory"
	"github.com/openshift/library-go/pkg/operator/csi/controller"
	"github.com/openshift/library-go/pkg/operator/events"
	"github.com/openshift/library-go/pkg/operator/loglevel"
	"github.com/openshift/library-go/pkg/operator/management"
	"github.com/openshift/library-go/pkg/operator/resource/resourceapply"
	"github.com/openshift/library-go/pkg/operator/staticresourcecontroller"
	"github.com/openshift/library-go/pkg/operator/v1helpers"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
)

type CSIDriverControllerOptions struct {
	dynamicClient       dynamic.Interface
	controllerManifest  string
	nodeManifest        string
	credentialsManifest string
}

type CSIDriverControllerOption func(*CSIDriverControllerOptions)

func WithControllerService(file string) CSIDriverControllerOption {
	return func(o *CSIDriverControllerOptions) {
		o.controllerManifest = file
	}
}

func WithNodeService(file string) CSIDriverControllerOption {
	return func(o *CSIDriverControllerOptions) {
		o.nodeManifest = file
	}
}

func WithCloudCredentials(dynamicClient dynamic.Interface, file string) CSIDriverControllerOption {
	return func(o *CSIDriverControllerOptions) {
		o.dynamicClient = dynamicClient
		o.credentialsManifest = file
	}
}

type ControllerSet struct {
	logLevelController        factory.Controller
	managementStateController factory.Controller
	staticResourcesController factory.Controller
	csiDriverController       *controller.Controller

	operatorClient v1helpers.OperatorClient
	eventRecorder  events.Recorder
}

func (c *ControllerSet) Run(ctx context.Context, workers int) {
	controllers := []factory.Controller{
		c.logLevelController,
		c.managementStateController,
		c.staticResourcesController,
	}

	for i := range controllers {
		if controllers[i] != nil {
			go controllers[i].Run(ctx, workers)
		}
	}

	go c.csiDriverController.Run(ctx, workers)
}

func (c *ControllerSet) WithLogLevelController() *ControllerSet {
	c.logLevelController = loglevel.NewClusterOperatorLoggingController(c.operatorClient, c.eventRecorder)
	return c
}

func (c *ControllerSet) WithManagementStateController(operandName string, supportsOperandRemoval bool) *ControllerSet {
	c.managementStateController = management.NewOperatorManagementStateController(operandName, c.operatorClient, c.eventRecorder)

	if supportsOperandRemoval {
		management.SetOperatorNotRemovable()
	}

	return c
}

func (c *ControllerSet) WithStaticResourcesController(
	name string,
	kubeClient kubernetes.Interface,
	kubeInformersForNamespace v1helpers.KubeInformersForNamespaces,
	manifests resourceapply.AssetFunc,
	files []string,
) *ControllerSet {
	c.staticResourcesController = staticresourcecontroller.NewStaticResourceController(
		name,
		manifests,
		files,
		(&resourceapply.ClientHolder{}).WithKubernetes(kubeClient),
		c.operatorClient,
		c.eventRecorder,
	).AddKubeInformers(kubeInformersForNamespace)

	return c
}

func (c *ControllerSet) WithCSIDriverController(
	name string,
	csiDriverName string,
	csiDriverNamespace string,
	assetFunc func(string) []byte,
	kubeClient kubernetes.Interface,
	namespacedInformerFactory informers.SharedInformerFactory,
	setters ...CSIDriverControllerOption,
) *ControllerSet {
	cdc := controller.New(
		name,
		csiDriverName,
		csiDriverNamespace,
		c.operatorClient,
		assetFunc,
		kubeClient,
		c.eventRecorder,
	)

	opts := &CSIDriverControllerOptions{}
	for _, setter := range setters {
		setter(opts)
	}

	if opts.controllerManifest != "" {
		cdc = cdc.WithControllerService(namespacedInformerFactory.Apps().V1().Deployments(), opts.controllerManifest)
	}

	if opts.nodeManifest != "" {
		cdc = cdc.WithNodeService(namespacedInformerFactory.Apps().V1().DaemonSets(), opts.nodeManifest)
	}

	if opts.credentialsManifest != "" {
		cdc = cdc.WithCloudCredentials(opts.dynamicClient, opts.credentialsManifest)

	}

	c.csiDriverController = cdc

	return c
}

func New(operatorClient v1helpers.OperatorClient, eventRecorder events.Recorder) *ControllerSet {
	return &ControllerSet{
		operatorClient: operatorClient,
		eventRecorder:  eventRecorder,
	}
}
