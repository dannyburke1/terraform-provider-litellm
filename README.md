# LiteLLM Terraform Provider

This Terraform provider allows you to manage LiteLLM resources through Infrastructure as Code. It provides support for managing models, teams, team members, and API keys via the LiteLLM REST API.

## Features

- Manage LiteLLM model configurations
- Create and manage teams
- Configure team members and their permissions
- Set usage limits and budgets
- Control access to specific models
- Specify model modes (e.g., completion, embedding, image generation)
- Manage API keys with fine-grained controls
- Support for reasoning effort configuration in the model resource
- Managed Users as an individual resource, without needing a tea

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 0.13.x
- [Go](https://golang.org/doc/install) >= 1.16 (for development)

## Using the Provider

To use the LiteLLM provider in your Terraform configuration, you need to declare it in the <code>terraform</code> block:

```hcl
terraform {
  required_providers {
    litellm = {
      source  = "dannyburke1/litellm"
      version = "~> 1.0.0"
    }
  }
}

provider "litellm" {
  api_base = var.litellm_api_base
  api_key  = var.litellm_api_key
}
```

Then, you can use the provider to manage LiteLLM resources. Here's an example of creating a model configuration:

```hcl
resource "litellm_model" "gpt4" {
  model_name          = "gpt-4-proxy"
  custom_llm_provider = "openai"
  model_api_key       = var.openai_api_key
  model_api_base      = "https://api.openai.com/v1"
  base_model          = "gpt-4"
  tier                = "paid"
  mode                = "chat"
  reasoning_effort    = "medium"  # Optional: "low", "medium", or "high"
  
  input_cost_per_million_tokens  = 30.0
  output_cost_per_million_tokens = 60.0
}
```

For full details on the <code>litellm_model</code> resource, see the [model resource documentation](docs/resources/model.md).

Here's an example of creating an API key with various options:

```hcl
resource "litellm_key" "example_key" {
  models               = ["gpt-4", "claude-3.5-sonnet"]
  max_budget           = 100.0
  user_id              = "user123"
  team_id              = "team456"
  max_parallel_requests = 5
  tpm_limit            = 1000
  rpm_limit            = 60
  budget_duration      = "30d"
  key_alias            = "prod-key-1"
  duration             = "30d"
  metadata             = {
    environment = "production"
  }
  allowed_cache_controls = ["no-cache", "max-age=3600"]
  soft_budget          = 80.0
  aliases              = {
    "gpt-4" = "gpt4"
  }
  config               = {
    default_model = "gpt-4"
  }
  permissions          = {
    can_create_keys = "true"
  }
  model_max_budget     = {
    "gpt-4" = 50.0
  }
  model_rpm_limit      = {
    "claude-3.5-sonnet" = 30
  }
  model_tpm_limit      = {
    "gpt-4" = 500
  }
  guardrails           = ["content_filter", "token_limit"]
  blocked              = false
  tags                 = ["production", "api"]
}
```

The <code>litellm_key</code> resource supports the following options:

- <code>models</code>: List of allowed models for this key
- <code>max_budget</code>: Maximum budget for the key
- <code>user_id</code> and <code>team_id</code>: Associate the key with a user and team
- <code>max_parallel_requests</code>: Limit concurrent requests
- <code>tpm_limit</code> and <code>rpm_limit</code>: Set tokens and requests per minute limits
- <code>budget_duration</code>: Specify budget duration (e.g., "30d", "30s", "7d", "10m")
- <code>key_alias</code>: Set a friendly name for the key
- <code>duration</code>: Set the key's validity period
- <code>metadata</code>: Add custom metadata to the key
- <code>allowed_cache_controls</code>: Specify allowed cache control directives
- <code>soft_budget</code>: Set a soft budget limit
- <code>aliases</code>: Define model aliases
- <code>config</code>: Set configuration options
- <code>permissions</code>: Specify key permissions
- <code>model_max_budget</code>, <code>model_rpm_limit</code>, <code>model_tpm_limit</code>: Set per-model limits
- <code>guardrails</code>: Apply specific guardrails to the key
- <code>blocked</code>: Flag to block/unblock the key
- <code>tags</code>: Add tags for organization and filtering

For full details on the <code>litellm_key</code> resource, see the [key resource documentation](docs/resources/key.md).

### Available Resources

- <code>litellm_model</code>: Manage model configurations. [Documentation](docs/resources/model.md)
- <code>litellm_team</code>: Manage teams. [Documentation](docs/resources/team.md)
- <code>litellm_team_member</code>: Manage team members. [Documentation](docs/resources/team_member.md)
- <code>litellm_key</code>: Manage API keys. [Documentation](docs/resources/key.md)

## Development

### Project Structure

The project is organized as follows:

```
terraform-provider-litellm/
├── litellm/
│   ├── provider.go
│   ├── resource_model.go
│   ├── resource_model_crud.go
│   ├── resource_team.go
│   ├── resource_team_member.go
│   ├── resource_key.go
│   ├── resource_key_utils.go
│   ├── resource_user.go
│   ├── types.go
│   └── utils.go
├── main.go
├── go.mod
├── go.sum
├── Makefile
└── ...
```

### Building the Provider

1. Clone the repository:
```sh
git clone https://github.com/your-username/terraform-provider-litellm.git
```

2. Enter the repository directory:
```sh
cd terraform-provider-litellm
```

3. Build and install the provider:
```sh
make install
```

### Development Commands

The Makefile provides several useful commands for development:

- `make build`: Builds the provider
- `make install`: Builds and installs the provider
- `make test`: Runs the test suite
- `make fmt`: Formats the code
- `make vet`: Runs go vet
- `make lint`: Runs golangci-lint
- `make clean`: Removes build artifacts and installed provider

### Testing

To run the tests:
```sh
make test
```

### Contributing

Contributions are welcome! Please read our [contributing guidelines](CONTRIBUTING.md) first.

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## Notes

- Always use environment variables or secure secret management solutions to handle sensitive information like API keys and AWS credentials.
- Refer to the `examples/` directory for more detailed usage examples.
- Make sure to keep your provider version updated for the latest features and bug fixes.
- v0.2.3 introduces support for the <code>reasoning_effort</code> attribute in the <code>litellm_model</code> resource. This attribute accepts "low", "medium", or "high" to control the model's reasoning effort.
