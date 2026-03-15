// Author: Harshad Khetpal <harshadkhetpal@gmail.com>
// MLTrainingJob Reconciler

package controllers

import (
	"context"
	"fmt"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	mlopsv1alpha1 "github.com/harshadkhetpal/kubernetes-ml-operator/api/v1alpha1"
)

type MLTrainingJobReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *MLTrainingJobReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	job := &mlopsv1alpha1.MLTrainingJob{}
	if err := r.Get(ctx, req.NamespacedName, job); err != nil {
		if errors.IsNotFound(err) { return ctrl.Result{}, nil }
		return ctrl.Result{}, err
	}

	logger.Info("Reconciling MLTrainingJob", "name", job.Name, "framework", job.Spec.Framework)

	if job.Status.Phase == "Succeeded" || job.Status.Phase == "Failed" {
		return ctrl.Result{}, nil
	}

	existing := &batchv1.Job{}
	err := r.Get(ctx, req.NamespacedName, existing)
	if errors.IsNotFound(err) {
		k8sJob := r.buildK8sJob(job)
		if err := r.Create(ctx, k8sJob); err != nil {
			return ctrl.Result{}, fmt.Errorf("creating K8s Job: %w", err)
		}
		job.Status.Phase = "Running"
		now := metav1.Now()
		job.Status.StartTime = &now
		r.Status().Update(ctx, job)
		logger.Info("Created K8s Job for MLTrainingJob", "job", job.Name)
		return ctrl.Result{}, nil
	}

	// Sync status from K8s Job
	if existing.Status.Succeeded > 0 {
		job.Status.Phase = "Succeeded"
		now := metav1.Now()
		job.Status.FinishTime = &now
	} else if existing.Status.Failed > 0 {
		job.Status.Phase = "Failed"
	}
	r.Status().Update(ctx, job)
	return ctrl.Result{}, nil
}

func (r *MLTrainingJobReconciler) buildK8sJob(mlJob *mlopsv1alpha1.MLTrainingJob) *batchv1.Job {
	env := []corev1.EnvVar{
		{Name: "MLFLOW_TRACKING_URI", Value: mlJob.Spec.Tracking.MLflowURI},
		{Name: "MLFLOW_EXPERIMENT_NAME", Value: mlJob.Spec.Tracking.Experiment},
		{Name: "DATASET_URI", Value: mlJob.Spec.Dataset},
	}
	for k, v := range mlJob.Spec.Hyperparameters {
		env = append(env, corev1.EnvVar{Name: "HP_" + k1, Value: v})
	}
	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: mlJob.Name, Namespace: mlJob.Namespace,
			Labels: map[string]string{"mlops.harshadkhetpal.dev/job": mlJob.Name},
		},
		Spec: batchv1.JobSpec{
			Parallelism: &mlJob.Spec.Replicas,
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:  "trainer",
						Image: mlJob.Spec.Image,
						Env:   env,
						Resources: mlJob.Spec.Resources,
					}},
					RestartPolicy: corev1.RestartPolicyNever,
				},
			},
		},
	}
}

func (r *MLTrainingJobReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).For(&mlopsv1alpha1.MLTrainingJob{}).Owns(&batchv1.Job{}).Complete(r)
}
