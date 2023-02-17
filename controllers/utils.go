package controllers

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

// NewHostPathVolume create a HostPath volume object
func newHostPathVolume(path string) corev1.Volume {
	v := corev1.Volume{
		Name: path,
		VolumeSource: corev1.VolumeSource{
			HostPath: &corev1.HostPathVolumeSource{
				Path: fmt.Sprintf("/%s", path),
			},
		},
	}

	hpType := corev1.HostPathUnset

	v.VolumeSource.HostPath.Type = &hpType

	return v
}

// Utility function to iterate over pods and return the names slice
func getPodNames(pods []corev1.Pod) []string {
	var podNames []string
	for _, pod := range pods {
		podNames = append(podNames, pod.Name)
	}
	return podNames
}

func generateResources(testnfspec string) corev1.ResourceRequirements {
	cpuRequest := "12"
	var memRequest string
	switch {
	case testnfspec == "large":
		memRequest = "32Gi"
	case testnfspec == "lite":
		memRequest = "16Gi"
	}

	requests := corev1.ResourceRequirements{
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:                           resource.MustParse(cpuRequest),
			corev1.ResourceMemory:                        resource.MustParse(memRequest),
			corev1.ResourceName("hugepages-1Gi"):         resource.MustParse("16Gi"),
			corev1.ResourceName("intel.com/sriovigbuio"): resource.MustParse("2"),
			corev1.ResourceName("intel.com/sriovpass"):   resource.MustParse("1"),
			corev1.ResourceName("intel.com/sriovvfio"):   resource.MustParse("1"),
		},
		Limits: corev1.ResourceList{
			corev1.ResourceMemory:                        resource.MustParse(memRequest),
			corev1.ResourceName("hugepages-1Gi"):         resource.MustParse("16Gi"),
			corev1.ResourceName("intel.com/sriovigbuio"): resource.MustParse("2"),
			corev1.ResourceName("intel.com/sriovpass"):   resource.MustParse("1"),
			corev1.ResourceName("intel.com/sriovvfio"):   resource.MustParse("1"),
		},
	}

	return requests
}

func generateImage(repository string, tag string) string {
	return fmt.Sprintf("%v:%v", repository, tag)
}

func newTrue() *bool {
	var t bool
	z := &t
	*z = true
	return z
}
