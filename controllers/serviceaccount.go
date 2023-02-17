package controllers

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"

	crdv1alpha1 "testnf-operator/api/v1alpha1"
)

func serviceAccountForTestnf(instance *crdv1alpha1.Testnf, scheme *runtime.Scheme) (*corev1.ServiceAccount, error) {
	saName := fmt.Sprintf("testnf-du-%s", instance.Spec.TestnfSpec)
	svca := &corev1.ServiceAccount{
		ObjectMeta: v1.ObjectMeta{
			Name:      saName,
			Namespace: instance.Namespace,
		},
	}

	if err := ctrl.SetControllerReference(instance, svca, scheme); err != nil {
		return nil, err
	}

	return svca, nil
}
