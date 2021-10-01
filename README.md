# ![ OpenSLO ](images/openslo.png)

## Table of Contents

- [OpenSLO](#openslo)
  - [Introduction](#introduction)
  - [Specification](#specification)
    - [Goals](#goals)
    - [Object Types](#object-types)
      - [General Schema](#general-schema)
        - [Notes](#notes)
      - [SLO](#slo)
        - [Notes](#notes-1)
        - [Objectives](#objectives)
      - [Service](#service)

## Introduction

The intent of this document is to outline the OpenSLO specification.

The goal of this project is to provide an open specification for defining and
interfacing with SLOs to allow for a common approach, giving a set vendor-agnostic
solution to defining and tracking SLOs. Platform specific implementation details
are purposefully excluded from the scope of this specification.

OpenSLO is an open specification-i.e., it is a specification created and
controlled, in an open and fair process, by an association or a standardization
body intending to achieve interoperability and interchangeability. An open
specification is not controlled by a single company or individual or by a group
with discriminatory membership criteria.  Additionally, this specification is
designed to be extended where needed to meet the needs of the implementation.

## Specification

### Goals

- Compliance with the Kubernetes YAML format.
- Vendor-agnostic
- Be flexible enough to be extended elsewhere

### Object Types

> 💡 **Note:** Specific attributes are described in detail in the **Notes** and
> under each integration section.

#### General Schema

```yaml
apiVersion: openslo/v1alpha
kind: SLO | Service
metadata:
  name: string
  displayName: string # optional
spec:
```

##### Notes

- **kind** *string* - required, either [SLO](#slo) or [Service](#service)
- **metadata.name:** *string* - required field, convention for naming object from
  [DNS RFC1123](https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names)
  `name` should:

  - contain at most 63 characters
  - contain only lowercase alphanumeric characters or `-`
  - start with an alphanumeric character
  - end with an alphanumeric character

---

#### SLO

A service level objective (SLO) is a target value or range of values for
a service level that is described by a service level indicator (SLI).

```yaml
apiVersion: openslo/v1alpha
kind: SLO
metadata:
  name: string
  displayName: string # optional
  labels: # optional, key <>value a pair of labels that can be used as metadata relevant to users
    userImpacting: "true"
    team: "identity"
    costCentre: "project1"
    serviceTier: "tier-1"
    tags:
      - auth
  annotations: map[string]string # optional, key <> value a pair of annotations that can be used as implementation metadata
    openslo.com/key1: value1
    fooimplementation.com/key2: value2
spec:
  description: string # optional
  service: [service name] # name of the service to associate this SLO with
  indicator: # represents the Service Level Indicator (SLI)
    thresholdMetric: # represents the metric used to inform the Service Level Object in the objectives stanza
      source: string # data source for the metric
      queryType: string # a name for the type of query to run on the data source
      query: string # the query to run to return the metric
      metadata: # optional, allows data source specific details to be passed
  timeWindows:
    # exactly one item, one of possible rolling time window or calendar aligned
    # rolling time window
    - unit: Second
      count: numeric
      isRolling: true
    # or
    # calendar aligned time window
    - unit: Year | Quarter | Month | Week | Day
      count: numeric # count of time units for example count: 7 and unit: Day means 7 days window
      calendar:
        startTime: 2020-01-21 12:30:00 # date with time in 24h format, format without time zone
        timeZone: America/New_York # name as in IANA Time Zone Database
      isRolling: false # false or not defined
  budgetingMethod: Occurrences | Timeslices
  objectives: # see objectives below for details
```

##### Notes (SLO)

- **metadata.labels:** *map[string]string|string[]* - optional field `key` <> `value`
  - the `key` segment is required and must contain at most 63 characters beginning and ending
     with an alphanumeric character `[a-z0-9A-Z]` with dashes `-`, underscores `_`, dots `.`
     and alphanumerics between.
  - the `value` of `key` segment can be a string or an array of strings
- **metadata.annotations:** *map[string]string* - optional field `key` <> `value`
  - `annotations` should be used to define implementation / system specific metadata about the SLO.
    For example, it can be metadata about a dashboard url, or how to name a metric created by the SLI, etc.
  - `key` have two segments: an optional `prefix` and `name`, separated by a slash `/`
  - the `name` segment is required and must contain at most 63 characters beginning and ending
    with an alphanumeric character `[a-z0-9A-Z]` with dashes `-`, underscores `_`, dots `.`
    and alphanumerics between.
  - the `prefix` is optional and must be a DNS subdomain: a series of DNS labels separated by dots `.`,
    it must contain at most 253 characters, followed by a slash `/`.
  - the `openslo.com/` is reserved for OpenSLO usage
- **indicator** optional, represents the Service Level Indicator (SLI).
  Currently this only supports one Metric, `thresholdMetric`, with `ratioMetric`
  supported in the [objectives](#objectives) stanza.
- **indicator.thresholdMetric** *Metric*, represents the query used for
  gathering data from metric sources. Raw data is used to compare objectives
  (threshold) values. If `thresholdMetric` is defined then `ratioMetrics`
  should be excluded in [objectives](#objectives).
- **timeWindows\[ \]** *TimeWindow* is a list but accepting only exactly one
  item, one of the rolling or calendar aligned
    time window:
  - Rolling time window. Minimum duration for rolling time window is 5
      minutes, maximum 31 days).

      ```yaml
      unit: Day | Hour | Minute
      count: numeric
      isRolling: true
      ```

  - Calendar Aligned time window. Minimum duration for calendar aligned time
    window is 1 day and maximum is 366 days.

      ```yaml
      unit: Year | Quarter | Month | Week | Day
      count: numeric
      calendar:
          startTime: 2020-01-21 12:30:00 # date with time in 24h format
          timeZone: America/New_York # name as in IANA Time Zone Database
      # isRolling: false # for calendar aligned set false value or not set
      ```

- **description** *string* optional field, contains at most 1050 characters

- **budgetingMethod** *enum(Occurrences \| Timeslices)*, required field
  - Occurrences method uses a ratio of counts of good events and total count of
    the event.
  - Timeslices method uses a ratio of good time slices vs. total time slices in
    a budgeting period.

- **objectives\[ \]** *Threshold*, required field, described in [Objectives](#objectives)
  section

##### Objectives

Objectives are the thresholds for your SLOs. You can use objectives to define
the tolerance levels for your metrics

```yaml
objectives:
  - displayName: string # optional
    op: lte | gte | lt | gt # conditional operator used to compare the SLI against the value. Only needed when using a thresholdMetric
    value: numeric # optional, value used to compare threshold metrics. Only needed when using a thresholdMetric
    target: numeric [0.0, 1.0) # budget target for given objective of the SLO
    timeSliceTarget: numeric (0.0, 1.0] # required only when budgetingMethod is set to TimeSlices
    # ratioMetric {good, total} should be defined only if thresholdMetric is not set.
    # ratioMetric good and total have to contain the same source type configuration (for example for prometheus).
    ratioMetric:
        incremental: true | false #todo: add description
        good: # the numerator
          source: string # data source for the "good" numerator
          queryType: string # a name for the type of query to run on the data source
          query: string # the query to run to return the numerator
          metadata: # optional, allows data source specific details to be passed
        total: # the denominator
          source: string # data source for the "total" denominator
          queryType: string # a name for the type of query to run on the data source
          query: string # the query to run to return the denominator
          metadata: # optional, allows data source specific details to be passed
```

Example:

```yaml
objectives:
  - displayName: Foo Total Errors
    target: 0.98
    ratioMetrics:
        incremental: true
        good:
          source: datadog
          queryType: query
          query: sum:requests.error{*}
        total:
          source: datadog
          queryType: query
          query: sum:requests.total{*}
```

The optional `metadata` key can be used to pass extraneous data to the data source,
for example, if your data source accepts variables, they could be passed via the
`metadata`.

An example is an internal tool which uses a templating feature to ease the maintenance
of SLOs by not repeating the same queries were there are only small differences.
The internal tool can then generate the final OpensLO specification based on the input described.

```yaml
objectives:
  - displayName: Foo Latency
    target: 0.98
    ratioMetrics:
        incremental: true
        good:
          source: prometheus
          queryType: query
          query: sum:requests.duration{*}
          metadata:
            bucket: "0.25"
            exclude_errors: "true"
        total:
          source: prometheus
          queryType: query
          query: sum:requests.duration{*}
          metadata:
            exclude_errors: "true"
```

##### Notes (Objectives)

- **objectives\[ \]** *Threshold*, required field. If `thresholdMetric` has
  been defined, only one Threshold can be defined. However if using `ratioMetric`
  then any number of Thresholds can be defined.

- **op** *enum(lte | gte | lt | gt)*, operator used to compare the SLI against
  the value. Only needed when using `thresholdMetric`

- **value numeric**, required field, used to compare values gathered from
  metric source. Only needed when using a `thresholdMetric`.

- **target numeric** *\[0.0, 1.0)*, required, budget target for given objective
  of the SLO

- **targetTimeSlices** *numeric* *\[0.0, 1.0\]*, required only when budgeting
  method is set to TimeSlices

- **indicator.ratioMetric** *Metric {Good, Total}*, if `ratioMetric` is defined
    then `thresholdMetric` should not be set in `indicator`

  - *Good* represents the query used for gathering data from metric sources used
   as the numerator. Received data is used to compare objectives (threshold)
   values to find good values.

  - *Total* represents the query used for gathering data from metric sources
    that is used as the denominator. Received data is used to compare objectives
    (threshold) values to find total number of metrics.

---

#### Service

A Service is a high-level grouping of SLO.

```yaml
apiVersion: openslo/v1alpha
kind: Service
metadata:
  name: string
  displayName: string # optional
spec:
  description: string # optional up to 1050 characters
```
