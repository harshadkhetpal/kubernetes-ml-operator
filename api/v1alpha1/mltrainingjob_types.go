// Author: Harshad Khetpal <harshadkhetpal@gmail.com>
// Kubernetes ML Operator — CRD Type Definitions

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/api/core/v1"
)

// MLTrainingJobSpec defines the desired state of MLTrainingJob
type MLTrainingJobSpec struct {
	// +kubebuilder:validation:Required
	Image string `json:"image"`

	// +kubebuilder:validation:Enum=pytorch;tensorflow;xgboost;sklearn
	Framework string `json:"framework"`

	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:default=1
	Replicas int32 `json:"replicas,omitempty"`

	// +kubebuilder:validation:Minimum=0
	GPUPerReplica int32 `json:"gpuPerReplica,omitempty"`

	Dataset string `json:"dataset"`

	Hyperparameters map[string]string `json:"hyperparameters,omitempty"`

	Resources corev1.ResourceRequirements `json:"resources,omitempty"`

	Tracking TrackingConfig `json:"tracking,omitempty"`
}

type TrackingConfig struct {
	MLflowURI  string `json:"mlflow_uri,omitempty"`
	Experiment string `json:"experiment,omitempty"`
	WandbProject string `json:"wandb_project,omitempty"`
}

// MLTrainingJobStatus defines the observed state of MLTrainingJob
type MLTrainingJobStatus struct {
	Phase      string      `json:"phase,omitempty"` // Pending, Running, Succeeded, Failed
	StartTime  *metav1.Time `json:"startTime,omitempty"`
	FinishTime *metav1.Time `json:"finishTime,omitempty"`
	Conditions []metav1.Condition `json:"conditions,omitempty"`
	RunID      string `json:"runId,omitempty"`
	ModelURI   string `json:"modelUri,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Phase",type=string,JSONPath=`.status.phase`
//+kubebuilder:printcolumn:name="Framework",type=string,JSONPath=`.spec.framework`
//+kubebuilder:printcolumn:name="Replicas",type=integer,JSONPath=`.spec.replicas`
//+kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`
type MLTrainingJob struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec   MLTrainingJobSpec   `json:"spec,omitempty"`
	Status MLTrainingJobStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true
type MLTrainingJobList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items []MLTrainingJob `json:"items"`
}

func init() { SchemeBuilder.Register(&MLTrainingJob{}, &MLTrainingJobList{}) }
