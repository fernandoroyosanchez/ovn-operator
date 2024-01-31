/*
Copyright 2022.

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

package v1beta1

import (
	"github.com/openstack-k8s-operators/lib-common/modules/common/condition"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// OvsvswitchdOvsContainerImage is the fall-back container image for OVNController ovs-*
	OvsvswitchdOvsContainerImage = "quay.io/podified-antelope-centos9/openstack-ovn-base:current-podified"
)

// OVSvswitchdSpec defines the desired state of OVSvswitchd
type OVSvswitchdSpec struct {
	// +kubebuilder:validation:Required
	// Image used for ovs-vswitchd containers (will be set to environmental default if empty)
	OvsContainerImage string `json:"ovsContainerImage"`

	// +kubebuilder:validation:Optional
	// Debug - enable debug for different deploy stages. If an init container is used, it runs and the
	// actual action pod gets started with sleep infinity
	Debug OVSvswitchdDebug `json:"debug,omitempty"`

	// +kubebuilder:validation:Optional
	// +optional
	NicMappings map[string]string `json:"nicMappings,omitempty"`

	// +kubebuilder:validation:Optional
	// NodeSelector to target subset of worker nodes running this service
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// +kubebuilder:validation:Optional
	// Resources - Compute Resources required by this service (Limits/Requests).
	// https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`

	// +kubebuilder:validation:Optional
	// NetworkAttachment is a NetworkAttachment resource name to expose the service to the given network.
	// If specified the IP address of this network is used as the OvnEncapIP.
	// Deprecated: superseded by NetworkAttachments
	NetworkAttachment string `json:"networkAttachment"`

	// +kubebuilder:validation:Optional
	// NetworkAttachments are NetworkAttachment resources used to expose the service to the specified networks.
	// If present, the IP of the attachment named "tenant", will be used as the OvnEncapIP.

	NetworkAttachments []string `json:"networkAttachments,omitempty"`
}

// OVSvswitchdDebug defines the observed state of OVSvswitchdDebug
type OVSvswitchdDebug struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=false
	// Service enable debug
	Service bool `json:"service"`
}

// OVSvswitchdStatus defines the observed state of OVSvswitchd
type OVSvswitchdStatus struct {
	// NumberReady of the OVNController instances
	NumberReady int32 `json:"numberReady,omitempty"`

	// DesiredNumberScheduled - total number of the nodes which should be running Daemon
	DesiredNumberScheduled int32 `json:"desiredNumberScheduled,omitempty"`

	// Conditions
	Conditions condition.Conditions `json:"conditions,omitempty" optional:"true"`

	// Map of hashes to track e.g. job status
	Hash map[string]string `json:"hash,omitempty"`

	// NetworkAttachments status of the deployment pods
	NetworkAttachments map[string][]string `json:"networkAttachments,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="NetworkAttachments",type="string",JSONPath=".status.networkAttachments",description="NetworkAttachments"
//+kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.conditions[0].status",description="Status"
//+kubebuilder:printcolumn:name="Message",type="string",JSONPath=".status.conditions[0].message",description="Message"

// OVSvswitchd is the Schema for the ovsvswitchds API
type OVSvswitchd struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OVSvswitchdSpec   `json:"spec,omitempty"`
	Status OVSvswitchdStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// OVSvswitchdList contains a list of OVSvswitchd
type OVSvswitchdList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OVSvswitchd `json:"items"`
}

func init() {
	SchemeBuilder.Register(&OVSvswitchd{}, &OVSvswitchdList{})
}

// IsReady - returns true if service is ready to server requests
func (instance OVSvswitchd) IsReady() bool {
	// Ready when:
	// OVSvswitchd is reconciled successfully
	return instance.Status.Conditions.IsTrue(condition.ReadyCondition)
}

// RbacConditionsSet - set the conditions for the rbac object
func (instance OVSvswitchd) RbacConditionsSet(c *condition.Condition) {
	instance.Status.Conditions.Set(c)
}

// RbacNamespace - return the namespace
func (instance OVSvswitchd) RbacNamespace() string {
	return instance.Namespace
}

// RbacResourceName - return the name to be used for rbac objects (serviceaccount, role, rolebinding)
func (instance OVSvswitchd) RbacResourceName() string {
	return "ovsvswitchd-" + instance.Name
}
