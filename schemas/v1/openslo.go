package v1

// The goal of OpenSLO is to provide an open specification for defining SLOs to enable a
// common, vendor–agnostic approach to tracking and interfacing with SLOs. Platform-specific
// implementation details are purposefully excluded from the scope of this specification.
type OpenSLO struct {
	APIVersion string      `json:"apiVersion"`// The version of specification format for this particular entity that this is written; against.
	Kind       OpenSLOKind `json:"kind"`      
	Metadata   Metadata    `json:"metadata"`  
	Spec       SpecSl      `json:"spec"`      
}

type Metadata struct {
	Annotations map[string]string `json:"annotations,omitempty"`// key <> value pairs which should be used to define implementation / system specific; metadata about the SLO. For example, it can be metadata about a dashboard url, or how to; name a metric created by the SLI, etc.
	DisplayName *string           `json:"displayName,omitempty"`
	Labels      map[string]*Label `json:"labels,omitempty"`     // key <> value pairs of labels that can be used as metadata relevant to users
	Name        string            `json:"name"`                 
}

// An Alert Condition allows you to define under which conditions an alert for an SLO needs
// to be triggered.
//
// An Alert Notification Target defines the possible targets where alert notifications
// should be delivered to. For example, this can be a web-hook, Slack or any other custom
// target.
//
// An Alert Policy allows you to define the alert conditions for an SLO.
//
// A DataSource represents connection details with a particular metric source.
//
// A Service is a high-level grouping of SLO.
//
// A service level indicator (SLI) represents how to read metrics from data sources.
//
// A service level objective (SLO) is a target value or a range of values for a service
// level that is described by a service level indicator (SLI).
type SpecSl struct {
	Condition           *ACCondition           `json:"condition,omitempty"`          
	Description         *string                `json:"description,omitempty"`        
	Severity            *string                `json:"severity,omitempty"`           
	Target              *string                `json:"target,omitempty"`             
	AlertWhenBreaching  *bool                  `json:"alertWhenBreaching,omitempty"` 
	AlertWhenNoData     *bool                  `json:"alertWhenNoData,omitempty"`    
	AlertWhenResolved   *bool                  `json:"alertWhenResolved,omitempty"`  
	Conditions          []ConditionElement     `json:"conditions,omitempty"`         
	NotificationTargets []NotificationTarget   `json:"notificationTargets,omitempty"`
	ConnectionDetails   map[string]interface{} `json:"connectionDetails,omitempty"`  // Fields used for creating a connection with particular datasource e.g. AccessKeys,; SecretKeys, etc. Everything that is valid YAML can be put here.
	Type                *string                `json:"type,omitempty"`               
	RatioMetric         *RatioMetric           `json:"ratioMetric,omitempty"`        
	ThresholdMetric     *ThresholdMetric       `json:"thresholdMetric,omitempty"`    
	AlertPolicies       []string               `json:"alertPolicies,omitempty"`      
	BudgetingMethod     *BudgetingMethod       `json:"budgetingMethod,omitempty"`    
	Indicator           *Indicator             `json:"indicator,omitempty"`          
	IndicatorRef        *string                `json:"indicatorRef,omitempty"`       
	Objectives          []Objective            `json:"objectives,omitempty"`         
	Service             *string                `json:"service,omitempty"`            
	TimeWindow          []TimeWindow           `json:"timeWindow,omitempty"`         
}

type ACCondition struct {
	AlertAfter     *string        `json:"alertAfter,omitempty"`    
	Kind           *ConditionKind `json:"kind,omitempty"`          
	LookbackWindow *string        `json:"lookbackWindow,omitempty"`
	Threshold      *float64       `json:"threshold,omitempty"`     
}

type ConditionElement struct {
	Metadata     *Metadata           `json:"metadata,omitempty"`    
	ConditionRef *string             `json:"conditionRef,omitempty"`
	Kind         *string             `json:"kind,omitempty"`        
	Spec         *SpecAlertCondition `json:"spec,omitempty"`        
}

// An Alert Condition allows you to define under which conditions an alert for an SLO needs
// to be triggered.
type SpecAlertCondition struct {
	Condition   ACCondition `json:"condition"`            
	Description *string     `json:"description,omitempty"`
	Severity    string      `json:"severity"`             
}

type Indicator struct {
	Metadata *Metadata `json:"metadata,omitempty"`
	Spec     *SpecSLI  `json:"spec,omitempty"`    
}

// A service level indicator (SLI) represents how to read metrics from data sources.
type SpecSLI struct {
	RatioMetric     *RatioMetric     `json:"ratioMetric,omitempty"`    
	ThresholdMetric *ThresholdMetric `json:"thresholdMetric,omitempty"`
}

type RatioMetric struct {
	Bad     *Bad  `json:"bad,omitempty"` 
	Counter bool  `json:"counter"`       
	Good    *Good `json:"good,omitempty"`
	Total   Total `json:"total"`         
}

type Bad struct {
	MetricSource MetricSource `json:"metricSource"`
}

type MetricSource struct {
	MetricSourceRef *string                `json:"metricSourceRef,omitempty"`
	Spec            map[string]interface{} `json:"spec"`                     
	Type            *string                `json:"type,omitempty"`           
}

type Good struct {
	MetricSource MetricSource `json:"metricSource"`
}

type Total struct {
	MetricSource MetricSource `json:"metricSource"`
}

type ThresholdMetric struct {
	MetricSource MetricSource `json:"metricSource"`
}

type NotificationTarget struct {
	Kind      *string                      `json:"kind,omitempty"`     
	Metadata  *NotificationTargetMetadata  `json:"metadata,omitempty"` 
	Spec      *SpecAlertNotificationTarget `json:"spec,omitempty"`     
	TargetRef *string                      `json:"targetRef,omitempty"`
}

type NotificationTargetMetadata struct {
	Metadata *Metadata `json:"metadata,omitempty"`
}

// An Alert Notification Target defines the possible targets where alert notifications
// should be delivered to. For example, this can be a web-hook, Slack or any other custom
// target.
type SpecAlertNotificationTarget struct {
	Description *string `json:"description,omitempty"`
	Target      string  `json:"target"`               
}

type Objective struct {
	DisplayName     *string          `json:"displayName,omitempty"`    
	Op              *Op              `json:"op,omitempty"`             
	Target          float64          `json:"target"`                   
	TimeSliceTarget *float64         `json:"timeSliceTarget,omitempty"`
	TimeSliceWindow *TimeSliceWindow `json:"timeSliceWindow"`          
	Value           *float64         `json:"value,omitempty"`          
}

type TimeWindow struct {
	Duration string `json:"duration"`
}

type OpenSLOKind string
const (
	AlertCondition OpenSLOKind = "AlertCondition"
	AlertNotificationTarget OpenSLOKind = "AlertNotificationTarget"
	AlertPolicy OpenSLOKind = "AlertPolicy"
	DataSource OpenSLOKind = "DataSource"
	SLI OpenSLOKind = "SLI"
	Service OpenSLOKind = "Service"
	Slo OpenSLOKind = "SLO"
)

type BudgetingMethod string
const (
	Occurrences BudgetingMethod = "Occurrences"
	Timeslices BudgetingMethod = "Timeslices"
)

type ConditionKind string
const (
	Burnrate ConditionKind = "burnrate"
)

type Op string
const (
	Gt Op = "gt"
	Gte Op = "gte"
	LTE Op = "lte"
	Lt Op = "lt"
)

type Label struct {
	String      *string
	StringArray []string
}

type TimeSliceWindow struct {
	Double *float64
	String *string
}
