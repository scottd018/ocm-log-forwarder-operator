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

// CreateDeploymentParentName creates the Deployment resource with name parent.Name.
func CreateDeploymentParentName(
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
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata": map[string]interface{}{
				// controlled by field:
				"name": parent.Name,
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
						// controlled by field:
						"serviceAccountName": parent.Name,
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
								// controlled by field: version
								//  OCM Log Forwarder version to use.  Any of the tags from the ocm-log-forwarder GitHub
								//  repo are supported here.
								//
								"image":           "ghcr.io/scottd018/ocm-log-forwarder:" + parent.Spec.Version + "",
								"imagePullPolicy": "Always",
								"env": []interface{}{
									// NOTE: present all config options here.  Use these as environment variables
									//       on the deployment so that changes here result in the app realizing
									//       those changes by restarting the managed pod.
									map[string]interface{}{
										"name": "OCM_CLUSTER_ID",
										// controlled by field: ocm.clusterId
										//  Cluster ID of the cluster to forward logs from.  This Cluster ID can be found in the OCM Console
										//  as part of the URL when selecting the cluster.  It shows up in a form such as
										//  '22tgckqk9c2ff3jd8ve62p0i2st14vrq'.
										//
										"value": parent.Spec.Ocm.ClusterId,
									},
									map[string]interface{}{
										"name": "OCM_SECRET_NAME",
										// controlled by field: ocm.secretRef
										//  The secret should contain the OCM JSON token obtained from OpenShift Cluster Manager.  It should
										//  have a single key/value pair with the form of clusterId=ocmTokenJson.  The clusterId
										//  should match the .spec.ocm.clusterId field, while the ocmTokenJson value should be a
										//  string form of the token obtained from OCM.
										//
										"value": parent.Spec.Ocm.SecretRef,
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
										"name": "OCM_POLL_INTERVAL_MINUTES",
										// controlled by field: ocm.pollInternalMinutes
										//  +kubebuilder:validation:Minimum=1
										//  +kubebuilder:validation:Maximum=1440
										//  How frequently, in minutes, the controller will poll the OpenShift Cluster Manager console for service logs.  Must
										//  be in the range of 1 minute to 1440 minutes (1 day).
										//
										"value": parent.Spec.Ocm.PollInternalMinutes,
									},
									map[string]interface{}{
										"name": "BACKEND_TYPE",
										// controlled by field: backend.type
										//  +kubebuilder:validation:Enum=elasticsearch
										//  Backend type where logs are sent and stored.  Only 'elasticsearch' supported at this time.  Requires
										//  backend.elasticSearch.url to be set.
										//
										"value": parent.Spec.Backend.Type,
									},
									map[string]interface{}{
										"name": "BACKEND_ES_URL",
										// controlled by field: backend.elasticSearch.url
										//  URL to which to ship logs when using the 'elasticsearch' as a backend in the .spec.backend.type
										//  field of this custom resource.
										//
										"value": parent.Spec.Backend.ElasticSearch.Url,
									},
									map[string]interface{}{
										"name": "BACKEND_ES_AUTH_TYPE",
										// controlled by field: backend.elasticSearch.authType
										//  +kubebuilder:validation:Enum=basic
										//  ElasticSearch authentication type to use.  Only 'basic' supported at this time.
										//
										//  * 'basic': For 'basic' authentication, the secret from .spec.backend.elasticSearch.secretRef should contain the
										//  basic authentication info for the ElasticSearch connection containing only a single key/value pair with
										//  the key as the username and the value as the password.
										//
										"value": parent.Spec.Backend.ElasticSearch.AuthType,
									},
									map[string]interface{}{
										"name": "BACKEND_ES_INDEX",
										// controlled by field: backend.elasticSearch.index
										//  +kubebuilder:validation:MaxLength=128
										//  Index name in ElasticSearch where service logs are sent.  Index name must be 128 characters or less.
										//
										"value": parent.Spec.Backend.ElasticSearch.Index,
									},
									map[string]interface{}{
										"name": "DEBUG",
										// controlled by field: debug
										//  Enable debug logging on the log forwarder.
										//
										"value": parent.Spec.Debug,
									},
									map[string]interface{}{
										"name": "BACKEND_ES_SECRET_NAME",
										// controlled by field: backend.elasticSearch.secretRef
										//  The secret should contain the authentication information for the ElasticSearch connection.  See
										//  .spec.backend.elasticSearch.authType for more information on secret requirements.  This secret
										//  should exist in the same namespace as the OCMLogForwarder resource.
										//
										"value": parent.Spec.Backend.ElasticSearch.SecretRef,
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

	return mutate.MutateDeploymentParentName(resourceObj, parent, reconciler, req)
}
