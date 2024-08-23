package openslo

//go:generate go-enum --marshal --nocomments --names --values

// Kind represents all the object kinds defined by OpenSLO specification.
// Keep in mind not all specification versions support every [Kind].
/* ENUM(
SLO
SLI
DataSource
Service
AlertPolicy
AlertCondition
AlertNotificationTarget
)*/
type Kind string
