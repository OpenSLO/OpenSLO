# OpenSLO

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
interfacing with SLOs to allow for a common, giving a set vendor-agnostic
solution to defining and tracking SLOs. Platform specific implimentation details
are purposefully excluded from the scope of this specification.

## Specification

### Goals

- Compliance with the Kubernetes YAML format.
- Vendor agnostic
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
  displayName: string   # optional
spec:
```

##### Notes

- **kind** *string* - required, either [SLO](#slo) or [Service](#service)
- **name:** *string* - required field, convention for naming object from
  [DNS RFC1123](https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names)
  `name` should:
  - contain at most 63 characters
  - contain only lowercase alphanumeric characters or `-`
  - start with an alphanumeric character
  - end with an alphanumeric character

---

#### SLO

A service level objective (SLO) is a target value or range of values for
a service level that is measured by a service level indicator (SLI).

```yaml
apiVersion: n9/v1alpha
kind: SLO
metadata:
  name: string
  displayName: string #optional
  project: string
spec:
  description: string #optional
  service: [service name] #name of the service to associate this SLO with
  timeWindows:
    # exactly one item, one of possible rolling time window or calendar aligned
    # rolling time window
    - unit: Day | Hour | Minute
      count: numeric
      isRolling: true
    # or
    # calendar aligned time window
    - unit: Year | Quarter | Month | Week | Day
      count: numeric # count of time units for example count: 7 and unit: Day means 7 days window
      calendar:
        startTime: 2020-01-21 12:30:00 # date with time in 24h format, format without time zone
        timeZone: America/New_York # name as in IANA Time Zone Database
      isRolling: false # false or or not defined
  budgetingMethod: Occurrences | Timeslices
  alertPolicies:
    - string # The name of the alert policy associated with this SLO
             # (alert policy only from the same project as SLO).
             # Allow to have 0 to 5 AlertPolicy per SLO.
  objectives:  # see objectives below for details
```

##### Notes

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
      # isRolling: false #for calendar aligned set false value or not set
      ```

- **description** *string* optional field, contains at most 1050 characters

- **alertPolicies\[ \]** *AlertPolicy*, optional field,
  **(count 0 up to 5 items)**, is a list of AlertPolicies names which can
    trigger alerts for given SLO when all conditions are met. AlertPolicies have
    to be created in the same project as SLO.

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

[TODO] This only covers ratio right now.  Since thresholds require the rawMetric,
which is part of indicator stanza, we need to find a way of handling that.

```yaml
objectives:
  - displayName: string # optional
    value: numeric # value used to compare metrics values. All objectives of the SLO need to have unique value.
    target: numeric [0.0, 1.0) #budget target for given objective of the SLO
    timeSliceTarget: numeric (0.0, 1.0] #required only when budgetingMethod is set to TimeSlices
    # countMetrics {good, total} should be defined only if raw metric is not set.
    # If rawMetric is defined on SLO level, raw data received from a metric source is compared with objective value.
    # Count metrics good and total have to contain the same source type configuration (for example for prometheus).
    countMetrics:
        incremental: true | false #todo: add description
        good: # the numerator
          source: # data source for the "good" numerator
          queryType: string # a name for the type of query to run on the data source
          query: string # the query to run to return the numerator
        total: # the denominator
          source: # data source for the "good" numerator
          queryType: string # a name for the type of query to run on the data source
          query: string # the query to run to return the numerator
```

Example:

```yaml
objectives:
  - displayName: Foo Total Errors
    value:  1
    target: 0.98
    countMetrics:
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

---

#### Service

A Service is a high-level grouping of SLO.

```yaml
apiVersion: openslo/v1alpha
kind: Service
metadata:
  name: string
  displayName: string  # optional
spec:
  description: string  # optional up to 1050 characters
```
