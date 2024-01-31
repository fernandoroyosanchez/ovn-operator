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
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// OVSDBServerDefaults -
type OVSDBServerDefaults struct {
	OvsContainerImageURL string
}

var ovsdbserverDefaults OVSDBServerDefaults

// log is for logging in this package.
var ovsdbserverlog = logf.Log.WithName("ovsdbserver-resource")

// SetupOVSDBServerDefaults - initialize OVSDBServer spec defaults for use with either internal or external webhooks
func SetupOVSDBServerDefaults(defaults OVSDBServerDefaults) {
	ovsdbserverDefaults = defaults
	ovsdbserverlog.Info("OVSDBServer defaults initialized", "defaults", defaults)
}

func (r *OVSDBServer) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-ovn-openstack-org-v1beta1-ovsdbserver,mutating=true,failurePolicy=fail,sideEffects=None,groups=ovn.openstack.org,resources=ovsdbservers,verbs=create;update,versions=v1beta1,name=movsdbserver.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &OVSDBServer{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *OVSDBServer) Default() {
	ovsdbserverlog.Info("default", "name", r.Name)

	r.Spec.Default()
}

// Default - set defaults for this OVSDBServer spec
func (spec *OVSDBServerSpec) Default() {
	if spec.OvsContainerImage == "" {
		spec.OvsContainerImage = ovsdbserverDefaults.OvsContainerImageURL
	}
}

//+kubebuilder:webhook:path=/validate-ovn-openstack-org-v1beta1-ovsdbserver,mutating=false,failurePolicy=fail,sideEffects=None,groups=ovn.openstack.org,resources=ovsdbservers,verbs=create;update,versions=v1beta1,name=vovsdbserver.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &OVSDBServer{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *OVSDBServer) ValidateCreate() error {
	ovsdbserverlog.Info("validate create", "name", r.Name)

	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *OVSDBServer) ValidateUpdate(old runtime.Object) error {
	ovsdbserverlog.Info("validate update", "name", r.Name)

	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *OVSDBServer) ValidateDelete() error {
	ovsdbserverlog.Info("validate delete", "name", r.Name)

	return nil
}
