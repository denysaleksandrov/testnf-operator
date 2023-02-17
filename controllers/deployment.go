package controllers

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	crdv1alpha1 "testnf-operator/api/v1alpha1"
)

func labels(v *crdv1alpha1.Testnf, tier string) map[string]string {
	// Fetches and sets labels

	return map[string]string{
		"app":  "testnf",
		"tier": tier,
	}
}

// ensureDeployment ensures Deployment resource presence in given namespace.
func (r *TestnfReconciler) ensureDeployment(request reconcile.Request,
	instance *crdv1alpha1.Testnf,
	dep *appsv1.Deployment,
	ctx context.Context,
) (*reconcile.Result, error) {

	// See if deployment already exists and create if it doesn't
	found := &appsv1.Deployment{}
	err := r.Get(ctx, types.NamespacedName{
		Name:      dep.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the deployment
		err = r.Create(ctx, dep)

		if err != nil {
			// Deployment failed
			return &reconcile.Result{}, err
		} else {
			// Deployment was successful
			return nil, nil
		}
	} else if err != nil {
		// Error that isn't due to the deployment not existing
		return &reconcile.Result{}, err
	}

	return nil, nil
}

// backendDeployment is a code for Creating Deployment
func (r *TestnfReconciler) testnfDeployment(instance *crdv1alpha1.Testnf) *appsv1.Deployment {

	labels := labels(instance, "testnf")

	// Define the liveness probe for the container
	livenessProbe := &corev1.Probe{
		ProbeHandler: corev1.ProbeHandler{
			HTTPGet: &corev1.HTTPGetAction{
				Path:   "/",
				Port:   intstr.FromInt(80),
				Scheme: corev1.URISchemeHTTP,
			},
		},
		PeriodSeconds:    int32(10),
		SuccessThreshold: int32(1),
		TimeoutSeconds:   int32(1),
	}

	// Define the security context for the container
	securityContext := &corev1.SecurityContext{
		Privileged: newTrue(),
	}

	// Define the volume mounts for the container
	volumeMounts := []corev1.VolumeMount{
		{
			Name:      "huge",
			MountPath: "/mnt/huge",
		},
		{
			Name:      "sys",
			MountPath: "/sys",
		},
		{
			Name:      "dev",
			MountPath: "/dev",
		},
		{
			Name:      "proc",
			MountPath: "/host/proc",
		},
	}

	// Define the volumes
	volumes := []corev1.Volume{
		{
			Name: "huge",
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{
					Medium: "HugePages",
				},
			},
		},
	}
	volumes = append(volumes, newHostPathVolume("sys"))
	volumes = append(volumes, newHostPathVolume("dev"))
	volumes = append(volumes, newHostPathVolume("proc"))

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name,
			Namespace: instance.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: instance.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      labels,
					Annotations: instance.Spec.Annotations,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: fmt.Sprintf("testnf-du-%s", instance.Spec.TestnfSpec),
					RestartPolicy:      "Always",
					NodeSelector:       instance.Spec.NodeSelectors,
					Volumes:            volumes,
					Containers: []corev1.Container{{
						Image:           generateImage(instance.Spec.Image.Repository, instance.Spec.Image.Tag),
						ImagePullPolicy: corev1.PullPolicy(instance.Spec.Image.PullPolicy),
						Name:            instance.Name,
						Resources:       generateResources(instance.Spec.TestnfSpec),
						Env: []corev1.EnvVar{
							{
								Name: "NODE_NAME",
								ValueFrom: &corev1.EnvVarSource{
									FieldRef: &corev1.ObjectFieldSelector{
										APIVersion: "v1",
										FieldPath:  "spec.nodeName",
									},
								},
							},
						},
						LivenessProbe:   livenessProbe,
						SecurityContext: securityContext,
						VolumeMounts:    volumeMounts,
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: 80,
								Name:          "http",
								Protocol:      "TCP",
							},
							{
								ContainerPort: 5201,
								Name:          "tcp",
								Protocol:      "TCP",
							},
						},
					}},
				},
			},
		},
	}

	controllerutil.SetControllerReference(instance, dep, r.Scheme)
	return dep
}
