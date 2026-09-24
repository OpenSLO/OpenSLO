# Schema

Select an API version to view its objects, field descriptions, and validation rules.

{{ generate_version_links() }}

The reference uses the OpenSLO Go SDK's object descriptions and validation rules.
Shared metadata is defined in each version overview.
Properties marked **reference** link to their definition instead of repeating it.
Inline objects link to the corresponding object page.
Expand a property to see its description, allowed values, examples, and rules.

In property paths, `[*]` identifies each array element, `.*` identifies each map
value, and `.*~` identifies each map key.
An **Applies when** column lists the conditions for a validation rule.
If a rule lists multiple conditions, all must hold.
Required fields within an optional object apply when that object is present.
Reference properties retain the rules for their use in the containing object.
The SDK supplies the allowed values. Check conditional rules for further restrictions.
