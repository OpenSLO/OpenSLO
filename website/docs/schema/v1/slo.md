# SLO

{{ generate_object_description("openslo/v1", "SLO") }}

## Examples

=== "Basic"

    ```yaml
    apiVersion: openslo/v1
    kind: SLO
    metadata:
      name: main-page
      displayName: Main page availability
    spec:
      description: Our main page should be available always
      service: web-shop
      indicator:
        metadata:
          name: response-code
          displayName: Response codes of requests to main page
        spec:
          ratioMetric:
            good:
              metricSource:
                type: Any # Here put any service that holds information you need.
                spec: # Fields necessary to query service for the data.
                  query: Any # 'query' is just an example field.
            total:
              metricSource:
                type: Any # Here put any service that holds information you need.
                spec: # Fields necessary to query service for the data.
                  query: Any # 'query' is just an example field.
      timeWindow:
       - duration: 2w
         isRolling: true
      budgetingMethod: RatioTimeslices
      objectives:
        - displayName: Good
          timeSliceWindow: 1m
          target: 0.99
    ```

## Properties

{{ generate_object_properties("openslo/v1", "SLO") }}
