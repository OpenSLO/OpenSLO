# ![ OpenSLO ][image-1]

## Table of Contents

- [OpenSLO](#openslo)
  - [Introduction](#introduction)
  - [Specification](#specification)
    - [Goals](#goals)
    - [Object Types](#object-types)
      - [General Schema](#general-schema)
        - [Notes](#notes)
      - [DataSource](#datasource)
        - [Notes](#notes-datasource)
      - [SLO](#slo)
        - [Notes](#notes-slo)
        - [Objectives](#objectives)
      - [SLI](#sli)
        - [Notes](#notes-sli)
        - [Ratio Metric](#ratio-metric)
      - [Alerting](#alert-policy)
        - [Notes](#notes-alert-policy)
        - [Alert Conditions](#alert-condition)
        - [Notification Targets](#alert-notification-target)
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
apiVersion: openslo/v0.1.0-beta
kind: SLO | Service
metadata:
  name: string
  displayName: string # optional
spec:
```

##### Notes

- **kind** *string* - required, either [SLO][12] or [Service][13]
- **metadata.name:** *string* - required field, convention for naming object from
  [DNS RFC1123][14]
  `name` should:

  - contain at most 63 characters
  - contain only lowercase alphanumeric characters or `-`
  - start with an alphanumeric character
  - end with an alphanumeric character

---

#### DataSource

A DataSource represents connection details with a particular metric source.

```yaml
apiVersion: openslo/v0.1.0-beta
kind: DataSource
metadata:
  name: string
  displayName: string # optional
spec:
  type: string # predefined type e.g. Prometheus, Datadog, etc.
  connectionDetails:
    # fields used for creating a connection with particular datasource e.g. AccessKeys, SecretKeys, etc.
    # everything that is valid YAML can be put here

```

##### Notes (DataSource)

DataSource enables reusing one source between many SLOs and moving
connection specific details e.g. authentication away from SLO.

An example of the DataSource kind can be:

```yaml
apiVersion: openslo/v0.1.0-beta
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

A service level objective (SLO) is a target value or range of values for
a service level that is described by a service level indicator (SLI).

```yaml
apiVersion: openslo/v0.1.0-beta
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
  indicator: # see SLI below for details
  indicatorRef: string # name of the SLI. Required if indicator is not given.
  timeWindow:
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
  alertPolicies: # see alert policies below details
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
- **indicator** optional, represents the Service Level Indicator (SLI),
  described in [SLI](#sli) section.
- **indicatorRef** optional, this is the name of Service Level Indicator (SLI).
  One of `indicator` or `indicatorRef` must be given.
- **timeWindow[ ]** *TimeWindow* is a list but accepting only exactly one
  item, one of the rolling or calendar aligned time window:

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
  - Occurrences method uses a ratio of counts of good events and total count of the event.
  - Timeslices method uses a ratio of good time slices vs. total time slices in a budgeting period.

- **objectives[ ]** *Threshold*, required field, described in [Objectives][17]
  section. If `thresholdMetric` has been defined, only one Threshold can be defined.
  However if using `ratioMetric` then any number of Thresholds can be defined.

- **alertPolicies\[ \]** *AlertPolicy*, optional field, described in [Alert Policies](#alert-policy)
  section

##### Objectives

Objectives are the thresholds for your SLOs. You can use objectives to define
the tolerance levels for your metrics.

```yaml
objectives:
  - displayName: string # optional
    op: lte | gte | lt | gt # conditional operator used to compare the SLI against the value. Only needed when using a thresholdMetric
    value: numeric # optional, value used to compare threshold metrics. Only needed when using a thresholdMetric
    target: numeric [0.0, 1.0) # budget target for given objective of the SLO
    timeSliceTarget: numeric (0.0, 1.0] # required only when budgetingMethod is set to TimeSlices
    timeSliceWindow: number | duration-shorthand # required only when budgetingMethod is set to TimeSlices
```

Example:

```yaml
objectives:
  - displayName: Foo Total Errors
    target: 0.98
```

##### Notes (Objectives)

- **op** *enum(lte | gte | lt | gt)*, operator used to compare the SLI against
  the value. Only needed when using `thresholdMetric`

- **value numeric**, required field, used to compare values gathered from
  metric source. Only needed when using a `thresholdMetric`.

- **target numeric** *[0.0, 1.0)*, required, budget target for given objective
  of the SLO

- **timeSliceTarget** *numeric* *[0.0, 1.0]*, required only when budgeting
  method is set to TimeSlices

- **timeSliceWindow** *(numeric | duration-shorthand)*, required only when budgeting
  method is set to TimeSlices. Denotes the size of timeslice for which data will be
  evaluated e.g. 5, 1m, 10m, 2h, 1d. Also ascertains the frequency at which to run the
  queries. Default interpretation of unit if specified as a number is minutes.

---

#### SLI

A service level indicator (SLI) represents how to gather data from metric sources.

```yaml
apiVersion: openslo/v0.1.0-beta
kind: SLI
metadata:
  name: string
  displayName: string # optional
spec:
  thresholdMetric: # either thresholdMetric or ratioMetric should be provided
    metricSource:
      metricSourceRef: string # optional, this field can be used to refer to DataSource object
      type: string # optional, this field describes predefined metric source type e.g. Prometheus, Datadog, etc.
      spec:
        # arbitrary chosen fields for every data source type to make it comfortable to use
        # anything that is valid YAML can be put here.
  ratioMetric: # either thresholdMetric or ratioMetric should be provided
    counter: true | false # true if the metric is a monotonically increasing counter,
                          # or false, if it is a single number that can arbitrarily go up or down
    good: # the numerator
      metricSource:
        metricSourceRef: string # optional
        type: string # optional
        spec:
          # arbitrary chosen fields for every data source type to make it comfortable to use.
    bad: # the numerator, required when "good" is not set
      metricSource:
        metricSourceRef: string # optional
        type: string # optional
        spec:
          # arbitrary chosen fields for every data source type to make it comfortable to use.
    total: # the denominator
      metricSource:
        metricSourceRef: string # optional
        type: string # optional
        spec:
          # arbitrary chosen fields for every data source type to make it comfortable to use.
```

##### Notes(SLI)

When filling data in `metricSource`

Either `ratioMetric` or `thresholdMetric` should be set.

- **thresholdMetric** *Metric*, represents the query used for
  gathering data from metric sources. Raw data is used to compare objectives
  (threshold) values.

- **ratioMetric** *Metric {Good, Total} or {Bad, Total}*.

  - *Good* represents the query used for gathering data from metric sources used
   as the numerator. Received data is used to compare objectives (threshold)
   values to find good values. If `Bad` is defined then `Good` should not be set.

  - *Bad* represents the query used for gathering data from metric sources used
   as the numerator. Received data is used to compare objectives (threshold)
   values to find bad values. If `Good` is defined then `Bad` should not be set.

  - *Total* represents the query used for gathering data from metric sources
   that is used as the denominator. Received data is used to compare objectives
   (threshold) values to find total number of metrics.

An example of an SLO where SLI is inlined:

```yaml
apiVersion: openslo/v0.1.0-beta
kind: SLO
metadata:
  name: foo-slo
  displayName: Foo SLO
spec:
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
the indicator would be: `(100 - 1 )  = 0.99`. This represents 99% on a 0-100 scale
using the formula `0.99 * 100 = 99`.

> 💡 **Note:** : As you can see for both query combinations we end up with the same calculated
> value for the service level indicator.

The required `spec` key will be used to pass extraneous data to the data source. Goal of this approach
is to give the maximum flexibility when querying data from a particular source. In the following examples
we can see that it works fine for simple and more complex cases.

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

An example of **thresholdMetric** without specifying DataSource name and kind:

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

#### Alert Policy

An Alert Policy allows you to define the alert conditions for a SLO.

```yaml
apiVersion: openslo/v0.1.0-beta
kind: AlertPolicy
metadata:
  name: string
  displayName: string # optional
spec:
  description: string # optional
  alertWhenNoData: boolean
  alertWhenResolved: boolean
  alertWhenBreaching: boolean
  conditions: # list of alert conditions
    - conditionRef: # required when alert condition is not inlined
  notificationTargets:
  - targetRef: # required when alert notification target is not inlined
```

#### Notes (Alert Policy)

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
  A condition can be inline defined or can refer to external Alert Notification Target
  defined in this case the following are required:
  - **targetRef** *string*: this is the name or path the Alert Notification Target

> 💡 **Note:**: The `conditions`-field is of the type `array` of *AlertCondition*
> but only allows one single condition to be defined.
> The use of an array is for future-proofing purposes.

An example of a Alert policy which refers to another Alert Condition:

```yaml
apiVersion: openslo/v0.1.0-beta
kind: AlertPolicy
metadata:
  name: AlertPolicy
  displayName: Alert Policy
spec:
  description: Alert policy for cpu usage breaches, notifies on-call devops via email
  alertWhenBreaching: true
  alertWhenResolved: false
  conditions:
    - operator: and
      conditionRef: cpu-usage-breach
  notificationTargets:
    - targetRef: OnCallDevopsMailNotification
```

An example of a Alert Policy were the Alert Condition is inlined:

```yaml
apiVersion: openslo/v0.1.0-beta
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
          threshold: 2
          lookbackWindow: 1h
          alertAfter: 5m
  notificationTargets:
    - targetRef: OnCallDevopsMailNotification
```

---

#### Alert Condition

An Alert Condition allows you to define in which conditions a alert of SLO
needs to be triggered.

```yaml
apiVersion: openslo/v0.1.0-beta
kind: AlertCondition
metadata:
  name: string
  displayName: string # optional
spec:
  description: string # optional
  severity: string # required
  condition: # optional
    kind: string
    threshold: number
    lookbackWindow: number
    alertAfter: number
```

#### Notes (Alert Condition)

- **severity** *string* , required field describing the severity level of the alert (ex. "sev1", "page", etc.)
- **condition**, required field. Defines the conditions of the alert
  - **kind** *enum(burnrate)* the kind of alerting condition thats checked, defaults to `burnrate`

If the kind is `burnrate` the following fields are required:

- **threshold** *number*, required field, the threshold that you want alert on
- **lookbackWindow** *number*, required field, the time-frame for which to calculate the threshold
- **alertAfter** *number*: required field, the duration the condition needs to be valid, defaults `0m`

If the alert condition is breaching, and the alert policy has `alertWhenBreaching` set to `true`
the alert will be triggered

If the alert condition is resolved, and the alert policy has `alertWhenResolved` set to `true`
the alert will be triggered

If the *service level objective* associated with the alert condition returns
no value for the burn rate, for example, due to the service level indicators
missing data (e.g. no time series being returned) and the `alertWhenNoData`
is set  to `true` the alert will be triggered.

> 💡 **Note:**: The `alertWhenBreaching` and `alertWhenResolved`, `alertWhenNoData` can be combined,
> if you want an alert to trigger when in all circumstances or for each separately.

---

An example of a alert condition is the following:

```yaml
apiVersion: openslo/v0.1.0-beta
kind: AlertCondition
metadata:
  name: cpu-usage-breach
  displayName: CPU usage breach
spec:
  description: If the CPU usage is too high for given period then it should alert
  severity: page
  condition:
    kind: burnrate
    threshold: 2
    lookbackWindow: 1h
    alertAfter: 5m
```

---

#### Alert Notification Target

An Alert Notification Target defines the possible targets where alert notifications
should be delivered to. For example, this can be a web-hook, Slack or any other
custom target.

```yaml
apiVersion: openslo/v0.1.0-beta
kind: AlertNotificationTarget
metadata:
  name: string
  displayName: string # optional, human readable name
spec:
  target: # required
  description: # optional
```

An example of the Alert Notification Target can be:

```yaml
apiVersion: openslo/v0.1.0-beta
kind: AlertNotificationTarget
metadata:
  name: OnCallDevopsMailNotification
spec:
  description: Notifies by a mail message to the on-call devops mailing group
  target: email
```

Alternatively, a similar notification target can be defined for Slack in this example

```yaml
apiVersion: openslo/v0.1.0-beta
kind: AlertNotificationTarget
metadata:
  name: OnCallDevopsSlackNotification
spec:
  description: "Sends P1 alert notifications to the slack channel"
  target: slack
```

##### Notes (Alert Notification Target)

- **name** *string*, required, the name of the notification target
- **metadata.labels:** *map[string]string|string[]* - optional field `key` <> `value`
  - the `key` segment is required and must contain at most 63 characters beginning and ending
     with an alphanumeric character `[a-z0-9A-Z]` with dashes `-`, underscores `_`, dots `.`
     and alphanumerics between.
  - the `value` of `key` segment can be a string or an array of strings
- **target** *string*, describes the target of the notification, e.g. Slack, email, web-hook, Opsgenie etc
- **description** *string*, optional description about the notification target

> 💡 **Note:**: The way the alert notification targets are is an implementation detail of the
> system that consumes the OpenSLO specification.

For example, if the OpenSLO is consumed by a solution that generates Prometheus recording rules,
and alerts, you can imagine that the name of the alert notification gets passed as label
to Alertmanager which then can be routed accordingly based on this label.

---

#### Service

A Service is a high-level grouping of SLO.

```yaml
apiVersion: openslo/v0.1.0-beta
kind: Service
metadata:
  name: string
  displayName: string # optional
spec:
  description: string # optional up to 1050 characters
```

[1]: #openslo
[2]: #introduction
[3]: #specification
[4]: #goals
[5]: #object-types
[6]: #general-schema
[7]: #notes
[8]: #slo
[9]: #notes-1
[10]: #objectives
[11]: #service
[12]: #slo
[13]: #service
[14]: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
[15]: #objectives
[16]: #objectives
[17]: #objectives

[image-1]: images/openslo.png
