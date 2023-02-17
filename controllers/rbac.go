package controllers

import (
	"fmt"
	crdv1alpha1 "testnf-operator/api/v1alpha1"

	rbacv1 "k8s.io/api/rbac/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func clusterRoleBindingForTestnf(instance *crdv1alpha1.Testnf, scheme *runtime.Scheme) *rbacv1.ClusterRoleBinding {
	sa := fmt.Sprintf("testnf-du-%s", instance.Spec.TestnfSpec)
	crb := &rbacv1.ClusterRoleBinding{
		ObjectMeta: v1.ObjectMeta{
			Name: sa,
		},
		RoleRef: rbacv1.RoleRef{
			Kind:     "ClusterRole",
			Name:     "cluster-admin",
			APIGroup: "rbac.authorization.k8s.io",
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      sa,
				Namespace: instance.Namespace,
			},
		},
	}
	controllerutil.SetControllerReference(instance, crb, scheme)

	return crb
}

func subjectForServiceAccount(namespace, name string) rbacv1.Subject {
	sa := rbacv1.Subject{
		Kind:      "ServiceAccount",
		Name:      name,
		Namespace: namespace,
	}
	return sa
}
