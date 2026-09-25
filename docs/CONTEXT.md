# Domain context

This guide defines the terms used in OpenSLO contributor documentation and agent instructions.
The [specification](../website/docs/specification.md) defines object fields and requirements.

## Reliability concepts

**Service level indicator (SLI)** is a quantitative measure of some aspect of a service's behavior, such as availability or latency.

**Service level objective (SLO)** is a target value or range for a service level measured by an SLI.

**Service level agreement (SLA)** is an agreement with users that defines consequences for meeting or missing its SLOs.

**Time window** is the period over which an SLO measures reliability.
OpenSLO supports rolling and calendar-aligned windows.

**Error budget** is the amount of unreliability that an SLO allows within its time window.
Failures consume that budget.

**Budgeting method** determines how measurements contribute to an SLO's success rate and error budget.
The [v1 specification](../website/docs/specification.md#notes-slo) defines three methods:

- `Occurrences` uses the number of good events divided by the total number of events.
  For example, at least 95% of requests over two weeks have latency below 100 ms.
- `Timeslices` uses the number of good time slices divided by the total number of time slices.
  Each slice is an interval classified as good or bad against the objective's `timeSliceTarget`.
  For example, classify each one-minute slice by its fraction of requests with latency below 100 ms.
  The SLO counts the good slices over two weeks.
- `RatioTimeslices` uses the average success ratio across all time slices in the budgeting period.
  For example, each one-minute slice measures the fraction of requests with latency below 100 ms.
  The SLO averages those fractions over two weeks.

**Ratio metric** expresses an SLI through a ratio of good or bad events to total events, or through a precomputed ratio.
OpenSLO converts a bad-event ratio to a success ratio when calculating the indicator value.

**Threshold metric** supplies values that an objective compares with a threshold.
OpenSLO uses the operators `lt`, `lte`, `gt`, and `gte` for these comparisons.
The [SLI specification](../website/docs/specification.md#sli) describes both metric forms.

**Composite SLO** combines multiple objectives, each of which can have its own queries, data sources, and targets.
The composite consumes error budget when any of its objectives consumes error budget.
All objectives use the same budgeting method.
That method determines how failures and objective weights affect the combined budget.
See the [composite SLO rules](../website/docs/specification.md#notes-composite-slo) for the calculations.

## OpenSLO objects and versions

**Specification** is the written definition of the OpenSLO format and its semantics.
It describes service level objectives independently of a particular vendor or implementation.
Draft proposals under [enhancements/](../enhancements/) describe proposed changes, such as [v2](../enhancements/v2alpha.md).

**API version** identifies the version of the format used by an object through its `apiVersion` field.
Examples include `openslo/v1` and the draft `openslo.com/v2alpha`.
Labels such as `v1` and `v2alpha` are shorthand used in documentation and schema page paths.
An SDK release has its own version, separate from the API versions it supports.

**Object** is one instance of an OpenSLO definition, such as a named SLO, typically expressed in YAML.
Its top-level fields are `apiVersion`, `kind`, `metadata`, and `spec`.
See the [general schema](../website/docs/specification.md#general-schema) for this structure.

**Object kind** is the category selected by an object's `kind` field, such as `SLO` or `DataSource`.
The API version determines the available kinds and their fields.
The v1 specification defines these kinds:

| Kind | Meaning |
| --- | --- |
| `Service` | A grouping that multiple SLOs can reference. |
| `SLO` | An object that configures a service level objective. |
| `SLI` | An object that configures metric queries for a service level indicator. |
| `DataSource` | Connection details for a metric source. |
| `AlertPolicy` | Alert conditions and notification targets associated with an SLO. |
| `AlertCondition` | Conditions that determine when an alert triggers. |
| `AlertNotificationTarget` | A destination for alert notifications. |

**Metadata** identifies and describes an object through fields such as `metadata.name`.
The API version defines the supported metadata fields.

**Object specification (`spec`)** is the part of an object that contains the settings for its kind.
For example, an SLO's `spec` contains its objectives.
Use “specification” for the format's written definition and `spec` for the object field.

**Object reference** identifies another OpenSLO object by name in fields such as `indicatorRef` or `conditionRef` in v1.
Some fields also allow an **inline object**, which embeds a definition inside the containing object.
The selected API version and field determine the supported forms.

## SDK and schema terminology

**Go SDK** is the implementation in the separate [OpenSLO/go-sdk repository](https://github.com/OpenSLO/go-sdk).
It provides object types, decoding, and validation.
The website uses its documentation output, and the example checker uses its decoder and validator.
SDK validation can differ from the specification, as the [SLO validation note](../website/docs/specification.md#slo) explains.

**Schema manifest**, or **manifest** in the contributor and website guides, is the JSON document produced by the SDK's documentation generator.
The SDK writes `docs/manifest.json`, which this repository imports as [website/api.json](../website/api.json).
It groups object descriptions, properties, examples, and validation rules by API version and object kind.
Use “OpenSLO object” or “YAML example” for a user-authored object definition to distinguish it from this documentation input.

**Schema reference** is the website documentation rendered from the manifest.
It describes the imported SDK types and validation rules.
Authored pages can add examples around the generated content.
See the [schema guide](../website/docs/schema.md) for how to read the reference.

**Property path** locates a field within an object in the manifest and schema reference.
The manifest uses `$` for the object root and paths such as `$.metadata.name` for fields.
Paths use `[*]` for each array element, `.*` for each map value, and `.*~` for each map key.

**Validation rule** describes a constraint on a value, sometimes subject to conditions.
A **validation plan** is a structured description of validation rules that the SDK exports for documentation.
The manifest associates rules with property paths.
Its `typeInfo.kind` describes a Go type category, such as `struct` or `string`.
Use “object kind” for the OpenSLO category selected by the object's `kind` field.

## Relationships

An SLO uses an SLI, a time window, and a budgeting method to express its reliability target.
The OpenSLO objects configure these concepts and their connections to data sources, services, and alerts.

The specification and draft proposals describe the format.
The Go SDK implements object types and validation, then exports a schema manifest for documentation.
This repository imports that manifest and renders the website's schema reference from it.
Schema pages for a draft API version document that draft's implementation without making the proposal stable.

The imported manifest and the [example checker's SDK dependency](../website/tools/example-check/go.mod) are separate inputs.
They can reflect different SDK revisions.
The [specification change workflow](specification-changes.md) explains how to keep them aligned.
The [website guide](../website/README.md#schema-reference) covers manifest generation, import, and rendering.

## Further reading

- [Site Reliability Engineering: How Google Runs Production Systems](https://sre.google/sre-book/table-of-contents/)
- [Implementing Service Level Objectives](https://www.oreilly.com/library/view/implementing-service-level/9781492076803/)
