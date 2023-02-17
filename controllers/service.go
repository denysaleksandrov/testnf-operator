package controllers

import (
	"context"
	crdv1alpha1 "testnf-operator/api/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// ensureService ensures Service is Running in a namespace.
func (r *TestnfReconciler) ensureService(request reconcile.Request, instance *crdv1alpha1.Testnf, service *corev1.Service, ctx context.Context) (*reconcile.Result, error) {

	// See if service already exists and create if it doesn't
	found := &corev1.Service{}
	err := r.Get(ctx, types.NamespacedName{Name: service.Name, Namespace: instance.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		// Create the service
		err = r.Create(ctx, service)

		if err != nil {
			// Service creation failed
			return &reconcile.Result{}, err
		} else {
			// Service creation was successful
			return nil, nil
		}
	} else if err != nil {
		// Error that isn't due to the service not existing
		return &reconcile.Result{}, err
	}

	return nil, nil
}

// testnfService is a code for creating a Service
func (r *TestnfReconciler) testnfService(v *crdv1alpha1.Testnf) *corev1.Service {
	labels := labels(v, "testnf")

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testnf-service",
			Namespace: v.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{{
				Name:       "http",
				Protocol:   corev1.ProtocolTCP,
				Port:       80,
				TargetPort: intstr.FromInt(80),
			},
				{
					Name:       "tcp",
					Protocol:   corev1.ProtocolTCP,
					Port:       5201,
					TargetPort: intstr.FromInt(5201),
				}},
			Type: corev1.ServiceTypeClusterIP,
		},
	}

	controllerutil.SetControllerReference(v, service, r.Scheme)
	return service
}
