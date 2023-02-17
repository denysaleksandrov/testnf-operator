/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// TestnfSpec defines the desired state of Testnf
type TestnfSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// The image of the Ingress Controller.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Image Image `json:"image"`

	// Replicas indicate the replicas to mantain
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=1
	Replicas *int32 `json:"replicas"`

	// Array of node selectors
	// +kubebuilder:validation:Optional
	NodeSelectors map[string]string `json:"nodeSelectors,omitempty"`

	// Specifies extra annotations of the service.
	// +kubebuilder:validation:Optional
	Annotations map[string]string `json:"annotations,omitempty"`

	// Spec defines Testnf size: flex, du, lite
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=large;lite
	// +kubebuilder:default:=large
	TestnfSpec string `json:"spec"`
}

// Image defines the Repository, Tag and ImagePullPolicy of the Controller Image.
type Image struct {
	// The repository of the image.
	Repository string `json:"repository"`
	// The tag (version) of the image.
	Tag string `json:"tag"`
	// The ImagePullPolicy of the image.
	// +kubebuilder:validation:Enum=Never;Always;IfNotPresent
	// +kubebuilder:default:=Always
	PullPolicy string `json:"pullPolicy"`
}

// TestnfStatus defines the observed state of Testnf
type TestnfStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Pods are the name of the Pods hosting the App
	Pods []string `json:"pods"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Testnf is the Schema for the testnfs API
type Testnf struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TestnfSpec   `json:"spec,omitempty"`
	Status TestnfStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// TestnfList contains a list of Testnf
type TestnfList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Testnf `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Testnf{}, &TestnfList{})
}
