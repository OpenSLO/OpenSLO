package v1

import (
	"github.com/OpenSLO/OpenSLO/pkg/openslo"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = openslo.Object(SLO{})

// SLOStatus defines the observed state of SLO
type SLOStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Conditions         []metav1.Condition `json:"conditions,omitempty"`
	CurrentSLO         string             `json:"currentSLO,omitempty"`
	LastEvaluationTime metav1.Time        `json:"lastEvaluationTime,omitempty"`
	Ready              string             `json:"ready,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:storageversion
//+kubebuilder:printcolumn:name="Ready",type=string,JSONPath=.status.ready,description="The reason for the current status of the SLO resource"
//+kubebuilder:printcolumn:name="Window",type=string,JSONPath=.spec.timeWindow[0].duration,description="The time window for the SLO resource"
//+kubebuilder:printcolumn:name="Age",type=date,JSONPath=.metadata.creationTimestamp,description="The time when the SLO resource was created"

type SLO struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SLOSpec   `json:"spec,omitempty"`
	Status SLOStatus `json:"status,omitempty"`
}

func (s SLO) GetVersion() openslo.Version {
	return APIVersion
}

func (s SLO) GetKind() openslo.Kind {
	return openslo.KindSLO
}

func (s SLO) GetName() string {
	return s.ObjectMeta.Name
}

func (s SLO) Validate() error {
	return nil
}

type SLOSpec struct {
	Description     string        `json:"description,omitempty"`
	Service         string        `json:"service"`
	Indicator       *SLOIndicator `json:"indicator,omitempty"`
	IndicatorRef    *string       `json:"indicatorRef,omitempty"`
	BudgetingMethod string        `json:"budgetingMethod"`
	TimeWindow      []TimeWindow  `json:"timeWindow,omitempty"`
	Objectives      []Objective   `json:"objectives"`
	// We don't make clear in the spec if this is a ref or inline.
	// We will make it a ref for now.
	// https://github.com/OpenSLO/OpenSLO/issues/133
	AlertPolicies []string `json:"alertPolicies,omitempty"`
}

type SLOIndicator struct {
	Metadata Metadata `json:"metadata"`
	Spec     SLISpec  `json:"spec"`
}

type Objective struct {
	DisplayName     string  `json:"displayName,omitempty"`
	Op              string  `json:"op,omitempty"`
	Value           float64 `json:"value,omitempty"`
	Target          float64 `json:"target"`
	TimeSliceTarget float64 `json:"timeSliceTarget,omitempty"`
	TimeSliceWindow string  `json:"timeSliceWindow,omitempty"`
}

type TimeWindow struct {
	Duration string `json:"duration"`
	// +kubebuilder:default=true
	IsRolling bool      `json:"isRolling"`
	Calendar  *Calendar `json:"calendar,omitempty"`
}

type Calendar struct {
	StartTime string `json:"startTime"`
	TimeZone  string `json:"timeZone"`
}

//+kubebuilder:object:root=true

// SLOList contains a list of SLO
type SLOList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SLO `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SLO{}, &SLOList{})
}
