package always200

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	examplev1alpha1 "always200/pkg/apis/example/v1alpha1"

	routev1 "github.com/openshift/api/route/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_always200")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Always200 Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileAlways200{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("always200-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Always200
	err = c.Watch(&source.Kind{Type: &examplev1alpha1.Always200{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner Always200
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &examplev1alpha1.Always200{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileAlways200 implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileAlways200{}

// ReconcileAlways200 reconciles a Always200 object
type ReconcileAlways200 struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Always200 object and makes changes based on the state read
// and what is in the Always200.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileAlways200) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Always200")

	// Fetch the Always200 instance
	instance := &examplev1alpha1.Always200{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	deployment, err := r.checkCreateDeployment (request, instance, r.always200Deployment(instance))
	if err != nil{
		return reconcile.Result{}, err
	}

	// If the spec.size in the CR changes, update the deployment number of replicas
	if deployment.Spec.Replicas != &instance.Spec.Size {
		controllerutil.CreateOrUpdate(context.TODO(), r.client, deployment, func() error {
			deployment.Spec.Replicas = &instance.Spec.Size
			return nil
		} )
	}

	err = r.checkCreateService(request, instance, r.always200Service(instance))
	if err != nil{
		return reconcile.Result{}, err
	}

	err = r.checkCreateRoute(request, instance, r.always200Route(instance))
	if err != nil{
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

// don't need this for such a small deployment
// never underestimate a programmers need to complicate simple things
func labels(cr *examplev1alpha1.Always200, tier string) map[string]string {
	// Fetches and sets labels

	return map[string]string{
		"app":             "always200",
		"always200_cr": cr.Name,
		"tier":            tier,
	}
}

// This is the equivalent of creating a deployment yaml and returning it
// It doesnt create anything on cluster
func (r *ReconcileAlways200) always200Deployment(cr *examplev1alpha1.Always200) (*appsv1.Deployment){
	// Build a Deployment
	labels := labels(cr, "backend-always200")
    size := cr.Spec.Size
	always200Deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "always200",
			Namespace: cr.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &size,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image: cr.Spec.Image,
						ImagePullPolicy: corev1.PullAlways,
						Name: "always200-pod",
						Ports: []corev1.ContainerPort{{
							ContainerPort: 8080,
							Name: "always200",
						}},
					}},
				},
			},
		},
	}

	// sets the this controller as owner
	controllerutil.SetControllerReference(cr, always200Deployment, r.scheme)
	return always200Deployment
}

// This is the equivalent of creating a service yaml and returning it
// It doesnt create anything on cluster
func (r ReconcileAlways200) always200Service(cr *examplev1alpha1.Always200) *corev1.Service {
	labels := labels(cr, "backend-always200")

	always200Service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:                       "always200-service",
			Namespace:                  cr.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{{
				Protocol: corev1.ProtocolTCP,
				Port: 8080,
				TargetPort: intstr.FromInt(8080),
			}},
		},
	}

	controllerutil.SetControllerReference(cr,always200Service,r.scheme)
	return always200Service
}

// This is the equivalent of creating a route yaml file and returning it
// It doesn't create anything on cluster
func (r ReconcileAlways200) always200Route(cr *examplev1alpha1.Always200) *routev1.Route {
	labels := labels(cr, "backend-always200")

	always200Route := &routev1.Route{
		ObjectMeta: metav1.ObjectMeta{
			Name: "always200",
			Namespace: cr.Namespace,
			Labels: labels,
		},
		Spec: routev1.RouteSpec{
			To: routev1.RouteTargetReference{
				Kind:   "Service",
				Name:   "always200-service",
			},
			Port: &routev1.RoutePort{
				TargetPort: intstr.FromInt(8080),
			},


		},
	}

	return always200Route
}


// check for a deployment if it doesn't exist it creates one on cluster using the deployment created in always200Deployment
func (r ReconcileAlways200) checkCreateDeployment(request reconcile.Request, instance *examplev1alpha1.Always200, always200Deployment *appsv1.Deployment) (*appsv1.Deployment,error){
    // check for a deployment in the namespace
	found := &appsv1.Deployment{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name:always200Deployment.Name,Namespace: instance.Namespace}, found)
	if err != nil {
		log.Info("Creating deployment")
		err = r.client.Create(context.TODO(), always200Deployment)
		if err !=nil{
			log.Error(err, "Failed to create deployment")
			return found, err
		} 
	}
	return found ,nil
}

// check for a service if it doesn't exist it creates one on cluster using the service created in always200Service
func (r ReconcileAlways200) checkCreateService(request reconcile.Request, instance *examplev1alpha1.Always200, always200Servcie *corev1.Service) error{
	// check for a deployment in the namespace
	found := &corev1.Service{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name:always200Servcie.Name,Namespace: instance.Namespace}, found)
	if err != nil {
		log.Info("Creating Service")
		err = r.client.Create(context.TODO(), always200Servcie)
		if err !=nil{
			log.Error(err, "Failed to create Service")
			return err
		} 
	}
	return nil
}

// check for a route if it doesn't exist it creates one on cluster using the route created in always200Route
func (r ReconcileAlways200) checkCreateRoute(request reconcile.Request, instance *examplev1alpha1.Always200, always200Route *routev1.Route) error{
	// check for a deployment in the namespace
	found := &routev1.Route{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name:always200Route.Name,Namespace: instance.Namespace}, found)
	if err != nil {
		log.Info("Creating Route")
		err = r.client.Create(context.TODO(), always200Route)
		if err !=nil{
			log.Error(err, "Failed to create Route")
			return err
		} 
	}
	return nil
}


