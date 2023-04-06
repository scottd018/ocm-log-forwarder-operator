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

package mutate

import (
	"fmt"
	"strconv"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/nukleros/operator-builder-tools/pkg/controller/workload"
	"github.com/nukleros/operator-builder-tools/pkg/resources"

	appsv1alpha1 "github.com/scottd018/ocm-log-forwarder-operator/apis/apps/v1alpha1"
)

const (
	logForwarderContainerName = "forwarder"
)

// MutateDeploymentForwarderNamespaceParentName mutates the Deployment resource with name parent.Name.
func MutateDeploymentForwarderNamespaceParentName(
	original client.Object,
	parent *appsv1alpha1.OCMLogForwarder, collection *appsv1alpha1.OCMLogForwarderConfig,
	reconciler workload.Reconciler, req *workload.Request,
) ([]client.Object, error) {
	// if either the reconciler or request are found to be nil, return the base object.
	if reconciler == nil || req == nil {
		return []client.Object{original}, nil
	}

	// get the unstructured object.  we need to use unstructured here because our typed object
	// is not quite valid yet.  we need to convert the invalid object into a valid typed object.
	object, err := resources.ToUnstructured(original)
	if err != nil {
		return returnError("unable to convert object to unstructured", original)
	}
	objectMap := object.Object

	// retrieve and validate the container specification
	containers, found, err := unstructured.NestedSlice(objectMap, "spec", "template", "spec", "containers")
	if err != nil || !found {
		return returnError("unable to find container specification", original)
	}

	// convert all environment values to strings
	for i, container := range containers {
		containerMap, ok := container.(map[string]interface{})
		if !ok {
			return returnError("unable to convert container object", original)
		}

		// retrieve the environment from the container spec
		environment, found, err := unstructured.NestedSlice(containerMap, "env")
		if err != nil || !found {
			return returnError(fmt.Sprintf("unable to find environment from container %s", containerMap["name"]), original)
		}

		for j := range environment {
			envMap, ok := environment[j].(map[string]interface{})
			if !ok {
				return returnError(fmt.Sprintf("unable to convert environment for container %s", containerMap["name"]), original)
			}

			for name, value := range envMap {
				switch t := value.(type) {
				case int:
					envMap[name] = strconv.FormatInt(int64(t), 10)
				case int64:
					envMap[name] = strconv.FormatInt(t, 10)
				case bool:
					envMap[name] = strconv.FormatBool(t)
				}
			}

			// replace the container environment with the updated map
			environment[j] = envMap
			containerMap["env"] = environment
			containers[i] = containerMap
		}
	}

	// update the containers
	if err := unstructured.SetNestedSlice(objectMap, containers, "spec", "template", "spec", "containers"); err != nil {
		return returnError("unable to set containers", original)
	}
	object.Object = objectMap

	// attempt to create a typed deployment object before continuing
	deployment := &appsv1.Deployment{}
	if err := resources.ToTyped(deployment, object); err != nil {
		return returnError("unable to convert object to deployment", object)
	}

	return []client.Object{deployment}, nil
}

func returnError(message string, object client.Object) ([]client.Object, error) {
	return []client.Object{object}, fmt.Errorf(
		"%s for object [%s/%s] of kind [%s]",
		message,
		object.GetNamespace(),
		object.GetName(),
		object.GetObjectKind().GroupVersionKind().Kind,
	)
}
