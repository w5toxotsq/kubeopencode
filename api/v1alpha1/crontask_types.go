// Copyright Contributors to the KubeOpenCode project

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ConcurrencyPolicy describes how the CronTask controller handles
// concurrent executions of Tasks created by a CronTask.
// +kubebuilder:validation:Enum=Allow;Forbid;Replace
type ConcurrencyPolicy string

const (
	// AllowConcurrent allows concurrent Task executions.
	AllowConcurrent ConcurrencyPolicy = "Allow"

	// ForbidConcurrent skips the new Task creation if the previous Task is still running.
	ForbidConcurrent ConcurrencyPolicy = "Forbid"

	// ReplaceConcurrent stops the currently running Task and creates a new one.
	ReplaceConcurrent ConcurrencyPolicy = "Replace"
)

const (
	// CronTaskLabelKey is the label key added to Tasks created by a CronTask.
	// The value is the CronTask name.
	CronTaskLabelKey = "kubeopencode.io/crontask"

	// CronTaskTriggerAnnotation is the annotation key used to manually trigger
	// a CronTask to create a Task immediately.
	CronTaskTriggerAnnotation = "kubeopencode.io/trigger"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope="Namespaced",shortName=ct
// +kubebuilder:printcolumn:JSONPath=`.spec.schedule`,name="Schedule",type=string
// +kubebuilder:printcolumn:JSONPath=`.spec.suspend`,name="Suspend",type=boolean
// +kubebuilder:printcolumn:JSONPath=`.status.active`,name="Active",type=integer,priority=1
// +kubebuilder:printcolumn:JSONPath=`.status.lastScheduleTime`,name="Last Schedule",type=date
// +kubebuilder:printcolumn:JSONPath=`.status.nextScheduleTime`,name="Next Schedule",type=date,priority=1
// +kubebuilder:printcolumn:JSONPath=`.metadata.creationTimestamp`,name="Age",type=date

// CronTask defines a scheduled Task execution.
// CronTask is a Task factory that creates Task objects on a cron schedule.
// It is analogous to Kubernetes CronJob, but creates Tasks instead of Jobs.
type CronTask struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec defines the desired state of CronTask
	Spec CronTaskSpec `json:"spec"`

	// Status represents the current status of the CronTask
	// +optional
	Status CronTaskStatus `json:"status,omitempty"`
}

// CronTaskSpec defines the CronTask configuration.
type CronTaskSpec struct {
	// Schedule in Cron format.
	// Uses standard 5-field cron syntax: minute hour day-of-month month day-of-week.
	//
	// Examples:
	//   "0 9 * * 1-5"   → Every weekday at 09:00
	//   "*/30 * * * *"  → Every 30 minutes
	//   "0 0 1 * *"     → First day of every month at midnight
	// +required
	Schedule string `json:"schedule"`

	// TimeZone is the IANA timezone for the cron schedule.
	// If not specified, UTC is used.
	//
	// Example: "Asia/Shanghai", "America/New_York", "Europe/London"
	// +optional
	TimeZone *string `json:"timeZone,omitempty"`

	// ConcurrencyPolicy specifies how to handle concurrent executions
	// of Tasks created by this CronTask.
	//
	// - "Allow": allows concurrent Task executions (multiple Tasks can run simultaneously)
	// - "Forbid": skips the new Task if the previous one is still running (default)
	// - "Replace": stops the running Task and creates a new one
	// +optional
	// +kubebuilder:default=Forbid
	ConcurrencyPolicy ConcurrencyPolicy `json:"concurrencyPolicy,omitempty"`

	// Suspend tells the controller to suspend subsequent executions.
	// Existing running Tasks are not affected. Defaults to false.
	// +optional
	Suspend *bool `json:"suspend,omitempty"`

	// StartingDeadlineSeconds is the deadline in seconds for starting
	// the Task if it misses its scheduled time for any reason.
	// Missed Tasks past this deadline are counted as missed and skipped.
	// If not specified, there is no deadline.
	// +optional
	// +kubebuilder:validation:Minimum=0
	StartingDeadlineSeconds *int64 `json:"startingDeadlineSeconds,omitempty"`

	// MaxRetainedTasks is the maximum number of child Tasks (in any state)
	// that can exist for this CronTask. When the limit is reached, new
	// scheduled Tasks are skipped until existing Tasks are cleaned up
	// (by global KubeOpenCodeConfig cleanup or manual deletion).
	//
	// This acts as a safety valve to prevent unbounded Task accumulation.
	// Defaults to 10.
	// +optional
	// +kubebuilder:default=10
	// +kubebuilder:validation:Minimum=1
	MaxRetainedTasks *int32 `json:"maxRetainedTasks,omitempty"`

	// TaskTemplate is the template for Tasks created by this CronTask.
	// +required
	TaskTemplate TaskTemplateSpec `json:"taskTemplate"`
}

// TaskTemplateSpec describes the Task that will be created by a CronTask.
type TaskTemplateSpec struct {
	// Metadata is the metadata to apply to created Tasks.
	// Only labels and annotations are supported.
	// +optional
	Metadata TaskTemplateMeta `json:"metadata,omitempty"`

	// Spec is the TaskSpec for created Tasks.
	// +required
	Spec TaskSpec `json:"spec"`
}

// TaskTemplateMeta is a subset of metav1.ObjectMeta for Task templates.
// Only labels and annotations are supported; other metadata fields are
// controlled by the CronTask controller.
type TaskTemplateMeta struct {
	// Labels to apply to the created Task.
	// +optional
	Labels map[string]string `json:"labels,omitempty"`

	// Annotations to apply to the created Task.
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
}

// CronTaskStatus defines the observed state of CronTask.
type CronTaskStatus struct {
	// Active is the number of currently active (Pending/Queued/Running) Tasks
	// created by this CronTask.
	// +optional
	Active int32 `json:"active,omitempty"`

	// ActiveRefs is the list of references to currently active Tasks.
	// +optional
	ActiveRefs []corev1.ObjectReference `json:"activeRefs,omitempty"`

	// LastScheduleTime is the last time a Task was successfully scheduled
	// (created) by this CronTask.
	// +optional
	LastScheduleTime *metav1.Time `json:"lastScheduleTime,omitempty"`

	// LastSuccessfulTime is the last time a Task created by this CronTask
	// completed successfully (phase = Completed).
	// +optional
	LastSuccessfulTime *metav1.Time `json:"lastSuccessfulTime,omitempty"`

	// NextScheduleTime is the next calculated schedule time.
	// +optional
	NextScheduleTime *metav1.Time `json:"nextScheduleTime,omitempty"`

	// TotalExecutions is the total number of Tasks created by this CronTask
	// since its creation.
	// +optional
	TotalExecutions int64 `json:"totalExecutions,omitempty"`

	// Conditions represent the latest available observations of the CronTask's state.
	// +optional
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CronTaskList contains a list of CronTask
type CronTaskList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CronTask `json:"items"`
}
