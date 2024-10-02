package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	sigscheme "sigs.k8s.io/controller-runtime/pkg/scheme"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

const numThermostats = 4

// ThermostatSpec defines the desired state of Thermostat
type ThermostatSpec struct {
	DesiredTemperature int `json:"desiredTemperature"`
	CurrentTemperature int `json:"currentTemperature"`
}

// Thermostat is the Schema for the thermostats API
type Thermostat struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ThermostatSpec `json:"spec,omitempty"`
}

// ThermostatList contains a list of Thermostat
type ThermostatList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Thermostat `json:"items"`
}

func (t *Thermostat) GetObjectKind() schema.ObjectKind { return &t.TypeMeta }
func (t *Thermostat) DeepCopyObject() runtime.Object {
	return &Thermostat{
		TypeMeta:   t.TypeMeta,
		ObjectMeta: *t.ObjectMeta.DeepCopy(),
		Spec:       t.Spec,
	}
}

func (tl *ThermostatList) GetObjectKind() schema.ObjectKind { return &tl.TypeMeta }
func (tl *ThermostatList) DeepCopyObject() runtime.Object {
	return &ThermostatList{
		TypeMeta: tl.TypeMeta,
		ListMeta: *tl.ListMeta.DeepCopy(),
		Items:    append([]Thermostat(nil), tl.Items...),
	}
}

// ThermostatReconciler reconciles a Thermostat Custom Resource
type ThermostatReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *ThermostatReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := ctrl.LoggerFrom(ctx)

	// Fetch the Thermostat instance
	thermostat := &Thermostat{}
	err := r.Get(ctx, req.NamespacedName, thermostat)
	if err != nil {
		if client.IgnoreNotFound(err) == nil {
			log.Info("Thermostat CR not found, assuming it has been deleted")

			// Check if any Thermostat resources exist
			var thermostatList ThermostatList
			if listErr := r.List(ctx, &thermostatList); listErr != nil {
				log.Error(listErr, "unable to list Thermostat resources")
				return ctrl.Result{}, listErr
			}
			if len(thermostatList.Items) == 0 {
				log.Info("No Thermostat resources left, shutting down the reconciler")
				// Trigger a shutdown signal
				go func() {
					log.Info("Shutting down manager")
					// Allow time for logs to flush
					time.Sleep(2 * time.Second)
					os.Exit(0)
				}()
			}

			return ctrl.Result{}, nil
		}
		log.Error(err, "unable to fetch Thermostat")
		return ctrl.Result{}, err
	}

	// Business logic: adjust the current temperature
	if thermostat.Spec.CurrentTemperature < thermostat.Spec.DesiredTemperature {
		thermostat.Spec.CurrentTemperature++
	} else if thermostat.Spec.CurrentTemperature > thermostat.Spec.DesiredTemperature {
		thermostat.Spec.CurrentTemperature--
	} else {
		// todo delete CR
		log.Info(fmt.Sprintf(
			"current temperature is equal to desired temperature, current Temp: %d, desiredTemp: %d",
			thermostat.Spec.CurrentTemperature,
			thermostat.Spec.DesiredTemperature,
		))
		log.Info("deleting Thermostat")
		r.Delete(ctx, thermostat)
		time.Sleep(time.Second)
		return ctrl.Result{}, nil
	}

	// Update the Thermostat status
	log.Info(fmt.Sprintf(
		"updating thermostat with new current temperature, currentTemp: %d, desiredTemp: %d",
		thermostat.Spec.CurrentTemperature,
		thermostat.Spec.DesiredTemperature,
	))
	err = r.Update(ctx, thermostat)
	if err != nil {
		log.Error(err, "unable to update Thermostat")
		return ctrl.Result{}, err
	}

	time.Sleep(time.Second)
	return ctrl.Result{}, nil
	// return ctrl.Result{RequeueAfter: time.Second * 30}, nil
}

func main() {
	var apiServerURL, kubeconfigPath string
	flag.StringVar(&apiServerURL, "api-server-url", "http://127.0.0.1:20888", "The address of the Kubernetes API server")
	flag.StringVar(&kubeconfigPath, "kube-config-path", "./kubeconfig.yaml", "The path to the kubeconfig file")

	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	config, err := clientcmd.BuildConfigFromFlags(apiServerURL, kubeconfigPath)
	if err != nil {
		setupLog.Error(err, "unable to set up client config")
		os.Exit(1)
	}

	mgr, err := ctrl.NewManager(config, ctrl.Options{
		Scheme: scheme,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if err = (&ThermostatReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create thermostat controller")
		os.Exit(1)
	}

	for i := 1; i <= numThermostats; i++ {
		createCR(mgr)
		time.Sleep(time.Second * 2)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

func (r *ThermostatReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&Thermostat{}).
		Complete(r)
}

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	groupVersion := schema.GroupVersion{Group: "thermostats.example.com", Version: "v1"}
	schemeBuilder := &sigscheme.Builder{GroupVersion: groupVersion}
	schemeBuilder.Register(&Thermostat{}, &ThermostatList{})
	utilruntime.Must(schemeBuilder.AddToScheme(scheme))
}

func createCR(mgr manager.Manager) {
	// Create the Thermostat CR
	name := "my-thermostat" + time.Now().Format("20060102150405")
	desiredTemp := 70
	// random int from 50 to 70
	currentTemp := rand.Intn(41) + 50

	setupLog.Info(fmt.Sprintf("Create the Thermostat CR '%s' with Desired Temp %d and Current Temp %d", name, desiredTemp, currentTemp))

	thermostat := &Thermostat{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Thermostat",
			APIVersion: "thermostats.example.com/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "default",
		},
		Spec: ThermostatSpec{
			DesiredTemperature: desiredTemp,
			CurrentTemperature: currentTemp,
		},
	}

	if err := mgr.GetClient().Create(context.Background(), thermostat); err != nil {
		setupLog.Error(err, fmt.Sprintf("unable to create Thermostat CR %s", name))
		os.Exit(1)
	}
}
