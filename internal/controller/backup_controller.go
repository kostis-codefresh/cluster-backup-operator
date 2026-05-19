package controller

import (
	"context"

	"github.com/kostis-codefresh/cluster-backup-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// BackupReconciler reconciles a Backup object
type BackupReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=backups.cloudnativedays.operators.com,resources=backups,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=backups.cloudnativedays.operators.com,resources=backups/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=backups.cloudnativedays.operators.com,resources=backups/finalizers,verbs=update

// Reconcile handles the backup resource
func (r *BackupReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("Starting reconciliation", "request", req.NamespacedName)

	// Fetch the Backup instance
	backup := &v1alpha1.Backup{}
	err := r.Get(ctx, req.NamespacedName, backup)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted
			logger.Info("Backup resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request
		logger.Error(err, "Failed to get Backup")
		return ctrl.Result{}, err
	}

	// Log the backup message - this is our main "action"
	logger.Info("*** OPERATOR MESSAGE ***", "message", backup.Spec.Message, "name", backup.Name)

	// Update status to show we processed it (if not already processed)
	if !backup.Status.Processed {
		logger.Info("Updating backup status", "name", backup.Name)
		backup.Status.Processed = true
		err = r.Status().Update(ctx, backup)
		if err != nil {
			logger.Error(err, "Failed to update Backup status")
			return ctrl.Result{}, err
		}
		logger.Info("Successfully updated backup status to processed")
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager
func (r *BackupReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Backup{}).
		Complete(r)
}
