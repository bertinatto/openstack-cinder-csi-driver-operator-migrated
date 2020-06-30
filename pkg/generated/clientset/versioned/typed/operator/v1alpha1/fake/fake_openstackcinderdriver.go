// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1alpha1 "github.com/openshift/openstack-cinder-csi-driver-operator/pkg/apis/operator/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeOpenStackCinderDrivers implements OpenStackCinderDriverInterface
type FakeOpenStackCinderDrivers struct {
	Fake *FakeCsiV1alpha1
}

var openstackcinderdriversResource = schema.GroupVersionResource{Group: "csi.openshift.io", Version: "v1alpha1", Resource: "openstackcinderdrivers"}

var openstackcinderdriversKind = schema.GroupVersionKind{Group: "csi.openshift.io", Version: "v1alpha1", Kind: "OpenStackCinderDriver"}

// Get takes name of the openStackCinderDriver, and returns the corresponding openStackCinderDriver object, and an error if there is any.
func (c *FakeOpenStackCinderDrivers) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.OpenStackCinderDriver, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(openstackcinderdriversResource, name), &v1alpha1.OpenStackCinderDriver{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.OpenStackCinderDriver), err
}

// List takes label and field selectors, and returns the list of OpenStackCinderDrivers that match those selectors.
func (c *FakeOpenStackCinderDrivers) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.OpenStackCinderDriverList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(openstackcinderdriversResource, openstackcinderdriversKind, opts), &v1alpha1.OpenStackCinderDriverList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.OpenStackCinderDriverList{ListMeta: obj.(*v1alpha1.OpenStackCinderDriverList).ListMeta}
	for _, item := range obj.(*v1alpha1.OpenStackCinderDriverList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested openStackCinderDrivers.
func (c *FakeOpenStackCinderDrivers) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(openstackcinderdriversResource, opts))
}

// Create takes the representation of a openStackCinderDriver and creates it.  Returns the server's representation of the openStackCinderDriver, and an error, if there is any.
func (c *FakeOpenStackCinderDrivers) Create(ctx context.Context, openStackCinderDriver *v1alpha1.OpenStackCinderDriver, opts v1.CreateOptions) (result *v1alpha1.OpenStackCinderDriver, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(openstackcinderdriversResource, openStackCinderDriver), &v1alpha1.OpenStackCinderDriver{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.OpenStackCinderDriver), err
}

// Update takes the representation of a openStackCinderDriver and updates it. Returns the server's representation of the openStackCinderDriver, and an error, if there is any.
func (c *FakeOpenStackCinderDrivers) Update(ctx context.Context, openStackCinderDriver *v1alpha1.OpenStackCinderDriver, opts v1.UpdateOptions) (result *v1alpha1.OpenStackCinderDriver, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(openstackcinderdriversResource, openStackCinderDriver), &v1alpha1.OpenStackCinderDriver{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.OpenStackCinderDriver), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeOpenStackCinderDrivers) UpdateStatus(ctx context.Context, openStackCinderDriver *v1alpha1.OpenStackCinderDriver, opts v1.UpdateOptions) (*v1alpha1.OpenStackCinderDriver, error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateSubresourceAction(openstackcinderdriversResource, "status", openStackCinderDriver), &v1alpha1.OpenStackCinderDriver{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.OpenStackCinderDriver), err
}

// Delete takes name of the openStackCinderDriver and deletes it. Returns an error if one occurs.
func (c *FakeOpenStackCinderDrivers) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(openstackcinderdriversResource, name), &v1alpha1.OpenStackCinderDriver{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeOpenStackCinderDrivers) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(openstackcinderdriversResource, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.OpenStackCinderDriverList{})
	return err
}

// Patch applies the patch and returns the patched openStackCinderDriver.
func (c *FakeOpenStackCinderDrivers) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.OpenStackCinderDriver, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(openstackcinderdriversResource, name, pt, data, subresources...), &v1alpha1.OpenStackCinderDriver{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.OpenStackCinderDriver), err
}