# Service

{{ generate_object_description("openslo/v1", "Service") }}

## Examples

=== "Basic"

    ```yaml
    apiVersion: openslo/v1
    kind: Service
    metadata:
      name: example-service
      displayName: Example Service
    spec:
      description: Example service description
    ```

=== "Full"

    ```yaml
    apiVersion: openslo/v1
    kind: Service
    metadata:
      annotations:
        openslo.com/service-folder: ./my/directory
      labels:
        env:
          - dev
        team:
          - team-a
          - team-b
      name: example-service
      displayName: Example Service
    spec:
      description: Example service description
    ```

## Properties

{{ generate_object_properties("openslo/v1", "Service") }}
