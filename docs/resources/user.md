# litellm_user Resource

Manages a user in LiteLLM. Users are created with a username, password, role and their own budget and keys.

## Example Usage

```hcl
resource "litellm_user" "this" {
  user_id         = "foo.bar"
  user_email      = foo.bar@foobar.com
  role            = "proxy_admin"
  models          = ["gpt-4-proxy", "claude-2"]

  tpm_limit       = 500000
  rpm_limit       = 5000
  max_budget      = 1000.0
  budget_duration = "30d"
  auto_create_keys = true
}
```

## Argument Reference

The following arguments are supported:

* `user_id` - (Required) A human-readable identifier for the user.

* `user_email` - (Required) The email tied to the user. Sends invitation email if SMTP is configured.  See `send_user_invite` below.

* `role` - (Optional) - The default role to assign to the user. Valid values are:
  * `proxy_admin`
  * `proxy_admin_viewer`
  * `internal_user`
  * `internal_user_viewer`
  * `team`
  * `customer`

* `models` - (Optional) List of model names that this user can access.


* `tpm_limit` - (Optional) User tokens per minute limit.

* `rpm_limit` - (Optional) User requests per minute limit.

* `max_budget` - (Optional) Maximum budget allocated to the user.

* `budget_duration` - (Optional) Duration for the budget cycle. Valid values are:
  * `daily`
  * `weekly`
  * `monthly`
  * `yearly`

* `auto_create_keys` - (Optional) Boolean value to enable or disable API keys on user creation.

* `send_email_invite` - (Optional) Boolean value to send users an email invitation for the LiteLLM instance

## Attribute Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The unique identifier for the user.

## Import

Users can be imported using the user ID:

```shell
terraform import litellm_user.user <user-id>
```

Note: The user ID is generated when the user is created and is different from the `user_alias`.