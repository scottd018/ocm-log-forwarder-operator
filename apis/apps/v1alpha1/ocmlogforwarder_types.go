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
	"errors"

	"github.com/nukleros/operator-builder-tools/pkg/controller/workload"
	"github.com/nukleros/operator-builder-tools/pkg/status"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var ErrUnableToConvertOCMLogForwarder = errors.New("unable to convert to OCMLogForwarder")

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// OCMLogForwarderSpec defines the desired state of OCMLogForwarder.
type OCMLogForwarderSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:validation:Optional
	Ocm OCMLogForwarderSpecOcm `json:"ocm,omitempty"`

	// +kubebuilder:validation:Optional
	Backend OCMLogForwarderSpecBackend `json:"backend,omitempty"`

	// +kubebuilder:default="latest"
	// +kubebuilder:validation:Optional
	// (Default: "latest")
	//
	//	OCM Log Forwarder version to use.  Any of the tags from the ocm-log-forwarder GitHub
	//	repo are supported here.
	Version string `json:"version,omitempty"`

	// +kubebuilder:default=false
	// +kubebuilder:validation:Optional
	// (Default: false)
	//
	//	Enable debug logging on the log forwarder.
	Debug bool `json:"debug,omitempty"`
}

type OCMLogForwarderSpecOcm struct {
	// +kubebuilder:default="ocm-token"
	// +kubebuilder:validation:Optional
	// (Default: "ocm-token")
	//
	//	The secret should contain the OCM JSON token obtained from OpenShift Cluster Manager.  It should
	//	have a single key/value pair with the form of clusterId=ocmTokenJson.  The clusterId
	//	should match the .spec.ocm.clusterId field, while the ocmTokenJson value should be a
	//	string form of the token obtained from OCM.
	SecretRef string `json:"secretRef,omitempty"`

	// +kubebuilder:validation:Required
	//
	//	Cluster ID of the cluster to forward logs from.  This Cluster ID can be found in the OCM Console
	//	as part of the URL when selecting the cluster.  It shows up in a form such as
	//	'22tgckqk9c2ff3jd8ve62p0i2st14vrq'.
	ClusterId string `json:"clusterId,omitempty"`

	// +kubebuilder:default=5
	// +kubebuilder:validation:Optional
	// (Default: 5)
	//
	//	+kubebuilder:validation:Minimum=1
	//	+kubebuilder:validation:Maximum=1440
	//	How frequently, in minutes, the controller will poll the OpenShift Cluster Manager console for service logs.  Must
	//	be in the range of 1 minute to 1440 minutes (1 day).
	PollInternalMinutes int `json:"pollInternalMinutes,omitempty"`
}

type OCMLogForwarderSpecBackend struct {
	// +kubebuilder:validation:Optional
	Elasticsearch OCMLogForwarderSpecBackendElasticsearch `json:"elasticsearch,omitempty"`

	// +kubebuilder:default="elasticsearch"
	// +kubebuilder:validation:Optional
	// (Default: "elasticsearch")
	//
	//	+kubebuilder:validation:Enum=elasticsearch
	//	Backend type where logs are sent and stored.  Only 'elasticsearch' supported at this time.  Requires
	//	backend.elasticsearch.url to be set.
	Type string `json:"type,omitempty"`
}

type OCMLogForwarderSpecBackendElasticsearch struct {
	// +kubebuilder:default="elastic-auth"
	// +kubebuilder:validation:Optional
	// (Default: "elastic-auth")
	//
	//	The secret should contain the authentication information for the ElasticSearch connection.  See
	//	.spec.backend.elasticsearch.authType for more information on secret requirements.  This secret
	//	should exist in the same namespace as the OCMLogForwarder resource.
	SecretRef string `json:"secretRef,omitempty"`

	// +kubebuilder:default="https://elasticsearch-es-http.elastic-system.svc.cluster.local:9200"
	// +kubebuilder:validation:Optional
	// (Default: "https://elasticsearch-es-http.elastic-system.svc.cluster.local:9200")
	//
	//	URL to which to ship logs when using the 'elasticsearch' as a backend in the .spec.backend.type
	//	field of this custom resource.
	Url string `json:"url,omitempty"`

	// +kubebuilder:default="basic"
	// +kubebuilder:validation:Optional
	// (Default: "basic")
	//
	//	+kubebuilder:validation:Enum=basic
	//	ElasticSearch authentication type to use.  Only 'basic' supported at this time.
	//
	//	* 'basic': For 'basic' authentication, the secret from .spec.backend.elasticsearch.secretRef should contain the
	//	basic authentication information for the ElasticSearch connection containing only a single key/value pair with
	//	the key as the username and the value as the password.
	AuthType string `json:"authType,omitempty"`

	// +kubebuilder:default="ocm_service_logs"
	// +kubebuilder:validation:Optional
	// (Default: "ocm_service_logs")
	//
	//	+kubebuilder:validation:MaxLength=128
	//	Index name in ElasticSearch where service logs are sent.  Index name must be 128 characters or less.
	Index string `json:"index,omitempty"`
}

// OCMLogForwarderStatus defines the observed state of OCMLogForwarder.
type OCMLogForwarderStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Created               bool                     `json:"created,omitempty"`
	DependenciesSatisfied bool                     `json:"dependenciesSatisfied,omitempty"`
	Conditions            []*status.PhaseCondition `json:"conditions,omitempty"`
	Resources             []*status.ChildResource  `json:"resources,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// OCMLogForwarder is the Schema for the ocmlogforwarders API.
type OCMLogForwarder struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              OCMLogForwarderSpec   `json:"spec,omitempty"`
	Status            OCMLogForwarderStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// OCMLogForwarderList contains a list of OCMLogForwarder.
type OCMLogForwarderList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OCMLogForwarder `json:"items"`
}

// interface methods

// GetReadyStatus returns the ready status for a component.
func (component *OCMLogForwarder) GetReadyStatus() bool {
	return component.Status.Created
}

// SetReadyStatus sets the ready status for a component.
func (component *OCMLogForwarder) SetReadyStatus(ready bool) {
	component.Status.Created = ready
}

// GetDependencyStatus returns the dependency status for a component.
func (component *OCMLogForwarder) GetDependencyStatus() bool {
	return component.Status.DependenciesSatisfied
}

// SetDependencyStatus sets the dependency status for a component.
func (component *OCMLogForwarder) SetDependencyStatus(dependencyStatus bool) {
	component.Status.DependenciesSatisfied = dependencyStatus
}

// GetPhaseConditions returns the phase conditions for a component.
func (component *OCMLogForwarder) GetPhaseConditions() []*status.PhaseCondition {
	return component.Status.Conditions
}

// SetPhaseCondition sets the phase conditions for a component.
func (component *OCMLogForwarder) SetPhaseCondition(condition *status.PhaseCondition) {
	for i, currentCondition := range component.GetPhaseConditions() {
		if currentCondition.Phase == condition.Phase {
			component.Status.Conditions[i] = condition

			return
		}
	}

	// phase not found, lets add it to the list.
	component.Status.Conditions = append(component.Status.Conditions, condition)
}

// GetResources returns the child resource status for a component.
func (component *OCMLogForwarder) GetChildResourceConditions() []*status.ChildResource {
	return component.Status.Resources
}

// SetResources sets the phase conditions for a component.
func (component *OCMLogForwarder) SetChildResourceCondition(resource *status.ChildResource) {
	for i, currentResource := range component.GetChildResourceConditions() {
		if currentResource.Group == resource.Group && currentResource.Version == resource.Version && currentResource.Kind == resource.Kind {
			if currentResource.Name == resource.Name && currentResource.Namespace == resource.Namespace {
				component.Status.Resources[i] = resource

				return
			}
		}
	}

	// phase not found, lets add it to the collection
	component.Status.Resources = append(component.Status.Resources, resource)
}

// GetDependencies returns the dependencies for a component.
func (*OCMLogForwarder) GetDependencies() []workload.Workload {
	return []workload.Workload{}
}

// GetComponentGVK returns a GVK object for the component.
func (*OCMLogForwarder) GetWorkloadGVK() schema.GroupVersionKind {
	return GroupVersion.WithKind("OCMLogForwarder")
}

func init() {
	SchemeBuilder.Register(&OCMLogForwarder{}, &OCMLogForwarderList{})
}
