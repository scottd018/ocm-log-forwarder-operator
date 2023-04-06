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

package apps

import (
	v1alpha1apps "github.com/scottd018/ocm-log-forwarder-operator/apis/apps/v1alpha1"
	//+kubebuilder:scaffold:operator-builder:imports

	"k8s.io/apimachinery/pkg/runtime/schema"
)

// OCMLogForwarderConfigGroupVersions returns all group version objects associated with this kind.
func OCMLogForwarderConfigGroupVersions() []schema.GroupVersion {
	return []schema.GroupVersion{
		v1alpha1apps.GroupVersion,
		//+kubebuilder:scaffold:operator-builder:groupversions
	}
}
