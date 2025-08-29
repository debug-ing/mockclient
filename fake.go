package mockclient

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type MockClient struct {
	client.Client
	handlers map[schema.GroupVersionKind]map[string]func() error
}

func NewMockClient(cl client.Client) *MockClient {
	return &MockClient{
		Client:   cl,
		handlers: make(map[schema.GroupVersionKind]map[string]func() error),
	}
}

func (c *MockClient) Inject(obj client.Object, verb string, handler func() error) {
	gvk := obj.GetObjectKind().GroupVersionKind()
	if _, ok := c.handlers[gvk]; !ok {
		c.handlers[gvk] = make(map[string]func() error)
	}
	c.handlers[gvk][verb] = handler
}

func (c *MockClient) Create(ctx context.Context, obj client.Object, opts ...client.CreateOption) error {
	gvk := obj.GetObjectKind().GroupVersionKind()
	if verbs, ok := c.handlers[gvk]; ok {
		if handler, ok := verbs["create"]; ok {
			if err := handler(); err != nil {
				return err
			}
			return nil
		}
	}
	return c.Client.Create(ctx, obj, opts...)
}

func (c *MockClient) Delete(ctx context.Context, obj client.Object, opts ...client.DeleteOption) error {
	gvk := obj.GetObjectKind().GroupVersionKind()
	if verbs, ok := c.handlers[gvk]; ok {
		if handler, ok := verbs["delete"]; ok {
			if err := handler(); err != nil {
				return err
			}
			return nil
		}
	}
	return c.Client.Delete(ctx, obj, opts...)
}

func (c *MockClient) Update(ctx context.Context, obj client.Object, opts ...client.UpdateOption) error {
	gvk := obj.GetObjectKind().GroupVersionKind()
	if verbs, ok := c.handlers[gvk]; ok {
		if handler, ok := verbs["update"]; ok {
			if err := handler(); err != nil {
				return err
			}
			return nil
		}
	}
	return c.Client.Update(ctx, obj, opts...)
}

func (c *MockClient) Get(ctx context.Context, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
	gvk := obj.GetObjectKind().GroupVersionKind()
	if verbs, ok := c.handlers[gvk]; ok {
		if handler, ok := verbs["get"]; ok {
			if err := handler(); err != nil {
				return err
			}
			return nil
		}
	}
	return c.Client.Get(ctx, key, obj, opts...)
}

func (c *MockClient) List(ctx context.Context, list client.ObjectList, opts ...client.ListOption) error {
	gvk := list.GetObjectKind().GroupVersionKind()
	if verbs, ok := c.handlers[gvk]; ok {
		if handler, ok := verbs["list"]; ok {
			if err := handler(); err != nil {
				return err
			}
			return nil
		}
	}
	return c.Client.List(ctx, list, opts...)
}
