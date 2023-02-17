package controllers

// import (
// 	"encoding/json"
// 	crdv1alpha1 "testnf-operator/api/v1alpha1"

// 	k8scni "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"
// 	ctrl "sigs.k8s.io/controller-runtime"

// 	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// 	"k8s.io/apimachinery/pkg/runtime"
// )

// func nadForTestnfController(instance *crdv1alpha1.Testnf, scheme *runtime.Scheme) *k8scni.NetworkAttachmentDefinition {
// 	// Define the network configuration as a JSON object
// 	plugins := []map[string]interface{}{}
// 	sriov := map[string]interface{}{
// 		"type": "sriov",
// 		"ipam": map[string]interface{}{
// 			"type":      "whereabouts",
// 			"datastore": "kubernetes",
// 			"range":     "192.167.3.0/24",
// 			"kubeconfig": map[string]string{
// 				"kubeconfig": "/etc/cni/net.d/whereabouts.d/whereabouts.kubeconfig",
// 			},
// 			"exclude": []string{
// 				"192.167.3.0/27",
// 			},
// 		},
// 	}
// 	sbr := map[string]interface{}{
// 		"type": "sbr",
// 	}
// 	tunning := map[string]interface{}{
// 		"type": "tuning",
// 		"sysctl": map[string]string{
// 			"net.core.somaxconn": "500",
// 		},
// 		"promisc": true,
// 		"mtu":     1200,
// 	}
// 	plugins = append(plugins, sriov)
// 	plugins = append(plugins, sbr)
// 	plugins = append(plugins, tunning)
// 	config := map[string]interface{}{
// 		"cniVersion": "0.3.1",
// 		"name":       "sriov-pass-network",
// 		"plugins":    plugins,
// 	}

// 	// Convert the configuration to a JSON string
// 	configJSON, err := json.Marshal(config)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// Define the NetworkAttachmentDefinition object
// 	nad := &k8scni.NetworkAttachmentDefinition{
// 		TypeMeta: metav1.TypeMeta{
// 			APIVersion: "k8s.cni.cncf.io/v1",
// 			Kind:       "NetworkAttachmentDefinition",
// 		},
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:        "sriov-pass",
// 			Annotations: instance.Spec.Annotations,
// 			Namespace:   instance.Namespace,
// 		},
// 		Spec: k8scni.NetworkAttachmentDefinitionSpec{
// 			Config: string(configJSON),
// 		},
// 	}

// 	ctrl.SetControllerReference(instance, nad, scheme)
// 	return nad
// }
