#

<!-- markdownlint-disable MD033-->
<picture>
  <source media="(prefers-color-scheme: dark)" srcset="images/openslo_light.png">
  <img alt="OpenSLO light theme" src="images/openslo.png">
</picture>
<!-- markdownlint-enable MD033-->

## Table of Contents

- [Introduction](#introduction)
- [Specification](#specification)
  - [Goals](#goals)
  - [General Schema](#general-schema)
    - [Notes (General Schema)](#notes-general-schema)
  - [Custom Data Types](#custom-data-types)
    - [duration-shorthand](#duration-shorthand)
  - [Object Types](#object-types)
    - [DataSource](#datasource)
      - [Notes (DataSource)](#notes-datasource)
    - [SLO](#slo)
      - [Notes (SLO)](#notes-slo)
      - [Objectives](#objectives)
        - [Notes (Objectives)](#notes-objectives)
    - [SLI](#sli)
      - [Notes (SLI)](#notes-sli)
      - [Ratio Metric](#ratio-metric)
    - [AlertPolicy](#alertpolicy)
      - [Notes (AlertPolicy)](#notes-alertpolicy)
    - [AlertCondition](#alertcondition)
      - [Notes (AlertCondition)](#notes-alertcondition)
    - [AlertNotificationTarget](#alertnotificationtarget)
      - [Notes (AlertNotificationTarget)](#notes-alertnotificationtarget)
    - [Service](#service)
- [Examples](examples/README.md)
- Work in progress for future versions
  - [v2alpha1](enhancements/v2alpha1.md)

## Introduction

The intent of this document is to outline the OpenSLO specification.

The goal of this project is to provide an open specification for defining SLOs
to enable a common, vendor–agnostic approach to tracking and interfacing with
SLOs. Platform-specific implementation details are purposefully excluded from
the scope of this specification.

OpenSLO is an open specification i.e., it is a specification created and
controlled, in an open and fair process, by an association or a standardization
body intending to achieve interoperability and interchangeability. An open
specification is not controlled by a single company or individual or by a group
with discriminatory membership criteria. Additionally, this specification is
designed to be extended where needed to meet the needs of the implementation.

Before making a contribute please read our [contribution guideline](https://github.com/OpenSLO/OpenSLO/blob/main/CONTRIBUTING.md).

## Specification

### Goals

- Compliance with the Kubernetes YAML format
- Vendor-agnostic
- Be flexible enough to be extended elsewhere

### General Schema

```yaml
apiVersion: openslo/v1
kind: DataSource | SLO | SLI | AlertPolicy | AlertCondition | AlertNotificationTarget | Service
metadata:
  name: string
  displayName: string # optional
  labels: # optional, key <>value a pair of labels that can be used as metadata relevant to users
    # Example labels
    userImpacting: "true"
    team: "identity"
    costCentre: "project1"
    serviceTier: "tier-1"
    tags:
      - auth
  annotations: map[string]string # optional, key <> value a pair of annotations that can be used as implementation metadata
    # Example annotations
    openslo.com/key1: value1
    fooimplementation.com/key2: value2
spec:
```

#### Notes (General Schema)

- **kind** *string* - required, one of: [DataSource](#datasource), [SLO](#slo),
  [SLI](#sli), [AlertPolicy](#alertpolicy), [AlertCondition](#alertcondition),
  [AlertNotificationTarget](#alertnotificationtarget), [Service](#service)
- **metadata.name:** *string* - required field
  - all implementations must at least support object names that follow [RFC1123][rfc1123-names]:
    - are up to 63 characters in length
    - contain lowercase alphanumeric characters or `-`
    - start with an alphanumeric character
    - end with an alphanumeric character
  - implementations are additionally encouraged to support names that:
    - are up to 255 characters in length
    - contain lowercase alphanumeric characters or `-`, `.`, `|`, `/`, `\`
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

### Custom Data Types

#### duration-shorthand

The duration shorthand is specified as a single–word string (no whitespaces) consisting
of a positive integer `number` followed by a case–sensitive single–character `postfix`.

Allowed postfixes are:

- *m* – minutes
- *h* – hours
- *d* – days
- *w* – weeks
- *M* – months
- *Q* – quarters
- *Y* – years

Examples: `12h`, `4w`, `1M`, `1Q`, `365d`, `1Y`.

This specification does not put requirements on how (or whether) to implement each
postfix, therefore implementers are free to pick an implementation that best suits
their environments.

There is however the possibility that future versions of this spec will take a more
prescriptive stance on this issue.

### Object Types

> 💡 **Note:** Specific attributes are described in detail in the **Notes**
> subsection of each object type's section.

#### DataSource

A DataSource represents connection details with a particular metric source.

> [Check work in progress for v2.](enhancements/v2alpha1.md#datasource)

```yaml
apiVersion: openslo/v1
kind: DataSource
metadata:
  name: string
  displayName: string # optional
spec:
  description: string # optional up to 1050 characters
  type: string # predefined type e.g. Prometheus, Datadog, etc.
  connectionDetails:
    # fields used for creating a connection with particular datasource e.g. AccessKeys, SecretKeys, etc.
    # everything that is valid YAML can be put here
```

##### Notes (DataSource)

DataSource enables reusing one source between many SLOs and moving
connection specific details (e.g. authentication) away from SLO definitions.

This spec does not enforce naming conventions for data source types, however
the OpenSLO project will publish guidelines in the form of supplementary materials
once common patterns start emerging from implementations.

An example of the DataSource kind can be:

```yaml
apiVersion: openslo/v1
kind: DataSource
metadata:
  name: string
  displayName: string # optional
spec:
  type: CloudWatch
  connectionDetails:
    accessKeyID: accessKey
    secretAccessKey: secretAccessKey
```

---

#### SLO

A service level objective (SLO) is a target value or a range of values for
a service level that is described by a service level indicator (SLI).

> [Check work in progress for v2.](enhancements/v2alpha1.md#slo)

```yaml
apiVersion: openslo/v1
kind: SLO
metadata:
  name: string
  displayName: string # optional
spec:
  description: string # optional up to 1050 characters
  service: string # name of the service to associate this SLO with, may refer (depends on implementation) to existing object Kind: Service
  indicator: # see SLI below for details
  indicatorRef: string # name of the SLI. Required if indicator is not given.
  timeWindow:
    # exactly one item; one of possible: rolling or calendar–aligned time window
    ## rolling time window
    - duration: duration-shorthand # duration of the window eg 1d, 4w
      isRolling: true
    # or
    ## calendar–aligned time window
    - duration: duration-shorthand # duration of the window eg 1M, 1Q, 1Y
      calendar:
        startTime: 2020-01-21 12:30:00 # date with time in 24h format, format without time zone
        timeZone: America/New_York # name as in IANA Time Zone Database
      isRolling: false # if omitted assumed `false` if `calendar:` is present
  budgetingMethod: Occurrences | Timeslices | RatioTimeslices
  objectives: # see objectives below for details
  alertPolicies: # see alert policies below for details
```

##### Notes (SLO)

- **indicator** optional, represents the Service Level Indicator (SLI),
  described in [SLI](#sli) section.
  One of `indicator` or `indicatorRef` must be given.
- **indicatorRef** optional, this is the name of Service Level Indicator (SLI).
  One of `indicator` or `indicatorRef` must be given.
- **timeWindow[ ]** optional, *TimeWindow* is a list but accepting only exactly one
  item, one of the rolling or calendar aligned time window:
  - Rolling time window. Duration should be provided in shorthand format
    e.g. 5m, 4w, 31d.
  - Calendar Aligned time window. Duration should be provided in shorthand format
    eg. 1d, 2M, 1Q, 366d.

- **description** *string* optional field, contains at most 1050 characters

- **budgetingMethod** *enum(Occurrences \| Timeslices \| RatioTimeslices)*, required field
  - Occurrences method uses a ratio of counts of good events to the total count of the events.
  - Timeslices method uses a ratio of good time slices to total time slices in a budgeting period.
  - RatioTimeslices method uses an average of all time slices' success ratios in a budgeting period.

- **objectives[ ]** *Threshold*, required field, described in [Objectives](#objectives)
  section. If `thresholdMetric` has been defined, only one Threshold can be defined.
  However if using `ratioMetric` then any number of Thresholds can be defined.

- **alertPolicies\[ \]** *AlertPolicy*, optional field.
  section. An alert policy can be defined inline or can refer to an [Alert Policies](#alertpolicy) object,
  in which case the following are required:
  - **alertPolicyRef** *string*: this is the name or path to the AlertPolicy

##### Objectives

Objectives are the thresholds for your SLOs. You can use objectives to define
the tolerance levels for your metrics.

```yaml
objectives:
  - displayName: string # optional
    op: lte | gte | lt | gt # conditional operator used to compare the SLI against the value. Only needed when using a thresholdMetric
    value: numeric # optional, value used to compare threshold metrics. Only needed when using a thresholdMetric
    target: numeric [0.0, 1.0) # budget target for given objective of the SLO, can't be used with targetPercent
    targetPercent: numeric [0.0, 100) # budget target for given objective of the SLO, can't be used with target
    timeSliceTarget: numeric (0.0, 1.0] # required only when budgetingMethod is set to TimeSlices
    timeSliceWindow: number | duration-shorthand # required only when budgetingMethod is set to TimeSlices or RatioTimeslices
```

Example:

```yaml
objectives:
  - displayName: Foo Total Errors
    target: 0.98
  - displayName: Bar Total Errors
    targetPercent: 99.99
```

###### Notes (Objectives)

- **op** *enum( lte | gte | lt | gt )*, operator used to compare the SLI against the value. Only needed when using a `thresholdMetric`

- **value** *numeric*, required field, used to compare values gathered from
  metric source. Only needed when using a `thresholdMetric`.

Either `target` or `targetPercent` must be used.

- **target** *numeric [0.0, 1.0)*, optional, but either this or `targetPercent` must
  be used. Budget target for a given objective of the SLO. A `target: 0.9995` is
  equivalent to `targetPercent: 99.95`.

- **targetPercent**: *numeric [0.0, 100)*, optional, but either this or `target` must
  be used. Budget target for a given objective of the SLO. A `targetPercent: 99.95`
  is equivalent to `target: 0.9995`.

- **timeSliceTarget** *numeric [0.0, 1.0]*, required only when budgeting
  method is set to TimeSlices

- **timeSliceWindow** *(numeric | duration-shorthand)*, required only when budgeting
  method is set to TimeSlices or RatioTimeslices. Denotes the size of a time slice for
  which data will be evaluated e.g. 5, 1m, 10m, 2h, 1d. Also ascertains the frequency
  at which to run the queries. Default interpretation of unit if specified as a number
  in minutes.

---

#### SLI

A service level indicator (SLI) represents how to read metrics from data sources.

> [Check work in progress for v2.](enhancements/v2alpha1.md#sli)

```yaml
apiVersion: openslo/v1
kind: SLI
metadata:
  name: string
  displayName: string # optional
spec:
  description: string # optional up to 1050 characters
  thresholdMetric: # either thresholdMetric or ratioMetric must be provided
    metricSource:
      metricSourceRef: string # optional, this field can be used to refer to DataSource object
      type: string # optional, this field describes predefined metric source type e.g. Prometheus, Datadog, etc.
      spec:
        # arbitrary chosen fields for every data source type to make it comfortable to use
        # anything that is valid YAML can be put here.
  ratioMetric: # either thresholdMetric or ratioMetric must be provided
    counter: true | false # true if the metric is a monotonically increasing counter,
                          # or false, if it is a single number that can arbitrarily go up or down
                          # ignored when using "raw"
    good: # the numerator, either "good" or "bad" must be provided if "total" is used
      metricSource:
        metricSourceRef: string # optional
        type: string # optional
        spec:
          # arbitrary chosen fields for every data source type to make it comfortable to use.
    bad: # the numerator, either "good" or "bad" must be provided if "total" is used
      metricSource:
        metricSourceRef: string # optional
        type: string # optional
        spec:
          # arbitrary chosen fields for every data source type to make it comfortable to use.
    total: # the denominator used with either "good" or "bad", either this or "raw" must be used
      metricSource:
        metricSourceRef: string # optional
        type: string # optional
        spec:
          # arbitrary chosen fields for every data source type to make it comfortable to use.

    rawType: success | failure # required with "raw", indicates how the stored ratio was calculated:
                               #  success – good/total
                               #  failure – bad/total
    raw: # the precomputed ratio stored as a metric, can't be used together with good/bad/total
      metricSource:
        metricSourceRef: string # optional
        type: string # optional
        spec:
          # arbitrary chosen fields for every data source type to make it comfortable to use.
```

##### Notes (SLI)

- **description** *string* optional field, contains at most 1050 characters

Either `ratioMetric` or `thresholdMetric` must be used.

- **thresholdMetric** *Metric*, represents the query used for
  gathering data from metric sources. Raw data is used to compare objectives
  (threshold) values.

- **ratioMetric** *Metric {good, total}, {bad, total} or raw*.

  - **counter** *enum(true \| false)*, specifies whether the metric is a monotonically
   increasing counter. Has no effect when using a `raw` query.

  - **good** represents the query used for gathering data from metric sources used
   as the numerator. Received data is used to compare objectives (threshold)
   values to find good values. If `bad` is defined then `good` must not be set.

  - **bad** represents the query used for gathering data from metric sources used
   as the numerator. Received data is used to compare objectives (threshold)
   values to find bad values. If `good` is defined then `bad` must not be set.

  - **total** represents the query used for gathering data from metric sources
   that is used as the denominator. Received data is used to compare objectives
   (threshold) values to find total number of metrics.

  - **rawType** *enum(success \| failure)*, required when using `raw`, specifies
   whether the ratios represented by the "raw" ratio metric are of successes or failures.
   Not to be used with `good` and `bad` as picking one of those determines the type of
   ratio.

  - **raw** represents the query used for gathering already precomputed ratios.
   The type of ratio (*success* or *failure*) is specified using `rawType`.

An example of an SLO where SLI is inlined:

```yaml
apiVersion: openslo/v1
kind: SLO
metadata:
  name: foo-slo
  displayName: Foo SLO
spec:
  service: foo
  indicator:
    metadata:
      name: foo-error
      displayName: Foo Error
    spec:
      ratioMetric:
        counter: true
        good:
          metricSource:
            metricSourceRef: datadog-datasource
            type: Datadog
            spec:
              query: sum:trace.http.request.hits.by_http_status{http.status_code:200}.as_count()
        total:
          metricSource:
            metricSourceRef: datadog-datasource
            type: Datadog
            spec:
              query: sum:trace.http.request.hits.by_http_status{*}.as_count()
  objectives:
    - displayName: Foo Total Errors
      target: 0.98
```

##### Ratio Metric

If a service level indicator has `ratioMetric` defined, the following maths can
be used to calculate the value of the SLI. Below we describe the advised formulas
for calculating the indicator value.

*Good-Total queries*
If the `good` and `total` queries are given then following formula can be used
to calculate the value:

```text
indicatorValue = good / total
```

If we have 99 good requests out of a total of 100 requests, the calculated value
for the indicator would be: `99 / 100  = 0.99`. This represents 99% on a 0-100 scale
using the formula `0.99 * 100 = 99`.

*Bad-Total queries*
If the `bad` and `total` queries are given then following formula can be used
to calculate the value:

```text
indicatorValue = ( total - bad ) / total
```

If we have 1 error out of a total of 100 requests, the calculated value for
the indicator would be: `(100 - 1) = 0.99`. This represents 99% on a 0-100 scale
using the formula `0.99 * 100 = 99`.

> 💡 **Note:** As you can see for both query combinations we end up with the same calculated
> value for the service level indicator.

The required `spec` key will be used to pass extraneous data to the data source. The goal of this approach
is to provide maximum flexibility when querying data from a particular source. In the following examples
we can see that this works fine for both simple and more complex cases.

An example of **ratioMetric**:

```yaml
ratioMetric:
  counter: true
  good:
    metricSource:
      type: Prometheus
      metricSourceRef: prometheus-datasource
      spec:
        query: sum(localhost_server_requests{code=~"2xx|3xx",host="*",instance="127.0.0.1:9090"})
  total:
    metricSource:
      type: Prometheus
      metricSourceRef: prometheus-datasource
      spec:
        query: localhost_server_requests{code="total",host="*",instance="127.0.0.1:9090"}
```

An example of **thresholdMetric**:

```yaml
thresholdMetric:
  metricSource:
    metricSourceRef: redshift-datasource
    spec:
      region: eu-central-1
      clusterId: metrics-cluster
      databaseName: metrics-db
      query: SELECT value, timestamp FROM metrics WHERE timestamp BETWEEN :date_from AND :date_to
```

Field `type` can be omitted because the type will be inferred from the DataSource when `metricSourceRef` is specified.

An example **thresholdMetric** that does not reference a defined DataSource:

```yaml
thresholdMetric:
  metricSource:
    type: Redshift
    spec:
      region: eu-central-1
      clusterId: metrics-cluster
      databaseName: metrics-db
      query: SELECT value, timestamp FROM metrics WHERE timestamp BETWEEN :date_from AND :date_to
      accessKeyID: accessKey
      secretAccessKey: secretAccessKey
 ```

 Field `type` can't be omitted because the reference to an existing DataSource is not specified.

---

#### AlertPolicy

An Alert Policy allows you to define the alert conditions for an SLO.

```yaml
apiVersion: openslo/v1
kind: AlertPolicy
metadata:
  name: string
  displayName: string # optional
spec:
  description: string # optional up to 1050 characters
  alertWhenNoData: boolean
  alertWhenResolved: boolean
  alertWhenBreaching: boolean
  conditions: # list of alert conditions
    - conditionRef: # required when alert condition is not inlined
  notificationTargets:
  - targetRef: # required when alert notification target is not inlined
```

##### Notes (AlertPolicy)

- **description** *string*, optional description about the alert policy, contains at most 1050 characters
- **alertWhenBreaching** *boolean*, `true`, `false`, whether the alert should be triggered
  when the condition is breaching
- **alertWhenResolved** *boolean*, `true`, `false`, whether the alert should be triggered
  when the condition is resolved
- **alertWhenNoData** *boolean*, `true`, `false`, whether the alert should be triggered
  when the condition indicates that no data is available
- **conditions\[ \]** *Alert Condition*, an array, (max of one condition), required field.
  A condition can be defined inline or can refer to external Alert condition defined in this case the following are required:
  - **conditionRef** *string*: this is the name or path the Alert condition
- **notificationTargets\[ \]** *Alert Notification Target*, required field.
  A condition can be defined inline or can refer to an [AlertNotificationTarget](#alertnotificationtarget)
  object, in which case the following are required:
  - **targetRef** *string*: this is the name or path to the AlertNotificationTarget

> 💡 **Note:** The `conditions` field is of the type `array` of *AlertCondition*
> but only allows one single condition to be defined.
> The use of an array is for future-proofing purposes.

An example of an Alert policy which refers to another Alert Condition:

```yaml
apiVersion: openslo/v1
kind: AlertPolicy
metadata:
  name: AlertPolicy
  displayName: Alert Policy
spec:
  description: Alert policy for cpu usage breaches, notifies on-call devops via email
  alertWhenBreaching: true
  alertWhenResolved: false
  conditions:
    - conditionRef: cpu-usage-breach
  notificationTargets:
    - targetRef: OnCallDevopsMailNotification
```

An example of an Alert Policy were the Alert Condition is inlined:

```yaml
apiVersion: openslo/v1
kind: AlertPolicy
metadata:
  name: AlertPolicy
  displayName: Alert Policy
spec:
  description: Alert policy for cpu usage breaches, notifies on-call devops via email
  alertWhenBreaching: true
  alertWhenResolved: false
  conditions:
    - kind: AlertCondition
      metadata:
        name: cpu-usage-breach
        displayName: CPU Usage breaching
      spec:
        description: SLO burn rate for cpu-usage-breach exceeds 2
        severity: page
        condition:
          kind: burnrate
          op: lte
          threshold: 2
          lookbackWindow: 1h
          alertAfter: 5m
  notificationTargets:
    - targetRef: OnCallDevopsMailNotification
```

---

#### AlertCondition

An Alert Condition allows you to define under which conditions an alert for an SLO
needs to be triggered.

```yaml
apiVersion: openslo/v1
kind: AlertCondition
metadata:
  name: string
  displayName: string # optional
spec:
  description: string # optional up to 1050 characters
  severity: string # required
  condition: # required
    kind: string
    op: enum
    threshold: number
    lookbackWindow: duration-shorthand
    alertAfter: duration-shorthand
```

##### Notes (AlertCondition)

- **description** *string*, optional description about the alert condition, contains at most 1050 characters
- **severity** *string*, required field describing the severity level of the alert (ex. "sev1", "page", etc.)
- **condition**, required field. Defines the conditions of the alert
  - **kind** *enum(burnrate)* the kind of alerting condition thats checked, defaults to `burnrate`

If the kind is `burnrate` the following fields are required:

- **op** *enum(lte | gte | lt | gt)*, required field, the conditional operator used to compare against the threshold
- **threshold** *number*, required field, the threshold that you want alert on
- **lookbackWindow** *duration-shorthand*, required field, the time-frame for which to calculate the threshold e.g. `5m`
- **alertAfter** *duration-shorthand*: required field, the duration the condition needs to be valid for before alerting, defaults to `0m`

If the alert condition is breaching, and the alert policy has `alertWhenBreaching` set to `true`
the alert will be triggered

If the alert condition is resolved, and the alert policy has `alertWhenResolved` set to `true`
the alert will be triggered

If the *service level objective* associated with the alert condition returns
no value for the burn rate, for example, due to the service level indicators
missing data (e.g. no time series being returned) and the `alertWhenNoData`
is set  to `true` the alert will be triggered.

> 💡 **Note:** The `alertWhenBreaching` and `alertWhenResolved`, `alertWhenNoData` can be combined,
> if you want an alert to trigger whenever at least one of these conditions is true.

---

An example of an alert condition:

```yaml
apiVersion: openslo/v1
kind: AlertCondition
metadata:
  name: cpu-usage-breach
  displayName: CPU usage breach
spec:
  description: If the CPU usage is too high for given period then it should alert
  severity: page
  condition:
    kind: burnrate
    op: lte
    threshold: 2
    lookbackWindow: 1h
    alertAfter: 5m
```

---

#### AlertNotificationTarget

An Alert Notification Target defines the possible targets where alert notifications
should be delivered to. For example, this can be a web-hook, Slack or any other
custom target.

```yaml
apiVersion: openslo/v1
kind: AlertNotificationTarget
metadata:
  name: string
  displayName: string # optional, human readable name
spec:
  target: # required
  description: # optional
```

An example Alert Notification Target:

```yaml
apiVersion: openslo/v1
kind: AlertNotificationTarget
metadata:
  name: OnCallDevopsMailNotification
spec:
  description: Notifies by a mail message to the on-call devops mailing group
  target: email
```

Alternatively, a similar notification target can be defined for Slack like this:

```yaml
apiVersion: openslo/v1
kind: AlertNotificationTarget
metadata:
  name: OnCallDevopsSlackNotification
spec:
  description: "Sends P1 alert notifications to the slack channel"
  target: slack
```

##### Notes (AlertNotificationTarget)

- **target** *string*, describes the target of the notification, e.g. Slack, email, web-hook, Opsgenie etc
- **description** *string*, optional description about the notification target, contains at most 1050 characters

> 💡 **Note:** The way the alert notification targets are is an implementation detail of the
> system that consumes the OpenSLO specification.
>
> For example, if the OpenSLO is consumed by a solution that generates Prometheus recording rules,
> and alerts, you can imagine that the name of the alert notification gets passed as a label
> to Alertmanager which then can be routed accordingly based on this label.

---

#### Service

A Service is a high-level grouping of SLO. It may be defined before creating SLO to be able to refer to it in SLO's `spec.service`.
Multiple SLOs can refer to the same Service.

```yaml
apiVersion: openslo/v1
kind: Service
metadata:
  name: string
  displayName: string # optional
spec:
  description: string # optional up to 1050 characters
```

[rfc1123-names]: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#dns-label-names
