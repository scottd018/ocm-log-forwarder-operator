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

// +kubebuilder:rbac:groups=core,resources=serviceaccounts,verbs=get;list;watch;create;update;patch;delete

// CreateServiceAccountForwarderNamespaceParentName creates the ServiceAccount resource with name parent.Name.
func CreateServiceAccountForwarderNamespaceParentName(
	parent *appsv1alpha1.OCMLogForwarderConfig,
	reconciler workload.Reconciler,
	req *workload.Request,
) ([]client.Object, error) {

	var resourceObj = &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion":                   "v1",
			"kind":                         "ServiceAccount",
			"automountServiceAccountToken": true,
			"metadata": map[string]interface{}{
				// controlled by field:
				"name": parent.Name,
				// controlled by field: forwarderNamespace
				"namespace": parent.Spec.ForwarderNamespace,
			},
		},
	}

	return mutate.MutateServiceAccountForwarderNamespaceParentName(resourceObj, parent, reconciler, req)
}

// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=roles,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=secrets,verbs=get;watch;list

// CreateRoleForwarderNamespaceParentNameOcm creates the Role resource with name parent.name + -ocm.
func CreateRoleForwarderNamespaceParentNameOcm(
	parent *appsv1alpha1.OCMLogForwarderConfig,
	reconciler workload.Reconciler,
	req *workload.Request,
) ([]client.Object, error) {

	var resourceObj = &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "rbac.authorization.k8s.io/v1",
			"kind":       "Role",
			"metadata": map[string]interface{}{
				// controlled by field:
				"name": "" + parent.Name + "-ocm",
				// controlled by field: forwarderNamespace
				"namespace": parent.Spec.ForwarderNamespace,
			},
			"rules": []interface{}{
				map[string]interface{}{
					"apiGroups": []interface{}{
						"",
					},
					"resources": []interface{}{
						"secrets",
					},
					"verbs": []interface{}{
						"get",
						"watch",
						"list",
					},
					"resourceNames": []interface{}{
						// NOTE: this secret must pre-exist for this operator to work.  The key should contain the
						//       string value of the JSON OCM API key.
						"ocm-token",
					},
				},
			},
		},
	}

	return mutate.MutateRoleForwarderNamespaceParentNameOcm(resourceObj, parent, reconciler, req)
}

// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=roles,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=secrets,verbs=get;watch;list

// CreateRoleForwarderNamespaceParentNameElastic creates the Role resource with name parent.name + -elastic.
func CreateRoleForwarderNamespaceParentNameElastic(
	parent *appsv1alpha1.OCMLogForwarderConfig,
	reconciler workload.Reconciler,
	req *workload.Request,
) ([]client.Object, error) {

	if parent.Spec.Backend.Type != "elasticsearch" {
		return []client.Object{}, nil
	}

	var resourceObj = &unstructured.Unstructured{
		Object: map[string]interface{}{
			// +operator-builder:resource:field=backend.type,value=elasticsearch,include=true
			"apiVersion": "rbac.authorization.k8s.io/v1",
			"kind":       "Role",
			"metadata": map[string]interface{}{
				// controlled by field:
				"name": "" + parent.Name + "-elastic",
				// controlled by field: forwarderNamespace
				"namespace": parent.Spec.ForwarderNamespace,
			},
			"rules": []interface{}{
				map[string]interface{}{
					"apiGroups": []interface{}{
						"",
					},
					"resources": []interface{}{
						"secrets",
					},
					"verbs": []interface{}{
						"get",
						"watch",
						"list",
					},
					"resourceNames": []interface{}{
						// NOTE: this secret must pre-exist for this operator to work.  The key should contain the
						//       string value of the JSON OCM API key.
						"elastic-auth",
					},
				},
			},
		},
	}

	return mutate.MutateRoleForwarderNamespaceParentNameElastic(resourceObj, parent, reconciler, req)
}

// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=rolebindings,verbs=get;list;watch;create;update;patch;delete

// CreateRoleBindingForwarderNamespaceParentNameOcm creates the RoleBinding resource with name parent.name + -ocm.
func CreateRoleBindingForwarderNamespaceParentNameOcm(
	parent *appsv1alpha1.OCMLogForwarderConfig,
	reconciler workload.Reconciler,
	req *workload.Request,
) ([]client.Object, error) {

	var resourceObj = &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "rbac.authorization.k8s.io/v1",
			"kind":       "RoleBinding",
			"metadata": map[string]interface{}{
				// controlled by field:
				"name": "" + parent.Name + "-ocm",
				// controlled by field: forwarderNamespace
				"namespace": parent.Spec.ForwarderNamespace,
			},
			"roleRef": map[string]interface{}{
				"apiGroup": "rbac.authorization.k8s.io",
				"kind":     "Role",
				// controlled by field:
				"name": "" + parent.Name + "-ocm",
			},
			"subjects": []interface{}{
				map[string]interface{}{
					"kind": "ServiceAccount",
					"name": "ocm-log-forwarder",
				},
			},
		},
	}

	return mutate.MutateRoleBindingForwarderNamespaceParentNameOcm(resourceObj, parent, reconciler, req)
}

// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=rolebindings,verbs=get;list;watch;create;update;patch;delete

// CreateRoleBindingForwarderNamespaceParentNameElastic creates the RoleBinding resource with name parent.name + -elastic.
func CreateRoleBindingForwarderNamespaceParentNameElastic(
	parent *appsv1alpha1.OCMLogForwarderConfig,
	reconciler workload.Reconciler,
	req *workload.Request,
) ([]client.Object, error) {

	if parent.Spec.Backend.Type != "elasticsearch" {
		return []client.Object{}, nil
	}

	var resourceObj = &unstructured.Unstructured{
		Object: map[string]interface{}{
			// +operator-builder:resource:field=backend.type,value=elasticsearch,include=true
			"apiVersion": "rbac.authorization.k8s.io/v1",
			"kind":       "RoleBinding",
			"metadata": map[string]interface{}{
				// controlled by field:
				"name": "" + parent.Name + "-elastic",
				// controlled by field: forwarderNamespace
				"namespace": parent.Spec.ForwarderNamespace,
			},
			"roleRef": map[string]interface{}{
				"apiGroup": "rbac.authorization.k8s.io",
				"kind":     "Role",
				// controlled by field:
				"name": "" + parent.Name + "-elastic",
			},
			"subjects": []interface{}{
				map[string]interface{}{
					"kind": "ServiceAccount",
					"name": "ocm-log-forwarder",
				},
			},
		},
	}

	return mutate.MutateRoleBindingForwarderNamespaceParentNameElastic(resourceObj, parent, reconciler, req)
}
