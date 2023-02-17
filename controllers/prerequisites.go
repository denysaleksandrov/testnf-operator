package controllers

import (
	"context"
	"fmt"
	crdv1alpha1 "testnf-operator/api/v1alpha1"

	"github.com/go-logr/logr"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
)

// create common resources
func (r *TestnfReconciler) createCommonResources(log logr.Logger, instance *crdv1alpha1.Testnf, ctx context.Context) error {
	// Create ClusterRoleBinding for all the Testnf resources.
	var err error

	clusterRoleName := fmt.Sprintf("testnf-du-%s", instance.Spec.TestnfSpec)
	crb := clusterRoleBindingForTestnf(instance, r.Scheme)

	err = r.Get(ctx, types.NamespacedName{Name: clusterRoleName, Namespace: v1.NamespaceAll}, crb)
	if err != nil && errors.IsNotFound(err) {
		log.Info("no previous ClusterRoleBinding found, creating a new one.")
		err = r.Create(ctx, crb)
	}

	if err != nil {
		return fmt.Errorf("error creating ClusterRoleBinding: %w", err)
	}

	// nad := nadForTestnfController(instance, r.Scheme)
	// fmt.Printf("%+v\n", r.Scheme)
	// err = r.Get(ctx, types.NamespacedName{Name: nad.Name, Namespace: instance.Namespace}, nad)
	// if err != nil && errors.IsNotFound(err) {
	// 	log.Info(fmt.Sprintf("no previous NetAttachDef %s found, creating a new one.", nad.Name))
	// 	err = r.Create(ctx, nad)
	// }

	// if err != nil {
	// 	return fmt.Errorf("error creating NetAttachDef: %w", err)
	// }

	return nil
}

// checkPrerequisites creates all necessary objects before the deployment
func (r *TestnfReconciler) checkPrerequisites(log logr.Logger, instance *crdv1alpha1.Testnf, ctx context.Context) error {
	sa, err := serviceAccountForTestnf(instance, r.Scheme)
	if err != nil {
		return err
	}
	existed, err := r.createIfNotExists(sa)
	if err != nil {
		return err
	}

	if !existed {
		log.Info("ServiceAccount created", "ServiceAccount.Namespace", sa.Namespace, "ServiceAccount.Name", sa.Name)
	}

	// Assign this new ServiceAccount to the ClusterRoleBinding (if is not present already)
	clusterRoleName := fmt.Sprintf("testnf-du-%s", instance.Spec.TestnfSpec)
	crb := clusterRoleBindingForTestnf(instance, r.Scheme)

	err = r.Get(ctx, types.NamespacedName{Name: clusterRoleName, Namespace: v1.NamespaceAll}, crb)
	if err != nil {
		return err
	}

	subject := subjectForServiceAccount(sa.Namespace, sa.Name)
	found := false
	for _, s := range crb.Subjects {
		if s.Name == subject.Name && s.Namespace == subject.Namespace {
			found = true
			break
		}
	}

	if !found {
		crb.Subjects = append(crb.Subjects, subject)

		err = r.Update(ctx, crb)
		if err != nil {
			return err
		}
	}

	return nil
}
