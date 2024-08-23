package openslo

//go:generate go-enum --marshal --nocomments --names --values

// Version represents a version of the OpenSLO specification.
/* ENUM(
v1alpha = openslo/v1alpha
v1 = openslo/v1
v2alpha = openslo.com/v2alpha
)*/
type Version string
