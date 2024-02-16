/*
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

package ovsdbserver

import (
	"context"
	"fmt"

	"github.com/fernandoroyosanchez/ovn-operator/api/v1beta1"
	"github.com/openstack-k8s-operators/lib-common/modules/common/env"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func getOVSDBServerPods(
	ctx context.Context,
	k8sClient client.Client,
	instance *v1beta1.OVSDBServer,
) (*corev1.PodList, error) {

	podList := &corev1.PodList{}
	podListOpts := &client.ListOptions{
		Namespace: instance.Namespace,
	}
	client.MatchingLabels{
		"service": ServiceName,
	}.ApplyToList(podListOpts)

	if err := k8sClient.List(ctx, podList, podListOpts); err != nil {
		err = fmt.Errorf("error listing pods for instance %s: %w", instance.Name, err)
		return podList, err
	}

	return podList, nil
}

// EnvDownwardAPI - set env from FieldRef->FieldPath, e.g. status.podIP
func EnvDownwardAPI(field string) env.Setter {
	return func(env *corev1.EnvVar) {
		if env.ValueFrom == nil {
			env.ValueFrom = &corev1.EnvVarSource{}
		}
		env.Value = ""

		if env.ValueFrom.FieldRef == nil {
			env.ValueFrom.FieldRef = &corev1.ObjectFieldSelector{}
		}

		env.ValueFrom.FieldRef.FieldPath = field
	}
}
