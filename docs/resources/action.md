---
page_title: "zitadel_action Resource - terraform-provider-zitadel"
subcategory: ""
description: |-
  Resource representing an action belonging to an organization.
---

# zitadel_action (Resource)

Resource representing an action belonging to an organization.

## Example Usage

```terraform
resource zitadel_action action {
  depends_on = [zitadel_org.org]
  provider   = zitadel

  org_id          = zitadel_org.org.id
  name            = "actionname"
  script          = "testscript"
  timeout         = "10s"
  allowed_to_fail = "true"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `allowed_to_fail` (Boolean) when true, the next action will be called even if this action fails
- `name` (String)
- `org_id` (String) ID of the organization
- `script` (String)
- `timeout` (String) after which time the action will be terminated if not finished

### Read-Only

- `id` (String) The ID of this resource.
- `state` (Number) the state of the action