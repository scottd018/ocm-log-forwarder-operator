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

// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete

// CreateDeploymentForwarderNamespaceParentName creates the Deployment resource with name parent.Name.
func CreateDeploymentForwarderNamespaceParentName(
	parent *appsv1alpha1.OCMLogForwarder,
	collection *appsv1alpha1.OCMLogForwarderConfig,
	reconciler workload.Reconciler,
	req *workload.Request,
) ([]client.Object, error) {

	if collection.Spec.Backend.Type != "elasticsearch" {
		return []client.Object{}, nil
	}

	var resourceObj = &unstructured.Unstructured{
		Object: map[string]interface{}{
			// +operator-builder:resource:collectionField=backend.type,value=elasticsearch,include=true
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata": map[string]interface{}{
				// controlled by field:
				"name": parent.Name,
				// controlled by collection field: forwarderNamespace
				"namespace": collection.Spec.ForwarderNamespace,
				"labels": map[string]interface{}{
					// controlled by field:
					"app.kubernetes.io/name": parent.Name,
				},
			},
			"spec": map[string]interface{}{
				"replicas": 1,
				"selector": map[string]interface{}{
					"matchLabels": map[string]interface{}{
						// controlled by field:
						"app.kubernetes.io/name": parent.Name,
					},
				},
				"template": map[string]interface{}{
					"metadata": map[string]interface{}{
						"labels": map[string]interface{}{
							// controlled by field:
							"app.kubernetes.io/name": parent.Name,
						},
					},
					"spec": map[string]interface{}{
						// controlled by collection field:
						"serviceAccountName": collection.Name,
						"nodeSelector": map[string]interface{}{
							"kubernetes.io/os": "linux",
						},
						"affinity": map[string]interface{}{
							"podAntiAffinity": map[string]interface{}{
								"preferredDuringSchedulingIgnoredDuringExecution": []interface{}{
									map[string]interface{}{
										"weight": 100,
										"podAffinityTerm": map[string]interface{}{
											"topologyKey": "kubernetes.io/hostname",
											"labelSelector": map[string]interface{}{
												"matchExpressions": []interface{}{
													map[string]interface{}{
														"key":      "app.kubernetes.io/name",
														"operator": "In",
														"values": []interface{}{
															// controlled by field:
															parent.Name,
														},
													},
												},
											},
										},
									},
								},
							},
						},
						"containers": []interface{}{
							map[string]interface{}{
								"name": "forwarder",
								// controlled by field: forwarderVersion
								//  OCM Log Forwarder version to use.  Any of the tags from the ocm-log-forwarder GitHub
								//  repo are supported here.
								//
								"image":           "ghcr.io/scottd018/ocm-log-forwarder:" + parent.Spec.ForwarderVersion + "",
								"imagePullPolicy": "Always",
								"env": []interface{}{
									// NOTE: present all config options here.  Use these as environment variables
									//       on the deployment so that changes here result in the app realizing
									//       those changes by restarting the managed pod.
									map[string]interface{}{
										"name": "OCM_CLUSTER_ID",
										// controlled by field: ocm.clusterId
										//  +kubebuilder:validation:Required
										//  Cluster ID of the cluster to forward logs from.  This Cluster ID can be found in the OCM Console
										//  as part of the URL when selecting the cluster.  It shows up in a form such as
										//  '22tgckqk9c2ff3jd8ve62p0i2st14vrq'.
										//
										"value": parent.Spec.Ocm.ClusterId,
									},
									map[string]interface{}{
										"name": "OCM_POLL_INTERVAL_MINUTES",
										// controlled by field: ocm.pollInternalMinutes
										//  +kubebuilder:validation:Minimum=1
										//  +kubebuilder:validation:Maximum=1440
										//  How frequently, in minutes, the controller will poll the OpenShift Cluster Manager console.  Must
										//  be in the range of 1 minute to 1440 minutes (1 day).
										//
										"value": parent.Spec.Ocm.PollInternalMinutes,
									},
									map[string]interface{}{
										"name": "BACKEND_TYPE",
										// controlled by collection field: backend.type
										//  +kubebuilder:validation:Enum=elasticsearch
										//  Backend type where logs are sent and stored.  Only 'elasticsearch' supported at this time.  Requires
										//  backend.elasticSearch.url to be set.
										//
										"value": collection.Spec.Backend.Type,
									},
									map[string]interface{}{
										"name": "BACKEND_ES_URL",
										// controlled by collection field: backend.elasticSearch.url
										//  URL to which to ship logs when using the 'elasticsearch' as a backend in the .spec.backend.type
										//  field of this custom resource.
										//
										"value": collection.Spec.Backend.ElasticSearch.Url,
									},
									map[string]interface{}{
										"name": "BACKEND_ES_AUTH_TYPE",
										// controlled by collection field: backend.elasticSearch.authType
										//  +kubebuilder:validation:Enum=basic
										//  ElasticSearch authentication type to use.  Only 'basic' supported at this time.  Requires
										//  a single key/value pair stored in a secret named 'elastic-auth' which contains the
										//  basic authentication info for the ElasticSearch connection.  Secret must exist within the same namespace
										//  where the OCM Log Collector is deployed to.
										//
										"value": collection.Spec.Backend.ElasticSearch.AuthType,
									},
									map[string]interface{}{
										"name": "BACKEND_ES_INDEX",
										// controlled by collection field: backend.elasticSearch.index
										//  +kubebuilder:validation:MaxLength=128
										//  Index name in ElasticSearch where service logs are sent.  Index name must be 128 characters or less.
										//
										"value": collection.Spec.Backend.ElasticSearch.Index,
									},
									map[string]interface{}{
										"name": "DEBUG",
										// controlled by field: debug
										//  Enable debug logging on the log forwarder.
										//
										"value": parent.Spec.Debug,
									},
									// WARN: do not change these as changing these have affect of conflicting with
									//       RBAC permissions.
									map[string]interface{}{
										"name":  "OCM_SECRET_NAME",
										"value": "ocm-token",
									},
									map[string]interface{}{
										"name": "OCM_SECRET_NAMESPACE",
										"valueFrom": map[string]interface{}{
											"fieldRef": map[string]interface{}{
												"fieldPath": "metadata.namespace",
											},
										},
									},
									map[string]interface{}{
										"name":  "BACKEND_ES_SECRET_NAME",
										"value": "elastic-auth",
									},
									map[string]interface{}{
										"name": "BACKEND_ES_SECRET_NAMESPACE",
										"valueFrom": map[string]interface{}{
											"fieldRef": map[string]interface{}{
												"fieldPath": "metadata.namespace",
											},
										},
									},
								},
								"securityContext": map[string]interface{}{
									"allowPrivilegeEscalation": false,
									"readOnlyRootFilesystem":   true,
									"capabilities": map[string]interface{}{
										"drop": []interface{}{
											"ALL",
										},
									},
									"runAsNonRoot": true,
									"runAsGroup":   0,
									"seccompProfile": map[string]interface{}{
										"type": "RuntimeDefault",
									},
								},
								"resources": map[string]interface{}{
									"requests": map[string]interface{}{
										"cpu":    "25m",
										"memory": "32Mi",
									},
									"limits": map[string]interface{}{
										"cpu":    "50m",
										"memory": "64Mi",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	resourceObj.SetNamespace(parent.Namespace)

	return mutate.MutateDeploymentForwarderNamespaceParentName(resourceObj, parent, collection, reconciler, req)
}
