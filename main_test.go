package main

import (
	"context"
	"errors"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestThermostatListDeepCopyObject(t *testing.T) {
	remaining := int64(1)
	original := &ThermostatList{
		ListMeta: metav1.ListMeta{RemainingItemCount: &remaining},
		Items: []Thermostat{{
			ObjectMeta: metav1.ObjectMeta{
				Labels:      map[string]string{"room": "kitchen"},
				Annotations: map[string]string{"note": "original"},
				Finalizers:  []string{"thermostats.example.com/cleanup"},
			},
			Spec: ThermostatSpec{DesiredTemperature: 70},
		}},
	}

	copied := original.DeepCopyObject().(*ThermostatList)
	copied.Items[0].Labels["room"] = "bedroom"
	copied.Items[0].Annotations["note"] = "copy"
	copied.Items[0].Finalizers[0] = "other"
	copied.Items[0].Spec.DesiredTemperature = 80
	*copied.RemainingItemCount = 2

	if original.Items[0].Labels["room"] != "kitchen" ||
		original.Items[0].Annotations["note"] != "original" ||
		original.Items[0].Finalizers[0] != "thermostats.example.com/cleanup" ||
		original.Items[0].Spec.DesiredTemperature != 70 ||
		*original.RemainingItemCount != 1 {
		t.Fatal("mutating the copied list changed the original")
	}
}

func TestThermostatListDeepCopyPreservesNilItems(t *testing.T) {
	copied := (&ThermostatList{}).DeepCopyObject().(*ThermostatList)
	if copied.Items != nil {
		t.Fatal("copy of a nil Items slice must remain nil")
	}
}

type deleteFailingClient struct {
	client.Client
	deleteErr error
}

func (c *deleteFailingClient) Get(_ context.Context, _ types.NamespacedName, obj client.Object, _ ...client.GetOption) error {
	obj.(*Thermostat).Spec = ThermostatSpec{DesiredTemperature: 70, CurrentTemperature: 70}
	return nil
}

func (c *deleteFailingClient) Delete(context.Context, client.Object, ...client.DeleteOption) error {
	return c.deleteErr
}

func TestReconcileReturnsDeleteError(t *testing.T) {
	deleteErr := errors.New("delete failed")
	reconciler := &ThermostatReconciler{Client: &deleteFailingClient{deleteErr: deleteErr}}

	_, err := reconciler.Reconcile(context.Background(), ctrl.Request{})
	if !errors.Is(err, deleteErr) {
		t.Fatalf("expected deletion error, got %v", err)
	}
}
