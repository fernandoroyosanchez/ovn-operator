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

// OVSvsswitchdDefaults -
type OVSvsswitchdDefaults struct {
	OvsContainerImageURL string
}

var ovsvswitchdDefaults OVSvsswitchdDefaults

// log is for logging in this package.
var ovsvswitchdlog = logf.Log.WithName("ovsvswitchd-resource")

// SetupOVSvsswitchdDefaults - initialize SetupOVSvsswitchd spec defaults for use with either internal or external webhooks
func SetupOVSvsswitchdDefaults(defaults OVSvsswitchdDefaults) {
	ovsvswitchdDefaults = defaults
	ovsvswitchdlog.Info("OVSvsswitchd defaults initialized", "defaults", defaults)
}

func (r *OVSvswitchd) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-ovn-openstack-org-v1beta1-ovsvswitchd,mutating=true,failurePolicy=fail,sideEffects=None,groups=ovn.openstack.org,resources=ovsvswitchds,verbs=create;update,versions=v1beta1,name=movsvswitchd.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &OVSvswitchd{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *OVSvswitchd) Default() {
	ovsvswitchdlog.Info("default", "name", r.Name)

	r.Spec.Default()
}

// Default - set defaults for this OVSvswitchd spec
func (spec *OVSvswitchdSpec) Default() {
	if spec.OvsContainerImage == "" {
		spec.OvsContainerImage = ovsvswitchdDefaults.OvsContainerImageURL
	}
}

//+kubebuilder:webhook:path=/validate-ovn-openstack-org-v1beta1-ovsvswitchd,mutating=false,failurePolicy=fail,sideEffects=None,groups=ovn.openstack.org,resources=ovsvswitchds,verbs=create;update,versions=v1beta1,name=vovsvswitchd.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &OVSvswitchd{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *OVSvswitchd) ValidateCreate() error {
	ovsvswitchdlog.Info("validate create", "name", r.Name)

	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *OVSvswitchd) ValidateUpdate(old runtime.Object) error {
	ovsvswitchdlog.Info("validate update", "name", r.Name)

	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *OVSvswitchd) ValidateDelete() error {
	ovsvswitchdlog.Info("validate delete", "name", r.Name)

	return nil
}
