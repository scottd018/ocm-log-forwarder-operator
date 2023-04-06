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
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"

	"github.com/nukleros/operator-builder-tools/pkg/controller/workload"

	appsv1alpha1 "github.com/scottd018/ocm-log-forwarder-operator/apis/apps/v1alpha1"
)

// sampleOCMLogForwarder is a sample containing all fields
const sampleOCMLogForwarder = `apiVersion: apps.dustinscott.io/v1alpha1
kind: OCMLogForwarder
metadata:
  name: ocmlogforwarder-sample
  namespace: default
spec:
  #collection:
    #name: "ocmlogforwarderconfig-sample"
    #namespace: ""
  forwarderVersion: "latest"
  ocm:
    clusterId: "22tgckqk9c2ff3jd8ve62p0i2st14vrq"
    pollInternalMinutes: 5
  debug: false
`

// sampleOCMLogForwarderRequired is a sample containing only required fields
const sampleOCMLogForwarderRequired = `apiVersion: apps.dustinscott.io/v1alpha1
kind: OCMLogForwarder
metadata:
  name: ocmlogforwarder-sample
  namespace: default
spec:
  #collection:
    #name: "ocmlogforwarderconfig-sample"
    #namespace: ""
  ocm:
    clusterId: "22tgckqk9c2ff3jd8ve62p0i2st14vrq"
`

// Sample returns the sample manifest for this custom resource.
func Sample(requiredOnly bool) string {
	if requiredOnly {
		return sampleOCMLogForwarderRequired
	}

	return sampleOCMLogForwarder
}

// Generate returns the child resources that are associated with this workload given
// appropriate structured inputs.
func Generate(
	workloadObj appsv1alpha1.OCMLogForwarder,
	collectionObj appsv1alpha1.OCMLogForwarderConfig,
	reconciler workload.Reconciler,
	req *workload.Request,
) ([]client.Object, error) {
	resourceObjects := []client.Object{}

	for _, f := range CreateFuncs {
		resources, err := f(&workloadObj, &collectionObj, reconciler, req)

		if err != nil {
			return nil, err
		}

		resourceObjects = append(resourceObjects, resources...)
	}

	return resourceObjects, nil
}

// GenerateForCLI returns the child resources that are associated with this workload given
// appropriate YAML manifest files.
func GenerateForCLI(workloadFile []byte, collectionFile []byte) ([]client.Object, error) {
	var workloadObj appsv1alpha1.OCMLogForwarder
	if err := yaml.Unmarshal(workloadFile, &workloadObj); err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml into workload, %w", err)
	}

	if err := workload.Validate(&workloadObj); err != nil {
		return nil, fmt.Errorf("error validating workload yaml, %w", err)
	}

	var collectionObj appsv1alpha1.OCMLogForwarderConfig
	if err := yaml.Unmarshal(collectionFile, &collectionObj); err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml into collection, %w", err)
	}

	if err := workload.Validate(&collectionObj); err != nil {
		return nil, fmt.Errorf("error validating collection yaml, %w", err)
	}

	return Generate(workloadObj, collectionObj, nil, nil)
}

// CreateFuncs is an array of functions that are called to create the child resources for the controller
// in memory during the reconciliation loop prior to persisting the changes or updates to the Kubernetes
// database.
var CreateFuncs = []func(
	*appsv1alpha1.OCMLogForwarder,
	*appsv1alpha1.OCMLogForwarderConfig,
	workload.Reconciler,
	*workload.Request,
) ([]client.Object, error){
	CreateDeploymentForwarderNamespaceParentName,
}

// InitFuncs is an array of functions that are called prior to starting the controller manager.  This is
// necessary in instances which the controller needs to "own" objects which depend on resources to
// pre-exist in the cluster. A common use case for this is the need to own a custom resource.
// If the controller needs to own a custom resource type, the CRD that defines it must
// first exist. In this case, the InitFunc will create the CRD so that the controller
// can own custom resources of that type.  Without the InitFunc the controller will
// crash loop because when it tries to own a non-existent resource type during manager
// setup, it will fail.
var InitFuncs = []func(
	*appsv1alpha1.OCMLogForwarder,
	*appsv1alpha1.OCMLogForwarderConfig,
	workload.Reconciler,
	*workload.Request,
) ([]client.Object, error){}

func ConvertWorkload(component, collection workload.Workload) (
	*appsv1alpha1.OCMLogForwarder,
	*appsv1alpha1.OCMLogForwarderConfig,
	error,
) {
	p, ok := component.(*appsv1alpha1.OCMLogForwarder)
	if !ok {
		return nil, nil, appsv1alpha1.ErrUnableToConvertOCMLogForwarder
	}

	c, ok := collection.(*appsv1alpha1.OCMLogForwarderConfig)
	if !ok {
		return nil, nil, appsv1alpha1.ErrUnableToConvertOCMLogForwarderConfig
	}

	return p, c, nil
}
