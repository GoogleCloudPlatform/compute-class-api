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

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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

// SubnetPriority defines a priority for a subnet to be used for the primary network interface.
type SubnetPriority struct {
	// Name defines the name of the subnetwork.
	//
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=63
	// +kubebuilder:validation:Pattern=`^[a-z]([-a-z0-9]*[a-z0-9])?$`
	Name string `json:"name" protobuf:"bytes,1,name=name"`

	// PodRange defines the name of the secondary subnet range reserved for pod IPs.
	//
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=63
	// +kubebuilder:validation:Pattern=`^[a-z]([-a-z0-9]*[a-z0-9])?$`
	PodRange string `json:"podRange" protobuf:"bytes,2,name=podRange"`
}

// NetworkConfig defines network-related settings for the ComputeClass.
type NetworkConfig struct {
	// SubnetPriorities is an ordered list of subnets to fall back through.
	// TODO(b/552484145): Increase the max items to 5 once the API is updated to support multiple subnets.
	//
	// +kubebuilder:listType=atomic
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=1
	// +optional
	SubnetPriorities []SubnetPriority `json:"subnetPriorities,omitempty" protobuf:"bytes,1,rep,name=subnetPriorities"`
}

// MinimumCapacity defines the minimum capacity required for a given
// compute class or priority. It allows managing statically sized infrastructure.
type MinimumCapacity struct {
	// TargetNodeCount defines a minimum number of nodes that should be present in the cluster.
	// If the active node count falls below the defined threshold,
	// Cluster Autoscaler will proactively provision capacity to satisfy the requirement.
	// +optional
	TargetNodeCount *int `json:"targetNodeCount,omitempty" protobuf:"bytes,1,opt,name=targetNodeCount"`
}

// ComputeClassSpec is a specification of provisioning priorities and
// other autoscaling settings.
//
// +kubebuilder:validation:XValidation:rule="!has(oldSelf.autopilot) || has(self.autopilot)", message="Autopilot is required once set"
// +kubebuilder:validation:XValidation:rule="(has(self.autopilot) && self.autopilot.enabled) ? !self.priorities.exists(priority, has(priority.nodepools)) : true", message="Nodepools priority cannot be used when Autopilot is enabled"
// +kubebuilder:validation:XValidation:rule="(has(self.autopilot) && self.autopilot.enabled) ? !(has(self.nodePoolAutoCreation) && !self.nodePoolAutoCreation.enabled) : true", message="NodePoolAutoCreation cannot be disabled when Autopilot is enabled"
// +kubebuilder:validation:XValidation:rule="(has(self.autopilot) && self.autopilot.enabled) ? (!has(self.nodePoolConfig) || !has(self.nodePoolConfig.imageType) || self.nodePoolConfig.imageType == \"cos_containerd\") : true", message="Only cos_containerd image type can be used when Autopilot is enabled"
// +kubebuilder:validation:XValidation:rule="(has(self.autopilot) && self.autopilot.enabled) ? (!has(self.nodePoolConfig) || !has(self.nodePoolConfig.loggingConfig) || !has(self.nodePoolConfig.loggingConfig.loggingVariantConfig) || !has(self.nodePoolConfig.loggingConfig.loggingVariantConfig.variant) || self.nodePoolConfig.loggingConfig.loggingVariantConfig.variant == \"DEFAULT\") : true", message="Only DEFAULT logging variant can be used when Autopilot is enabled"
// +kubebuilder:validation:XValidation:rule="(has(self.nodePoolConfig) && has(self.nodePoolConfig.workloadType) && !has(self.nodePoolGroup)) ? self.nodePoolConfig.workloadType == \"HIGH_AVAILABILITY\" : true", message="If NodePoolGroup is not specified NodePoolConfig.WorkloadType can only be HIGH_AVAILABILITY if set"
// +kubebuilder:validation:XValidation:rule="self.priorities.exists(priority, has(priority.podFamily)) ? (has(self.autopilot) && self.autopilot.enabled) : true", message="In GKE Standard, pod family can be used only if Autopilot is enabled"
// +kubebuilder:validation:XValidation:rule="(has(self.nodePoolConfig) && has(self.nodePoolConfig.confidentialNodeType)) ? self.priorities.all(priority, has(priority.machineFamily) || has(priority.machineType)) : true", message="If using NodePoolConfig.ConfidentialNodeType, each priority must specify either MachineFamily or MachineType."
// +kubebuilder:validation:XValidation:rule="(has(self.nodePoolConfig) && has(self.nodePoolConfig.confidentialNodeType) && self.nodePoolConfig.confidentialNodeType == \"SEV\") ? self.priorities.all(priority, ((has(priority.machineFamily) && priority.machineFamily in ['n2d', 'c2d', 'c3d', 'c4d', 'g4']) || (has(priority.machineType) && priority.machineType.split('-')[0] in ['n2d', 'c2d', 'c3d', 'c4d', 'g4']))) : true", message="ConfidentialNodeType SEV only supports N2D, C2D, C3D, C4D, G4"
// +kubebuilder:validation:XValidation:rule="(has(self.nodePoolConfig) && has(self.nodePoolConfig.confidentialNodeType) && self.nodePoolConfig.confidentialNodeType == \"SEV\") ? self.priorities.all(priority, (has(priority.machineFamily) && priority.machineFamily == 'g4') ? false : (has(priority.machineType) && priority.machineType.startsWith('g4-')) ? (priority.machineType == 'g4-standard-48' && (!has(priority.gpu) || has(priority.gpu) && (!has(priority.gpu.type) || priority.gpu.type == 'nvidia-rtx-pro-6000'))) : true) : true", message="ConfidentialNodeType SEV on G4 only supports g4-standard-48 machine type and nvidia-rtx-pro-6000 GPU type"
// +kubebuilder:validation:XValidation:rule="(has(self.nodePoolConfig) && has(self.nodePoolConfig.confidentialNodeType) && self.nodePoolConfig.confidentialNodeType == \"SEV_SNP\") ? self.priorities.all(priority, ((has(priority.machineFamily) && priority.machineFamily in ['n2d']) || (has(priority.machineType) && priority.machineType.split('-')[0] in ['n2d']))) : true", message="ConfidentialNodeType SEV_SNP only supports N2D"
// TDX should only be enabled on C3, C4, A3 and A4 machine families. This is because TDX has specific hardware requirements that are only met by these machine families.
// +kubebuilder:validation:XValidation:rule="(has(self.nodePoolConfig) && has(self.nodePoolConfig.confidentialNodeType) && self.nodePoolConfig.confidentialNodeType == \"TDX\") ? self.priorities.all(priority, (has(priority.machineFamily) && priority.machineFamily in ['c3', 'c4', 'a3', 'a4']) || (has(priority.machineType) && (priority.machineType.startsWith('c3-standard-') || priority.machineType.startsWith('c4-standard-') || priority.machineType == 'a3-highgpu-1g' || priority.machineType == 'a4-highgpu-8g'))) : true", message="ConfidentialNodeType TDX only supports C3 standard, C4 standard, A3 and A4"
// +kubebuilder:validation:XValidation:rule="(has(self.nodePoolConfig) && has(self.nodePoolConfig.confidentialNodeType) && self.nodePoolConfig.confidentialNodeType == \"TDX\") ? self.priorities.all(priority, (has(priority.machineFamily) && priority.machineFamily == 'c3' || has(priority.machineType) && priority.machineType.startsWith('c3-standard-')) ? (!has(priority.gpu) || has(priority.gpu) && (!has(priority.gpu.type) || priority.gpu.type == 'nvidia-h100-80gb')) : true) : true", message="ConfidentialNodeType TDX on C3 only supports c3-standard- machine type and nvidia-h100-80gb GPU type"
// +kubebuilder:validation:XValidation:rule="(has(self.nodePoolConfig) && has(self.nodePoolConfig.confidentialNodeType) && self.nodePoolConfig.confidentialNodeType == \"TDX\") ? self.priorities.all(priority, (has(priority.machineFamily) && priority.machineFamily == 'c4' || has(priority.machineType) && priority.machineType.startsWith('c4-standard-')) ? (!has(priority.gpu) || has(priority.gpu) && (!has(priority.gpu.type) || priority.gpu.type == 'nvidia-h100-80gb')) : true) : true", message="ConfidentialNodeType TDX on C4 only supports c4-standard- machine type and nvidia-h100-80gb GPU type"
// +kubebuilder:validation:XValidation:rule="(has(self.nodePoolConfig) && has(self.nodePoolConfig.confidentialNodeType) && self.nodePoolConfig.confidentialNodeType == \"TDX\") ? self.priorities.all(priority, (has(priority.machineFamily) && priority.machineFamily == 'a3' || has(priority.machineType) && priority.machineType == 'a3-highgpu-1g') ? (!has(priority.gpu) || has(priority.gpu) && (!has(priority.gpu.type) || priority.gpu.type == 'nvidia-h100-80gb')) : true) : true", message="ConfidentialNodeType TDX on A3 only supports a3-highgpu-1g machine type and nvidia-h100-80gb GPU type"
// +kubebuilder:validation:XValidation:rule="(has(self.nodePoolConfig) && has(self.nodePoolConfig.confidentialNodeType) && self.nodePoolConfig.confidentialNodeType == \"TDX\") ? self.priorities.all(priority, (has(priority.machineFamily) && priority.machineFamily == 'a4' || has(priority.machineType) && priority.machineType == 'a4-highgpu-8g') ? (!has(priority.gpu) || has(priority.gpu) && (!has(priority.gpu.type) || priority.gpu.type == 'nvidia-b200')) : true) : true", message="ConfidentialNodeType TDX on A4 only supports a4-highgpu-8g machine type and nvidia-b200 GPU type"
// +kubebuilder:validation:XValidation:rule="(has(self.nodePoolConfig) && has(self.nodePoolConfig.confidentialNodeType) && self.priorities.exists(priority, has(priority.gpu))) ? (has(self.priorityDefaults) && has(self.priorityDefaults.location) && has(self.priorityDefaults.location.zones)) || self.priorities.all(priority, has(priority.location) && has(priority.location.zones)) : true", message="When using confidential GPUs you must specify location.zones"
// +kubebuilder:validation:XValidation:rule="self.priorities.all(p, has(p.priorityScore)) || self.priorities.all(p, !has(p.priorityScore))", message="PriorityScore must be set for all priorities or for none of them"
// +kubebuilder:validation:XValidation:rule="self.priorities.all(p, (has(p.gpu) && has(p.gpu.topology)) ? (((has(p.machineFamily) && p.machineFamily == 'a4x') || (has(p.gpu.type) && p.gpu.type == 'nvidia-gb200')) && has(p.placement) && has(p.placement.policyName)) : true)", message="GPU Topology is supported only for A4X machine family or nvidia-gb200 GPU type together with placement (workload) policy"
// +kubebuilder:validation:XValidation:rule="self.priorities.all(p, (has(p.spot) && p.spot)) || !has(self.priorityDefaults) || !has(self.priorityDefaults.nodeSystemConfig) || !has(self.priorityDefaults.nodeSystemConfig.kubeletConfig) || !has(self.priorityDefaults.nodeSystemConfig.kubeletConfig.shutdownGracePeriodSeconds)", message="shutdownGracePeriodSeconds is only supported for Spot"
// +kubebuilder:validation:XValidation:rule="has(self.minimumCapacity) && has(self.minimumCapacity.targetNodeCount) ? self.priorities.all(p, has(p.machineType) || has(p.gpu) || has(p.tpu) || (has(p.reservations) && p.reservations.affinity == 'Specific')) : true",message="Spec-level MinimumCapacity requires all priorities to have machine specifications (machineType, gpu, tpu, or specific reservation)"
// +kubebuilder:validation:XValidation:rule="(has(self.autopilot) && self.autopilot.enabled) ? ((!has(self.priorityDefaults) || !has(self.priorityDefaults.nodeSystemConfig) || !has(self.priorityDefaults.nodeSystemConfig.linuxNodeConfig) || !has(self.priorityDefaults.nodeSystemConfig.linuxNodeConfig.sysctls) || !has(self.priorityDefaults.nodeSystemConfig.linuxNodeConfig.sysctls.net__dot__ipv4__dot__tcp_congestion_control)) && self.priorities.all(p, !has(p.nodeSystemConfig) || !has(p.nodeSystemConfig.linuxNodeConfig) || !has(p.nodeSystemConfig.linuxNodeConfig.sysctls) || !has(p.nodeSystemConfig.linuxNodeConfig.sysctls.net__dot__ipv4__dot__tcp_congestion_control))) : true", message="Sysctl config net.ipv4.tcp_congestion_control cannot be set when Autopilot is enabled"
// +kubebuilder:validation:XValidation:rule="(has(self.autopilot) && self.autopilot.enabled) ? (!has(self.nodePoolConfig) || !has(self.nodePoolConfig.networkTags) || size(self.nodePoolConfig.networkTags) == 0) : true", message="networkTags cannot be used when Autopilot is enabled"
// +kubebuilder:validation:XValidation:rule="(has(self.nodePoolConfig) && has(self.nodePoolConfig.sandbox) && has(self.nodePoolConfig.sandbox.type) && self.nodePoolConfig.sandbox.type == 'microvm') ? (has(self.priorities) && size(self.priorities) > 0) : true", message="ComputeClass with sandbox type 'microvm' requires at least one priority to be specified"
// +kubebuilder:validation:XValidation:rule="(has(self.nodePoolConfig) && has(self.nodePoolConfig.sandbox) && has(self.nodePoolConfig.sandbox.type) && self.nodePoolConfig.sandbox.type == 'microvm') ? (has(self.priorities) ? self.priorities.all(p, has(p.nodepools) || (has(p.enableNestedVirtualization) ? p.enableNestedVirtualization : (has(self.priorityDefaults) && has(self.priorityDefaults.enableNestedVirtualization) && self.priorityDefaults.enableNestedVirtualization))) : true) : true", message="ComputeClass with sandbox type 'microvm' requires enableNestedVirtualization to be true for all priorities that are not preexisting nodepools"
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

	// MinimumCapacity defines declarative minimum node preprovisioning requirements
	// for the entire ComputeClass.
	// +optional
	MinimumCapacity *MinimumCapacity `json:"minimumCapacity,omitempty" protobuf:"bytes,11,opt,name=minimumCapacity"`

	// AllocationStrategyDefaults define the default allocation strategies for different provisioning models.
	//
	// +optional
	AllocationStrategyDefaults *AllocationStrategyDefaults `json:"allocationStrategyDefaults,omitempty" protobuf:"bytes,12,opt,name=allocationStrategyDefaults"`

	// NetworkConfig defines network-related settings for the ComputeClass.
	//
	// +optional
	NetworkConfig *NetworkConfig `json:"networkConfig,omitempty" protobuf:"bytes,13,opt,name=networkConfig"`
}

type NetworkingDra struct {
	// +optional
	// +kubebuilder:default=false
	Enabled bool `json:"enabled,omitempty" protobuf:"bytes,1,opt,name=enabled"`
}

// Dra represents a set of settings related to dynamic resource allocation
type Dra struct {
	Networking NetworkingDra `json:"networking,omitempty" protobuf:"bytes,1,opt,name=networking"`
}

// TpuDriverMode is an enumeration of supported Google TPU driver modes.
type TpuDriverMode string

const (
	// TpuDriverModeDevicePlugin enables managed device plugin mode for Google TPU driver.
	TpuDriverModeDevicePlugin TpuDriverMode = "DevicePlugin"
	// TpuDriverModeDynamicResourceAllocation enables managed DRA mode for Google TPU driver.
	TpuDriverModeDynamicResourceAllocation TpuDriverMode = "DynamicResourceAllocation"
)

// GoogleTpu describes how Google TPU should be functioning on the node
type GoogleTpu struct {
	// DriverMode determines the behaviour of the Google TPU driver.
	//
	// +kubebuilder:validation:Enum=DevicePlugin;DynamicResourceAllocation
	// +kubebuilder:default=DevicePlugin
	// +optional
	DriverMode TpuDriverMode `json:"driverMode,omitempty" protobuf:"bytes,1,opt,name=driverMode"`
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

	// ConfigDrift describes whether drifted nodes should be replaced to match the ComputeClass config.
	// A node is considered "drifted" if its current state diverges from the ComputeClass's configuration
	// within `nodePoolConfig` (e.g., nodeVersion, labels, taints, instanceMetadata).
	// When disabled, drift is ignored.
	//
	// +optional
	ConfigDrift *bool `json:"configDrift,omitempty" protobuf:"bytes,3,opt,name=configDrift"`

	// ReconciliationPolicy defines how nodes should be migrated.
	//
	// +optional
	ReconciliationPolicy *ReconciliationPolicy `json:"reconciliationPolicy,omitempty" protobuf:"bytes,4,opt,name=reconciliationPolicy"`
}

// ReconciliationPolicy describes how nodes should be migrated.
type ReconciliationPolicy struct {
	// Strategy defines the migration strategy to use.
	// It is only applicable to config drift active migration.
	// Optimize rule priority and ensure all daemon set pods running always use CreateBeforeDelete.
	// Supported values:
	// * CreateBeforeDelete
	// * DeleteBeforeCreate
	//
	// +optional
	// +kubebuilder:default=CreateBeforeDelete
	Strategy MigrationStrategy `json:"strategy,omitempty" protobuf:"bytes,1,opt,name=strategy"`

	// MaxNodeDisruption defines the maximum number of nodes that can be deleted at the same time during drift migration.
	//
	// +optional
	// +kubebuilder:validation:Minimum=1
	MaxNodeDisruption *int32 `json:"maxNodeDisruption,omitempty" protobuf:"bytes,2,opt,name=maxNodeDisruption"`

	// AtomicGroupLabels defines a list of node label keys used to group drifted nodes.
	// Nodes are only grouped together if they share the exact same values for ALL specified labels.
	//
	// +optional
	AtomicGroupLabels []string `json:"atomicGroupLabels,omitempty" protobuf:"bytes,3,rep,name=atomicGroupLabels"`
}

// MigrationStrategy defines the strategy used for active migration.
//
// +kubebuilder:validation:Enum=CreateBeforeDelete;DeleteBeforeCreate
type MigrationStrategy string

const (
	// MigrationStrategyCreateBeforeDelete creates new nodes before deleting old ones.
	MigrationStrategyCreateBeforeDelete MigrationStrategy = "CreateBeforeDelete"
	// MigrationStrategyDeleteBeforeCreate deletes old nodes before creating new ones.
	MigrationStrategyDeleteBeforeCreate MigrationStrategy = "DeleteBeforeCreate"
)

// ShieldedInstanceConfig defines the shielded instance configuration for auto-created node pools.
type ShieldedInstanceConfig struct {
	// EnableSecureBoot defines whether secure boot is enabled.
	// +optional
	// +kubebuilder:default=false
	EnableSecureBoot *bool `json:"enableSecureBoot,omitempty" protobuf:"bytes,1,opt,name=enableSecureBoot"`
	// EnableIntegrityMonitoring defines whether integrity monitoring is enabled.
	// +optional
	// +kubebuilder:default=true
	EnableIntegrityMonitoring *bool `json:"enableIntegrityMonitoring,omitempty" protobuf:"bytes,2,opt,name=enableIntegrityMonitoring"`
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

	// ShieldedInstanceConfig defines the shielded instance configuration for auto-created node pools.
	// +optional
	ShieldedInstanceConfig *ShieldedInstanceConfig `json:"shieldedInstanceConfig,omitempty" protobuf:"bytes,4,opt,name=shieldedInstanceConfig"`
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

type MaintenanceExclusionType string

const (
	// MaintenanceExclusionUntilEndOfSupport denotes an exclusion until the end of support of the nodepool's minor version.
	MaintenanceExclusionUntilEndOfSupport MaintenanceExclusionType = "UNTIL_END_OF_SUPPORT"
)

// NodePoolConfig defines required node pool configuration. Existing node pools will be matched with the ComputeClass
// only if their configuration match this field. Auto-provisioned node pools will be created with this configuration.
// +kubebuilder:validation:XValidation:rule="has(self.imageType) && self.imageType == 'custom_containerd' ? has(self.customImageConfig) : true", message="customImageConfig must be specified when imageType is custom_containerd"
// +kubebuilder:validation:XValidation:rule="has(self.customImageConfig) ? has(self.imageType) && self.imageType == 'custom_containerd' : true", message="customImageConfig can only be set when imageType is custom_containerd"
type NodePoolConfig struct {
	// ServiceAccount used by the node pool.
	//
	// +optional
	ServiceAccount string `json:"serviceAccount,omitempty" protobuf:"string,1,name=serviceAccount"`

	// Image type used by nodes in the node pool.
	//
	// +kubebuilder:validation:Enum=cos_containerd;ubuntu_containerd;custom_containerd
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

	// Dra describes settings related to dynamic resource allocation
	// and its integration with autoprovisioning
	//
	// +optional
	// +kubebuilder:validation:Optional
	Dra Dra `json:"dra,omitempty" protobuf:"bytes,13,opt,name=dra"`

	// IPType specifies whether the nodes in the node pool use public or private IP addresses.
	// Possible values are "public" or "private".
	// An empty string indicates the default IP type.
	// This setting corresponds to the presence and value of the cloud.google.com/private-node node selector.
	//
	// +optional
	// +kubebuilder:validation:Enum=public;private
	IPType string `json:"ipType,omitempty" protobuf:"string,14,opt,name=ipType"`

	// NodeVersion defines the GKE version to be used for the node pool.
	// If unspecified, the GKE cluster server will automatically pick a version
	// as per https://cloud.google.com/kubernetes-engine/versioning#specifying_node_version.
	//
	// +optional
	NodeVersion string `json:"nodeVersion,omitempty" protobuf:"string,15,opt,name=nodeVersion"`

	// Tpu defines node pool configuration for Google TPU.
	//
	// +optional
	Tpu GoogleTpu `json:"tpu,omitempty" protobuf:"bytes,16,opt,name=tpu"`

	// Sandbox contains sandbox configuration.
	//
	// +optional
	Sandbox *Sandbox `json:"sandbox,omitempty" protobuf:"bytes,17,opt,name=sandbox"`

	// WorkloadMetadata specifies how node metadata is exposed to the workload.
	// Possible values are "GCE_METADATA" or "GKE_METADATA".
	//
	// +optional
	// +kubebuilder:validation:Enum=GCE_METADATA;GKE_METADATA
	WorkloadMetadata *string `json:"workloadMetadata,omitempty" protobuf:"bytes,18,opt,name=workloadMetadata"`

	// InstanceMetadata is a map of custom key-value pairs to be injected into the underlying Compute Engine instances.
	// +optional
	// +kubebuilder:validation:MaxProperties=32
	// +kubebuilder:validation:XValidation:rule="self.all(k, k.matches('^[a-zA-Z0-9_-]+$') && size(k) < 128)", message="Metadata keys must be alphanumeric with dashes/underscores and less than 128 characters"
	// +kubebuilder:validation:XValidation:rule="self.all(k, size(self[k]) <= 32768)", message="Metadata values cannot exceed 32768 characters"
	// +kubebuilder:validation:XValidation:rule="self.all(k, !(k in ['cluster-location', 'cluster-name', 'cluster-uid', 'configure-sh', 'containerd-configure-sh', 'enable-os-login', 'gci-ensure-gke-docker', 'gci-metrics-enabled', 'gci-update-strategy', 'instance-template', 'kube-env', 'startup-script', 'user-data', 'disable-address-manager', 'windows-startup-script-ps1', 'common-psm1', 'k8s-node-setup-psm1', 'install-ssh-psm1', 'user-profile-psm1']))", message="Reserved metadata keys are not allowed"
	InstanceMetadata map[string]string `json:"instanceMetadata,omitempty" protobuf:"bytes,19,rep,name=instanceMetadata"`

	// TaintConfig contains node pool taint configuration.
	//
	// +optional
	TaintConfig *NodePoolTaintConfig `json:"taintConfig,omitempty" protobuf:"bytes,20,opt,name=taintConfig"`

	// MaintenanceExclusion defines the type of exclusion policy applied to node pools.
	// UNTIL_END_OF_SUPPORT - will not be upgraded until end of support of the nodepool's minor version
	//
	// +optional
	MaintenanceExclusion *MaintenanceExclusionType `json:"maintenanceExclusion,omitempty" protobuf:"bytes,21,opt,name=maintenanceExclusion"`

	// NodeDrainConfig contains node drain related configurations for node pool.
	//
	// +optional
	NodeDrainConfig *NodeDrainConfig `json:"nodeDrainConfig,omitempty" protobuf:"bytes,22,opt,name=nodeDrainConfig"`

	// Custom node image configuration used by nodes in the node pool.
	//
	// +optional
	CustomImageConfig *CustomImageConfig `json:"customImageConfig,omitempty" protobuf:"string,23,name=customImageConfig"`

	// ContainerdConfig defines customization for containerd.
	//
	// +optional
	ContainerdConfig *ContainerdConfig `json:"containerdConfig,omitempty" protobuf:"bytes,24,opt,name=containerdConfig"`

	// NetworkTags specifies network firewall tags assigned to all nodes created within the pool.
	// Note: Network tags are legacy and will be deprecated. It is recommended to use spec.nodePoolConfig.resourceManagerTags instead.
	//
	// +kubebuilder:validation:MaxItems=64
	// +optional
	NetworkTags []string `json:"networkTags,omitempty" protobuf:"bytes,25,rep,name=networkTags"`
}

type CustomImageConfig struct {
	// Image used by nodes in the node pool.
	//
	// +required
	ImageName string `json:"imageName,omitempty" protobuf:"string,1,name=imageName"`

	// Image project for the image used by nodes in the node pool.
	//
	// +required
	ImageProjectId string `json:"imageProjectId,omitempty" protobuf:"string,2,name=imageProjectId"`
}

// NodeDrainConfig contains node drain related configurations for node pool.
type NodeDrainConfig struct {
	// PdbTimeoutDuration specifies the duration of the PDB timeout period for node drain.
	// Must be a duration string ending with 's' unit/suffix (e.g., "60s", "100.5s").
	// Other time units such as 'm' or 'h' are not supported.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Pattern=`^[0-9]+(\.[0-9]{1,9})?s$`
	// +kubebuilder:validation:XValidation:rule="duration(self) >= duration('0s')",message="PDB timeout duration must be a non-negative duration"
	// +kubebuilder:validation:XValidation:rule="duration(self) <= duration('168h')",message="PDB timeout duration must be less than or equal to 168h"
	PdbTimeoutDuration *string `json:"pdbTimeoutDuration,omitempty" protobuf:"bytes,1,opt,name=pdbTimeoutDuration"`

	// GraceTerminationDuration specifies the duration of the grace termination period for node drain.
	// Must be a duration string ending with 's' unit/suffix (e.g., "60s", "100.5s").
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Pattern=`^[0-9]+(\.[0-9]{1,9})?s$`
	// +kubebuilder:validation:XValidation:rule="duration(self) >= duration('0s')",message="Graceful termination duration must be a non-negative duration"
	// +kubebuilder:validation:XValidation:rule="duration(self) <= duration('24h')",message="Graceful termination duration must be less than or equal to 24h"
	GraceTerminationDuration *string `json:"graceTerminationDuration,omitempty" protobuf:"bytes,2,opt,name=graceTerminationDuration"`

	// RespectPdbDuringNodePoolDeletion specifies whether to respect PDB during node pool deletion.
	//
	// +optional
	// +kubebuilder:validation:Optional
	RespectPdbDuringNodePoolDeletion *bool `json:"respectPdbDuringNodePoolDeletion,omitempty" protobuf:"bytes,3,opt,name=respectPdbDuringNodePoolDeletion"`
}

// NodePoolTaintConfig contains node pool taint configuration.
type NodePoolTaintConfig struct {
	// ArchitectureTaintBehavior specifies the behavior of architecture taint.
	// If set to NONE, architecture taint will not be applied to the nodes.
	// If set to ARM, it will be applied to ARM64 nodes.
	// Default is ARM if unspecified.
	//
	// +optional
	// +kubebuilder:validation:Enum=NONE;ARM
	ArchitectureTaintBehavior string `json:"architectureTaintBehavior,omitempty" protobuf:"string,1,opt,name=architectureTaintBehavior"`
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

// Sandbox stores sandbox configuration for nodepools.
type Sandbox struct {
	// Type defines the sandbox type (e.g., gvisor) for all nodes managed by this class.
	// +optional
	// +kubebuilder:validation:Enum=gvisor;microvm
	Type string `json:"type,omitempty" protobuf:"string,1,opt,name=type"`
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
//
// +kubebuilder:validation:XValidation:rule="!has(self.bootDiskStoragePools) || !has(self.bootDiskType) || self.bootDiskType == 'hyperdisk-balanced'", message="bootDiskStoragePools requires bootDiskType to be 'hyperdisk-balanced' or omitted"
// +kubebuilder:validation:XValidation:rule="has(self.bootDiskProvisionedIops) == has(self.bootDiskProvisionedThroughput)", message="bootDiskProvisionedIops and bootDiskProvisionedThroughput must be specified together"
// +kubebuilder:validation:XValidation:rule="(!has(self.bootDiskProvisionedIops) && !has(self.bootDiskProvisionedThroughput)) || (has(self.bootDiskType) && self.bootDiskType == 'hyperdisk-balanced')", message="bootDiskProvisionedIops and bootDiskProvisionedThroughput can only be specified for a Hyperdisk bootDiskType"
// +kubebuilder:validation:XValidation:rule="!has(self.localSsdEncryptionMode) || (has(self.localSSDCount) && self.localSSDCount > 0)", message="localSsdEncryptionMode can only be specified when localSSDCount is greater than 0"
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

	// BootDiskStoragePools defines a list of storage pools to be used for the boot disk.
	//
	// +kubebuilder:listType=atomic
	// +kubebuilder:validation:MinItems=1
	// +optional
	BootDiskStoragePools []BootDiskStoragePool `json:"bootDiskStoragePools,omitempty" protobuf:"bytes,6,rep,name=bootDiskStoragePools"`

	// BootDiskProvisionedIops defines the provisioned IOPS for the boot disk.
	// Only supported when bootDiskType is a Hyperdisk type (e.g. hyperdisk-balanced).
	// When custom performance is configured for Hyperdisk Balanced, both bootDiskProvisionedIops and bootDiskProvisionedThroughput must be specified together.
	//
	// +optional
	// +kubebuilder:validation:Minimum=2000
	// +kubebuilder:validation:Maximum=160000
	BootDiskProvisionedIops *int64 `json:"bootDiskProvisionedIops,omitempty" protobuf:"varint,7,opt,name=bootDiskProvisionedIops"`

	// BootDiskProvisionedThroughput defines the provisioned throughput in MB/s for the boot disk.
	// Only supported when bootDiskType is a Hyperdisk type (e.g. hyperdisk-balanced).
	// When custom performance is configured for Hyperdisk Balanced, both bootDiskProvisionedIops and bootDiskProvisionedThroughput must be specified together.
	//
	// +optional
	// +kubebuilder:validation:Minimum=140
	// +kubebuilder:validation:Maximum=2400
	BootDiskProvisionedThroughput *int64 `json:"bootDiskProvisionedThroughput,omitempty" protobuf:"varint,8,opt,name=bootDiskProvisionedThroughput"`

	// LocalSSDEncryptionMode specifies the encryption strategy for local SSD storage attached to instances in the priority level.
	// Currently supported modes:
	// * STANDARD_ENCRYPTION
	// * EPHEMERAL_KEY_ENCRYPTION
	//
	// +kubebuilder:validation:Enum=STANDARD_ENCRYPTION;EPHEMERAL_KEY_ENCRYPTION
	// +optional
	LocalSSDEncryptionMode *string `json:"localSsdEncryptionMode,omitempty" protobuf:"bytes,9,opt,name=localSsdEncryptionMode"`
}

// BootDiskStoragePool represents a storage pool configuration for a boot disk.
type BootDiskStoragePool struct {
	// Name defines the name of the storage pool.
	//
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Pattern=`^[a-z0-9-]+$`
	Name string `json:"name" protobuf:"bytes,1,name=name"`

	// Zone defines the GCP zone where the storage pool resides.
	//
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Pattern=`^[a-z0-9-]+$`
	Zone string `json:"zone" protobuf:"bytes,2,name=zone"`

	// Project defines the GCP project where the storage pool resides.
	// If omitted, defaults to the cluster's project.
	//
	// +optional
	// +kubebuilder:validation:Pattern=`^[a-z0-9-]+$`
	Project string `json:"project,omitempty" protobuf:"bytes,3,opt,name=project"`
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
	// Zones is a list of GCE zones where reservations are to be consumed.
	//
	// +kubebuilder:listType=atomic
	// +kubebuilder:validation:MinItems=1
	// +optional
	Zones []string `json:"zones,omitempty" protobuf:"bytes,4,opt,name=zones"`
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
// +kubebuilder:validation:Enum=Specific;AnyBestEffort;None;AnyThenFail
type ReservationAffinity string

const (
	// SpecificAffinity affinity allows to consume only specific reservations.
	SpecificAffinity ReservationAffinity = "Specific"
	// AnyBestEffortAffinity affinity allows to consume any reservation with a possibility to fallback to on demand.
	AnyBestEffortAffinity ReservationAffinity = "AnyBestEffort"
	// NoneAffinity prevents reservations from being used.
	NoneAffinity ReservationAffinity = "None"
	// AnyThenFail affinity allows to consume any reservation without a possibility to fallback to on demand.
	AnyThenFail ReservationAffinity = "AnyThenFail"
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
	// "None" affinity would prevent reservations from being used.
	// "AnyThenFail" affinity would consider any non-specific reservation available
	// to be claimed without a fallback to on-demand nodes in case of none claimable,
	// that means, node creation will fail if no reservation is available.
	//
	// +required
	Affinity ReservationAffinity `json:"affinity" protobuf:"bytes,2,name=affinity"`
}

// Priority is a specification of preferred machine characteristics.
//
// +kubebuilder:validation:MinProperties=1
// +kubebuilder:validation:XValidation:rule="has(self.nodepools) ? (size(dyn(self)) == 1 + (has(self.allocationStrategy) ? 1 : 0) + (has(self.priorityScore) ? 1 : 0)) : true", message="Nodepool field cannot be set along with other nodepool configuration fields"
// +kubebuilder:validation:XValidation:rule="!(has(self.machineFamily) && has(self.machineType))",message="MachineFamily and MachineType cannot be set together"
// +kubebuilder:validation:XValidation:rule="!(has(self.machineType) && (has(self.minCores) || has(self.minMemoryGb)))",message="MachineType cannot be set together with MinCores/MinMemoryGb"
// +kubebuilder:validation:XValidation:rule="!(has(self.machineFamily) && self.machineFamily == 'ek')", message="MachineFamily cannot be equal to 'ek'"
// +kubebuilder:validation:XValidation:rule="!(has(self.machineType) && self.machineType.startsWith('ek'))", message="MachineType cannot start with 'ek' prefix"
// +kubebuilder:validation:XValidation:rule="!(has(self.flexStart) && has(self.spot) && self.spot == true && self.flexStart.enabled == true)", message="Flex Start provisioning model is incompatible with Spot"
// +kubebuilder:validation:XValidation:rule="!has(self.capacityCheckWaitTimeSeconds) || has(self.tpu) || (has(self.flexStart) && self.flexStart.enabled)", message="capacityCheckWaitTimeSeconds is only supported for Flex Start and for multi-host TPUs"
// +kubebuilder:validation:XValidation:rule="(has(self.spot) && self.spot) || !has(self.nodeSystemConfig) || !has(self.nodeSystemConfig.kubeletConfig) || !has(self.nodeSystemConfig.kubeletConfig.shutdownGracePeriodSeconds)", message="shutdownGracePeriodSeconds is only supported for Spot"
// +kubebuilder:validation:XValidation:rule="!(has(self.gpuDirect) && self.gpuDirect == 'rdma') || has(self.acceleratorNetworkProfile)", message="acceleratorNetworkProfile must be specified when gpuDirect is 'rdma'"
// +kubebuilder:validation:XValidation:rule="has(self.minimumCapacity) && has(self.minimumCapacity.targetNodeCount) ? (has(self.machineType) || has(self.gpu) || has(self.tpu) || (has(self.reservations) && self.reservations.affinity == 'Specific')) : true",message="Priority-level MinimumCapacity requires a machineType, gpu, tpu, or specific reservation"
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
	// +kubebuilder:validation:Enum=general-purpose;general-purpose-arm
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
	// +kubebuilder:validation:MaxProperties=100
	NodeLabels map[string]string `json:"nodeLabels,omitempty" protobuf:"bytes,20,opt,name=nodeLabels"`

	// Taints is used to add user defined Kubernetes taints to all nodes in the new node pool.
	// These taints are applied to the Kubernetes API node object and can be used in tolerations for pod scheduling.
	//
	// +optional
	// +kubebuilder:validation:MaxItems=100
	Taints []TaintConfig `json:"taints,omitempty" protobuf:"bytes,21,opt,name=taints"`
	// A higher value is treated as a higher priority.
	// Priorities with the same priorityScore value are treated equally.
	// Not more than 3 priorities can have the same priorityScore.
	//
	// +optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=1000
	PriorityScore *int `json:"priorityScore,omitempty" protobuf:"bytes,22,opt,name=priorityScore"`

	// AcceleratorNetworkProfile defines the type of automated accelerator network provisioning to use.
	// Possible values:
	// "auto": Enables automatic ANP configuration based on the machine type.
	// "auto-<suffix>": Enables automatic ANP with a custom network profile suffix.
	// "<profile-name>": References a user-managed AcceleratorNetworkProfile resource name.
	// +kubebuilder:validation:MaxLength=63
	// +kubebuilder:validation:XValidation:rule="self.matches('^[a-z]([a-z0-9-]*[a-z0-9])?$') && size(self) <= 63",message="acceleratorNetworkProfile must be a valid resource ID up to 63 characters (lowercase letters, numbers, hyphens)"
	// +optional
	AcceleratorNetworkProfile *string `json:"acceleratorNetworkProfile,omitempty" protobuf:"bytes,24,opt,name=acceleratorNetworkProfile"`

	// GpuDirect defines the gpu direct strategy.
	// Possible values:
	// "rdma"
	// +kubebuilder:validation:XValidation:rule="self == 'rdma'",message="gpuDirect must be 'rdma'"
	// +optional
	GpuDirect string `json:"gpuDirect,omitempty" protobuf:"bytes,25,name=gpuDirect"`

	// MinimumCapacity defines declarative minimum node preprovisioning requirements
	// for this specific priority.
	// +optional
	MinimumCapacity *MinimumCapacity `json:"minimumCapacity,omitempty" protobuf:"bytes,26,opt,name=minimumCapacity"`

	// InstanceMetadata is a map of custom key-value pairs to be injected into the underlying Compute Engine instances.
	// Overrides conflicting keys defined in NodePoolConfig.InstanceMetadata.
	// +optional
	// +kubebuilder:validation:MaxProperties=32
	// +kubebuilder:validation:XValidation:rule="self.all(k, k.matches('^[a-zA-Z0-9_-]+$') && size(k) < 128)", message="Metadata keys must be alphanumeric with dashes/underscores and less than 128 characters"
	// +kubebuilder:validation:XValidation:rule="self.all(k, size(self[k]) <= 32768)", message="Metadata values cannot exceed 32768 characters"
	// +kubebuilder:validation:XValidation:rule="self.all(k, !(k in ['cluster-location', 'cluster-name', 'cluster-uid', 'configure-sh', 'containerd-configure-sh', 'enable-os-login', 'gci-ensure-gke-docker', 'gci-metrics-enabled', 'gci-update-strategy', 'instance-template', 'kube-env', 'startup-script', 'user-data', 'disable-address-manager', 'windows-startup-script-ps1', 'common-psm1', 'k8s-node-setup-psm1', 'install-ssh-psm1', 'user-profile-psm1']))", message="Reserved metadata keys are not allowed"
	InstanceMetadata map[string]string `json:"instanceMetadata,omitempty" protobuf:"bytes,27,rep,name=instanceMetadata"`

	// AllocationStrategy defines the allocation strategy for a node pool.
	//
	// +optional
	AllocationStrategy *AllocationStrategy `json:"allocationStrategy,omitempty" protobuf:"bytes,28,opt,name=allocationStrategy"`

	// EnableNestedVirtualization specifies whether to enable nested virtualization on the nodes.
	//
	// +optional
	EnableNestedVirtualization *bool `json:"enableNestedVirtualization,omitempty" protobuf:"bytes,29,opt,name=enableNestedVirtualization"`

	// PerformanceMonitoringUnit defines the virtualized performance monitoring unit configuration for hardware execution profiling.
	// Currently supported values:
	// * ARCHITECTURAL
	// * ENHANCED
	// * STANDARD
	//
	// +kubebuilder:validation:Enum=ARCHITECTURAL;ENHANCED;STANDARD
	// +optional
	PerformanceMonitoringUnit *string `json:"performanceMonitoringUnit,omitempty" protobuf:"bytes,30,opt,name=performanceMonitoringUnit"`
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
	// +kubebuilder:validation:Enum=default;latest;autoinstall-disabled
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

// ZoneType is an enumeration of supported zone types.
//
// +kubebuilder:validation:Enum=STANDARD;AI;CLUSTER_DEFAULT
type ZoneType string

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

	// ZoneTypes specifies sets of zones used for provisioning.
	// STANDARD zone type designates the core Google Cloud zones within a region.
	// AI zone type designates specialized zones optimized for AI.
	// CLUSTER_DEFAULT zone type designate zones specified in the cluster's autoprovisioningLocations or cluster’s locations if autoprovisioningLocations is empty.
	//
	// +optional
	// +kubebuilder:validation:UniqueItems:true
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=3
	ZoneTypes []ZoneType `json:"zoneTypes,omitempty" protobuf:"bytes,3,opt,name=zoneTypes"`
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

	// EnableNestedVirtualization specifies whether to enable nested virtualization on the nodes.
	//
	// +optional
	EnableNestedVirtualization *bool `json:"enableNestedVirtualization,omitempty" protobuf:"bytes,3,opt,name=enableNestedVirtualization"`
}

// AllocationStrategy is an enumeration of supported allocation strategies.
//
// +kubebuilder:validation:Enum=lowest-cost;fleet-efficiency
type AllocationStrategy string

const (
	// AllocationStrategyLowestCost uses the existing lowest price strategy.
	AllocationStrategyLowestCost AllocationStrategy = "lowest-cost"
	// AllocationStrategyFleetEfficiency prefers instances that would keep GCE fleet most efficient.
	AllocationStrategyFleetEfficiency AllocationStrategy = "fleet-efficiency"
)

// AllocationStrategyDefaults defines the default allocation strategies for different provisioning models.
// These can be overridden at priority level.
type AllocationStrategyDefaults struct {
	// OnDemand defines the default allocation strategy for on-demand provisioning model.
	//
	// +optional
	OnDemand *AllocationStrategy `json:"onDemand,omitempty" protobuf:"bytes,1,opt,name=onDemand"`

	// Spot defines the default allocation strategy for spot provisioning model.
	//
	// +optional
	Spot *AllocationStrategy `json:"spot,omitempty" protobuf:"bytes,2,opt,name=spot"`

	// FlexStart defines the default allocation strategy for flex start provisioning model.
	//
	// +optional
	FlexStart *AllocationStrategy `json:"flexStart,omitempty" protobuf:"bytes,3,opt,name=flexStart"`
}

// TaintConfig applies the given kubernetes taints on all nodes in the new node pool, which can be used with tolerations for pod scheduling.
// Any workload that does not tolerate the taints specified in this object will not be scheduled to the node pool.
// More info: https://cloud.google.com/sdk/gcloud/reference/container/node-pools/create#--node-taints
type TaintConfig struct {

	// Node taint key. The key must conform to syntax described in https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#syntax-and-character-set.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MaxLength=320
	Key string `json:"key,omitempty" protobuf:"bytes,1,opt,name=key"`

	// The value that matches the specified taint key.
	// +kubebuilder:validation:Pattern=`^([a-z0-9][-A-Za-z0-9_.]{1,61})?[A-Za-z0-9]$`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MaxLength=63
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
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty" protobuf:"bytes,1,rep,name=conditions"`

	// PriorityStatuses represent the statuses of Priorities within a given ComputeClass.
	// +optional
	PriorityStatuses []PriorityStatus `json:"priorityStatuses,omitempty" protobuf:"bytes,2,rep,name=priorityStatuses"`

	// ResourceInfo represents the current information about resource allocation and usage within the Compute Class.
	// +optional
	ResourceInfo []ResourceInfo `json:"resourceInfo,omitempty" protobuf:"bytes,3,rep,name=resourceInfo"`
}

// PriorityStatus describes a Status of ComputeClass priority.
type PriorityStatus struct {
	// Identifier represents the identifier of priority this PriorityStatus refers to.
	// If WhenUnsatisfiable is set to "ScaleUpAnyway", there will be an additional PriorityStatus with the identifier "ScaleUpAnyway",
	// and it will contain information about capacity provisioned as part of the implicit "ScaleUpAnyway" rule.
	Identifier string `json:"identifier,omitempty" protobuf:"bytes,1,opt,name=identifier"`

	// Conditions represent the observations of a priority current state.
	// +optional
	Conditions []metav1.Condition `json:"conditions" protobuf:"bytes,2,rep,name=conditions"`

	// ResourceInfo represents the current information about resource allocation and usage within the priority.
	// +optional
	ResourceInfo []ResourceInfo `json:"resourceInfo" protobuf:"bytes,3,rep,name=resourceInfo"`

	// ScalingEventsHistory represents the aggregated information about scaling events.
	// +optional
	ScalingEventsHistory *ScalingEventsHistory `json:"scalingEventsHistory,omitempty" protobuf:"bytes,4,opt,name=scalingEventsHistory"`

	// ConfigHash represents the combined hash of the global configuration and this specific priority.
	// This hash is also applied to node pools, enabling comparison to determine whether a node pool was created with the current configuration.
	// +optional
	ConfigHash string `json:"configHash,omitempty" protobuf:"bytes,5,opt,name=configHash"`
}

// ResourceName represents the resource a given ResourceInfo applies to. Can be one of "cpu", "memory", "ephemeral-storage", "nvidia.com/gpu", or "google.com/tpu".
// +kubebuilder:validation:Enum=cpu;memory;ephemeral-storage;nvidia.com/gpu;google.com/tpu
type ResourceName string

// ResourceUnit specifies the unit used to measure a resource.
// +kubebuilder:validation:Enum=Cores;GiB;Cards
type ResourceUnit string

// ResourceInfo describes current usage of resources.
type ResourceInfo struct {
	// Name is the name of a given resource measured in this ResourceInfo.
	Name *ResourceName `json:"name,omitempty" protobuf:"bytes,1,opt,name=name"`

	// Unit is a unit a given resource was measured in.
	Unit *ResourceUnit `json:"unit,omitempty" protobuf:"bytes,2,opt,name=unit"`

	// TargetCount represents the target count of a given resource within a priority. Can be lower than current count if there is ongoing node consolidation or higher, if there is ongoing node provisioning event.
	TargetCount *int `json:"targetCount,omitempty" protobuf:"bytes,3,opt,name=targetCount"`

	// CurrentCount represents the current count of a given resource.
	CurrentCount *int `json:"currentCount,omitempty" protobuf:"bytes,4,opt,name=currentCount"`

	// CurrentUtilizationPercentage represents the percentage of utilization for the resource `Name` at the `MeasuredAt` timestamp.
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=100
	CurrentUtilizationPercentage *int `json:"currentUtilizationPercentage,omitempty" protobuf:"bytes,5,opt,name=currentUtilizationPercentage"`

	// MeasuredAt represents the timestamp at which the resource information was measured.
	MeasuredAt *metav1.Time `json:"measuredAt,omitempty" protobuf:"bytes,6,opt,name=measuredAt"`
}

// ScalingEventsHistory represents the aggregated information about scaling events.
type ScalingEventsHistory struct {
	// ConsolidatedNodesCount represents how many nodes in this priority were consolidated.
	ConsolidatedNodesCount *int `json:"consolidatedNodesCount,omitempty" protobuf:"bytes,1,opt,name=consolidatedNodesCount"`

	// ProvisionedNodesCount represents how many nodes in this priority were added.
	ProvisionedNodesCount *int `json:"provisionedNodesCount,omitempty" protobuf:"bytes,2,opt,name=provisionedNodesCount"`

	// MigratedNodesCount represents how many nodes in this priority were removed as part of high priority migration.
	MigratedNodesCount *int `json:"migratedNodesCount,omitempty" protobuf:"bytes,3,opt,name=migratedNodesCount"`

	// MeasuredAt represents a timestamp at which the data was gathered.
	MeasuredAt *metav1.Time `json:"measuredAt,omitempty" protobuf:"bytes,4,opt,name=measuredAt"`

	// MeasuredSince represents a timestamp at which data started being collected.
	MeasuredSince *metav1.Time `json:"measuredSince,omitempty" protobuf:"bytes,5,opt,name=measuredSince"`
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

// ContainerdConfig defines customization for containerd.
type ContainerdConfig struct {
	// PrivateRegistryAccessConfig defines access configuration for private container registries.
	// +optional
	PrivateRegistryAccessConfig *PrivateRegistryAccessConfig `json:"privateRegistryAccessConfig,omitempty" protobuf:"bytes,1,opt,name=privateRegistryAccessConfig"`

	// WritableCgroups defines writable cgroups configuration for the node pool.
	// +optional
	WritableCgroups *WritableCgroups `json:"writableCgroups,omitempty" protobuf:"bytes,2,opt,name=writableCgroups"`

	// RegistryHosts defines containerd registry host configuration.
	// +kubebuilder:validation:MaxItems=25
	// +optional
	RegistryHosts []*RegistryHostConfig `json:"registryHosts,omitempty" protobuf:"bytes,3,rep,name=registryHosts"`
}

// PrivateRegistryAccessConfig defines access configuration for private container registries.
type PrivateRegistryAccessConfig struct {
	// Enabled is a boolean flag to enable or disable private registry access.
	// +optional
	Enabled *bool `json:"enabled,omitempty" protobuf:"varint,1,opt,name=enabled"`

	// CertificateAuthorityDomainConfig configures domain-specific certificate authorities for secure communication with private registries.
	// +optional
	CertificateAuthorityDomainConfig []*CertificateAuthorityDomainConfig `json:"certificateAuthorityDomainConfig,omitempty" protobuf:"bytes,2,rep,name=certificateAuthorityDomainConfig"`
}

// CertificateAuthorityDomainConfig configures domain-specific certificate authorities for secure communication with private registries.
type CertificateAuthorityDomainConfig struct {
	// FQDNs is a list of fully qualified domain names associated with the certificate authority.
	// +optional
	FQDNs []string `json:"fqdns,omitempty" protobuf:"bytes,1,rep,name=fqdns"`

	// GCPSecretManagerCertificateConfig specifies the location of CA certificates stored in Google Cloud Secret Manager.
	// +optional
	GCPSecretManagerCertificateConfig *GCPSecretManagerCertificateConfig `json:"gcpSecretManagerCertificateConfig,omitempty" protobuf:"bytes,2,opt,name=gcpSecretManagerCertificateConfig"`
}

// GCPSecretManagerCertificateConfig specifies the location of CA certificates stored in Google Cloud Secret Manager.
type GCPSecretManagerCertificateConfig struct {
	// SecretURI specifies the location of the secret in Google Cloud Secret Manager.
	//
	// +optional
	// +kubebuilder:validation:Pattern=`^projects/[^/]+(/locations/[^/]+)?/secrets/[^/]+/versions/[^/]+$`
	SecretURI *string `json:"secretURI,omitempty" protobuf:"bytes,1,opt,name=secretURI"`
}

// WritableCgroups defines writable cgroups configuration.
type WritableCgroups struct {
	// Enabled is a boolean flag to enable writable cgroups for the containerd runtime.
	// +optional
	Enabled *bool `json:"enabled,omitempty" protobuf:"varint,1,opt,name=enabled"`
}

// RegistryHostConfig defines containerd registry host configuration.
type RegistryHostConfig struct {
	// Server is the FQDN of the primary registry server.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=256
	Server string `json:"server,omitempty" protobuf:"bytes,1,opt,name=server"`

	// Hosts is the list of host configs for the registry server.
	// +kubebuilder:validation:MaxItems=10
	// +optional
	Hosts []*HostConfig `json:"hosts,omitempty" protobuf:"bytes,2,rep,name=hosts"`
}

// HostConfig defines mirror configurations for the primary registry server.
type HostConfig struct {
	// Host is the FQDN of the mirror host.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=256
	Host string `json:"host,omitempty" protobuf:"bytes,1,opt,name=host"`

	// Capabilities are permissions for the mirror (e.g., HOST_CAPABILITY_PULL, HOST_CAPABILITY_RESOLVE, HOST_CAPABILITY_PUSH).
	// +kubebuilder:validation:MaxItems=3
	// +optional
	// +kubebuilder:validation:items:Enum=HOST_CAPABILITY_PULL;HOST_CAPABILITY_RESOLVE;HOST_CAPABILITY_PUSH
	Capabilities []string `json:"capabilities,omitempty" protobuf:"bytes,2,rep,name=capabilities"`

	// OverridePath determines if the default path should be overridden.
	// +optional
	OverridePath *bool `json:"overridePath,omitempty" protobuf:"varint,3,opt,name=overridePath"`

	// Header is a list of custom HTTP headers for registry requests.
	// +kubebuilder:validation:MaxItems=5
	// +optional
	Header []*HostHeader `json:"header,omitempty" protobuf:"bytes,4,rep,name=header"`

	// CA is the list of certificate authorities to use for the registry.
	// +kubebuilder:validation:MaxItems=5
	// +optional
	CA []*RegistryHostCertificateConfig `json:"ca,omitempty" protobuf:"bytes,5,rep,name=ca"`

	// Client is the list of client-side TLS configurations (certificate and key).
	// +kubebuilder:validation:MaxItems=5
	// +optional
	Client []*RegistryHostClientCertificateConfig `json:"client,omitempty" protobuf:"bytes,6,rep,name=client"`

	// DialTimeout is the maximum amount of time a dial will wait for a connect to complete.
	// +optional
	// +kubebuilder:validation:Pattern=`^([0-9]+([.][0-9]+)?(ns|us|µs|ms|s|m|h))+$`
	DialTimeout *string `json:"dialTimeout,omitempty" protobuf:"bytes,7,opt,name=dialTimeout"`
}

// HostHeader defines custom HTTP headers for registry requests.
type HostHeader struct {
	// Key is the header key.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=256
	Key string `json:"key,omitempty" protobuf:"bytes,1,opt,name=key"`

	// Value is the list of header values.
	// +kubebuilder:validation:MaxItems=5
	// +kubebuilder:validation:items:MaxLength=256
	// +optional
	Value []string `json:"value,omitempty" protobuf:"bytes,2,rep,name=value"`
}

// RegistryHostCertificateConfig configures certificate for the registry.
type RegistryHostCertificateConfig struct {
	// GcpSecretManagerSecretUri specifies the secret in Google Cloud Secret Manager.
	//
	// +optional
	// +kubebuilder:validation:Pattern=`^projects/[^/]+(/locations/[^/]+)?/secrets/[^/]+/versions/[^/]+$`
	// +kubebuilder:validation:MaxLength=256
	GcpSecretManagerSecretUri *string `json:"gcpSecretManagerSecretUri,omitempty" protobuf:"bytes,1,opt,name=gcpSecretManagerSecretUri"`
}

// RegistryHostClientCertificateConfig configures pairs of certificates and keys for the registry client.
type RegistryHostClientCertificateConfig struct {
	// Cert specifies the client certificate.
	// +optional
	Cert *RegistryHostCertificateConfig `json:"cert,omitempty" protobuf:"bytes,1,opt,name=cert"`

	// Key specifies the client key.
	// +optional
	Key *RegistryHostCertificateConfig `json:"key,omitempty" protobuf:"bytes,2,opt,name=key"`
}
