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

package ocmlogforwarderconfig

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/nukleros/operator-builder-tools/pkg/controller/workload"

	appsv1alpha1 "github.com/scottd018/ocm-log-forwarder-operator/apis/apps/v1alpha1"
	"github.com/scottd018/ocm-log-forwarder-operator/apis/apps/v1alpha1/ocmlogforwarderconfig/mutate"
)

// +kubebuilder:rbac:groups=core,resources=namespaces,verbs=get;list;watch;create;update;patch;delete

// CreateNamespaceForwarderNamespace creates the Namespace resource with name parent.Spec.ForwarderNamespace.
func CreateNamespaceForwarderNamespace(
	parent *appsv1alpha1.OCMLogForwarderConfig,
	reconciler workload.Reconciler,
	req *workload.Request,
) ([]client.Object, error) {

	var resourceObj = &unstructured.Unstructured{
		Object: map[string]interface{}{
			// controlled by field: provider
			//  +kubebuilder:validation:Enum=rosa
			//  Provider to which the OCM Log Forwarder is provisioned to.  Only 'rosa' supported at this time.
			//
			"apiVersion": "v1",
			"kind":       "Namespace",
			"metadata": map[string]interface{}{
				// controlled by field: forwarderNamespace
				//  Namespace used to provision the OCM Log Forwarder workloads.  Each workload provisoned will
				//  end up in this namespace.
				//
				"name": parent.Spec.ForwarderNamespace,
			},
		},
	}

	return mutate.MutateNamespaceForwarderNamespace(resourceObj, parent, reconciler, req)
}
