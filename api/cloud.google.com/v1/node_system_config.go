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

	SwapConfig *SwapConfig `json:"swapConfig,omitempty" protobuf:"bytes,5,opt,name=swapConfig"`

	// Additional entries to be added to /etc/hosts.
	// +optional
	AdditionalEtcHosts []*EtcHostsEntry `json:"additionalEtcHosts,omitempty" protobuf:"bytes,6,rep,name=additionalEtcHosts"`

	// Additional entries to be added to /etc/resolv.conf.
	// +optional
	AdditionalEtcResolvConf []*ResolvedConfEntry `json:"additionalEtcResolvConf,omitempty" protobuf:"bytes,7,rep,name=additionalEtcResolvConf"`

	// Additional entries to be added to /etc/systemd/resolved.conf.
	// +optional
	AdditionalEtcSystemdResolvedConf []*ResolvedConfEntry `json:"additionalEtcSystemdResolvedConf,omitempty" protobuf:"bytes,8,rep,name=additionalEtcSystemdResolvedConf"`

	// Support for running custom init code while bootstrapping nodes.
	// +optional
	CustomNodeInit *CustomNodeInit `json:"customNodeInit,omitempty" protobuf:"bytes,9,opt,name=customNodeInit"`

	// Parameters that can be configured on the kernel.
	// +optional
	KernelOverrides *KernelOverrides `json:"kernelOverrides,omitempty" protobuf:"bytes,10,opt,name=kernelOverrides"`

	// Configures the timezone of the node.
	// +kubebuilder:validation:MaxLength=256
	// +optional
	TimeZone *string `json:"timeZone,omitempty" protobuf:"bytes,11,opt,name=timeZone"`

	// AccurateTimeConfig defines accurate time configuration for a node.
	// +optional
	AccurateTimeConfig *AccurateTimeConfig `json:"accurateTimeConfig,omitempty" protobuf:"bytes,12,opt,name=accurateTimeConfig"`

	// Contains VFIO-related configurations for this node.
	// +optional
	NodeVfioConfig *NodeVfioConfig `json:"nodeVfioConfig,omitempty" protobuf:"bytes,13,opt,name=nodeVfioConfig"`

	// Controls the configuration for the disk IO scheduler.
	// +optional
	DiskIoScheduler *DiskIoScheduler `json:"diskIoScheduler,omitempty" protobuf:"bytes,14,opt,name=diskIoScheduler"`
}

// AccurateTimeConfig defines accurate time configuration for a node.
type AccurateTimeConfig struct {
	// EnablePtpKvmTimeSync controls whether to enable accurate time synchronization with PTP-KVM.
	// +optional
	EnablePtpKvmTimeSync *bool `json:"enablePtpKvmTimeSync,omitempty" protobuf:"bytes,1,opt,name=enablePtpKvmTimeSync"`
}

// NodeVfioConfig defines configuration settings for VFIO on a node.
type NodeVfioConfig struct {
	// Specifies the maximum number of DMA entries (pages) that can be mapped
	// by the VFIO IOMMU type 1 driver for a container. This limit affects the
	// total amount of host memory that can be pinned for direct device access,
	// which is often critical for high-performance devices like TPUs and GPUs.
	// This setting corresponds to the kernel parameter at:
	// /sys/module/vfio_iommu_type1/parameters/dma_entry_limit
	// The default value in the kernel is 65535. Higher values may be
	// needed for workloads mapping large memory regions.
	//
	// +kubebuilder:validation:Minimum=65535
	// +kubebuilder:validation:Maximum=4194304
	// +optional
	DmaEntryLimit *int32 `json:"dmaEntryLimit,omitempty" protobuf:"bytes,1,opt,name=dmaEntryLimit"`
}

// DiskIoScheduler contains the configuration for the disk IO scheduler.
type DiskIoScheduler struct {
	// Configures the IO scheduler for the boot disk or ephemeral lssd that runs
	// node system workloads.
	//
	// +kubebuilder:validation:Enum=mq-deadline;bfq;kyber;none
	// +optional
	NodeSystemIoScheduler string `json:"nodeSystemIoScheduler,omitempty" protobuf:"bytes,1,opt,name=nodeSystemIoScheduler"`

	// Configures the IO scheduler for the attached disks.
	//
	// +kubebuilder:validation:Enum=mq-deadline;bfq;kyber;none
	// +optional
	NodeAttachedDiskIoScheduler string `json:"nodeAttachedDiskIoScheduler,omitempty" protobuf:"bytes,2,opt,name=nodeAttachedDiskIoScheduler"`
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
// +kubebuilder:validation:XValidation:rule="!has(self.shutdownGracePeriodCriticalPodsSeconds) || (has(self.shutdownGracePeriodSeconds) && self.shutdownGracePeriodCriticalPodsSeconds <= self.shutdownGracePeriodSeconds)", message="ShutdownGracePeriodCriticalPodsSeconds must be less than or equal to ShutdownGracePeriodSeconds and requires ShutdownGracePeriodSeconds to be set"
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
	// ShutdownGracePeriodSeconds is the maximum allowed grace period
	// (in seconds) that the node should delay the shutdown during a graceful shutdown.
	//
	// +kubebuilder:validation:Enum=0;30;120
	// +kubebuilder:validation:Optional
	ShutdownGracePeriodSeconds *int32 `json:"shutdownGracePeriodSeconds,omitempty" protobuf:"bytes,20,opt,name=shutdownGracePeriodSeconds"`
	// ShutdownGracePeriodCriticalPodsSeconds is the maximum allowed grace period
	// (in seconds) that is used to terminate critical pods during a node shutdown.
	// This value should be <= ShutdownGracePeriodSeconds, and is only valid if ShutdownGracePeriodSeconds is set.
	//
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=120
	// +kubebuilder:validation:Optional
	ShutdownGracePeriodCriticalPodsSeconds *int32 `json:"shutdownGracePeriodCriticalPodsSeconds,omitempty" protobuf:"bytes,21,opt,name=shutdownGracePeriodCriticalPodsSeconds"`

	// CrashLoopBackOff contains the configuration to modify node level parameters
	// for container restart behavior.
	//
	// +kubebuilder:validation:Optional
	CrashLoopBackOff *CrashLoopBackOff `json:"crashLoopBackOff,omitempty" protobuf:"bytes,22,opt,name=crashLoopBackOff"`

	// ReservedResourcesConfig contains the configuration for the reserved resources on the node.
	// +kubebuilder:validation:Optional
	ReservedResourcesConfig *ReservedResourcesConfig `json:"reservedResourcesConfig,omitempty" protobuf:"bytes,23,opt,name=reservedResourcesConfig"`
	// This setting defines the maximum number of concurrent workers to rotate container log files.
	// The value must be an integer between 1 and 10, inclusive.
	// Default is 1 if unspecified.
	//
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=10
	// +kubebuilder:validation:Optional
	ContainerLogMaxWorkers *int64 `json:"containerLogMaxWorkers,omitempty" protobuf:"bytes,24,opt,name=containerLogMaxWorkers"`
	// This setting specifies the interval at which the container logs are monitored for performing the log rotate operation.
	// The value must be a positive duration between 3s and 300s, inclusive.
	// Default is "10s" if unspecified.
	//
	// +kubebuilder:validation:Pattern=`^[0-9]+(\.[0-9]{1,9})?s$`
	// +kubebuilder:validation:XValidation:rule="duration(self) >= duration('3s')",message="containerLogMonitorInterval must be greater than or equal to 3s"
	// +kubebuilder:validation:XValidation:rule="duration(self) <= duration('300s')",message="containerLogMonitorInterval must be less than or equal to 300s"
	// +kubebuilder:validation:Optional
	ContainerLogMonitorInterval *string `json:"containerLogMonitorInterval,omitempty" protobuf:"bytes,25,opt,name=containerLogMonitorInterval"`
}

// CrashLoopBackOff contains the configuration to modify node level parameters
// for container restart behavior.
type CrashLoopBackOff struct {
	// MaxContainerRestartPeriod is the maximum duration the backoff delay can
	// accrue to for container restarts. If not set, defaults to the internal
	// crashloopbackoff maximum.
	// The value must be a duration string. Eg. "30s", "1m30s".
	// The value must be positive and less than '5m'.
	// +kubebuilder:validation:Pattern=`^([0-9]+([.][0-9]+)?(ns|us|µs|ms|s|m|h))+$`
	// +optional
	MaxContainerRestartPeriod *string `json:"maxContainerRestartPeriod,omitempty" protobuf:"bytes,1,opt,name=maxContainerRestartPeriod"`
}

// ReservedResourcesConfig contains the configuration for the reserved resources on the node.
type ReservedResourcesConfig struct {
	// CpuReservedMillicore is the amount of CPU to reserve for system daemons.
	// This is a user-specified value. If unspecified, GKE decides the default.
	// +optional
	CpuReservedMillicore *int64 `json:"cpuReservedMillicore,omitempty" protobuf:"varint,1,opt,name=cpuReservedMillicore"`

	// MemoryReservedMib is the amount of memory to reserve for system daemons (in MiB).
	// This is a user-specified value. If unspecified, GKE decides the default.
	// +optional
	MemoryReservedMib *int64 `json:"memoryReservedMib,omitempty" protobuf:"varint,2,opt,name=memoryReservedMib"`
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
	// Controls TCP Packetization-Layer Path MTU Discovery. Supported values:
	// 0: Disabled
	// 1: Disabled by default, enabled when an ICMP black hole detected
	// 2: Always enabled, use initial MSS of tcp_base_mss.
	//
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=2
	// +kubebuilder:validation:Optional
	Net_ipv4_tcp_mtu_probing *int64 `json:"net.ipv4.tcp_mtu_probing,omitempty" protobuf:"bytes,41,opt,name=net.ipv4.tcp_mtu_probing"`
	// Maximal number of timewait sockets held by system simultaneously. If this number is exceeded time-wait socket is immediately destroyed and warning is printed.
	//
	// +kubebuilder:validation:Minimum=4096
	// +kubebuilder:validation:Maximum=2147483647
	// +kubebuilder:validation:Optional
	Net_ipv4_tcp_max_tw_buckets *int64 `json:"net.ipv4.tcp_max_tw_buckets,omitempty" protobuf:"bytes,42,opt,name=net.ipv4.tcp_max_tw_buckets"`
	// Number of times initial SYNs for an active TCP connection attempt will be retransmitted.
	//
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=127
	// +kubebuilder:validation:Optional
	Net_ipv4_tcp_syn_retries *int64 `json:"net.ipv4.tcp_syn_retries,omitempty" protobuf:"bytes,43,opt,name=net.ipv4.tcp_syn_retries"`
	// Control use of Explicit Congestion Notification (ECN) by TCP. ECN is used only when both ends of the TCP connection indicate support for it.
	//
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=2
	// +kubebuilder:validation:Optional
	Net_ipv4_tcp_ecn *int64 `json:"net.ipv4.tcp_ecn,omitempty" protobuf:"bytes,44,opt,name=net.ipv4.tcp_ecn"`
	// Set the congestion control algorithm to be used for new connections. The algorithm “reno” is always available, but additional choices may be available based on kernel configuration. Default is set as part of kernel configuration. For passive connections, the listener congestion control choice is inherited.
	//
	// +kubebuilder:validation:MaxLength=10
	// +kubebuilder:validation:Pattern=`^[a-zA-Z0-9_]+$`
	// +kubebuilder:validation:Optional
	Net_ipv4_tcp_congestion_control *string `json:"net.ipv4.tcp_congestion_control,omitempty" protobuf:"bytes,45,opt,name=net.ipv4.tcp_congestion_control"`
	// Controls use of the performance events system by unprivileged users (without CAP_PERFMON). The default value is 2 in kernel.
	//
	// +kubebuilder:validation:Minimum=-1
	// +kubebuilder:validation:Maximum=3
	// +kubebuilder:validation:Optional
	Kernel_perf_event_paranoid *int64 `json:"kernel.perf_event_paranoid,omitempty" protobuf:"bytes,46,opt,name=kernel.perf_event_paranoid"`
	// A global limit on how much time real-time scheduling may use.
	//
	// +kubebuilder:validation:Minimum=-1
	// +kubebuilder:validation:Maximum=1000000
	// +kubebuilder:validation:Optional
	Kernel_sched_rt_runtime_us *int64 `json:"kernel.sched_rt_runtime_us,omitempty" protobuf:"bytes,47,opt,name=kernel.sched_rt_runtime_us"`
	// Control whether the kernel panics when a soft lockup is detected.
	//
	// +kubebuilder:validation:Optional
	Kernel_softlockup_panic *bool `json:"kernel.softlockup_panic,omitempty" protobuf:"bytes,48,opt,name=kernel.softlockup_panic"`
	// Defines the scope and restrictions for the ptrace() system call, impacting process debugging and tracing. Supported values:
	// 0: Classic ptrace permissions.
	// 1: Restricted ptrace (default in many distributions) - only child processes or CAP_SYS_PTRACE.
	// 2: Admin-only ptrace - only processes with CAP_SYS_PTRACE.
	// 3: No ptrace - ptrace calls are disallowed.
	//
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=3
	// +kubebuilder:validation:Optional
	Kernel_yama_ptrace_scope *int64 `json:"kernel.yama.ptrace_scope,omitempty" protobuf:"bytes,49,opt,name=kernel.yama.ptrace_scope"`
	// Indicates whether restrictions are placed on exposing kernel addresses via /proc and other interfaces.
	//
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=2
	// +kubebuilder:validation:Optional
	Kernel_kptr_restrict *int64 `json:"kernel.kptr_restrict,omitempty" protobuf:"bytes,50,opt,name=kernel.kptr_restrict"`
	// Indicates whether unprivileged users are prevented from using dmesg(8) to view messages from the kernel’s log buffer.
	//
	// +kubebuilder:validation:Optional
	Kernel_dmesg_restrict *bool `json:"kernel.dmesg_restrict,omitempty" protobuf:"bytes,51,opt,name=kernel.dmesg_restrict"`
	// Controls the functions allowed to be invoked via the SysRq key. List of possible values:
	// 0: Disables sysrq completely.
	// 1: Enables all sysrq functions.
	// >1 - bitmask of allowed sysrq functions. More details in https://docs.kernel.org/admin-guide/sysrq.html.
	//
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=511
	// +kubebuilder:validation:Optional
	Kernel_sysrq *int64 `json:"kernel.sysrq,omitempty" protobuf:"bytes,52,opt,name=kernel.sysrq"`
	// Contains the amount of dirty memory at which the background kernel flusher threads will start writeback.
	// Note: Vm_dirty_background_bytes is the counterpart of Vm_dirty_background_ratio. Only one of them may be specified at a time.
	//
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=68719476736
	// +kubebuilder:validation:Optional
	Vm_dirty_background_bytes *int64 `json:"vm.dirty_background_bytes,omitempty" protobuf:"bytes,53,opt,name=vm.dirty_background_bytes"`
	// Contains the amount of dirty memory at which a process generating disk writes will itself start writeback.
	// Note: Vm_dirty_bytes is the counterpart of Vm_dirty_ratio. Only one of them may be specified at a time.
	// Note: the minimum value allowed for Vm_dirty_bytes is two pages (in bytes); any value lower than this limit will be ignored and the old configuration will be retained.
	//
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=68719476736
	// +kubebuilder:validation:Optional
	Vm_dirty_bytes *int64 `json:"vm.dirty_bytes,omitempty" protobuf:"bytes,54,opt,name=vm.dirty_bytes"`
	// Defines the core dump pattern for the kernel.
	// Only absolute paths are supported. Piping and relative paths are not allowed.
	//
	// +kubebuilder:validation:MaxLength=128
	// +kubebuilder:validation:Pattern="^/[a-zA-Z0-9/._%-]+$"
	// +kubebuilder:validation:Optional
	Kernel_core_pattern *string `json:"kernel.core_pattern,omitempty" protobuf:"bytes,55,opt,name=kernel.core_pattern"`
	// Controls the maximum number of keys that a nonroot user may own.
	//
	// +kubebuilder:validation:Minimum=200
	// +kubebuilder:validation:Maximum=1048576
	// +kubebuilder:validation:Optional
	Kernel_keys_maxkeys *int64 `json:"kernel.keys.maxkeys,omitempty" protobuf:"bytes,56,opt,name=kernel.keys.maxkeys"`
	// Represents the maximum number of bytes that a nonroot user can hold in the payload section of all their keys.
	//
	// +kubebuilder:validation:Minimum=20000
	// +kubebuilder:validation:Maximum=2097152
	// +kubebuilder:validation:Optional
	Kernel_keys_maxbytes *int64 `json:"kernel.keys.maxbytes,omitempty" protobuf:"bytes,57,opt,name=kernel.keys.maxbytes"`
	// Tells the garbage collector the minimum number of network entries that can sit in cache (floor).
	//
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=262144
	// +kubebuilder:validation:Optional
	Net_ipv4_neigh_default_gc_thresh1 *int64 `json:"net.ipv4.neigh.default.gc_thresh1,omitempty" protobuf:"bytes,58,opt,name=net.ipv4.neigh.default.gc_thresh1"`
	// Acts as a soft limit to the number of network device entries stored in cache (soft ceiling).
	//
	// +kubebuilder:validation:Minimum=512
	// +kubebuilder:validation:Maximum=524288
	// +kubebuilder:validation:Optional
	Net_ipv4_neigh_default_gc_thresh2 *int64 `json:"net.ipv4.neigh.default.gc_thresh2,omitempty" protobuf:"bytes,59,opt,name=net.ipv4.neigh.default.gc_thresh2"`
	// Sets a hard ceiling (absolute maximum) for the network neighbor cache.
	//
	// +kubebuilder:validation:Minimum=1024
	// +kubebuilder:validation:Maximum=1048576
	// +kubebuilder:validation:Optional
	Net_ipv4_neigh_default_gc_thresh3 *int64 `json:"net.ipv4.neigh.default.gc_thresh3,omitempty" protobuf:"bytes,60,opt,name=net.ipv4.neigh.default.gc_thresh3"`
	// Configure IPv6 forwarding on default network interface.
	//
	// +kubebuilder:validation:Optional
	Net_ipv6_conf_default_forwarding *bool `json:"net.ipv6.conf.default.forwarding,omitempty" protobuf:"bytes,61,opt,name=net.ipv6.conf.default.forwarding"`
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

// SwapConfig specifies the swap memory configuration for a node pool.
// +kubebuilder:validation:XValidation:rule="(has(self.bootDiskProfile) ? 1 : 0) + (has(self.ephemeralLocalSsdProfile) ? 1 : 0) + (has(self.dedicatedLocalSsdProfile) ? 1 : 0) <= 1",message="only one of bootDiskProfile, ephemeralLocalSsdProfile, or dedicatedLocalSsdProfile may be set"
type SwapConfig struct {
	// Enables or disables swap for the node pool. Default to false.
	Enabled bool `json:"enabled,omitempty" protobuf:"bytes,1,opt,name=enabled"`

	// If omitted, swap space is encrypted by default.
	// +optional
	EncryptionConfig *SwapConfigEncryptionConfig `json:"encryptionConfig,omitempty" protobuf:"bytes,2,opt,name=encryptionConfig"`

	// --- Performance Profile (oneof) ---
	// Only ONE of the following profile fields should be set (non-nil).

	// Use the node's boot disk for swap.
	// +optional
	BootDiskProfile *SwapConfigBootDiskProfile `json:"bootDiskProfile,omitempty" protobuf:"bytes,3,opt,name=bootDiskProfile,oneof=performanceProfile"`
	// Use the local SSD (shared with ephemeral storage) for swap.
	// +optional
	EphemeralLocalSsdProfile *SwapConfigEphemeralLocalSsdProfile `json:"ephemeralLocalSsdProfile,omitempty" protobuf:"bytes,4,opt,name=ephemeralLocalSsdProfile,oneof=performanceProfile"`
	// Provision a new, separate local NVMe SSD exclusively for swap.
	// +optional
	DedicatedLocalSsdProfile *SwapConfigDedicatedLocalSsdProfile `json:"dedicatedLocalSsdProfile,omitempty" protobuf:"bytes,5,opt,name=dedicatedLocalSsdProfile,oneof=performanceProfile"`
}

// SwapConfigEncryptionConfig defines encryption settings for the swap space.
type SwapConfigEncryptionConfig struct {
	// If true, swap space will NOT be encrypted.
	// Defaults to false, swap space is encrypted by default.
	Disabled bool `json:"disabled,omitempty" protobuf:"bytes,1,opt,name=disabled,proto3"`
}

// SwapConfigBootDiskProfile defines swap on the node's boot disk.
// +kubebuilder:validation:XValidation:rule="(has(self.swapSizeGib) ? 1 : 0) + (has(self.swapSizePercent) ? 1 : 0) <= 1",message="only one of swapSizeGib or swapSizePercent may be set"
type SwapConfigBootDiskProfile struct {
	// --- Swap Size (oneof) ---
	// Only one of the following size fields should be set.

	// The size of the swap space in GiB.
	// +kubebuilder:validation:Minimum=1
	// +optional
	SwapSizeGib *int64 `json:"swapSizeGib,omitempty" protobuf:"bytes,1,opt,name=swapSizeGib,oneof=swapSize"`
	// The size of the swap space as a percentage of the node's boot disk.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=50
	// +optional
	SwapSizePercent *int32 `json:"swapSizePercent,omitempty" protobuf:"bytes,2,opt,name=swapSizePercent,oneof=swapSize"`
}

// SwapConfigEphemeralLocalSsdProfile defines swap on the local SSD.
// +kubebuilder:validation:XValidation:rule="(has(self.swapSizeGib) ? 1 : 0) + (has(self.swapSizePercent) ? 1 : 0) <= 1",message="only one of swapSizeGib or swapSizePercent may be set"
type SwapConfigEphemeralLocalSsdProfile struct {
	// --- Swap Size (oneof) ---
	// Only one of the following size fields should be set.

	// The size of the swap space in GiB.
	// +kubebuilder:validation:Minimum=1
	// +optional
	SwapSizeGib *int64 `json:"swapSizeGib,omitempty" protobuf:"bytes,1,opt,name=swapSizeGib,oneof=swapSize"`
	// The size of the swap space as a percentage of the node's ephemeral storage local SSDs.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=80
	// +optional
	SwapSizePercent *int32 `json:"swapSizePercent,omitempty" protobuf:"bytes,2,opt,name=swapSizePercent,oneof=swapSize"`
}

// SwapConfigDedicatedLocalSsdProfile provisions a new local SSD for swap.
type SwapConfigDedicatedLocalSsdProfile struct {
	// +kubebuilder:validation:Minimum=1
	// The number of physical local NVMe SSD disks to attach.
	DiskCount int64 `json:"diskCount,omitempty" protobuf:"bytes,1,opt,name=diskCount"`
}

// EtcHostsEntry defines an entry in /etc/hosts.
type EtcHostsEntry struct {
	// The IPv4 or IPv6 address of the host.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	Ip string `json:"ip,omitempty" protobuf:"bytes,1,opt,name=ip"`
	// The hostname of the host.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	Host string `json:"host,omitempty" protobuf:"bytes,2,opt,name=host"`
}

// ResolvedConfEntry defines an entry in resolved.conf.
type ResolvedConfEntry struct {
	// The key of resolved.conf
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	Key string `json:"key,omitempty" protobuf:"bytes,1,opt,name=key"`
	// The value of resolved.conf
	// +kubebuilder:validation:MaxItems=256
	// +optional
	Value []string `json:"value,omitempty" protobuf:"bytes,2,rep,name=value"`
}

// CustomNodeInit defines the init script to be executed on the node.
type CustomNodeInit struct {
	// The init script to be executed on the node.
	// +optional
	InitScript *InitScript `json:"initScript,omitempty" protobuf:"bytes,1,opt,name=initScript"`
}

// InitScript defines the init script source and arguments.
type InitScript struct {
	// The Cloud Storage URI for storing the init script.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	GcsUri *string `json:"gcsUri,omitempty" protobuf:"bytes,1,opt,name=gcsUri"`
	// The generation of the init script stored in GCS.
	// +optional
	GcsGeneration *int64 `json:"gcsGeneration,omitempty" protobuf:"varint,2,opt,name=gcsGeneration"`
	// Optional arguments to be passed to the init script.
	// +kubebuilder:validation:MaxItems=50
	// +kubebuilder:validation:items:MinLength=1
	// +kubebuilder:validation:items:MaxLength=512
	// +optional
	Args []string `json:"args,omitempty" protobuf:"bytes,3,rep,name=args"`
	// The resource name of the secret manager secret hosting the init script.
	// +kubebuilder:validation:MaxLength=256
	// +optional
	GcpSecretManagerSecretUri *string `json:"gcpSecretManagerSecretUri,omitempty" protobuf:"bytes,4,opt,name=gcpSecretManagerSecretUri"`
}

// KernelOverrides defines kernel parameters.
type KernelOverrides struct {
	// Optional kernel command line arguments overrides.
	// +optional
	KernelCommandlineOverrides *KernelCommandlineOverrides `json:"kernelCommandlineOverrides,omitempty" protobuf:"bytes,1,opt,name=kernelCommandlineOverrides"`
	// LRU Gen (Multi-Gen LRU) options.
	// +optional
	LruGen *LRUGen `json:"lruGen,omitempty" protobuf:"bytes,2,opt,name=lruGen"`
}

// KernelCommandlineOverrides defines kernel command line argument overrides.
type KernelCommandlineOverrides struct {
	// Defines the change of spec_rstack_overflow.
	// +kubebuilder:validation:Enum=SPEC_RSTACK_OVERFLOW_UNSPECIFIED;SPEC_RSTACK_OVERFLOW_OFF
	// +optional
	SpecRstackOverflow *string `json:"specRstackOverflow,omitempty" protobuf:"bytes,1,opt,name=specRstackOverflow"`
	// Defines the change of init_on_alloc.
	// +kubebuilder:validation:Enum=INIT_ON_ALLOC_UNSPECIFIED;INIT_ON_ALLOC_OFF
	// +optional
	InitOnAlloc *string `json:"initOnAlloc,omitempty" protobuf:"bytes,2,opt,name=initOnAlloc"`
}

// LRUGen defines Multi-Gen LRU options.
type LRUGen struct {
	// Enable LRU Gen.
	// +optional
	Enabled *bool `json:"enabled,omitempty" protobuf:"varint,1,opt,name=enabled"`
	// Prevent working set of N milliseconds from getting evicted.
	// +optional
	MinTtlMs *int32 `json:"minTtlMs,omitempty" protobuf:"varint,2,opt,name=minTtlMs"`
}
