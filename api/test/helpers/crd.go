/*
Copyright 2024 Red Hat
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

package helpers

import (
	"context"
	"time"

	ovnv1 "github.com/fernandoroyosanchez/ovn-operator/api/v1beta1"
	"github.com/go-logr/logr"
	"github.com/google/uuid"
	"github.com/onsi/gomega"
	"github.com/openstack-k8s-operators/lib-common/modules/common/condition"
	k8s_errors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	base "github.com/openstack-k8s-operators/lib-common/modules/common/test/helpers"
)

// TestHelper is a collection of helpers for testing operators. It extends the
// generic TestHelper from modules/test.
type TestHelper struct {
	*base.TestHelper
}

// NewTestHelper returns a TestHelper
func NewTestHelper(
	ctx context.Context,
	k8sClient client.Client,
	timeout time.Duration,
	interval time.Duration,
	logger logr.Logger,
) *TestHelper {
	helper := &TestHelper{}
	helper.TestHelper = base.NewTestHelper(ctx, k8sClient, timeout, interval, logger)
	return helper
}

// CreateOVNNorthd creates a new OVNNorthd instance with the specified
// namespace in the Kubernetes cluster.
//
// Example usage:
//
//	ovnNorthd := th.CreateOVNNorthd(namespace, spec)
//	DeferCleanup(th.DeleteOVNNorthd, ovnNorthd)
func (th *TestHelper) CreateOVNNorthd(namespace string, spec ovnv1.OVNNorthdSpec) types.NamespacedName {
	name := "ovnnorthd-" + uuid.New().String()
	ovnnorthd := &ovnv1.OVNNorthd{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "ovn.openstack.org/v1beta1",
			Kind:       "OVNNorthd",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: spec,
	}

	gomega.Expect(th.K8sClient.Create(th.Ctx, ovnnorthd)).Should(gomega.Succeed())
	th.Logger.Info("OVNNorthd created", "OVNNorthd", name)
	return types.NamespacedName{Namespace: namespace, Name: name}
}

// DeleteOVNNorthd deletes a OVNNorthd resource from the Kubernetes cluster.
//
// After the deletion, the function checks again if the OVNNorthd is
// successfully deleted.
//
// Example usage:
//
//	ovnNorthd := th.CreateOVNNorthd(namespace, spec)
//	DeferCleanup(th.DeleteOVNNorthd, ovnNorthd)
func (th *TestHelper) DeleteOVNNorthd(name types.NamespacedName) {
	gomega.Eventually(func(g gomega.Gomega) {
		ovnnorthd := &ovnv1.OVNNorthd{}
		err := th.K8sClient.Get(th.Ctx, name, ovnnorthd)
		// if it is already gone that is OK
		if k8s_errors.IsNotFound(err) {
			return
		}
		g.Expect(err).NotTo(gomega.HaveOccurred())

		g.Expect(th.K8sClient.Delete(th.Ctx, ovnnorthd)).Should(gomega.Succeed())

		err = th.K8sClient.Get(th.Ctx, name, ovnnorthd)
		g.Expect(k8s_errors.IsNotFound(err)).To(gomega.BeTrue())
	}, th.Timeout, th.Interval).Should(gomega.Succeed())
}

// GetOVNNorthd retrieves a OVNNorthd resource.
//
// The function returns a pointer to the retrieved OVNNorthd resource.
//
// Example usage:
//
//	ovnNorthdName := th.CreateOVNNorthd(namespace, spec)
//	ovnNorthd := th.GetOVNNorthd(ovnNorthdName)
func (th *TestHelper) GetOVNNorthd(name types.NamespacedName) *ovnv1.OVNNorthd {
	instance := &ovnv1.OVNNorthd{}
	gomega.Eventually(func(g gomega.Gomega) {
		g.Expect(th.K8sClient.Get(th.Ctx, name, instance)).Should(gomega.Succeed())
	}, th.Timeout, th.Interval).Should(gomega.Succeed())
	return instance
}

// SimulateOVNNorthdReady simulates the readiness of a OVNNorthd resource by
// setting the Ready condition of the OVNNorthd to true.
//
// Example usage:
// th.SimulateOVNNorthdReady(ovnNorthdName)
func (th *TestHelper) SimulateOVNNorthdReady(name types.NamespacedName) {
	gomega.Eventually(func(g gomega.Gomega) {
		service := th.GetOVNNorthd(name)
		service.Status.Conditions.MarkTrue(condition.ReadyCondition, "Ready")
		g.Expect(th.K8sClient.Status().Update(th.Ctx, service)).To(gomega.Succeed())
	}, th.Timeout, th.Interval).Should(gomega.Succeed())
	th.Logger.Info("Simulated GetOVNNorthd ready", "on", name)
}

// CreateOVNDBCluster creates a new OVNDBCluster instance with the specified
// namespace in the Kubernetes cluster.
//
// Example usage:
//
//	ovnDBCluster := th.CreateOVNDBCluster(namespace, spec)
//	DeferCleanup(th.DeleteOVNDBCluster, ovnDBCluster)
func (th *TestHelper) CreateOVNDBCluster(namespace string, spec ovnv1.OVNDBClusterSpec) types.NamespacedName {
	name := "ovndbcluster-" + uuid.New().String()
	ovnDBCluster := &ovnv1.OVNDBCluster{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "ovn.openstack.org/v1beta1",
			Kind:       "OVNDBCluster",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: spec,
	}

	gomega.Expect(th.K8sClient.Create(th.Ctx, ovnDBCluster)).Should(gomega.Succeed())
	th.Logger.Info("OVNDBCluster created", "OVNDBCluster", name)
	return types.NamespacedName{Namespace: namespace, Name: name}
}

// DeleteOVNDBCluster deletes a OVNDBCluster resource from the Kubernetes cluster.
//
// After the deletion, the function checks again if the OVNDBCluster is
// successfully deleted.
//
// Example usage:
//
//	ovnDBCluster := th.CreateOVNDBCluster(namespace, spec)
//	DeferCleanup(th.DeleteOVNDBCluster, ovnDBCluster)
func (th *TestHelper) DeleteOVNDBCluster(name types.NamespacedName) {
	gomega.Eventually(func(g gomega.Gomega) {
		ovnDBCluster := &ovnv1.OVNDBCluster{}
		err := th.K8sClient.Get(th.Ctx, name, ovnDBCluster)
		// if it is already gone that is OK
		if k8s_errors.IsNotFound(err) {
			return
		}
		g.Expect(err).NotTo(gomega.HaveOccurred())

		g.Expect(th.K8sClient.Delete(th.Ctx, ovnDBCluster)).Should(gomega.Succeed())

		err = th.K8sClient.Get(th.Ctx, name, ovnDBCluster)
		g.Expect(k8s_errors.IsNotFound(err)).To(gomega.BeTrue())
	}, th.Timeout, th.Interval).Should(gomega.Succeed())
}

// GetOVNDBCluster retrieves a OVNDBCluster resource.
//
// The function returns a pointer to the retrieved OVNDBCluster resource.
//
// Example usage:
//
//	ovnDBClusterName := th.CreateOVNDBCluster(namespace, spec)
//	ovnDBCluster := th.GetOVNDBCluster(ovnDBClusterName)
func (th *TestHelper) GetOVNDBCluster(name types.NamespacedName) *ovnv1.OVNDBCluster {
	instance := &ovnv1.OVNDBCluster{}
	gomega.Eventually(func(g gomega.Gomega) {
		g.Expect(th.K8sClient.Get(th.Ctx, name, instance)).Should(gomega.Succeed())
	}, th.Timeout, th.Interval).Should(gomega.Succeed())
	return instance
}

// SimulateOVNDBClusterReady simulates the readiness of a OVNDBCluster resource by
// setting the Ready condition of the OVNDBCluster to true.
//
// Example usage:
// th.SimulateOVNDBClusterReady(ovnDBClusterName)
func (th *TestHelper) SimulateOVNDBClusterReady(name types.NamespacedName) {
	gomega.Eventually(func(g gomega.Gomega) {
		service := th.GetOVNDBCluster(name)
		service.Status.Conditions.MarkTrue(condition.ReadyCondition, "Ready")
		g.Expect(th.K8sClient.Status().Update(th.Ctx, service)).To(gomega.Succeed())
	}, th.Timeout, th.Interval).Should(gomega.Succeed())
	th.Logger.Info("Simulated GetOVNDBCluster ready", "on", name)
}

// CreateOVNController creates a new OVNController instance with the specified
// namespace in the Kubernetes cluster.
//
// Example usage:
//
//	ovnController := th.CreateOVNController(namespace, spec)
//	DeferCleanup(th.DeleteOVNController, ovnController)
func (th *TestHelper) CreateOVNController(namespace string, spec ovnv1.OVNControllerSpec) types.NamespacedName {
	name := "ovncontroller-" + uuid.New().String()
	ovnController := &ovnv1.OVNController{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "ovn.openstack.org/v1beta1",
			Kind:       "OVNController",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: spec,
	}

	gomega.Expect(th.K8sClient.Create(th.Ctx, ovnController)).Should(gomega.Succeed())
	th.Logger.Info("OVNController created", "OVNController", name)
	return types.NamespacedName{Namespace: namespace, Name: name}
}

// DeleteOVNController deletes a OVNController resource from the Kubernetes cluster.
//
// After the deletion, the function checks again if the OVNController is
// successfully deleted.
//
// Example usage:
//
//	ovnController := th.CreateOVNController(namespace, spec)
//	DeferCleanup(th.DeleteOVNController, ovnController)
func (th *TestHelper) DeleteOVNController(name types.NamespacedName) {
	gomega.Eventually(func(g gomega.Gomega) {
		ovnController := &ovnv1.OVNController{}
		err := th.K8sClient.Get(th.Ctx, name, ovnController)
		// if it is already gone that is OK
		if k8s_errors.IsNotFound(err) {
			return
		}
		g.Expect(err).NotTo(gomega.HaveOccurred())

		g.Expect(th.K8sClient.Delete(th.Ctx, ovnController)).Should(gomega.Succeed())

		err = th.K8sClient.Get(th.Ctx, name, ovnController)
		g.Expect(k8s_errors.IsNotFound(err)).To(gomega.BeTrue())
	}, th.Timeout, th.Interval).Should(gomega.Succeed())
}

// GetOVNController retrieves a OVNController resource.
//
// The function returns a pointer to the retrieved OVNController resource.
//
// Example usage:
//
//	ovnControllerName := th.CreateOVNController(namespace, spec)
//	ovnController := th.GetOVNController(ovnControllerName)
func (th *TestHelper) GetOVNController(name types.NamespacedName) *ovnv1.OVNController {
	instance := &ovnv1.OVNController{}
	gomega.Eventually(func(g gomega.Gomega) {
		g.Expect(th.K8sClient.Get(th.Ctx, name, instance)).Should(gomega.Succeed())
	}, th.Timeout, th.Interval).Should(gomega.Succeed())
	return instance
}

// SimulateOVNControllerReady simulates the readiness of a OVNController resource by
// setting the Ready condition of the OVNController to true.
//
// Example usage:
// th.SimulateOVNControllerReady(ovnControllerName)
func (th *TestHelper) SimulateOVNControllerReady(name types.NamespacedName) {
	gomega.Eventually(func(g gomega.Gomega) {
		service := th.GetOVNController(name)
		service.Status.Conditions.MarkTrue(condition.ReadyCondition, "Ready")
		g.Expect(th.K8sClient.Status().Update(th.Ctx, service)).To(gomega.Succeed())
	}, th.Timeout, th.Interval).Should(gomega.Succeed())
	th.Logger.Info("Simulated GetOVNController ready", "on", name)
}

// CreateOVSDBServer creates a new OVSDBServer instance with the specified
// namespace in the Kubernetes cluster.
//
// Example usage:
//
//	ovsdbserver := th.CreateOVSDBServer (namespace, spec)
//	DeferCleanup(th.DeleteOVSDBServer, ovsdbserver)
func (th *TestHelper) CreateOVSDBServer(namespace string, spec ovnv1.OVSDBServerSpec) types.NamespacedName {
	name := "ovsdbserver-" + uuid.New().String()
	ovsdbserver := &ovnv1.OVSDBServer{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "ovn.openstack.org/v1beta1",
			Kind:       "OVSDBServer",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: spec,
	}

	gomega.Expect(th.K8sClient.Create(th.Ctx, ovsdbserver)).Should(gomega.Succeed())
	th.Logger.Info("OVSDBServer created", "OVSDBServer", name)
	return types.NamespacedName{Namespace: namespace, Name: name}
}

// DeleteOVSDBServer deletes a OVSDBServer resource from the Kubernetes cluster.
//
// After the deletion, the function checks again if the OVSDBServer is
// successfully deleted.
//
// Example usage:
//
//	ovsdbserver := th.CreateOVSDBServer(namespace, spec)
//	DeferCleanup(th.DeleteOVSDBServer, ovsdbserver)
func (th *TestHelper) DeleteOVSDBServer(name types.NamespacedName) {
	gomega.Eventually(func(g gomega.Gomega) {
		ovsdbserver := &ovnv1.OVSDBServer{}
		err := th.K8sClient.Get(th.Ctx, name, ovsdbserver)
		// if it is already gone that is OK
		if k8s_errors.IsNotFound(err) {
			return
		}
		g.Expect(err).NotTo(gomega.HaveOccurred())

		g.Expect(th.K8sClient.Delete(th.Ctx, ovsdbserver)).Should(gomega.Succeed())

		err = th.K8sClient.Get(th.Ctx, name, ovsdbserver)
		g.Expect(k8s_errors.IsNotFound(err)).To(gomega.BeTrue())
	}, th.Timeout, th.Interval).Should(gomega.Succeed())
}

// GetOVSDBServer retrieves a OVSDBServer resource.
//
// The function returns a pointer to the retrieved OVSDBServer resource.
//
// Example usage:
//
//	ovsdbserverName := th.CreateOVSDBServer(namespace, spec)
//	ovsdbserver := th.GetOVSDBServer(ovsdbserverName)
func (th *TestHelper) GetOVSDBServer(name types.NamespacedName) *ovnv1.OVSDBServer {
	instance := &ovnv1.OVSDBServer{}
	gomega.Eventually(func(g gomega.Gomega) {
		g.Expect(th.K8sClient.Get(th.Ctx, name, instance)).Should(gomega.Succeed())
	}, th.Timeout, th.Interval).Should(gomega.Succeed())
	return instance
}

// SimulateOVSDBServerReady simulates the readiness of a OVSDBServer resource by
// setting the Ready condition of the OVSDBServer to true.
//
// Example usage:
// th.SimulateOVSDBServerReady(ovsdbserverName)
func (th *TestHelper) SimulateOVSDBServerReady(name types.NamespacedName) {
	gomega.Eventually(func(g gomega.Gomega) {
		service := th.GetOVSDBServer(name)
		service.Status.Conditions.MarkTrue(condition.ReadyCondition, "Ready")
		g.Expect(th.K8sClient.Status().Update(th.Ctx, service)).To(gomega.Succeed())
	}, th.Timeout, th.Interval).Should(gomega.Succeed())
	th.Logger.Info("Simulated GetOVSDBServer ready", "on", name)
}

// CreateOVSvswitchd creates a new OVSvswitchd instance with the specified
// namespace in the Kubernetes cluster.
//
// Example usage:
//
//	ovsvswitchd := th.CreateOVSvswitchd (namespace, spec)
//	DeferCleanup(th.DeleteOVSvswitchd, ovsvswitchd)
func (th *TestHelper) CreateOVSvswitchd(namespace string, spec ovnv1.OVSvswitchdSpec) types.NamespacedName {
	name := "ovsvswitchd-" + uuid.New().String()
	ovsvswitchd := &ovnv1.OVSvswitchd{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "ovn.openstack.org/v1beta1",
			Kind:       "OVSvswitchd",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: spec,
	}

	gomega.Expect(th.K8sClient.Create(th.Ctx, ovsvswitchd)).Should(gomega.Succeed())
	th.Logger.Info("OVSvswitchd created", "OVSvswitchd", name)
	return types.NamespacedName{Namespace: namespace, Name: name}
}

// DeleteOVSvswitchd deletes a OVSvswitchd resource from the Kubernetes cluster.
//
// After the deletion, the function checks again if the OVSvswitchd is
// successfully deleted.
//
// Example usage:
//
//	ovsvswitchd := th.CreateOVSvswitchd(namespace, spec)
//	DeferCleanup(th.DeleteOVSvswitchd, ovsvswitchd)
func (th *TestHelper) DeleteOVSvswitchd(name types.NamespacedName) {
	gomega.Eventually(func(g gomega.Gomega) {
		ovsvswitchd := &ovnv1.OVSvswitchd{}
		err := th.K8sClient.Get(th.Ctx, name, ovsvswitchd)
		// if it is already gone that is OK
		if k8s_errors.IsNotFound(err) {
			return
		}
		g.Expect(err).NotTo(gomega.HaveOccurred())

		g.Expect(th.K8sClient.Delete(th.Ctx, ovsvswitchd)).Should(gomega.Succeed())

		err = th.K8sClient.Get(th.Ctx, name, ovsvswitchd)
		g.Expect(k8s_errors.IsNotFound(err)).To(gomega.BeTrue())
	}, th.Timeout, th.Interval).Should(gomega.Succeed())
}

// GetOVSvswitchd retrieves a OVSvswitchd resource.
//
// The function returns a pointer to the retrieved OVSvswitchd resource.
//
// Example usage:
//
//	ovsvswitchdName := th.CreateOVSvswitchd(namespace, spec)
//	ovsvswitchd := th.GetOVSvswitchd(ovsvswitchdName)
func (th *TestHelper) GetOVSvswitchd(name types.NamespacedName) *ovnv1.OVSvswitchd {
	instance := &ovnv1.OVSvswitchd{}
	gomega.Eventually(func(g gomega.Gomega) {
		g.Expect(th.K8sClient.Get(th.Ctx, name, instance)).Should(gomega.Succeed())
	}, th.Timeout, th.Interval).Should(gomega.Succeed())
	return instance
}

// SimulateOVSvswitchdReady simulates the readiness of a OVSvswitchd resource by
// setting the Ready condition of the OVSvswitchd to true.
//
// Example usage:
// th.SimulateOVSvswitchdReady(ovsvswitchdName)
func (th *TestHelper) SimulateOVSvswitchdReady(name types.NamespacedName) {
	gomega.Eventually(func(g gomega.Gomega) {
		service := th.GetOVSvswitchd(name)
		service.Status.Conditions.MarkTrue(condition.ReadyCondition, "Ready")
		g.Expect(th.K8sClient.Status().Update(th.Ctx, service)).To(gomega.Succeed())
	}, th.Timeout, th.Interval).Should(gomega.Succeed())
	th.Logger.Info("Simulated GetOVSvswitchd ready", "on", name)
}
