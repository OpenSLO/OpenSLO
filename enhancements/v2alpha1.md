# v2alpha1

This is the place for refining ideas for new versions of OpenSLO spec. It's not supposed to be stable, this is a living document

## [General Schema](https://github.com/OpenSLO/OpenSLO?tab=readme-ov-file#general-schema)

Since the goal of the OpenSLO spec is to be compatible with Kubernetes, we should make a couple of fixes
in the specification to reach that goal.

1. Use `openslo.com/<version>` instead of `openslo/<version>`
2. Remove `displayName` from the [ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#objectmeta-v1-meta)
3. Adhere to [metadata.labels](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/) and
   [metadata.annotations](https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/) Kubernetes standards.

## [DataSource](../README.md#datasource)

**Rationale:** Simplify syntax. Avoid being needlessly verbose without sacrificing flexibility and readability.

```yaml
apiVersion: openslo.com/v2alpha1
kind: DataSource
metadata:
  name: string
  labels: object # optional
spec:
  description: string # optional up to 1050 characters
  <<dataSourceName>>: # e.g. cloudWatch, datadog, prometheus (arbitrary chosen, implementor decision)
  # fields used for creating a connection with particular datasource e.g. AccessKeys, SecretKeys, etc.
  # everything that is valid YAML can be put here
```

An example of the DataSource kind can be:

```yaml
apiVersion: openslo.com/v2alpha1
kind: DataSource
metadata:
  name: string
spec:
  cloudWatch:
    accessKeyID: accessKey
    secretAccessKey: secretAccessKey
```

## [SLO](../README.md#slo)

**Rationale:** Make names more straightforward and aligned with others. Change field indicator to `sli` and `indicatorRef` to `sliRef`
it tells which kind of object should be referred there. This change should also apply to each objective.

```yaml
apiVersion: openslo.com/v2alpha1
kind: SLO
metadata:
  name: string
  labels: object # optional
spec:
  description: string # optional up to 1050 characters
  service: string # name of the service to associate this SLO with, may refer (depends on implementation) to existing object Kind: Service
  sli: # see SLI below for details
  sliRef: string # name of the SLI, required if indicator is not given and you want to reference to existing SLI
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

### Objectives

Objectives are the thresholds for your SLOs. You can use objectives to define
the tolerance levels for your metrics.

**Note:** While in the previous version only one objective was allowed for Threshold metrics, with `v2alpha1` we'll allow any number of objectives.

```yaml
objectives:
  - displayName: string # optional
    labels: object # optional
    op: lte | gte | lt | gt # conditional operator used to compare the SLI against the value. Only needed when using a thresholdMetric
    value: numeric # optional, value used to compare threshold metrics. Only needed when using a thresholdMetric
    target: numeric [0.0, 1.0) # budget target for given objective of the SLO, can't be used with targetPercent
    targetPercent: numeric [0.0, 100) # budget target for given objective of the SLO, can't be used with target
    timeSliceTarget: numeric (0.0, 1.0] # required only when budgetingMethod is set to TimeSlices
    timeSliceWindow: number | duration-shorthand # required only when budgetingMethod is set to TimeSlices or RatioTimeslices
    sli: # required only when creating composite SLO, see SLI below for more details
    sliRef: string # required only when creating composite SLO, required if sli is not given.
    compositeWeight: numeric (0.0, inf+] # optional, supported only when declaring multiple objectives, default value 1.
```

## [SLI](../README.md#sli)

**Rationale:** Get rid of `metricSource` (reduce the level of indentation), and use the new syntax of `DataSource` directly.

```yaml
apiVersion: openslo.com/v2alpha1
kind: SLI
metadata:
  name: string
  labels: object # optional
spec:
  description: string # optional up to 1050 characters
  thresholdMetric: # either thresholdMetric or ratioMetric must be provided
    # either dataSourceRef or <<dataSourceName>> must be provided
    dataSourceRef: string # refer to already defined DataSource object
    <<dataSourceName>>: # inline whole DataSource e.g. cloudWatch, datadog, prometheus (arbitrary chosen, implementor decision)
  # fields used for creating a connection with particular datasource e.g. AccessKeys, SecretKeys, etc.
  # everything that is valid YAML can be put here
    spec:
      # arbitrary chosen fields for every DataSource type to make it comfortable to use
      # anything that is valid YAML can be put here.
  ratioMetric: # either thresholdMetric or ratioMetric must be provided
    counter:
      true | false # true if the metric is a monotonically increasing counter,
      # or false, if it is a single number that can arbitrarily go up or down
      # ignored when using "raw"
    good: # the numerator, either "good" or "bad" must be provided if "total" is used
      # either dataSourceRef or <<dataSourceName>> must be provided
      dataSourceRef: string # refer to already defined DataSource object
      <<dataSourceName>>: # inline whole DataSource e.g. cloudWatch, datadog, prometheus (arbitrary chosen, implementor decision)
   # fields used for creating a connection with particular datasource e.g. AccessKeys, SecretKeys, etc.
   # everything that is valid YAML can be put here
      spec:
        # arbitrary chosen fields for every DataSource type to make it comfortable to use
        # anything that is valid YAML can be put here.
    bad: # the numerator, either "good" or "bad" must be provided if "total" is used
      # either dataSourceRef or <<dataSourceName>> must be provided
      dataSourceRef: string # refer to already defined DataSource object
      <<dataSourceName>>: # inline whole DataSource e.g. cloudWatch, datadog, prometheus (arbitrary chosen, implementor decision)
   # fields used for creating a connection with particular datasource e.g. AccessKeys, SecretKeys, etc.
   # everything that is valid YAML can be put here
      spec:
        # arbitrary chosen fields for every DataSource type to make it comfortable to use
        # anything that is valid YAML can be put here
    total: # the denominator used with either "good" or "bad", either this or "raw" must be used
      # either dataSourceRef or <<dataSourceName>> must be provided
      dataSourceRef: string # refer to already defined DataSource object
      <<dataSourceName>>: # inline whole DataSource e.g. cloudWatch, datadog, prometheus (arbitrary chosen, implementor decision)
   # fields used for creating a connection with particular datasource e.g. AccessKeys, SecretKeys, etc.
   # everything that is valid YAML can be put here
      spec:
        # arbitrary chosen fields for every DataSource type to make it comfortable to use
        # anything that is valid YAML can be put here

    rawType:
      success | failure # required with "raw", indicates how the stored ratio was calculated:
      #  success – good/total
      #  failure – bad/total
    raw: # the precomputed ratio stored as a metric, can't be used together with good/bad/total
      # either dataSourceRef or <<dataSourceName>> must be provided
      dataSourceRef: string # refer to already defined DataSource object
      <<dataSourceName>>:# inline whole DataSource e.g. cloudWatch, datadog, prometheus (arbitrary chosen, implementor decision)
        # fields used for creating a connection with particular datasource e.g. AccessKeys, SecretKeys, etc.
        # everything that is valid YAML can be put here
      spec:
        # arbitrary chosen fields for every DataSource type to make it comfortable to use
        # anything that is valid YAML can be put here
```

An example of an SLO where SLI is inlined:

```yaml
apiVersion: openslo.com/v2alpha1
kind: SLO
metadata:
  name: foo-slo
spec:
  service: foo
  sli:
    metadata:
      name: foo-error
    spec:
      ratioMetric:
        counter: true
        good:
          dataSourceRef: datadog-datasource
          spec:
            query: sum:trace.http.request.hits.by_http_status{http.status_code:200}.as_count()
        total:
          dataSourceRef: datadog-datasource
          spec:
            query: sum:trace.http.request.hits.by_http_status{*}.as_count()
  objectives:
    - displayName: Foo Total Errors
      target: 0.98
```

An example of **ratioMetric**:

```yaml
ratioMetric:
  counter: true
  good:
    dataSourceRef: prometheus-datasource
    spec:
      query: sum(localhost_server_requests{code=~"2xx|3xx",host="*",instance="127.0.0.1:9090"})
  total:
    dataSourceRef: prometheus-datasource
    spec:
      query: localhost_server_requests{code="total",host="*",instance="127.0.0.1:9090"}
```

An example of **thresholdMetric**:

```yaml
thresholdMetric:
  dataSourceRef: redshift-datasource
  spec:
    region: eu-central-1
    clusterId: metrics-cluster
    databaseName: metrics-db
    query: SELECT value, timestamp FROM metrics WHERE timestamp BETWEEN :date_from AND :date_to
```

An example **thresholdMetric** that does not reference a defined DataSource (it has `DataSource` inlined):

```yaml
thresholdMetric:
  redshift:
    accessKeyID: accessKey
    secretAccessKey: secretAccessKey
  spec:
    region: eu-central-1
    clusterId: metrics-cluster
    databaseName: metrics-db
    query: SELECT value, timestamp FROM metrics WHERE timestamp BETWEEN :date_from AND :date_to
```
