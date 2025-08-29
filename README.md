# Mock Kubernetes Client

A lightweight **mock client wrapper** around [controller-runtime’s client](https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/client) that allows you to inject custom handlers for CRUD operations (`Create`, `Update`, `Delete`, `Get`, `List`).  

This is useful for **unit tests** where you want to simulate success/failure scenarios without depending on a real Kubernetes cluster.


## Features

- 🔧 Inject custom behavior (error or success) for specific verbs.
- ✅ Fallback to the real `controller-runtime` fake client if no handler is injected.
- 🧪 Simple API for testing Kubernetes controllers or operators.


## Installation

```bash
go get github.com/debug-ing/mockclient
```

## Usage

### 1. Create a Mock Client
```go
scheme := runtime.NewScheme()
_ = corev1.AddToScheme(scheme)

realClient := fake.NewClientBuilder().WithScheme(scheme).Build()
mc := mockclient.NewMockClient(realClient)
```
### 2. Inject an Error for **Create**
```go
pod := &corev1.Pod{
    ObjectMeta: metav1.ObjectMeta{
        Name:      "mypod",
        Namespace: "default",
    },
}

expectedErr := errors.New("simulated create error")

mc.Inject(pod, "create", func() error {
    return expectedErr
})

err := mc.Create(context.Background(), pod)
if err != nil {
    fmt.Println("Create failed as expected:", err)
}
```

### 3. Inject a Success Handler for **Delete**
```go
called := false
mc.Inject(pod, "delete", func() error {
    called = true
    return nil
})

_ = mc.Delete(context.Background(), pod)
fmt.Println("Delete handler called?", called)
```
### 4. Fallback to Real Client
```go
err := mc.Create(context.Background(), pod)
if err != nil {
    panic(err)
}

got := &corev1.Pod{}
_ = mc.Get(context.Background(), client.ObjectKey{Name: "mypod", Namespace: "default"}, got)
fmt.Println("Pod created in fake client:", got.Name)
```
