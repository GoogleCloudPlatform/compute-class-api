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
}

// NodePoolAutoCreation defines node-pool autoprovisioning related settings.
type NodePoolAutoCreation struct {
	// Enabled indicates whether NodePoolAutoCreation is enabled for a given ComputeClass.
	//
	// +kubebuilder:validation:Required
	// +kubebuilder:default=false
	Enabled bool `json:"enabled" protobuf:"bytes,1,name=enabled"`
}

// Priority is a specification of preferred machine characteristics.
//
// +kubebuilder:validation:MinProperties=1
// +kubebuilder:validation:XValidation:rule="has(self.nodepools) ? (size(dyn(self)) == 1) : true", message="Nodepool field cannot be set along with other fields"
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
}

// ComputeClassStatus is the current status of the ComputeClass.
type ComputeClassStatus struct {
	// Conditions represent the observations of a ComputeClass's current state.
	//
	// +optional
	Conditions []metav1.Condition `json:"conditions" protobuf:"bytes,1,rep,name=conditions"`
}
