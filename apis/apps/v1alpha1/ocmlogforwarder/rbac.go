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

package ocmlogforwarder

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/nukleros/operator-builder-tools/pkg/controller/workload"

	appsv1alpha1 "github.com/scottd018/ocm-log-forwarder-operator/apis/apps/v1alpha1"
	"github.com/scottd018/ocm-log-forwarder-operator/apis/apps/v1alpha1/ocmlogforwarder/mutate"
)

// +kubebuilder:rbac:groups=core,resources=serviceaccounts,verbs=get;list;watch;create;update;patch;delete

// CreateServiceAccountParentName creates the ServiceAccount resource with name parent.Name.
func CreateServiceAccountParentName(
	parent *appsv1alpha1.OCMLogForwarder,
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
			},
		},
	}

	resourceObj.SetNamespace(parent.Namespace)

	return mutate.MutateServiceAccountParentName(resourceObj, parent, reconciler, req)
}

// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=roles,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=secrets,verbs=get;watch;list

// CreateRoleParentNameOcm creates the Role resource with name parent.name + -ocm.
func CreateRoleParentNameOcm(
	parent *appsv1alpha1.OCMLogForwarder,
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
						// controlled by field: ocm.secretRef
						parent.Spec.Ocm.SecretRef,
					},
				},
			},
		},
	}

	resourceObj.SetNamespace(parent.Namespace)

	return mutate.MutateRoleParentNameOcm(resourceObj, parent, reconciler, req)
}

// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=roles,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=secrets,verbs=get;watch;list

// CreateRoleParentNameElastic creates the Role resource with name parent.name + -elastic.
func CreateRoleParentNameElastic(
	parent *appsv1alpha1.OCMLogForwarder,
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
						// controlled by field: backend.elasticsearch.secretRef
						parent.Spec.Backend.Elasticsearch.SecretRef,
					},
				},
			},
		},
	}

	resourceObj.SetNamespace(parent.Namespace)

	return mutate.MutateRoleParentNameElastic(resourceObj, parent, reconciler, req)
}

// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=rolebindings,verbs=get;list;watch;create;update;patch;delete

// CreateRoleBindingParentNameOcm creates the RoleBinding resource with name parent.name + -ocm.
func CreateRoleBindingParentNameOcm(
	parent *appsv1alpha1.OCMLogForwarder,
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
					// controlled by field:
					"name": parent.Name,
				},
			},
		},
	}

	resourceObj.SetNamespace(parent.Namespace)

	return mutate.MutateRoleBindingParentNameOcm(resourceObj, parent, reconciler, req)
}

// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=rolebindings,verbs=get;list;watch;create;update;patch;delete

// CreateRoleBindingParentNameElastic creates the RoleBinding resource with name parent.name + -elastic.
func CreateRoleBindingParentNameElastic(
	parent *appsv1alpha1.OCMLogForwarder,
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
					// controlled by field:
					"name": parent.Name,
				},
			},
		},
	}

	resourceObj.SetNamespace(parent.Namespace)

	return mutate.MutateRoleBindingParentNameElastic(resourceObj, parent, reconciler, req)
}
