SHELL :=/bin/bash

all: build
.PHONY: all

# Include the library makefile
include $(addprefix ./vendor/github.com/openshift/build-machinery-go/make/, \
	golang.mk \
	targets/openshift/deps-gomod.mk \
	targets/openshift/images.mk \
	targets/openshift/bindata.mk \
	targets/openshift/codegen.mk \
)

# Codegen module needs setting these required variables
CODEGEN_OUTPUT_PACKAGE :=github.com/openshift/azure-disk-csi-driver-operator/pkg/generated
CODEGEN_API_PACKAGE :=github.com/openshift/azure-disk-csi-driver-operator/pkg/apis
CODEGEN_GROUPS_VERSION :=operator:v1alpha1

define run-codegen
	"$(SHELL)" \
        "$(CODEGEN_PKG)/generate-groups.sh" \
	"$(CODEGEN_GENERATORS)" \
	"$(CODEGEN_OUTPUT_PACKAGE)" \
	"$(CODEGEN_API_PACKAGE)" \
	"$(CODEGEN_GROUPS_VERSION)" \
    --output-base $(CODEGEN_OUTPUT_BASE) \
    --go-header-file $(CODEGEN_GO_HEADER_FILE) \
    $1
endef

# Run core verification and all self contained tests.
#
# Example:
#   make check
check: | verify test-unit
.PHONY: check

IMAGE_REGISTRY?=registry.svc.ci.openshift.org

# This will call a macro called "build-image" which will generate image specific targets based on the parameters:
# $0 - macro name
# $1 - target name
# $2 - image ref
# $3 - Dockerfile path
# $4 - context directory for image build
# It will generate target "image-$(1)" for building the image and binding it as a prerequisite to target "images".
$(call build-image,azure-disk-csi-driver-operator,$(IMAGE_REGISTRY)/ocp/4.5:azure-disk-csi-driver-operator,./Dockerfile.rhel7,.)

# generate bindata targets
# $0 - macro name
# $1 - target suffix
# $2 - input dirs
# $3 - prefix
# $4 - pkg
# $5 - output
$(call add-bindata,generated,./assets/...,assets,generated,pkg/generated/bindata.go)

clean:
	$(RM) azure-disk-csi-driver-operator
.PHONY: clean

GO_TEST_PACKAGES :=./pkg/... ./cmd/...

# Run e2e tests. Requires openshift-tests in $PATH.
#
# Example:
#   make test-e2e
test-e2e:
	hack/e2e.sh

.PHONY: test-e2e
