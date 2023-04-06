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

var ErrUnableToConvertOCMLogForwarderConfig = errors.New("unable to convert to OCMLogForwarderConfig")

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// OCMLogForwarderConfigSpec defines the desired state of OCMLogForwarderConfig.
type OCMLogForwarderConfigSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:default="rosa"
	// +kubebuilder:validation:Optional
	// (Default: "rosa")
	//  +kubebuilder:validation:Enum=rosa
	//  Provider to which the OCM Log Forwarder is provisioned to.  Only 'rosa' supported at this time.
	//
	Provider string `json:"provider,omitempty"`

	// +kubebuilder:default="ocm-log-forwarder"
	// +kubebuilder:validation:Optional
	// (Default: "ocm-log-forwarder")
	//  Namespace used to provision the OCM Log Forwarder workloads.  Each workload provisoned will
	//  end up in this namespace.
	//
	ForwarderNamespace string `json:"forwarderNamespace,omitempty"`

	// +kubebuilder:validation:Optional
	Backend OCMLogForwarderConfigSpecBackend `json:"backend,omitempty"`
}

type OCMLogForwarderConfigSpecBackend struct {
	// +kubebuilder:default="elasticsearch"
	// +kubebuilder:validation:Optional
	// (Default: "elasticsearch")
	//  +kubebuilder:validation:Enum=elasticsearch
	//  Backend type where logs are sent and stored.  Only 'elasticsearch' supported at this time.  Requires
	//  backend.elasticSearch.url to be set.
	//
	Type string `json:"type,omitempty"`

	// +kubebuilder:validation:Optional
	ElasticSearch OCMLogForwarderConfigSpecBackendElasticSearch `json:"elasticSearch,omitempty"`
}

type OCMLogForwarderConfigSpecBackendElasticSearch struct {
	// +kubebuilder:default="https://elasticsearch-es-http.elastic-system.svc.cluster.local:9200"
	// +kubebuilder:validation:Optional
	// (Default: "https://elasticsearch-es-http.elastic-system.svc.cluster.local:9200")
	//  URL to which to ship logs when using the 'elasticsearch' as a backend in the .spec.backend.type
	//  field of this custom resource.
	//
	Url string `json:"url,omitempty"`

	// +kubebuilder:default="basic"
	// +kubebuilder:validation:Optional
	// (Default: "basic")
	//  +kubebuilder:validation:Enum=basic
	//  ElasticSearch authentication type to use.  Only 'basic' supported at this time.  Requires
	//  a single key/value pair stored in a secret named 'elastic-auth' which contains the
	//  basic authentication info for the ElasticSearch connection.  Secret must exist within the same namespace
	//  where the OCM Log Collector is deployed to.
	//
	AuthType string `json:"authType,omitempty"`

	// +kubebuilder:default="ocm_service_logs"
	// +kubebuilder:validation:Optional
	// (Default: "ocm_service_logs")
	//  +kubebuilder:validation:MaxLength=128
	//  Index name in ElasticSearch where service logs are sent.  Index name must be 128 characters or less.
	//
	Index string `json:"index,omitempty"`
}

// OCMLogForwarderConfigStatus defines the observed state of OCMLogForwarderConfig.
type OCMLogForwarderConfigStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Created               bool                     `json:"created,omitempty"`
	DependenciesSatisfied bool                     `json:"dependenciesSatisfied,omitempty"`
	Conditions            []*status.PhaseCondition `json:"conditions,omitempty"`
	Resources             []*status.ChildResource  `json:"resources,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster

// OCMLogForwarderConfig is the Schema for the ocmlogforwarderconfigs API.
type OCMLogForwarderConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              OCMLogForwarderConfigSpec   `json:"spec,omitempty"`
	Status            OCMLogForwarderConfigStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// OCMLogForwarderConfigList contains a list of OCMLogForwarderConfig.
type OCMLogForwarderConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OCMLogForwarderConfig `json:"items"`
}

// interface methods

// GetReadyStatus returns the ready status for a component.
func (component *OCMLogForwarderConfig) GetReadyStatus() bool {
	return component.Status.Created
}

// SetReadyStatus sets the ready status for a component.
func (component *OCMLogForwarderConfig) SetReadyStatus(ready bool) {
	component.Status.Created = ready
}

// GetDependencyStatus returns the dependency status for a component.
func (component *OCMLogForwarderConfig) GetDependencyStatus() bool {
	return component.Status.DependenciesSatisfied
}

// SetDependencyStatus sets the dependency status for a component.
func (component *OCMLogForwarderConfig) SetDependencyStatus(dependencyStatus bool) {
	component.Status.DependenciesSatisfied = dependencyStatus
}

// GetPhaseConditions returns the phase conditions for a component.
func (component *OCMLogForwarderConfig) GetPhaseConditions() []*status.PhaseCondition {
	return component.Status.Conditions
}

// SetPhaseCondition sets the phase conditions for a component.
func (component *OCMLogForwarderConfig) SetPhaseCondition(condition *status.PhaseCondition) {
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
func (component *OCMLogForwarderConfig) GetChildResourceConditions() []*status.ChildResource {
	return component.Status.Resources
}

// SetResources sets the phase conditions for a component.
func (component *OCMLogForwarderConfig) SetChildResourceCondition(resource *status.ChildResource) {
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
func (*OCMLogForwarderConfig) GetDependencies() []workload.Workload {
	return []workload.Workload{}
}

// GetComponentGVK returns a GVK object for the component.
func (*OCMLogForwarderConfig) GetWorkloadGVK() schema.GroupVersionKind {
	return GroupVersion.WithKind("OCMLogForwarderConfig")
}

func init() {
	SchemeBuilder.Register(&OCMLogForwarderConfig{}, &OCMLogForwarderConfigList{})
}
