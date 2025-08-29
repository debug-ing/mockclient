package mockclient

import (
	"context"
	"errors"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func newTestClient() (*MockClient, *corev1.Pod) {
	scheme := runtime.NewScheme()
	_ = corev1.AddToScheme(scheme)
	realClient := fake.NewClientBuilder().WithScheme(scheme).Build()
	fc := NewMockClient(realClient)
	pod := &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "mypod", Namespace: "default"}}
	return fc, pod
}

func TestMockClient_Handlers(t *testing.T) {
	tests := []struct {
		name       string
		verb       string
		injectErr  error
		expectErr  bool
		expectCall bool
	}{
		{
			name:       "Create returns injected error",
			verb:       "create",
			injectErr:  errors.New("simulated create error"),
			expectErr:  true,
			expectCall: true,
		},
		{
			name:       "Delete handler called successfully",
			verb:       "delete",
			injectErr:  nil,
			expectErr:  false,
			expectCall: true,
		},
		{
			name:       "Update handler returns injected error",
			verb:       "update",
			injectErr:  errors.New("update failed"),
			expectErr:  true,
			expectCall: true,
		},
		{
			name:       "Get handler succeeds",
			verb:       "get",
			injectErr:  nil,
			expectErr:  false,
			expectCall: true,
		},
		{
			name:       "List handler returns injected error",
			verb:       "list",
			injectErr:  errors.New("list failed"),
			expectErr:  true,
			expectCall: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fc, pod := newTestClient()

			called := false
			fc.Inject(pod, tt.verb, func() error {
				called = true
				return tt.injectErr
			})

			var err error
			switch tt.verb {
			case "create":
				err = fc.Create(context.Background(), pod)
			case "delete":
				err = fc.Delete(context.Background(), pod)
			case "update":
				err = fc.Update(context.Background(), pod)
			case "get":
				err = fc.Get(context.Background(), client.ObjectKey{Name: "mypod", Namespace: "default"}, pod)
			case "list":
				list := &corev1.PodList{}
				err = fc.List(context.Background(), list)
			}

			if tt.expectErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if tt.expectCall && !called {
				t.Fatalf("expected handler to be called")
			}
		})
	}
}

func TestMockClient_FallbackToRealClient(t *testing.T) {
	fc, pod := newTestClient()
	if err := fc.Create(context.Background(), pod); err != nil {
		t.Fatalf("expected no error from realClient, got %v", err)
	}
	got := &corev1.Pod{}
	if err := fc.Get(context.Background(), client.ObjectKey{Name: "mypod", Namespace: "default"}, got); err != nil {
		t.Fatalf("expected pod in realClient, got %v", err)
	}
}
