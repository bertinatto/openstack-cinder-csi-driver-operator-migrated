module github.com/openshift/azure-disk-csi-driver-operator

go 1.13

require (
	github.com/evanphx/json-patch v4.5.0+incompatible // indirect
	github.com/go-bindata/go-bindata v3.1.2+incompatible
	github.com/googleapis/gnostic v0.3.1 // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/openshift/api v0.0.0-20200521101457-60c476765272
	github.com/openshift/build-machinery-go v0.0.0-20200424080330-082bf86082cc
	github.com/openshift/library-go v0.0.0
	github.com/prometheus/client_golang v1.4.1
	github.com/spf13/cobra v0.0.6
	github.com/spf13/pflag v1.0.5
	google.golang.org/genproto v0.0.0-20191220175831-5c49e3ecc1c1 // indirect
	k8s.io/apimachinery v0.18.3
	k8s.io/client-go v0.18.3
	k8s.io/code-generator v0.18.3
	k8s.io/component-base v0.18.3
	k8s.io/klog v1.0.0
)

replace github.com/openshift/library-go => ../library-go
