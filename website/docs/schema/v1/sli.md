# SLI

{{ generate_object_description("openslo/v1", "SLI") }}

## Examples

=== "Ratio Metric"

    ```yaml
    apiVersion: openslo/v1
    kind: SLI
    metadata:
      name: api-success-rate
      displayName: API Success Rate
    spec:
      description: Ratio of successful API requests to total requests
      ratioMetric:
        good:
          metricSource:
            type: Prometheus
            spec:
              query: sum(rate(http_requests_total{status=~"2.."}[5m]))
        total:
          metricSource:
            type: Prometheus
            spec:
              query: sum(rate(http_requests_total[5m]))
    ```

=== "Threshold Metric"

    ```yaml
    apiVersion: openslo/v1
    kind: SLI
    metadata:
      name: api-latency
      displayName: API Latency
    spec:
      description: 95th percentile API response time
      thresholdMetric:
        metricSource:
          type: Prometheus
          spec:
            query: histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))
    ```

=== "Inline in SLO"

    ```yaml
    apiVersion: openslo/v1
    kind: SLO
    metadata:
      name: api-availability
      displayName: API Availability
    spec:
      service: api-service
      indicator:
        metadata:
          name: api-success-rate
          displayName: API Success Rate
        spec:
          ratioMetric:
            good:
              metricSource:
                type: Prometheus
                spec:
                  query: sum(rate(http_requests_total{status=~"2.."}[5m]))
            total:
              metricSource:
                type: Prometheus
                spec:
                  query: sum(rate(http_requests_total[5m]))
      timeWindow:
        - duration: 30d
          isRolling: true
      budgetingMethod: Occurrences
      objectives:
        - displayName: 99.9% availability
          target: 0.999
    ```

## Properties

{{ generate_object_properties("openslo/v1", "SLI") }}
