/*
* Copyright 2025 Google LLC
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     https://www.apache.org/licenses/LICENSE-2.0
*
*     Unless required by applicable law or agreed to in writing, software
*     distributed under the License is distributed on an "AS IS" BASIS,
*     WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
*     See the License for the specific language governing permissions and
*     limitations under the License.
 */
package v1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:storageversions
// +kubebuilder:metadata:labels="addonmanager.kubernetes.io/mode=Reconcile"
// +kubebuilder:metadata:annotations="components.gke.io/layer=addon"
// +kubebuilder:resource:scope=Cluster,shortName=cc;ccs
// +kubebuilder:subresource:status

// ComputeClass is a way to impact Cluster Autoscaler scaling
// decisions based on user preferences. It gives control over preference of
// hardware to be selected by Cluster Autoscaler.
// Given ComputeClass affects only workloads using workload separation
// label equal to CCs name, except ComputeClass with name default
// which will be used for workloads not specifying any preferences.
type ComputeClass struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object metadata. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
	//
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Specification of the ComputeClass object.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#spec-and-status.
	// +required
	Spec ComputeClassSpec `json:"spec" protobuf:"bytes,2,name=spec"`
	// Status of the ComputeClass.
	//
	// +optional
	Status ComputeClassStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ComputeClassList is a list of ComputeClass objects.
type ComputeClassList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata.
	//
	// +optional
	metav1.ListMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`
	// Items, list of ComputeClass returned from API.
	//
	// +optional
	Items []ComputeClass `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// ComputeClassSpec is a specification of provisioning priorities and
// other autoscaling settings.
//
// +kubebuilder:validation:XValidation:rule="!has(oldSelf.autopilot) || has(self.autopilot)", message="Autopilot is required once set"
// +kubebuilder:validation:XValidation:rule="(has(self.autopilot) && self.autopilot.enabled) ? !self.priorities.exists(priority, has(priority.nodepools)) : true", message="Nodepools priority cannot be used when Autopilot is enabled"
// +kubebuilder:validation:XValidation:rule="(has(self.autopilot) && self.autopilot.enabled) ? !(has(self.nodePoolAutoCreation) && !self.nodePoolAutoCreation.enabled) : true", message="NodePoolAutoCreation cannot be disabled when Autopilot is enabled"
// +kubebuilder:validation:XValidation:rule="(has(self.autopilot) && self.autopilot.enabled) ? (!has(self.nodePoolConfig) || !has(self.nodePoolConfig.imageType) || self.nodePoolConfig.imageType == \"cos_containerd\") : true", message="Only cos_containerd image type can be used when Autopilot is enabled"
// +kubebuilder:validation:XValidation:rule="(has(self.nodePoolConfig) && has(self.nodePoolConfig.workloadType) && !has(self.nodePoolGroup)) ? self.nodePoolConfig.workloadType == \"HIGH_AVAILABILITY\" : true", message="If NodePoolGroup is not specified NodePoolConfig.WorkloadType can only be HIGH_AVAILABILITY if set"
// +kubebuilder:validation:XValidation:rule="self.priorities.exists(priority, has(priority.podFamily)) ? (has(self.autopilot) && self.autopilot.enabled) : true", message="In GKE Standard, pod family can be used only if Autopilot is enabled"
// +kubebuilder:validation:XValidation:rule="self.priorities.exists(priority, has(priority.podFamily)) ? !(has(self.nodePoolConfig)) : true", message="Pod family cannot be used with nodePoolConfig"
// +kubebuilder:validation:XValidation:rule="self.priorities.exists(priority, has(priority.podFamily)) ? !(has(self.priorityDefaults)) : true", message="Pod family cannot be used with priorityDefaults"
// +kubebuilder:validation:XValidation:rule="self.priorities.exists(priority, has(priority.podFamily)) ? (self.whenUnsatisfiable != \"ScaleUpAnyway\") : true", message="Pod family cannot be used with ScaleUpAnyway"
// +kubebuilder:validation:XValidation:rule="self.priorities.exists(priority, has(priority.podFamily)) ? self.priorities.all(priority, has(priority.podFamily)) : true", message="Pod family priorities cannot be mixed with non-pod family priorities"
// +kubebuilder:validation:XValidation:rule="(has(self.nodePoolConfig) && has(self.nodePoolConfig.confidentialNodeType)) ? self.priorities.all(priority, has(priority.machineFamily) || has(priority.machineType)) : true", message="If using NodePoolConfig.ConfidentialNodeType, each priority must specify either MachineFamily or MachineType."
// +kubebuilder:validation:XValidation:rule="(has(self.nodePoolConfig) && has(self.nodePoolConfig.confidentialNodeType) && self.nodePoolConfig.confidentialNodeType == \"SEV\") ? self.priorities.all(priority, ((has(priority.machineFamily) && priority.machineFamily in ['n2d', 'c2d', 'c3d', 'c4d']) || (has(priority.machineType) && priority.machineType.split('-')[0] in ['n2d', 'c2d', 'c3d', 'c4d']))) : true", message="ConfidentialNodeType SEV only supports N2D, C2D, C3D, C4D"
// +kubebuilder:validation:XValidation:rule="(has(self.nodePoolConfig) && has(self.nodePoolConfig.confidentialNodeType) && self.nodePoolConfig.confidentialNodeType == \"SEV_SNP\") ? self.priorities.all(priority, ((has(priority.machineFamily) && priority.machineFamily in ['n2d']) || (has(priority.machineType) && priority.machineType.split('-')[0] in ['n2d']))) : true", message="ConfidentialNodeType SEV_SNP only supports N2D"
// +kubebuilder:validation:XValidation:rule="(has(self.nodePoolConfig) && has(self.nodePoolConfig.confidentialNodeType) && self.nodePoolConfig.confidentialNodeType == \"TDX\") ? self.priorities.all(priority, ((has(priority.machineFamily) && priority.machineFamily in ['c3', 'a3']) || (has(priority.machineType) && (priority.machineType.startsWith('c3-standard-') || priority.machineType == 'a3-highgpu-1g')) || (has(priority.gpu) && has(priority.gpu.type) && priority.gpu.type == 'nvidia-h100-80gb'))) : true", message="ConfidentialNodeType TDX only supports C3 standard and A3 highgpu 1"
// +kubebuilder:validation:XValidation:rule="(has(self.nodePoolConfig) && has(self.nodePoolConfig.confidentialNodeType) && self.priorities.exists(priority, has(priority.gpu))) ? (has(self.priorityDefaults) && has(self.priorityDefaults.location) && has(self.priorityDefaults.location.zones)) || self.priorities.all(priority, has(priority.location) && has(priority.location.zones)) : true", message="When using confidential GPUs you must specify location.zones"
type ComputeClassSpec struct {
	// Priorities is a description of user preferences to be
	// used by a given ComputeClass.
	// +kubebuilder:listType=atomic
	// +kubebuilder:validation:MinItems=0
	// +kubebuilder:validation:MaxItems=1000
	// +kubebuilder:default={}
	// +optional
	Priorities []Priority `json:"priorities" protobuf:"bytes,1,name=priorities"`

	// NodePoolAutoCreation describes the auto provisioning settings for a given
	// ComputeClass.
	// +kubebuilder:default={enabled: false}
	// +optional
	NodePoolAutoCreation *NodePoolAutoCreation `json:"nodePoolAutoCreation,omitempty" protobuf:"bytes,2,opt,name=nodePoolAutoCreation"`

	// ActiveMigration describes settings related to active reconciliation of
	// a given ComputeClass.
	//
	// +optional
	ActiveMigration *ActiveMigration `json:"activeMigration,omitempty" protobuf:"bytes,3,opt,name=activeMigration"`

	// WhenUnsatisfiable describes autoscaler behaviour in case none
	// of the provided priorities is satisfiable.
	// Currently supported values:
	// * ScaleUpAnyway
	// * DoNotScaleUp
	//
	// +kubebuilder:validation:Enum=ScaleUpAnyway;DoNotScaleUp
	// +kubebuilder:default=DoNotScaleUp
	WhenUnsatisfiable string `json:"whenUnsatisfiable" protobuf:"bytes,4,name=whenUnsatisfiable"`

	// AutoscalingPolicy describes settings related to active reconciliation of
	// a given ComputeClass.
	// +optional
	AutoscalingPolicy *AutoscalingPolicy `json:"autoscalingPolicy,omitempty" protobuf:"bytes,5,opt,name=autoscalingPolicy"`

	// Autopilot describes the autopilot settings for a given ComputeClass.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="Autopilot is immutable"
	Autopilot *Autopilot `json:"autopilot,omitempty" protobuf:"bytes,6,opt,name=autopilot"`

	// NodePoolConfig defines required node pool configuration. Existing node pools will be matched with the ComputeClass
	// only if their configuration match this field. Auto-provisioned node pools will be created with this configuration.
	// +optional
	NodePoolConfig *NodePoolConfig `json:"nodePoolConfig,omitempty" protobuf:"bytes,7,opt,name=nodePoolConfig"`

	// NodePoolGroup defines required node pool configurations that are shared between a group of node pools.
	// Existing node pools will be matched with the ComputeClass only if their configuration matches this field.
	// Auto-provisioned node pools will be created with this configuration.
	// +optional
	NodePoolGroup *NodePoolGroup `json:"nodePoolGroup,omitempty" protobuf:"bytes,8,opt,name=nodePoolGroup"`

	// PriorityDefaults define the default rules for all priorities if the rule doesn't exist in some priority.
	// Note: PriorityDefaults doesn't apply to priorities with only Nodepools.
	//
	// +kubebuilder:validation:Optional
	PriorityDefaults *PriorityDefaults `json:"priorityDefaults,omitempty" protobuf:"bytes,9,opt,name=priorityDefaults"`

	// Description is an arbitrary string that usually provides guidelines on
	// when this compute class should be used.
	// +optional
	Description string `json:"description,omitempty" protobuf:"bytes,10,opt,name=description"`
}

// AutoscalingPolicy defines autoscaling related settings.
type AutoscalingPolicy struct {
	// ConsolidationDelayMinutes determines how long a node should be unneeded before it is eligible for scale down.
	// Minimum duration is 1 minute, maximum is 24 hours or 1440 minutes
	//
	// +kubebuilder:validation:Maximum=1440
	// +kubebuilder:validation:Minimum=1
	// +optional
	ConsolidationDelayMinutes *int `json:"consolidationDelayMinutes,omitempty" protobuf:"bytes,2,opt,name=consolidationDelayMinutes"`

	// ConsolidationThreshold determines resource utilization threshold below which a node can be considered for scale down.
	//
	// +kubebuilder:validation:Maximum=100
	// +kubebuilder:validation:Minimum=0
	// +optional
	ConsolidationThreshold *int `json:"consolidationThreshold,omitempty" protobuf:"bytes,3,opt,name=consolidationThreshold"`

	// GPUConsolidationThreshold determines GPU resource utilization threshold below which a node can be considered for scale down.
	// Utilization calculation only cares about GPU resource for accelerator node, CPU and memory utilization will be ignored.
	//
	// +kubebuilder:validation:Maximum=100
	// +kubebuilder:validation:Minimum=0
	// +optional
	GPUConsolidationThreshold *int `json:"gpuConsolidationThreshold,omitempty" protobuf:"bytes,4,opt,name=gpuConsolidationThreshold"`
}

// ActiveMigration describes if and what type of active migration
// should be performed.
type ActiveMigration struct {
	// OptimizeRulePriority defines whether workloads affected by given
	// ComputeClass should be migrated to nodepool defined by higher priority rule, if possible.
	//
	// +kubebuilder:validation:Required
	// +kubebuilder:default=false
	OptimizeRulePriority bool `json:"optimizeRulePriority" protobuf:"bytes,1,name=optimizeRulePriority"`

	// EnsureAllDaemonSetPodsRunning defines whether node pools should be migrated
	// to larger ones to ensure that all daemon sets are schedulable.
	//
	// +optional
	EnsureAllDaemonSetPodsRunning *bool `json:"ensureAllDaemonSetPodsRunning,omitempty" protobuf:"bytes,2,name=ensureAllDaemonSetPodsRunning"`
}

// NodePoolAutoCreation defines node-pool autoprovisioning related settings.
type NodePoolAutoCreation struct {
	// Enabled indicates whether NodePoolAutoCreation is enabled for a given ComputeClass.
	//
	// +kubebuilder:validation:Required
	// +kubebuilder:default=false
	Enabled bool `json:"enabled" protobuf:"bytes,1,name=enabled"`

	// DynamicMaxPodsPerNode if set to true specifies that max pods per node value for managed node pools will be selected
	// by Cluster Autoscaler automatically, based on the binpacking simulation results. It is ignored if there is a Priority.MaxPodsPerNode value specified.
	// If not specified the value defaults to being true for Compute Classes with Autopilot enabled.
	// If set to false cluster wide static value for max pods per node is used.
	//
	// +optional
	DynamicMaxPodsPerNode *bool `json:"dynamicMaxPodsPerNode,omitempty" protobuf:"bytes,2,opt,name=dynamicMaxPodsPerNode"`

	// DynamicBootDiskSize if set to true specifies that boot disk size value for managed node pools will be selected
	// by Cluster Autoscaler automatically, based on the binpacking simulation results. It is ignored if there is a Priority.Storage.BootDiskSize value specified.
	// If not specified the value defaults to being true for Compute Classes with Autopilot enabled.
	// If set to false cluster wide static value from AutoprovisioningNodePoolDefaults is used.
	//
	// +optional
	DynamicBootDiskSize *bool `json:"dynamicBootDiskSize,omitempty" protobuf:"bytes,3,opt,name=dynamicBootDiskSize"`
}

// Autopilot defines describes the autopilot settings for a given ComputeClass.
type Autopilot struct {
	// Enabled indicates whether nodes created for this compute class should be Autopilot managed.
	//
	// +kubebuilder:validation:Required
	// +kubebuilder:default=false
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="Autopilot is immutable"
	Enabled bool `json:"enabled" protobuf:"bytes,1,name=enabled"`
}

// NodePoolConfig defines required node pool configuration. Existing node pools will be matched with the ComputeClass
// only if their configuration match this field. Auto-provisioned node pools will be created with this configuration.
type NodePoolConfig struct {
	// ServiceAccount used by the node pool.
	//
	// +optional
	ServiceAccount string `json:"serviceAccount,omitempty" protobuf:"string,1,name=serviceAccount"`

	// Image type used by nodes in the node pool.
	//
	// +kubebuilder:validation:Enum=cos_containerd;ubuntu_containerd
	// +optional
	ImageType string `json:"imageType,omitempty" protobuf:"string,2,name=imageType"`

	// WorkloadType defines Collection or Goodput SLO for the workload. Currently
	// supported values:
	// * HIGH_AVAILABILITY - for Collection SLO
	// * HIGH_THROUGHPUT - for Goodput SLO
	// HIGH_AVAILABILITY is desired for running serving workloads which require
	// most of the infrastructure (slices) running all the time to achieve high
	// availability.
	// HIGH_THROUGHPUT is desired for running batch/training jobs
	// which require all underlying infrastructure (slices) running for most of
	// the time to make progress. HIGH_THROUGHPUT can be only set for a multi-host
	// scenario, that is, when NodePoolGroup is set.
	//
	// +optional
	// +kubebuilder:validation:Enum=HIGH_AVAILABILITY;HIGH_THROUGHPUT
	WorkloadType string `json:"workloadType,omitempty" protobuf:"bytes,3,opt,name=workloadType"`

	// NodeLabels is used to add user defined Kubernetes labels to all nodes in the new node pool.
	// These labels are applied to the Kubernetes API node object and can be used in nodeSelectors for pod scheduling.
	// Note: Node labels are distinct from GKE labels.
	// More info: https://cloud.google.com/sdk/gcloud/reference/container/node-pools/create#--node-labels
	//
	// +optional
	// +kubebuilder:validation:MaxProperties=100
	// +kubebuilder:validation:XValidation:rule="self.all(key, key.matches('^[a-zA-Z0-9][a-zA-Z0-9-._]{0,62}$'))",message="Node labels must begin with a letter or number, and may contain letters, numbers, hyphens, dots, and underscores, up to 63 characters each."
	// +kubebuilder:validation:XValidation:rule="self.all(key, !key.contains('cloud.google.com'))",message="Node labels cannot contain reserved `cloud.google.com` substring."
	// +kubebuilder:validation:XValidation:rule="self.all(key, !key.contains('kubernetes.io'))",message="Node labels cannot contain reserved `kubernetes.io` substring."
	// +kubebuilder:validation:XValidation:rule="self.all(key, !key.contains('gke.io'))",message="Node labels cannot contain reserved `gke.io` substring."
	// +kubebuilder:validation:XValidation:rule="self.all(key, !key.contains('k8s.io'))",message="Node labels cannot contain reserved `k8s.io` substring."
	NodeLabels map[string]string `json:"nodeLabels,omitempty" protobuf:"bytes,4,opt,name=nodeLabels"`

	// Taints is used to add user defined Kubernetes taints to all nodes in the new node pool.
	// These taints are applied to the Kubernetes API node object and can be used in tolerations for pod scheduling.
	//
	// +optional
	// +kubebuilder:validation:MaxItems=100
	Taints []TaintConfig `json:"taints,omitempty" protobuf:"bytes,5,opt,name=taints"`

	// ConfidentialNodeType: Defines the type of technology used by the
	// confidential node.
	//
	// Possible values:
	//   "CONFIDENTIAL_INSTANCE_TYPE_UNSPECIFIED" - No type specified. Do not use
	// this value.
	//   "SEV" - AMD Secure Encrypted Virtualization.
	//   "SEV_SNP" - AMD Secure Encrypted Virtualization - Secure Nested Paging.
	//   "TDX" - Intel Trust Domain eXtension.
	// +kubebuilder:validation:Enum=CONFIDENTIAL_INSTANCE_TYPE_UNSPECIFIED;SEV;SEV_SNP;TDX
	// +optional
	ConfidentialNodeType string `json:"confidentialNodeType,omitempty" protobuf:"string,6,opt,name=confidentialNodeType"`

	// AutoRepair if set to true specifies that a node pool should have auto repair enabled, disabled in case of being set to false.
	//
	// +optional
	AutoRepair *bool `json:"autoRepair,omitempty" protobuf:"bytes,7,opt,name=autoRepair"`

	// AutoUpgrade if set to true specifies that a node pool should have auto upgrade enabled, disabled in case of being set to false.
	//
	// +optional
	AutoUpgrade *bool `json:"autoUpgrade,omitempty" protobuf:"bytes,8,opt,name=autoUpgrade"`

	// ImageStreaming contains image streaming settings.
	//
	// +optional
	ImageStreaming *ImageStreaming `json:"imageStreaming,omitempty" protobuf:"bytes,9,opt,name=imageStreaming"`

	// ResourceManagerTags defines what existing GCE resource manager tag key/value pairs
	// with purpose GCE_FIREWALL to attach to all node pools.
	// Referenced Tags must be created beforehand via Resource Manager API.
	// +kubebuilder:validation:MaxItems=5
	// +optional
	ResourceManagerTags []Tags `json:"resourceManagerTags,omitempty" protobuf:"bytes,10,opt,name=resourceManagerTags"`

	// Gvnic contains Google Virtual NIC settings.
	// +optional
	Gvnic *Gvnic `json:"gvnic,omitempty" protobuf:"bytes,11,opt,name=gvnic"`

	// Contains logging configuration.
	// +optional
	LoggingConfig *NodePoolLoggingConfig `json:"loggingConfig,omitempty" protobuf:"bytes,12,opt,name=loggingConfig"`
}

// NodePoolLoggingConfig specifies logging configuration for nodepools.
type NodePoolLoggingConfig struct {
	// Logging variant configuration.
	// +optional
	LoggingVariantConfig *LoggingVariantConfig `json:"loggingVariantConfig,omitempty" protobuf:"bytes,1,opt,name=loggingVariantConfig"`
}

// LoggingVariantConfig specifies logging variant configuration.
type LoggingVariantConfig struct {
	// Logging variant deployed on nodes.
	// +optional
	// +kubebuilder:validation:Enum=DEFAULT;MAX_THROUGHPUT
	Variant string `json:"variant,omitempty" protobuf:"string,1,opt,name=variant"`
}

// Gvnic stores Google Virtual NIC settings.
type Gvnic struct {
	// Enabled indicates whether gVNIC is enabled on the node pool.
	//
	// +kubebuilder:validation:Required
	// +kubebuilder:default=false
	Enabled bool `json:"enabled" protobuf:"bytes,1,name=enabled"`
}

// ImageStreaming stores container image streaming settings. It is equivalent to `GcfsConfig` in GKE.
// https://cloud.google.com/kubernetes-engine/docs/reference/rest/v1/GcfsConfig
type ImageStreaming struct {
	// Enabled enables container image` streaming.
	//
	// +kubebuilder:validation:Required
	// +kubebuilder:default=false
	Enabled bool `json:"enabled" protobuf:"bytes,1,name=enabled"`
}

// NodePoolGroup defines required node pool configurations that are shared between a group of node pools. It is
// GKE equivalent of GCE's Multi-MIG. Existing node pools will be matched with the ComputeClass only if their configuration
// matches this field. Auto-provisioned node pools will be created with this configuration.
type NodePoolGroup struct {
	// Name defines the name of the node pool group, e.g. MultiMIG
	//
	// +required
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name" protobuf:"bytes,1,name=name"`
}

// Storage defines storage config per priority rule.
type Storage struct {
	// BootDiskSize defines the size of a disk attached to node, specified in GB.
	//
	// +optional
	// +kubebuilder:validation:Minimum=10
	BootDiskSize *int `json:"bootDiskSize,omitempty" protobuf:"bytes,1,opt,name=bootDiskSize"`

	// BootDiskType defines type of the disk attached to the node.
	// Note that available boot disk types depend on the machine family / machine type selected.
	// Currently supported types:
	// * pd-balanced
	// * pd-standard
	// * pd-ssd
	// * hyperdisk-balanced
	//
	// +kubebuilder:validation:Enum=pd-balanced;pd-standard;pd-ssd;hyperdisk-balanced
	// +optional
	BootDiskType *string `json:"bootDiskType,omitempty" protobuf:"bytes,2,opt,name=bootDiskType"`
	// BootDiskKMSKey defines a key used to encrypt the boot disk attached.
	//
	// +optional
	// +kubebuilder:validation:Pattern=projects/[^/]+/locations/[^/]+/keyRings/[^/]+/cryptoKeys/[^/]+
	BootDiskKMSKey *string `json:"bootDiskKMSKey,omitempty" protobuf:"bytes,3,opt,name=bootDiskKMSKey"`
	// LocalSSDCount defines a number of local SSDs attached to node.
	//
	// +optional
	// +kubebuilder:validation:Minimum=1
	LocalSSDCount *int `json:"localSSDCount,omitempty" protobuf:"bytes,4,opt,name=localSSDCount"`
	// SecondaryBootDisks represent persistent disks attached to a node with special configurations based on their modes.
	//
	// +optional
	SecondaryBootDisks []SecondaryBootDisk `json:"secondaryBootDisks,omitempty" protobuf:"bytes,5,opt,name=secondaryBootDisks"`
}

// SecondaryBootDisk represents a persistent disk attached to a node with special configurations based on its mode.
type SecondaryBootDisk struct {
	// The name of the disk image.
	//
	// +required
	DiskImageName string `json:"diskImageName" protobuf:"bytes,1,name=diskImageName"`
	// The name of the project that the disk image belongs to.
	//
	// +optional
	Project *string `json:"project,omitempty" protobuf:"bytes,2,opt,name=project"`
	// Currently supported modes:
	// * MODE_UNSPECIFIED - MODE_UNSPECIFIED is when mode is not set.
	// * CONTAINER_IMAGE_CACHE - it is for using the secondary boot disk as a container image cache.
	//
	// +optional
	// +kubebuilder:validation:Enum=MODE_UNSPECIFIED;CONTAINER_IMAGE_CACHE
	Mode *string `json:"mode,omitempty" protobuf:"bytes,3,opt,name=mode"`
}

// SpecificReservation defines a single specific reservation to be consumed by the created node.
type SpecificReservation struct {
	// Name of the reservation to be used.
	Name string `json:"name" protobuf:"bytes,1,name=name"`
	// Project is the project where the specific reservation lives.
	//
	// +optional
	Project string `json:"project,omitempty" protobuf:"bytes,2,opt,name=project"`
	// ReservationBlock is the block of the reservation.
	//
	// +optional
	ReservationBlock *ReservationBlock `json:"reservationBlock,omitempty" protobuf:"bytes,3,opt,name=reservationBlock"`
}

// ReservationBlock is the block of the reservation.
type ReservationBlock struct {
	// Name is the name of the block.
	//
	// +required
	Name string `json:"name" protobuf:"bytes,1,name=name"`
	// ReservationSubBlock is the subBlock of the reservation block.
	//
	// +optional
	ReservationSubBlock *ReservationSubBlock `json:"reservationSubBlock,omitempty" protobuf:"bytes,2,opt,name=reservationSubBlock"`
}

// ReservationSubBlock is the subBlock of the reservation block.
type ReservationSubBlock struct {
	// Name is the name of the subBlock.
	//
	// +required
	Name string `json:"name" protobuf:"bytes,1,name=name"`
}

// ReservationAffinity is an enumeration of supported reservation affinities
//
// +kubebuilder:validation:Enum=Specific;AnyBestEffort;None
type ReservationAffinity string

const (
	// SpecificAffinity affinity allows to consume only specific reservations.
	SpecificAffinity ReservationAffinity = "Specific"
	// AnyBestEffortAffinity affinity allows to consume any reservation with a possibility to fallback to on demand.
	AnyBestEffortAffinity ReservationAffinity = "AnyBestEffort"
	// NoneAffinity prevents reservations from being used.
	NoneAffinity ReservationAffinity = "None"
)

// Reservations define reservations configuration per priority rule.
//
// +kubebuilder:validation:XValidation:message="Unable to set specific reservations for non specific affinity",rule="has(self.specific) && self.specific.size() > 0 ? self.affinity == \"Specific\" : true"
// +kubebuilder:validation:XValidation:message="At least 1 specific reservation required for specific affinity",rule="self.affinity == \"Specific\" ? has(self.specific) && self.specific.size() > 0 : true"
type Reservations struct {
	// Specific is a non prioritized list of specific reservations to be considered by the priority rule.
	//
	// +kubebuilder:listType=atomic
	// +kubebuilder:validation:MinItems=0
	// +optional
	Specific []SpecificReservation `json:"specific,omitempty" protobuf:"bytes,1,opt,name=specific"`

	// ReservationAffinity affects reservations considered and the way how they are consumed.
	// "Specific" means that only specific reservations are considered with no fallback possible.
	// "AnyBestEffort" affinity would consider any non-specific reservation available
	// to be claimed with a fallback to on-demand nodes in case of none claimable.
	// "None" affinity would prevent reservations from being used
	//
	// +required
	Affinity ReservationAffinity `json:"affinity" protobuf:"bytes,2,name=affinity"`
}

// Priority is a specification of preferred machine characteristics.
//
// +kubebuilder:validation:MinProperties=1
// +kubebuilder:validation:XValidation:rule="has(self.nodepools) ? (size(dyn(self)) == 1) : true", message="Nodepool field cannot be set along with other fields"
// +kubebuilder:validation:XValidation:rule="!(has(self.machineFamily) && has(self.machineType))",message="MachineFamily and MachineType cannot be set together"
// +kubebuilder:validation:XValidation:rule="!(has(self.machineType) && (has(self.minCores) || has(self.minMemoryGb)))",message="MachineType cannot be set together with MinCores/MinMemoryGb"
// +kubebuilder:validation:XValidation:rule="has(self.reservations) && !has(self.tpu) ? has(self.machineFamily) || has(self.machineType) : true",message="Non-TPU reservations can be used only with machine type or machine family defined"
// +kubebuilder:validation:XValidation:rule="!(has(self.machineFamily) && self.machineFamily == 'ek')", message="MachineFamily cannot be equal to 'ek'"
// +kubebuilder:validation:XValidation:rule="!(has(self.machineType) && self.machineType.startsWith('ek'))", message="MachineType cannot start with 'ek' prefix"
// +kubebuilder:validation:XValidation:rule="!(has(self.flexStart) && has(self.spot) && self.spot == true && self.flexStart.enabled == true)", message="Flex Start provisioning model is incompatible with Spot"
// +kubebuilder:validation:XValidation:rule="has(self.podFamily) ? (has(self.spot) ? (size(dyn(self)) == 2) : (size(dyn(self)) == 1)) : true", message="Spot selection is the only configurable priority setting when using podFamily"
// +kubebuilder:validation:XValidation:rule="!has(self.capacityCheckWaitTimeSeconds) || has(self.tpu) || (has(self.flexStart) && self.flexStart.enabled)", message="capacityCheckWaitTimeSeconds is only supported for Flex Start and for multi-host TPUs"
type Priority struct {
	// Machine family describes preferred instance family for a node. If none is specified,
	// the default autoprovisioning machine family is used.
	//
	// +optional
	// +kubebuilder:validation:MaxLength=10
	MachineFamily *string `json:"machineFamily,omitempty" protobuf:"bytes,1,opt,name=machineFamily"`
	// Spot if set to true specifies that a node should be a spot instance, on-demand otherwise.
	//
	// +optional
	Spot *bool `json:"spot,omitempty" protobuf:"bytes,2,opt,name=spot"`
	// MinCores describes a minimum number of CPU cores of a node.
	//
	// +optional
	// +kubebuilder:validation:Minimum=0
	MinCores *int `json:"minCores,omitempty" protobuf:"bytes,3,opt,name=minCores"`
	// MinMemoryGb describes a minimum GBs of memory of a node.
	//
	// +optional
	// +kubebuilder:validation:Minimum=0
	MinMemoryGb *int `json:"minMemoryGb,omitempty" protobuf:"bytes,4,opt,name=minMemoryGb"`
	// Nodepools describes preference of specific, preexisting nodepools.
	//
	// +optional
	Nodepools []string `json:"nodepools,omitempty" protobuf:"bytes,5,opt,name=nodepools"`
	// Storage describes storage config of a node.
	//
	// +optional
	Storage *Storage `json:"storage,omitempty" protobuf:"bytes,6,opt,name=storage"`

	// MachineType defines preferred machine type for a node.
	//
	// +optional
	// +kubebuilder:validation:MaxLength=100
	MachineType *string `json:"machineType,omitempty" protobuf:"bytes,7,opt,name=machineType"`

	// Gpu defines preferred GPU config for a node.
	//
	// +optional
	Gpu *GPU `json:"gpu,omitempty" protobuf:"bytes,8,opt,name=gpu"`

	// Tpu defines preferred TPU config for a node.
	//
	// +optional
	Tpu *TPU `json:"tpu,omitempty" protobuf:"bytes,9,opt,name=tpu"`

	// Reservations defines reservations config for a node.
	//
	// +optional
	Reservations *Reservations `json:"reservations,omitempty" protobuf:"bytes,10,opt,name=reservations"`

	// MaxRunDurationSeconds defines the maximum duration for the nodes to exist. If unspecified, the nodes can exist indefinitely.
	//
	// +optional
	MaxRunDurationSeconds *int `json:"maxRunDurationSeconds,omitempty" protobuf:"bytes,11,opt,name=maxRunDurationSeconds"`

	// MaxPodsPerNode describes the maximum number of pods a node can accommodate.
	//
	// +optional
	// +kubebuilder:validation:Minimum=8
	// +kubebuilder:validation:Maximum=256
	MaxPodsPerNode *int `json:"maxPodsPerNode,omitempty" protobuf:"bytes,12,opt,name=maxPodsPerNode"`

	// NodeSystemConfig defines node system config for a node.
	//
	// +kubebuilder:validation:Optional
	NodeSystemConfig *NodeSystemConfig `json:"nodeSystemConfig,omitempty" protobuf:"bytes,13,opt,name=nodeSystemConfig"`

	// FlexStart defines Flex Start provisioning model.
	//
	// +kubebuilder:validation:Optional
	FlexStart *FlexStart `json:"flexStart,omitempty" protobuf:"bytes,14,opt,name=flexStart"`

	// PodFamily represents pod-based provisioning and billing config.
	//
	// +optional
	// +kubebuilder:validation:Enum=general-purpose
	PodFamily *string `json:"podFamily,omitempty" protobuf:"bytes,15,opt,name=podFamily"`

	// Location describes CCC zonal preferences config.
	//
	// +optional
	Location *Location `json:"location,omitempty" protobuf:"bytes,16,opt,name=location"`

	// Placement defines resource policy used for BYOPP and BYOWP
	//
	// +kubebuilder:validation:Optional
	Placement *Placement `json:"placement,omitempty" protobuf:"bytes,17,opt,name=placement"`

	// CapacityCheckWaitTimeSeconds defines for how long will this priority be attempted to scale up before moving on to the next priority.
	//
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=86400
	CapacityCheckWaitTimeSeconds *int `json:"capacityCheckWaitTimeSeconds,omitempty" protobuf:"bytes,18,opt,name=capacityCheckWaitTimeSeconds"`

	// MinCpuPlatform defines the minimum CPU platform for a node.
	//
	// +optional
	// +kubebuilder:validation:Enum={Intel Sandy Bridge,Intel Ivy Bridge,Intel Haswell,Intel Broadwell,Intel Skylake,Intel Cascade Lake,Intel Ice Lake,Intel Sapphire Rapids,Intel Emerald Rapids,Intel Granite Rapids,AMD Rome,AMD Milan,AMD Genoa,AMD Turin,Ampere Altra,Google Axion,Nvidia Grace}
	MinCpuPlatform *string `json:"minCpuPlatform,omitempty" protobuf:"bytes,19,opt,name=minCpuPlatform"`

	// NodeLabels is used to add user defined Kubernetes labels to all nodes in the new node pool.
	// These labels are applied to the Kubernetes API node object and can be used in nodeSelectors for pod scheduling.
	// Note: Node labels are distinct from GKE labels.
	// More info: https://cloud.google.com/sdk/gcloud/reference/container/node-pools/create#--node-labels
	//
	// +optional
	// +kubebuilder:validation:MaxProperties=20
	// +kubebuilder:validation:XValidation:rule="self.all(key, key.matches('^[a-zA-Z0-9][a-zA-Z0-9-._]{0,62}$'))",message="Node labels must begin with a letter or number, and may contain letters, numbers, hyphens, dots, and underscores, up to 63 characters each."
	// +kubebuilder:validation:XValidation:rule="self.all(key, !key.contains('cloud.google.com'))",message="Node labels cannot contain reserved `cloud.google.com` substring."
	// +kubebuilder:validation:XValidation:rule="self.all(key, !key.contains('kubernetes.io'))",message="Node labels cannot contain reserved `kubernetes.io` substring."
	// +kubebuilder:validation:XValidation:rule="self.all(key, !key.contains('gke.io'))",message="Node labels cannot contain reserved `gke.io` substring."
	// +kubebuilder:validation:XValidation:rule="self.all(key, !key.contains('k8s.io'))",message="Node labels cannot contain reserved `k8s.io` substring."
	NodeLabels map[string]string `json:"nodeLabels,omitempty" protobuf:"bytes,20,opt,name=nodeLabels"`

	// Taints is used to add user defined Kubernetes taints to all nodes in the new node pool.
	// These taints are applied to the Kubernetes API node object and can be used in tolerations for pod scheduling.
	//
	// +optional
	// +kubebuilder:validation:MaxItems=20
	Taints []TaintConfig `json:"taints,omitempty" protobuf:"bytes,21,opt,name=taints"`

	// HostMaintenancePolicy contains the desired maintenance policy for the
	// Google Compute Engine hosts.
	//
	// +kubebuilder:validation:Optional
	HostMaintenancePolicy *HostMaintenancePolicy `json:"hostMaintenancePolicy,omitempty" protobuf:"bytes,22,opt,name=hostMaintenancePolicy"`
}

// Placement describes preference of Resource Policy for BYOPP
type Placement struct {
	// PolicyName defines the name of the resource policy, e.g. my-resource-policy
	//
	// +required
	// +kubebuilder:validation:MinLength=1
	PolicyName string `json:"policyName" protobuf:"bytes,1,name=policyName"`
}

// GPU describes preference on given GPU config.
type GPU struct {
	// Type describes preferred GPU accelerator type for a node.
	Type string `json:"type,omitempty" protobuf:"bytes,1,name=type"`
	// Count describes preferred count of GPUs for a node.
	// +kubebuilder:validation:Minimum=0
	Count int64 `json:"count,omitempty" protobuf:"bytes,2,name=count"`
	// DriverVersion describes version of GPU driver for a node.
	// +kubebuilder:validation:Enum=default;latest
	// +kubebuilder:default=default
	// +optional
	DriverVersion string `json:"driverVersion,omitempty" protobuf:"bytes,3,name=driverVersion"`
	// The topology defines the physical arrangement of GPUs chips within a slice.
	// +optional
	Topology string `json:"topology,omitempty" protobuf:"bytes,4,name=topology"`

	// GpuSharing defines the way the nodes would share the GPU.
	// +optional
	GpuSharing *GpuSharing `json:"gpuSharing,omitempty" protobuf:"bytes,5,name=gpuSharing"`
}

// TPU describes preference on given TPU config.
type TPU struct {
	// Type describes preferred TPU type for a node.
	Type string `json:"type,omitempty" protobuf:"bytes,1,name=type"`
	// Count describes preferred count of TPU chips for a node.
	Count int64 `json:"count,omitempty" protobuf:"bytes,2,name=count"`
	// Topology describes preferred TPU topology of a node.
	Topology string `json:"topology,omitempty" protobuf:"bytes,3,name=topology"`
}

// FlexStart defines Flex Start provisioning model.
type FlexStart struct {
	// Enabled indicates whether Flex Start provisioning model is enabled.
	//
	// +kubebuilder:validation:Required
	// +kubebuilder:default=false
	Enabled bool `json:"enabled" protobuf:"bytes,1,name=enabled"`

	// NodeRecycling defines node recycling config.
	//
	// +kubebuilder:validation:Optional
	NodeRecycling *NodeRecyclingConfig `json:"nodeRecycling,omitempty" protobuf:"bytes,2,opt,name=nodeRecycling"`
}

// NodeRecyclingConfig defines node recycling config.
type NodeRecyclingConfig struct {
	// LeadTimeSeconds defines how much time before node termination timestamp CA should start looking for a replacement node.
	//
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=604800
	// +kubebuilder:validation:Required
	LeadTimeSeconds *int `json:"leadTimeSeconds" protobuf:"bytes,1,name=leadTimeSeconds"`
}

// NodeSystemConfig defines node system config for a node.
type NodeSystemConfig struct {
	LinuxNodeConfig *LinuxNodeConfig `json:"linuxNodeConfig,omitempty" protobuf:"bytes,1,opt,name=linuxNodeConfig"`
	KubeletConfig   *KubeletConfig   `json:"kubeletConfig,omitempty" protobuf:"bytes,2,opt,name=kubeletConfig"`
}

// LinuxNodeConfig defines linux node config for a node.
type LinuxNodeConfig struct {
	Sysctls   *SysctlsConfig   `json:"sysctls,omitempty" protobuf:"bytes,1,opt,name=sysctls"`
	Hugepages *HugepagesConfig `json:"hugepageConfig,omitempty" protobuf:"bytes,2,opt,name=hugepageConfig"`
	// Controls transparent hugepage support for anonymous memory. Currently supported values:
	// * TRANSPARENT_HUGEPAGE_ENABLED_ALWAYS: Transparent hugepage is enabled system wide.
	// * TRANSPARENT_HUGEPAGE_ENABLED_MADVISE: Transparent hugepage is enabled inside MADV_HUGEPAGE regions. This is the default kernel configuration.
	// * TRANSPARENT_HUGEPAGE_ENABLED_NEVER: Transparent hugepage is disabled.
	// * TRANSPARENT_HUGEPAGE_ENABLED_UNSPECIFIED: Default value. GKE will not modify the kernel configuration.
	//
	// +optional
	// +kubebuilder:validation:Enum=TRANSPARENT_HUGEPAGE_ENABLED_ALWAYS;TRANSPARENT_HUGEPAGE_ENABLED_MADVISE;TRANSPARENT_HUGEPAGE_ENABLED_NEVER;TRANSPARENT_HUGEPAGE_ENABLED_UNSPECIFIED
	TransparentHugepageEnabled *string `json:"transparentHugepageEnabled,omitempty" protobuf:"bytes,3,opt,name=transparentHugepageEnabled"`
	// Defines the transparent hugepage defrag configuration on the node. Currently supported values:
	// * TRANSPARENT_HUGEPAGE_DEFRAG_ALWAYS: An application requesting THP will stall on allocation failure and directly reclaim pages and compact memory in an effort to allocate a THP immediately.
	// * TRANSPARENT_HUGEPAGE_DEFRAG_DEFER: An application will wake kswapd in the background to reclaim pages and wake kcompactd to compact memory so that THP is available in the near future. It is the responsibility of khugepaged to then install the THP pages later.
	// * TRANSPARENT_HUGEPAGE_DEFRAG_DEFER_WITH_MADVISE: An application will enter direct reclaim and compaction like always, but only for regions that have used madvise(MADV_HUGEPAGE); all other regions will wake kswapd in the background to reclaim pages and wake kcompactd to compact memory so that THP is available in the near future.
	// * TRANSPARENT_HUGEPAGE_DEFRAG_MADVISE: An application will enter direct reclaim and compaction like always, but only for regions that have used madvise(MADV_HUGEPAGE); all other regions will wake kswapd in the background to reclaim pages and wake kcompactd to compact memory so that THP is available in the near future.
	// * TRANSPARENT_HUGEPAGE_DEFRAG_NEVER: An application will never enter direct reclaim or compaction.
	// * TRANSPARENT_HUGEPAGE_DEFRAG_UNSPECIFIED: Default value. GKE will not modify the kernel configuration.
	//
	// +optional
	// +kubebuilder:validation:Enum=TRANSPARENT_HUGEPAGE_DEFRAG_ALWAYS;TRANSPARENT_HUGEPAGE_DEFRAG_DEFER;TRANSPARENT_HUGEPAGE_DEFRAG_DEFER_WITH_MADVISE;TRANSPARENT_HUGEPAGE_DEFRAG_MADVISE;TRANSPARENT_HUGEPAGE_DEFRAG_NEVER;TRANSPARENT_HUGEPAGE_DEFRAG_UNSPECIFIED
	TransparentHugepageDefrag *string `json:"transparentHugepageDefrag,omitempty" protobuf:"bytes,4,opt,name=transparentHugepageDefrag"`
}

type TopologyManager struct {
	// Policy controls the Kubelet's Topology Manager policy.
	// Policies:
	// * none: (default) The Kubelet does not perform any topology alignment.
	// * best-effort: The Kubelet will attempt to align resources but will not fail pod admission.
	// * restricted: The Kubelet will reject pods that do not align to the minimal number of NUMA domains.
	// * single-numa-node: The Kubelet will reject pods that do not align to a single NUMA domain.
	//
	// +kubebuilder:validation:Enum=none;best-effort;restricted;single-numa-node
	// +kubebuilder:validation:Optional
	Policy *string `json:"policy,omitempty" protobuf:"bytes,1,opt,name=policy"`
	// Scope controls the Kubelet's Topology Manager scope.
	// Scopes:
	// * container: (default) The Kubelet performs topology alignment for each container in a pod.
	// * pod: The Kubelet performs topology alignment for the pod as a whole. This setting ensures pod-level topology alignment, where the Topology Manager treats all containers as a single unit to place them on a common set of NUMA nodes.
	//
	// +kubebuilder:validation:Enum=container;pod
	// +kubebuilder:validation:Optional
	Scope *string `json:"scope,omitempty" protobuf:"bytes,2,opt,name=scope"`
}

// MemoryManagerConfig defines the configuration for the Kubelet Memory Manager.
type MemoryManager struct {
	// Policy controls the Kubelet's Memory Manager policy.
	// The Static policy is required for the Topology Manager to perform memory affinity alignment.
	// Policies:
	// * None: (default) The Kubelet does not perform any memory alignment.
	// * Static: The Kubelet allows pods in the Guaranteed QoS class to be granted memory from a single NUMA node.
	//
	// +kubebuilder:validation:Enum=None;Static
	// +kubebuilder:validation:Optional
	Policy *string `json:"policy,omitempty" protobuf:"bytes,1,opt,name=policy"`
}

// EvictionSoft is a map of signal names to quantities that defines soft eviction thresholds.
// A soft eviction threshold pairs with a grace period. The kubelet does not evict pods until the grace period is exceeded.
// +kubebuilder:validation:Optional
type EvictionSoft struct {
	// MemoryAvailable is the soft eviction threshold for memory.available.
	// The value must be a quantity, e.g., "100Mi".
	// The value must be greater than the GKE default hard eviction threshold of 100Mi and less than 50% of machine memory.
	// +kubebuilder:validation:Pattern=`^[0-9]+(\.[0-9]+)?(Ki|Mi|Gi)$`
	// +optional
	MemoryAvailable *string `json:"memoryAvailable,omitempty" protobuf:"bytes,1,opt,name=memoryAvailable"`

	// NodefsAvailable is the soft eviction threshold for nodefs.available.
	// The value must be a percentage, e.g., "20%".
	// The value must be between 10% and 50% inclusive.
	// +kubebuilder:validation:Pattern=`^[0-9]+(\.[0-9]+)?%$`
	// +optional
	NodefsAvailable *string `json:"nodefsAvailable,omitempty" protobuf:"bytes,2,opt,name=nodefsAvailable"`
	// ImagefsAvailable is the soft eviction threshold for imagefs.available.
	// The value must be a percentage. Eg. "10%".
	// The value must be between 15% and 50% inclusive.
	// +kubebuilder:validation:Pattern=`^[0-9]+(\.[0-9]+)?%$`
	// +optional
	ImagefsAvailable *string `json:"imagefsAvailable,omitempty" protobuf:"bytes,3,opt,name=imagefsAvailable"`
	// ImagefsInodesFree is the soft eviction threshold for imagefs.inodesFree.
	// The value must be a percentage. Eg. "5%".
	// The value must be between 5% and 50% inclusive.
	// +kubebuilder:validation:Pattern=`^[0-9]+(\.[0-9]+)?%$`
	// +optional
	ImagefsInodesFree *string `json:"imagefsInodesFree,omitempty" protobuf:"bytes,4,opt,name=imagefsInodesFree"`
	// NodefsInodesFree is the soft eviction threshold for nodefs.inodesFree.
	// The value must be a percentage. Eg. "5%".
	// The value must be between 5% and 50% inclusive.
	// +kubebuilder:validation:Pattern=`^[0-9]+(\.[0-9]+)?%$`
	// +optional
	NodefsInodesFree *string `json:"nodefsInodesFree,omitempty" protobuf:"bytes,5,opt,name=nodefsInodesFree"`
	// PidAvailable is the soft eviction threshold for pid.available.
	// The value must be a percentage. Eg. "10%".
	// The value must be between 10% and 50% inclusive.
	// +kubebuilder:validation:Pattern=`^[0-9]+(\.[0-9]+)?%$`
	// +optional
	PidAvailable *string `json:"pidAvailable,omitempty" protobuf:"bytes,6,opt,name=pidAvailable"`
}

// EvictionSoftGracePeriod is a map of signal names to durations that defines grace periods for soft eviction thresholds.
// Each soft eviction threshold must have a corresponding grace period.
// +kubebuilder:validation:Optional
type EvictionSoftGracePeriod struct {
	// MemoryAvailable is the grace period for the memory.available soft eviction threshold.
	// The value must be a duration string. Eg. "30s", "1m30s".
	// The value must be positive and less than '5m'.
	// +kubebuilder:validation:Pattern=`^([0-9]+([.][0-9]+)?(ns|us|µs|ms|s|m|h))+$`
	// +optional
	MemoryAvailable *string `json:"memoryAvailable,omitempty" protobuf:"bytes,1,opt,name=memoryAvailable"`
	// NodefsAvailable is the grace period for the nodefs.available soft eviction threshold.
	// The value must be a duration string. Eg. "30s", "1m30s".
	// The value must be positive and less than '5m'.
	// +kubebuilder:validation:Pattern=`^([0-9]+([.][0-9]+)?(ns|us|µs|ms|s|m|h))+$`
	// +optional
	NodefsAvailable *string `json:"nodefsAvailable,omitempty" protobuf:"bytes,2,opt,name=nodefsAvailable"`
	// ImagefsAvailable is the grace period for the imagefs.available soft eviction threshold.
	// The value must be a duration string. Eg. "30s", "1m30s".
	// The value must be positive and less than '5m'.
	// +kubebuilder:validation:Pattern=`^([0-9]+([.][0-9]+)?(ns|us|µs|ms|s|m|h))+$`
	// +optional
	ImagefsAvailable *string `json:"imagefsAvailable,omitempty" protobuf:"bytes,3,opt,name=imagefsAvailable"`
	// ImagefsInodesFree is the grace period for the imagefs.inodesFree soft eviction threshold.
	// The value must be a duration string. Eg. "30s", "1m30s".
	// The value must be positive and less than '5m'.
	// +kubebuilder:validation:Pattern=`^([0-9]+([.][0-9]+)?(ns|us|µs|ms|s|m|h))+$`
	// +optional
	ImagefsInodesFree *string `json:"imagefsInodesFree,omitempty" protobuf:"bytes,4,opt,name=imagefsInodesFree"`
	// NodefsInodesFree is the grace period for the nodefs.inodesFree soft eviction threshold.
	// The value must be a duration string. Eg. "30s", "1m30s".
	// The value must be positive and less than '5m'.
	// +kubebuilder:validation:Pattern=`^([0-9]+([.][0-9]+)?(ns|us|µs|ms|s|m|h))+$`
	// +optional
	NodefsInodesFree *string `json:"nodefsInodesFree,omitempty" protobuf:"bytes,5,opt,name=nodefsInodesFree"`
	// PidAvailable is the grace period for the pid.available soft eviction threshold.
	// The value must be a duration string. Eg. "30s", "1m30s".
	// The value must be positive and less than '5m'.
	// +kubebuilder:validation:Pattern=`^([0-9]+([.][0-9]+)?(ns|us|µs|ms|s|m|h))+$`
	// +optional
	PidAvailable *string `json:"pidAvailable,omitempty" protobuf:"bytes,6,opt,name=pidAvailable"`
}

// EvictionMinimumReclaim is a map of signal names to quantities that defines minimum reclaims.
// It describes the minimum amount of a given resource the kubelet will reclaim when performing a pod eviction.
// By default, all values are 0 if unspecified.
// +kubebuilder:validation:Optional
type EvictionMinimumReclaim struct {
	// MemoryAvailable is the minimum reclaim for memory.available.
	// The value must be a percentage, e.g., "5%".
	// The value must be positive and less than 10%.
	// +kubebuilder:validation:Pattern=`^[0-9]+(\.[0-9]+)?%$`
	// +optional
	MemoryAvailable *string `json:"memoryAvailable,omitempty" protobuf:"bytes,1,opt,name=memoryAvailable"`
	// NodefsAvailable is the minimum reclaim for nodefs.available.
	// The value must be a percentage, e.g., "5%".
	// The value must be positive and less than 10%.
	// +kubebuilder:validation:Pattern=`^[0-9]+(\.[0-9]+)?%$`
	// +optional
	NodefsAvailable *string `json:"nodefsAvailable,omitempty" protobuf:"bytes,2,opt,name=nodefsAvailable"`
	// ImagefsAvailable is the minimum reclaim for imagefs.available.
	// The value must be a percentage, e.g., "5%".
	// The value must be positive and less than 10%.
	// +kubebuilder:validation:Pattern=`^[0-9]+(\.[0-9]+)?%$`
	// +optional
	ImagefsAvailable *string `json:"imagefsAvailable,omitempty" protobuf:"bytes,3,opt,name=imagefsAvailable"`
	// ImagefsInodesFree is the minimum reclaim for imagefs.inodesFree.
	// The value must be a percentage, e.g., "5%".
	// The value must be positive and less than 10%.
	// +kubebuilder:validation:Pattern=`^[0-9]+(\.[0-9]+)?%$`
	// +optional
	ImagefsInodesFree *string `json:"imagefsInodesFree,omitempty" protobuf:"bytes,4,opt,name=imagefsInodesFree"`
	// NodefsInodesFree is the minimum reclaim for nodefs.inodesFree.
	// The value must be a percentage, e.g., "5%".
	// The value must be positive and less than 10%.
	// +kubebuilder:validation:Pattern=`^[0-9]+(\.[0-9]+)?%$`
	// +optional
	NodefsInodesFree *string `json:"nodefsInodesFree,omitempty" protobuf:"bytes,5,opt,name=nodefsInodesFree"`
	// PidAvailable is the minimum reclaim for pid.available.
	// The value must be a percentage, e.g., "5%".
	// The value must be positive and less than 10%.
	// +kubebuilder:validation:Pattern=`^[0-9]+(\.[0-9]+)?%$`
	// +optional
	PidAvailable *string `json:"pidAvailable,omitempty" protobuf:"bytes,6,opt,name=pidAvailable"`
}

// KubeletConfig defines kubelet config for a node.
//
// +kubebuilder:validation:XValidation:rule="has(self.imageGcHighThresholdPercent)&&has(self.imageGcLowThresholdPercent) ? self.imageGcHighThresholdPercent>self.imageGcLowThresholdPercent : true", message="ImageGcLowThresholdPercent must be lower than imageGcHighThresholdPercent"
// +kubebuilder:validation:XValidation:rule="has(self.imageGcHighThresholdPercent)&&!has(self.imageGcLowThresholdPercent) ? self.imageGcHighThresholdPercent>80 : true", message="ImageGcHighThresholdPercent must be higher than 80 which is default value of imageGcLowThresholdPercent"
type KubeletConfig struct {
	// This setting enforces the Pod's CPU limit. Setting this value to false means that the CPU limits for Pods are ignored.
	// Ignoring CPU limits might be desirable in certain scenarios where Pods are sensitive to CPU limits.
	// The risk of disabling cpuCFSQuota is that a rogue Pod can consume more CPU resources than intended.
	//
	// +kubebuilder:validation:Optional
	CpuCfsQuota *bool `json:"cpuCfsQuota,omitempty" protobuf:"bytes,1,opt,name=cpuCfsQuota"`
	// This setting sets the CPU CFS quota period value, cpu.cfs_period_us, which specifies the period of how often a cgroup's access to CPU resources should be reallocated.
	// This option lets you tune the CPU throttling behavior. Value must be 1ms <= period <= 1s.
	//
	// +kubebuilder:validation:Pattern="^([1-9][0-9]*)m?s$"
	// +kubebuilder:validation:Optional
	CpuCfsQuotaPeriod *string `json:"cpuCfsQuotaPeriod,omitempty" protobuf:"bytes,2,opt,name=cpuCfsQuotaPeriod"`
	// This setting controls the kubelet's CPU Manager Policy. The default value is none which is the default CPU affinity scheme, providing no affinity beyond what the OS scheduler does automatically.
	// Setting this value to static allows Pods in the Guaranteed QoS class with integer CPU requests to be assigned exclusive use of CPUs.
	//
	// +kubebuilder:validation:Enum=none;static
	// +kubebuilder:validation:Optional
	CpuManagerPolicy *string `json:"cpuManagerPolicy,omitempty" protobuf:"bytes,3,opt,name=cpuManagerPolicy"`
	// This setting sets the maximum number of process IDs (PIDs) that each Pod can use.
	//
	// +kubebuilder:validation:Minimum=1024
	// +kubebuilder:validation:Maximum=4194304
	// +kubebuilder:validation:Optional
	PodPidsLimit *int64 `json:"podPidsLimit,omitempty" protobuf:"bytes,4,opt,name=podPidsLimit"`
	// This setting sets the percent of disk usage before which image garbage collection is never
	// run. Lowest disk usage to garbage collect to. The percent is calculated as
	// this field value out of 100. Default is 80 if unspecified.
	//
	// +kubebuilder:validation:Minimum=10
	// +kubebuilder:validation:Maximum=84
	// +kubebuilder:validation:Optional
	ImageGcLowThresholdPercent *int64 `json:"imageGcLowThresholdPercent,omitempty" protobuf:"bytes,5,opt,name=imageGcLowThresholdPercent"`
	// This setting sets the percent of disk usage after which image garbage collection is always
	// run. The percent is calculated as this field value out of 100. Default is 85 if unspecified.
	//
	// +kubebuilder:validation:Minimum=11
	// +kubebuilder:validation:Maximum=85
	// +kubebuilder:validation:Optional
	ImageGcHighThresholdPercent *int64 `json:"imageGcHighThresholdPercent,omitempty" protobuf:"bytes,6,opt,name=imageGcHighThresholdPercent"`
	// This setting sets the minimum age for an unused image before it is garbage collected.
	// The string must be a decimal number with a unit suffix, such as "300s", "1.5h", and "2h45m".
	// Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
	// The value must be a positive duration and less than or equal to 2 minutes.
	// Default is "2m" if unspecified.
	//
	// +kubebuilder:validation:Pattern=^([0-9]+([.][0-9]+)?(ns|us|µs|ms|s|m|h))+$
	// +kubebuilder:validation:Optional
	ImageMinimumGcAge *string `json:"imageMinimumGcAge,omitempty" protobuf:"bytes,7,opt,name=imageMinimumGcAge"`
	// This setting sets the maximum age an image can be unused before it is garbage collected.
	// The string must be a decimal number with a unit suffix, such as "300s", "1.5h", and "2h45m".
	// Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
	// The value must be a positive duration.
	// Default is "0s" if unspecified, which disables the field.
	//
	// +kubebuilder:validation:Pattern=^([0-9]+([.][0-9]+)?(ns|us|µs|ms|s|m|h))+$
	// +kubebuilder:validation:Optional
	ImageMaximumGcAge *string `json:"imageMaximumGcAge,omitempty" protobuf:"bytes,8,opt,name=imageMaximumGcAge"`
	// This setting sets the maximum size of the container log file before it is rotated.
	// Format: positive number + unit, Eg. 100Ki, 10Mi, 5Gi. Valid units are Ki,
	// Mi, Gi. The value must be between 10Mi and 500Mi. And the total
	// container log size (container_log_max_size * container_log_max_files)
	// cannot exceed 1% of the total storage of the node.
	// Default is 10Mi in OSS if unspecified.
	//
	// +kubebuilder:validation:Pattern="^([0-9]+([.][0-9]+)?(Ki|Mi|Gi))+$"
	// +kubebuilder:validation:Optional
	ContainerLogMaxSize *string `json:"containerLogMaxSize,omitempty" protobuf:"bytes,9,opt,name=containerLogMaxSize"`
	// This setting sets the maximum number of container log files that can be present for a
	// container. Default is 5 in OSS if unspecified.
	//
	// +kubebuilder:validation:Minimum=2
	// +kubebuilder:validation:Maximum=10
	// +kubebuilder:validation:Optional
	ContainerLogMaxFiles *int64 `json:"containerLogMaxFiles,omitempty" protobuf:"bytes,10,opt,name=containerLogMaxFiles"`
	// This setting defines a comma-separated allowlist of unsafe sysctls or sysctl patterns
	// (ending in `*`). The unsafe namespaced sysctl groups are `kernel.shm*`, `kernel.msg*`,
	// `kernel.sem`, `fs.mqueue.*`, and `net.*`. Leaving this allowlist empty means they cannot be set on Pods.
	//
	// +kubebuilder:listType=atomic
	// +kubebuilder:validation:MaxItems=100
	// +kubebuilder:validation:items:MinLength=1
	// +kubebuilder:validation:items:MaxLength=253
	// +kubebuilder:validation:items:Pattern="^([a-z0-9]([-_a-z0-9]*[a-z0-9])?[./])*([a-z0-9][-_a-z0-9]*)?[a-z0-9*]$"
	AllowedUnsafeSysctls []string `json:"allowedUnsafeSysctls,omitempty" protobuf:"bytes,11,opt,name=allowedUnsafeSysctls"`
	// This setting sets the maximum number of image pulls in parallel. Default is 2 or 3 depending on boot disk type.
	//
	// +kubebuilder:validation:Minimum=2
	// +kubebuilder:validation:Maximum=5
	// +kubebuilder:validation:Optional
	MaxParallelImagePulls *int64 `json:"maxParallelImagePulls,omitempty" protobuf:"bytes,12,opt,name=maxParallelImagePulls"`
	// This setting sets whether to enable single process OOM killer.
	// If set to true, the processes in a container will be OOM killed individually instead of as a group.
	//
	// +kubebuilder:validation:Optional
	SingleProcessOOMKill *bool `json:"singleProcessOOMKill,omitempty" protobuf:"bytes,13,opt,name=singleProcessOOMKill"`
	// EvictionSoft defines soft eviction thresholds.
	//
	// +kubebuilder:validation:Optional
	EvictionSoft *EvictionSoft `json:"evictionSoft,omitempty" protobuf:"bytes,14,opt,name=evictionSoft"`
	// EvictionSoftGracePeriod defines grace periods for soft eviction thresholds.
	//
	// +kubebuilder:validation:Optional
	EvictionSoftGracePeriod *EvictionSoftGracePeriod `json:"evictionSoftGracePeriod,omitempty" protobuf:"bytes,15,opt,name=evictionSoftGracePeriod"`
	// EvictionMinimumReclaim defines minimum reclaims.
	//
	// +kubebuilder:validation:Optional
	EvictionMinimumReclaim *EvictionMinimumReclaim `json:"evictionMinimumReclaim,omitempty" protobuf:"bytes,16,opt,name=evictionMinimumReclaim"`
	// EvictionMaxPodGracePeriodSeconds is the maximum allowed grace period
	// (in seconds) to use when terminating pods in response to a soft eviction
	// threshold being met.
	//
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=300
	// +kubebuilder:validation:Optional
	EvictionMaxPodGracePeriodSeconds *int64 `json:"evictionMaxPodGracePeriodSeconds,omitempty" protobuf:"bytes,17,opt,name=evictionMaxPodGracePeriodSeconds"`
	// TopologyManager contains the configuration for the Kubelet Topology Manager.
	//
	// +kubebuilder:validation:Optional
	TopologyManager *TopologyManager `json:"topologyManager,omitempty" protobuf:"bytes,18,opt,name=topologyManager"`
	// MemoryManager contains the configuration for the Kubelet Memory Manager.
	//
	// +kubebuilder:validation:Optional
	MemoryManager *MemoryManager `json:"memoryManager,omitempty" protobuf:"bytes,19,opt,name=memoryManager"`
}

// SysctlsConfig defines sysctls config for a node.
type SysctlsConfig struct {
	// Maximum number of packets, queued on the INPUT side, when the interface receives packets faster than kernel can process them.
	//
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=2147483647
	// +kubebuilder:validation:Optional
	Net_core_netdev_max_backlog *int64 `json:"net.core.netdev_max_backlog,omitempty" protobuf:"bytes,1,opt,name=net.core.netdev_max_backlog"`
	// The maximum receive socket buffer size in bytes.
	//
	// +kubebuilder:validation:Minimum=2304
	// +kubebuilder:validation:Maximum=2147483647
	// +kubebuilder:validation:Optional
	Net_core_rmem_max *int64 `json:"net.core.rmem_max,omitempty" protobuf:"bytes,2,opt,name=net.core.rmem_max"`
	// The default setting (in bytes) of the socket send buffer.
	//
	// +kubebuilder:validation:Minimum=4608
	// +kubebuilder:validation:Maximum=2147483647
	// +kubebuilder:validation:Optional
	Net_core_wmem_default *int64 `json:"net.core.wmem_default,omitempty" protobuf:"bytes,3,opt,name=net.core.wmem_default"`
	// The maximum send socket buffer size in bytes.
	//
	// +kubebuilder:validation:Minimum=4608
	// +kubebuilder:validation:Maximum=2147483647
	// +kubebuilder:validation:Optional
	Net_core_wmem_max *int64 `json:"net.core.wmem_max,omitempty" protobuf:"bytes,4,opt,name=net.core.wmem_max"`
	// Maximum ancillary buffer size allowed per socket. Ancillary data is a sequence of struct cmsghdr structures with appended data.
	//
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=2147483647
	// +kubebuilder:validation:Optional
	Net_core_optmem_max *int64 `json:"net.core.optmem_max,omitempty" protobuf:"bytes,5,opt,name=net.core.optmem_max"`
	// Limit of socket listen() backlog, known in userspace as SOMAXCONN. Defaults to 128. See also tcp_max_syn_backlog for additional tuning for TCP sockets.
	//
	// +kubebuilder:validation:Minimum=128
	// +kubebuilder:validation:Maximum=2147483647
	// +kubebuilder:validation:Optional
	Net_core_somaxconn *int64 `json:"net.core.somaxconn,omitempty" protobuf:"bytes,6,opt,name=net.core.somaxconn"`
	// Minimal size of receive buffer used by UDP sockets in moderation. Each UDP socket is able to use the size for receiving data, even if total pages of UDP sockets exceed udp_mem pressure. The unit is byte. Default: 1 page. The three values are: min, default, max. Eg. '4096 87380 6291456'.
	//
	// +kubebuilder:validation:Optional
	Net_ipv4_tcp_rmem *string `json:"net.ipv4.tcp_rmem,omitempty" protobuf:"bytes,7,opt,name=net.ipv4.tcp_rmem"`
	// Minimal size of send buffer used by UDP sockets in moderation. Each UDP socket is able to use the size for sending data, even if total pages of UDP sockets exceed udp_mem pressure. The unit is byte. Default: 1 page. The three values are: min, default, max. Eg. '4096 87380 6291456'.
	//
	// +kubebuilder:validation:Optional
	Net_ipv4_tcp_wmem *string `json:"net.ipv4.tcp_wmem,omitempty" protobuf:"bytes,8,opt,name=net.ipv4.tcp_wmem"`
	// Allow to reuse TIME-WAIT sockets for new connections when it is safe from protocol viewpoint. It should not be changed without advice/request of technical experts.
	//
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=2
	// +kubebuilder:validation:Optional
	Net_ipv4_tcp_tw_reuse *int64 `json:"net.ipv4.tcp_tw_reuse,omitempty" protobuf:"bytes,9,opt,name=net.ipv4.tcp_tw_reuse"`
	// Low latency busy poll timeout for poll and select. (needs CONFIG_NET_RX_BUSY_POLL) Approximate time in us to busy loop waiting for events.
	//
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=2147483647
	// +kubebuilder:validation:Optional
	Net_core_busy_poll *int64 `json:"net.core.busy_poll,omitempty" protobuf:"bytes,10,opt,name=net.core.busy_poll"`
	// Low latency busy poll timeout for socket reads. (needs CONFIG_NET_RX_BUSY_POLL) Approximate time in us to busy loop waiting for packets on the device queue.
	//
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=2147483647
	// +kubebuilder:validation:Optional
	Net_core_busy_read *int64 `json:"net.core.busy_read,omitempty" protobuf:"bytes,11,opt,name=net.core.busy_read"`
	// Changing this value is same as changing conf/default/disable_ipv6 setting and also all per-interface disable_ipv6 settings to the same value.
	//
	// +kubebuilder:validation:Optional
	Net_ipv6_conf_all_disable_ipv6 *bool `json:"net.ipv6.conf.all.disable_ipv6,omitempty" protobuf:"bytes,12,opt,name=net.ipv6.conf.all.disable_ipv6"`
	// Disable IPv6 operation.
	//
	// +kubebuilder:validation:Optional
	Net_ipv6_conf_default_disable_ipv6 *bool `json:"net.ipv6.conf.default.disable_ipv6,omitempty" protobuf:"bytes,13,opt,name=net.ipv6.conf.default.disable_ipv6"`
	// Maximum number of memory map areas a process may have.
	//
	// +kubebuilder:validation:Minimum=65536
	// +kubebuilder:validation:Maximum=2147483647
	// +kubebuilder:validation:Optional
	Vm_max_map_count *int64 `json:"vm.max_map_count,omitempty" protobuf:"bytes,14,opt,name=vm.max_map_count"`
	// The system-wide maximum number of shared memory segments.
	//
	// +kubebuilder:validation:Minimum=4096
	// +kubebuilder:validation:Maximum=32768
	// +kubebuilder:validation:Optional
	Kernel_shmmni *int64 `json:"kernel.shmmni,omitempty" protobuf:"bytes,15,opt,name=kernel.shmmni"`
	// The maximum size (in bytes) of a single shared memory segment allowed by the kernel.
	// Note that the actual range should be integer between 0 and 18446744073692774399, while kubebuilder would lose some precision on uint64 during the internal representation and parsing.
	//
	// +kubebuilder:validation:Pattern="^([0-9]+)$"
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=20
	// +kubebuilder:validation:Optional
	Kernel_shmall *string `json:"kernel.shmall,omitempty" protobuf:"bytes,16,opt,name=kernel.shmall"`
	// The total amount of shared memory pages that can be used on the system at one time.
	// Note that the actual range should be integer between 0 and 18446744073692774399, while kubebuilder would lose some precision on uint64 during the internal representation and parsing.
	//
	// +kubebuilder:validation:Pattern="^([0-9]+)$"
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=20
	// +kubebuilder:validation:Optional
	Kernel_shmmax *string `json:"kernel.shmmax,omitempty" protobuf:"bytes,17,opt,name=kernel.shmmax"`
	// The default receive socket buffer size in bytes.
	//
	// +kubebuilder:validation:Minimum=2304
	// +kubebuilder:validation:Maximum=2147483647
	// +kubebuilder:validation:Optional
	Net_core_rmem_default *int64 `json:"net.core.rmem_default,omitempty" protobuf:"bytes,18,opt,name=net.core.rmem_default"`
	// The size of connection tracking table.
	//
	// +kubebuilder:validation:Minimum=65536
	// +kubebuilder:validation:Maximum=4194304
	// +kubebuilder:validation:Optional
	Net_netfilter_nf_conntrack_max *int64 `json:"net.netfilter.nf_conntrack_max,omitempty" protobuf:"bytes,19,opt,name=net.netfilter.nf_conntrack_max"`
	// The size of hash table for connection tracking.
	// +kubebuilder:validation:Minimum=65536
	// +kubebuilder:validation:Maximum=524288
	// +kubebuilder:validation:Optional
	Net_netfilter_nf_conntrack_buckets *int64 `json:"net.netfilter.nf_conntrack_buckets,omitempty" protobuf:"bytes,20,opt,name=net.netfilter.nf_conntrack_buckets"`
	// Whether to enable connection tracking flow accounting.
	//
	// +kubebuilder:validation:Optional
	Net_netfilter_nf_conntrack_acct *bool `json:"net.netfilter.nf_conntrack_acct,omitempty" protobuf:"bytes,21,opt,name=net.netfilter.nf_conntrack_acct"`
	// The duration of dead connections before deleted automatically from connection tracking table.
	//
	// +kubebuilder:validation:Minimum=600
	// +kubebuilder:validation:Maximum=86400
	// +kubebuilder:validation:Optional
	Net_netfilter_nf_conntrack_tcp_timeout_established *int64 `json:"net.netfilter.nf_conntrack_tcp_timeout_established,omitempty" protobuf:"bytes,22,opt,name=net.netfilter.nf_conntrack_tcp_timeout_established"`
	// The period for which the TCP connections can remain in the CLOSE_WAIT state, and stay in the table.
	//
	// +kubebuilder:validation:Minimum=60
	// +kubebuilder:validation:Maximum=3600
	// +kubebuilder:validation:Optional
	Net_netfilter_nf_conntrack_tcp_timeout_close_wait *int64 `json:"net.netfilter.nf_conntrack_tcp_timeout_close_wait,omitempty" protobuf:"bytes,23,opt,name=net.netfilter.nf_conntrack_tcp_timeout_close_wait"`
	// The period for which the TCP connections can remain in the TIME_WAIT state, and stay in the table.
	//
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=600
	// +kubebuilder:validation:Optional
	Net_netfilter_nf_conntrack_tcp_timeout_time_wait *int64 `json:"net.netfilter.nf_conntrack_tcp_timeout_time_wait,omitempty" protobuf:"bytes,24,opt,name=net.netfilter.nf_conntrack_tcp_timeout_time_wait"`
	// The maximum number of file descriptors that can be opened by a process.
	//
	// +kubebuilder:validation:Minimum=1048576
	// +kubebuilder:validation:Maximum=2147483584
	// +kubebuilder:validation:Optional
	Fs_nr_open *int64 `json:"fs.nr_open,omitempty" protobuf:"bytes,25,opt,name=fs.nr_open"`
	// The maximum number of inotify watches that a user can create.
	//
	// +kubebuilder:validation:Minimum=8192
	// +kubebuilder:validation:Maximum=1048576
	// +kubebuilder:validation:Optional
	Fs_inotify_max_user_watches *int64 `json:"fs.inotify.max_user_watches,omitempty" protobuf:"bytes,26,opt,name=fs.inotify.max_user_watches"`
	// The maximum number of inotify instances that a user can create.
	//
	// +kubebuilder:validation:Minimum=8192
	// +kubebuilder:validation:Maximum=1048576
	// +kubebuilder:validation:Optional
	Fs_inotify_max_user_instances *int64 `json:"fs.inotify.max_user_instances,omitempty" protobuf:"bytes,27,opt,name=fs.inotify.max_user_instances"`
	// Determines the kernel's memory overcommit handling strategy.
	// Supported values:
	// 0: Rejects allocations that are obviously too large.
	// 1: Allows overcommit until memory is exhausted.
	// 2 (strict): Prevents overcommit beyond swap space plus a percentage of RAM defined by 'vm.overcommit_ratio'.
	//
	// +kubebuilder:validation:Enum=0;1;2
	// +kubebuilder:validation:Optional
	Vm_overcommit_memory *int64 `json:"vm.overcommit_memory,omitempty" protobuf:"bytes,28,opt,name=vm.overcommit_memory"`
	// Specifies the percentage of physical RAM allowed for overcommit when 'vm.overcommit_memory' is set to 2.
	// The total committed address space cannot exceed swap plus this RAM percentage.
	//
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=100
	// +kubebuilder:validation:Optional
	Vm_overcommit_ratio *int64 `json:"vm.overcommit_ratio,omitempty" protobuf:"bytes,29,opt,name=vm.overcommit_ratio"`
	// Adjusts the kernel's preference for reclaiming memory used for dentry (directory) and inode caches.
	//
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=100
	// +kubebuilder:validation:Optional
	Vm_vfs_cache_pressure *int64 `json:"vm.vfs_cache_pressure,omitempty" protobuf:"bytes,30,opt,name=vm.vfs_cache_pressure"`
	// Percentage of system memory that can be filled with dirty pages (modified but not yet written to disk) before background kernel flusher threads begin writeback.
	// This value should be less than 'vm.dirty_ratio'.
	//
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=100
	// +kubebuilder:validation:Optional
	Vm_dirty_background_ratio *int64 `json:"vm.dirty_background_ratio,omitempty" protobuf:"bytes,31,opt,name=vm.dirty_background_ratio"`
	// Percentage of system memory that can be filled with dirty pages before processes performing writes are forced to block and write out dirty data synchronously.
	// This value should be greater than 'vm.dirty_background_ratio'.
	//
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=100
	// +kubebuilder:validation:Optional
	Vm_dirty_ratio *int64 `json:"vm.dirty_ratio,omitempty" protobuf:"bytes,32,opt,name=vm.dirty_ratio"`
	// Maximum age (in hundredths of a second) that dirty data can remain in memory before kernel flusher threads write it to disk.
	// Lower values result in faster, more frequent writebacks.
	//
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=6000
	// +kubebuilder:validation:Optional
	Vm_dirty_expire_centisecs *int64 `json:"vm.dirty_expire_centisecs,omitempty" protobuf:"bytes,33,opt,name=vm.dirty_expire_centisecs"`
	// Interval (in hundredths of a second) at which kernel flusher threads wake up to write 'old' dirty data to disk.
	//
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=1000
	// +kubebuilder:validation:Optional
	Vm_dirty_writeback_centisecs *int64 `json:"vm.dirty_writeback_centisecs,omitempty" protobuf:"bytes,34,opt,name=vm.dirty_writeback_centisecs"`
	// Maximum number of file-handles that the Linux kernel will allocate.
	//
	// +kubebuilder:validation:Minimum=104857
	// +kubebuilder:validation:Maximum=67108864
	// +kubebuilder:validation:Optional
	Fs_file_max *int64 `json:"fs.file-max,omitempty" protobuf:"bytes,35,opt,name=fs.file-max"`
	// The maximum system-wide number of asynchronous io requests.
	//
	// +kubebuilder:validation:Minimum=65536
	// +kubebuilder:validation:Maximum=4194304
	// +kubebuilder:validation:Optional
	Fs_aio_max_nr *int64 `json:"fs.aio-max-nr,omitempty" protobuf:"bytes,36,opt,name=fs.aio-max-nr"`
	// Maximal number of TCP sockets not attached to any user file handle.

	// +kubebuilder:validation:Minimum=16384
	// +kubebuilder:validation:Maximum=262144
	// +kubebuilder:validation:Optional
	Net_ipv4_tcp_max_orphans *int64 `json:"net.ipv4.tcp_max_orphans,omitempty" protobuf:"bytes,37,opt,name=net.ipv4.tcp_max_orphans"`
	// Controls the tendency of the kernel to move processes out of physical memory and onto the swap disk.

	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=200
	// +kubebuilder:validation:Optional
	Vm_swappiness *int64 `json:"vm.swappiness,omitempty" protobuf:"bytes,38,opt,name=vm.swappiness"`
	// Controls the aggressiveness of kswapd. The flag defines the amount of memory left in a node before kswapd is woken up and how much memory needs to be freed before kswapd goes back to sleep.

	// +kubebuilder:validation:Minimum=10
	// +kubebuilder:validation:Maximum=3000
	// +kubebuilder:validation:Optional
	Vm_watermark_scale_factor *int64 `json:"vm.watermark_scale_factor,omitempty" protobuf:"bytes,39,opt,name=vm.watermark_scale_factor"`
	// Minimum free memory before OOM.

	// +kubebuilder:validation:Minimum=67584
	// +kubebuilder:validation:Maximum=1048576
	// +kubebuilder:validation:Optional
	Vm_min_free_kbytes *int64 `json:"vm.min_free_kbytes,omitempty" protobuf:"bytes,40,opt,name=vm.min_free_kbytes"`
}

// HugepagesConfig defines hugepages config for a node.
type HugepagesConfig struct {
	// Number of 1-gigabyte-sized huge pages to allocate.
	//
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Optional
	HugepageSize1g *int64 `json:"hugepage_size1g,omitempty" protobuf:"bytes,1,opt,name=hugepage_size1g"`
	// Number of 2-megabyte-sized huge pages to allocate.
	//
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Optional
	HugepageSize2m *int64 `json:"hugepage_size2m,omitempty" protobuf:"bytes,2,opt,name=hugepage_size2m"`
}

// Location describes CCC zonal preferences config.
type Location struct {
	// Zones lists zones considered for node autoprovisioning.
	//
	// +kubebuilder:listType=atomic
	// +kubebuilder:validation:MinItems=1
	// +optional
	Zones []string `json:"zones,omitempty" protobuf:"bytes,1,opt,name=zones"`

	// LocationPolicy specifies the strategy for selecting zones when scaling up a node
	// pool managed by this Compute Class. This setting controls the distribution of new
	// nodes across zones in the node pool's region and corresponds to the node pool
	// setting of the same name.
	// More info: https://cloud.google.com/sdk/gcloud/reference/container/node-pools/create#--location-policy
	// +optional
	// +kubebuilder:validation:Enum=ANY;BALANCED
	LocationPolicy *string `json:"locationPolicy,omitempty" protobuf:"bytes,2,opt,name=locationPolicy"`
}

// PriorityDefaults define the default rules for all priorities if the rule doesn't exist in some priority.
type PriorityDefaults struct {
	// NodeSystemConfig defines node system config for a node.
	//
	// +kubebuilder:validation:Optional
	NodeSystemConfig *NodeSystemConfig `json:"nodeSystemConfig,omitempty" protobuf:"bytes,1,opt,name=nodeSystemConfig"`
	// Location describes CCC zonal preferences config.
	//
	// +optional
	Location *Location `json:"location,omitempty" protobuf:"bytes,2,opt,name=location"`

	// HostMaintenancePolicy contains the desired maintenance policy for the
	// Google Compute Engine hosts.
	//
	// +kubebuilder:validation:Optional
	HostMaintenancePolicy *HostMaintenancePolicy `json:"hostMaintenancePolicy,omitempty" protobuf:"bytes,3,opt,name=hostMaintenancePolicy"`
}

// TaintConfig applies the given kubernetes taints on all nodes in the new node pool, which can be used with tolerations for pod scheduling.
// Any workload that does not tolerate the taints specified in this object will not be scheduled to the node pool.
// More info: https://cloud.google.com/sdk/gcloud/reference/container/node-pools/create#--node-taints
type TaintConfig struct {

	// Node taint key. The key must begin with a letter or number, and may contain letters, numbers, hyphens, dots, and underscores up to 253 characters.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^([a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*/)?(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])$`
	// +kubebuilder:validation:XValidation:rule="!self.contains('cloud.google.com')",message="Taint keys cannot contain reserved `cloud.google.com` substring."
	// +kubebuilder:validation:XValidation:rule="!self.contains('kubernetes.io')",message="Taint keys cannot contain reserved `kubernetes.io` substring."
	// +kubebuilder:validation:XValidation:rule="!self.contains('gke.io')",message="Taint keys cannot contain reserved `gke.io` substring."
	// +kubebuilder:validation:XValidation:rule="!self.contains('k8s.io')",message="Taint keys cannot contain reserved `k8s.io` substring."
	// +kubebuilder:validation:MaxLength=253
	Key string `json:"key,omitempty" protobuf:"bytes,1,opt,name=key"`

	// The value that matches the specified taint key.
	// +kubebuilder:validation:Pattern=`^([a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*/)?(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])$`
	// +kubebuilder:validation:Optional
	Value string `json:"value,omitempty" protobuf:"bytes,2,opt,name=value"`

	// It defines the taint's effect on pods that does not have the necessary toleration.
	// The following values are supported: NoSchedule, PreferNoSchedule, and NoExecute.
	// +kubebuilder:validation:Enum=NoSchedule;PreferNoSchedule;NoExecute
	// +kubebuilder:validation:Required
	Effect string `json:"effect,omitempty" protobuf:"bytes,3,opt,name=effect"`
}

// Tags define the key/value of resource manager tags.
// Tags must be in one of the following formats ([KEY]=[VALUE])
// 1. tagKeys/{tag_key_id}=tagValues/{tag_value_id}
// 2. {org_id}/{tag_key_name}={tag_value_name}
// 3. {project_id}/{tag_key_name}={tag_value_name}
type Tags struct {
	// +kubebuilder:validation:Required
	Key string `json:"key,omitempty" protobuf:"bytes,1,opt,name=key"`

	// +kubebuilder:validation:Required
	Value string `json:"value,omitempty" protobuf:"bytes,2,opt,name=value"`
}

// ComputeClassStatus is the current status of the ComputeClass.
type ComputeClassStatus struct {
	// Conditions represent the observations of a ComputeClass's current state.
	//
	// +optional
	Conditions []metav1.Condition `json:"conditions" protobuf:"bytes,1,rep,name=conditions"`
}

// HostMaintenancePolicy contains the maintenance policy for the hosts on which
// the GKE VMs run on.
type HostMaintenancePolicy struct {
	// Specifies the frequency of planned maintenance events.
	// Currently supported values:
	// * AS_NEEDED
	// * PERIODIC
	//
	// +optional
	// +kubebuilder:validation:Enum=AS_NEEDED;PERIODIC
	MaintenanceInterval string `json:"maintenanceInterval" protobuf:"bytes,1,opt,name=maintenanceInterval"`
}

// GpuSharing represents the GPU sharing configuration for
// Hardware Accelerators.
type GpuSharing struct {
	// SharingStrategy The type of GPU sharing strategy to enable on the GPU node.
	// Possible values:
	// * TIME_SHARING - GPUs are time-shared between containers.
	// * MPS - GPUs are shared between containers with NVIDIA MPS.
	// +kubebuilder:validation:Enum=MPS;TIME_SHARING
	// +optional
	SharingStrategy string `json:"sharingStrategy,omitempty" protobuf:"bytes,1,opt,name=sharingStrategy"`

	// MaxSharedClientsPerGPU describes the max number of containers that can
	// share a physical GPU.
	// +kubebuilder:validation:Minimum=0
	// +optional
	MaxSharedClientsPerGPU int64 `json:"maxSharedClientsPerGPU,omitempty" protobuf:"bytes,2,name=maxSharedClientsPerGPU"`

	// GpuPartitionSize is size of partitions to create on the GPU. Valid values are
	// described in the NVIDIA mig user guide. Example: "1g.5gb"
	// (https://docs.nvidia.com/datacenter/tesla/mig-user-guide/#partitioning).
	// +optional
	GpuPartitionSize string `json:"gpuPartitionSize,omitempty" protobuf:"bytes,3,name=gpuPartitionSize"`
}
