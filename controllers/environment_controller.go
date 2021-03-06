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

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "github.com/pavel-mikhalchuk/kubebuilder-playground/api/v1"
)

// EnvironmentReconciler reconciles a Environment object
type EnvironmentReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// Reconcile func docs are WIP
// +kubebuilder:rbac:groups=infra.mikhalchuk.k8s,resources=environments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infra.mikhalchuk.k8s,resources=environments/status,verbs=get;update;patch
func (r *EnvironmentReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("environment", req.NamespacedName)

	var env infrav1.Environment
	if err := r.Get(ctx, req.NamespacedName, &env); err != nil {
		log.Error(err, "unable to fetch Env")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	env.Status.NfsProvisioned = true

	if err := r.Status().Update(ctx, &env); err != nil {
		log.Error(err, "Unable to update Environment status")
	}

	return ctrl.Result{}, nil
}

// SetupWithManager adds EnvironmentReconciler to the manager
func (r *EnvironmentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrav1.Environment{}).
		Complete(r)
}
