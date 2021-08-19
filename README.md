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
spec:
  description: string # optional
  service: [service name] # name of the service to associate this SLO with
  indicator: # represents the Service Level Indicator (SLI)
    thresholdMetric: # represents the metric used to inform the Service Level Object in the objectives stanza
      source: string # data source for the metric
      queryType: string # a name for the type of query to run on the data source
      query: string # the query to run to return the metric
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
  alertPolicies: # see alert policies below details
```

##### Notes (SLO)

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

- **alertPolicies\[ \]** *AlertPolicy*, optional field, described in [Alert Policies](#alertpolicies)
  section

##### Objectives

Objectives are the thresholds for your SLOs. You can use objectives to define
the tolerance levels for your metrics

```yaml
objectives:
  - displayName: string # optional
    op: lte | gte | lt | gt # conditional operator used to compare the SLI against the value. Only needed when using a thresholdMetric
    value: numeric # value used to compare metrics values. All objectives of the SLO need to have a unique value.
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
        total: # the denominator
          source: string # data source for the "total" denominator
          queryType: string # a name for the type of query to run on the data source
          query: string # the query to run to return the denominator
```

Example:

```yaml
objectives:
  - displayName: Foo Total Errors
    value:  1
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

##### Notes (Objectives)

- **objectives\[ \]** *Threshold*, required field. If `thresholdMetric` has
  been defined, only one Threshold can be defined. However if using `ratioMetric`
  then any number of Thresholds can be defined.

- **op** *enum(lte | gte | lt | gt)*, operator used to compare the SLI against
  the value. Only needed when using `thresholdMetric`

- **value numeric**, required field, used to compare values gathered from
  metric source

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

#### Alert Policy

An Alert Policy allows you to define the alert conditions for a SLO.

```yaml
apiVersion: openslo/v1alpha
kind: AlertPolicy
metadata:
  name: string
  displayName: string # optional
spec:
  description: string # optional
  conditions: # list of alert conditions
    - conditionRef: # required when alert condition is not inlined
  notificationTargets:
  - targetRef: # required when alert notification target is not inlined

```

#### Notes (Alert Policy)

- **conditions\[ \]** *Alert Condition*, required field.
  A condition can be inline defined or can refer to external Alert condition defined
  in this case the following are required:
  - **conditionRef** *string*: this is the name or path the Alert condition
- **notificationTargets\[ \]** *Alert Notification Target*, required field.
  A condition can be inline defined or can refer to external Alert Notification Target
  defined in this case the following are required:
  - **targetRef** *string*: this is the name or path the Alert Notification Target

An example of a Alert policy which refers to another Alert Condition:

```yaml
apiVersion: openslo/v1alpha
kind: AlertPolicy
metadata:
  name: AlertPolicy
  displayName: Alert Policy
spec:
  description: Alert policy for cpu usage breaches, notifies on-call devops via email
  conditions:
    - conditionRef: cpu-usage-breach
  notificationTargets:
    - targetRef: OnCallDevopsMailNotification

```

---

#### Alert Condition

An Alert Condition allows you to define in which conditions a alert of SLO
needs to be triggered.

```yaml
apiVersion: openslo/v1alpha
kind: AlertCondition
metadata:
  name: string
  displayName: string # optional
spec:
  description: # optional
  severity: string # required (ticket or page)
  condition: # optional
    kind: string
    threshold: number
    lookbackWindow: number
    controlLookbackWindow: number # optional
    alertAfter: number
```

#### Notes (Alert Condition)

- **severity** *enum(ticket, page)*, required field. The severity level of the alert
- **condition**, required field. Defines the conditions of the alert
  - **kind** *enum(burnrate, guard, custom)* the kind of alerting condition thats checked, defaults to `burnrate`
  If the kind is `burnrate` the following fields are required:
  - **threshold** *number*, required field, the threshold that you want alert on
  - **lookbackWindow** *number*, required field, the time-frame for which to calculate the threshold
  - **controlLookbackWindow** *number*, optional field
  - **alertAfter** *number*: required field, the duration the condition needs to be valid, defaults `0m`
  If the kind is `custom` the following fields are required:
  - **threshold** *number*, required field, the threshold that you want alert on
  - **comparison** *enum(lt, lte, gt, gte, eq)*, optional field, defines
  how the threshold should be compared to meet the threshold, defaults to `gt`
  If the kind is `guard` the following fields are required, can  be used to
  define quality conditions:
  - **threshold** *number*, required field, the threshold that you want alert on
  - **criteriaType** *enum(pass, warning)*, required field, defines the criteria type of the condition
  - **weight** *number*, required field, the weight or importance of the condition

---

An example of a alert condition is the following:

```yaml
apiVersion: openslo/v1alpha
kind: AlertCondition
metadata:
  name: cpu-usage-breach
  displayName: CPU usage breach
spec:
  description: If the CPU usage is too high for given period then it should alert
  severity: page
  condition:
    kind: burnrate
    threshold: 0.9
    lookbackWindow: 1h
    alertAfter: 5m
```

---

#### Alert Notification Target

An Alert Notification Target defines the possible targets where alert notifications
should be delivered to. For example, this can be a web-hook, Slack or similar

```yaml
apiVersion: openslo/v1alpha
kind: AlertNotificationTarget
metadata:
  name:
  displayName: string # optional, human readable name
spec:
  target: # required
  description: # optional
  parameters: # required
    - name: string # required
      description: string # optional
      type: string # required, string, secret, number, url
      value: string # optional
      defaultValue: string # optional
```

An example of the Alert Notification Target can be:

```yaml
apiVersion: openslo/v1alpha
kind: AlertNotificationTarget
metadata:
  name: OnCallDevopsMailNotification
spec:
  description: Notifies by a mail message to the on-call devops mailing group
  target: email
  parameters:
    - name: emailAddress
      type: email
      value: on-call-devops@openslo.org
```

Alternatively, a similar notification target can be defined for Slack in this example

```yaml
apiVersion: openslo/v1alpha
kind: AlertNotificationTarget
metadata:
  name: OnCallDevopsSlackNotification
spec:
  description: "Sends P1 alert notifications to the #alerts channel"
  target: slack
  parameters:
    - name: channel
      type: string
      value: alerts
    - name: webhook
      type: url
      value: https://hooks.slack.com/services/XXX1XXXXX/XXXXXXX8X/acdifghf7Klxy8xZCEDXyOk5I
    - name: template
      type: string
      value: |-
        A priority one alert has been triggered for the service $service_name for SLO $slo_name,
        please extinguish the fire, thank you!
        Don't forget to schedule a meeting the company board via Executive Secretary
```

##### Notes (Alert Notification Target)

- **target** *string*, describes the target of the notification, e.g. Slack, email, web-hook, Opsgenie etc
- **description** *string*, optional description about the notification target
- **parameters**, defined all the available parameters for the target, e.g. username, password, Slack web-hook url
  - **name** *string*, unique name of the parameter
  - **description** *string* (optional), human readable description of the parameter
  - **type** *enum(string, secret, number, email, url)* the accepted type of
  value of the parameter, if not given it fallbacks to `string`
  - **value** *string*, a hard-coded value of the parameter name
    (e.g. web-hook url)
  - **defaultValue** *string*, a fallback or default value for the parameter

**Note**: The OpenSLO comes with a few so called template fields which can be used as part of the specification to
get the name of the relevant service, SLO name, SLO description and

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
