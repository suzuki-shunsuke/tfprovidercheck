# tfprovidercheck

[![DeepWiki](https://img.shields.io/badge/DeepWiki-suzuki--shunsuke%2Ftfprovidercheck-blue.svg?logo=data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACwAAAAyCAYAAAAnWDnqAAAAAXNSR0IArs4c6QAAA05JREFUaEPtmUtyEzEQhtWTQyQLHNak2AB7ZnyXZMEjXMGeK/AIi+QuHrMnbChYY7MIh8g01fJoopFb0uhhEqqcbWTp06/uv1saEDv4O3n3dV60RfP947Mm9/SQc0ICFQgzfc4CYZoTPAswgSJCCUJUnAAoRHOAUOcATwbmVLWdGoH//PB8mnKqScAhsD0kYP3j/Yt5LPQe2KvcXmGvRHcDnpxfL2zOYJ1mFwrryWTz0advv1Ut4CJgf5uhDuDj5eUcAUoahrdY/56ebRWeraTjMt/00Sh3UDtjgHtQNHwcRGOC98BJEAEymycmYcWwOprTgcB6VZ5JK5TAJ+fXGLBm3FDAmn6oPPjR4rKCAoJCal2eAiQp2x0vxTPB3ALO2CRkwmDy5WohzBDwSEFKRwPbknEggCPB/imwrycgxX2NzoMCHhPkDwqYMr9tRcP5qNrMZHkVnOjRMWwLCcr8ohBVb1OMjxLwGCvjTikrsBOiA6fNyCrm8V1rP93iVPpwaE+gO0SsWmPiXB+jikdf6SizrT5qKasx5j8ABbHpFTx+vFXp9EnYQmLx02h1QTTrl6eDqxLnGjporxl3NL3agEvXdT0WmEost648sQOYAeJS9Q7bfUVoMGnjo4AZdUMQku50McDcMWcBPvr0SzbTAFDfvJqwLzgxwATnCgnp4wDl6Aa+Ax283gghmj+vj7feE2KBBRMW3FzOpLOADl0Isb5587h/U4gGvkt5v60Z1VLG8BhYjbzRwyQZemwAd6cCR5/XFWLYZRIMpX39AR0tjaGGiGzLVyhse5C9RKC6ai42ppWPKiBagOvaYk8lO7DajerabOZP46Lby5wKjw1HCRx7p9sVMOWGzb/vA1hwiWc6jm3MvQDTogQkiqIhJV0nBQBTU+3okKCFDy9WwferkHjtxib7t3xIUQtHxnIwtx4mpg26/HfwVNVDb4oI9RHmx5WGelRVlrtiw43zboCLaxv46AZeB3IlTkwouebTr1y2NjSpHz68WNFjHvupy3q8TFn3Hos2IAk4Ju5dCo8B3wP7VPr/FGaKiG+T+v+TQqIrOqMTL1VdWV1DdmcbO8KXBz6esmYWYKPwDL5b5FA1a0hwapHiom0r/cKaoqr+27/XcrS5UwSMbQAAAABJRU5ErkJggg==)](https://deepwiki.com/suzuki-shunsuke/tfprovidercheck)

[Install](INSTALL.md) | [Usage](#usage) | [Config](#configuration)

Censor [Terraform Providers](https://developer.hashicorp.com/terraform/language/providers).

```console
# Only google provider and azurerm provider are allowed
$ cat .tfprovidercheck.yaml
providers:
  - name: registry.terraform.io/hashicorp/google
    version: ">= 4.0.0"
  - name: registry.terraform.io/hashicorp/azurerm

# tfprovidercheck fails because aws provider is disallowed
$ terraform version -json | tfprovidercheck
FATA[0000] tfprovidercheck failed                        error="this Terraform Provider is disallowed" program=tfprovidercheck provider_name=registry.terraform.io/hashicorp/aws tfprovidercheck_version=0.1.0
```

tfprovidercheck is a command line tool to execute Terraform security.
It prevents malicious Terraform Providers from being executed.
You can define the allow list of Terraform Providers and their versions, and check if disallowed providers aren't used.

## Usage

Please run `terraform init` in advance to update the list of Terraform Providers.

```sh
terraform version -json | tfprovidercheck [-c <configuration file path>]
```

To prevent malicious codes from being executed, you should run tfprovidercheck before running other Terraform commands such as `terraform validate`, `terraform plan`, and `terraform apply`.

```console
$ tfprovidercheck --help
tfprovidercheck - Censor Terraform Providers

https://github.com/suzuki-shunsuke/tfprovidercheck

Usage:
  tfprovidercheck [<options>]

Options:
  -help, -h     Show help
  -version, -v  Show version
  -config, -c   Configuration file path
```

## Configuration

There are several ways to configure tfprovidercheck.
In order of priority, they are as follows.

1. The command line option `-config [-c]`, which is the configuration file path
1. The environment variable `TFPROVIDERCHECK_CONFIG_BODY`, which is the configuration itself (YAML)
1. The environment variable `TFPROVIDERCHECK_CONFIG`, which is the configuration file path
1. The configuration file `.tfprovidercheck.yaml` on the current directory

The field `providers` lists allowed providers and their versions.

e.g.

```yaml
providers:
  - name: registry.terraform.io/hashicorp/aws
    version: ">= 3.0.0" # Quotes are necessary because '>' is a special character for YAML
  - name: registry.terraform.io/hashicorp/google
    # version is optional
```

- `name` (Required, string): `name` must be equal to the provider name. Regular expression and glob aren't supported
- `version` (Optional, string): The version constraint of Terraform Provider. `version` is evaluated as [hashicorp/go-version' Version Constraints](https://github.com/hashicorp/go-version#version-constraints). If `version` is empty, any version is allowed

## :bulb: Prevent configuration from being tampered

It's important to prevent configuration from being tamperd.
If you run tfprovidercheck on GitHub Actions, `pull_request_target` event is useful to prevent workflows from being tampered.

[Secure GitHub Actions by pull_request_target](https://dev.to/suzukishunsuke/secure-github-actions-by-pullrequesttarget-641)

tfprovidercheck supports configuring with the environment variable `TFPROVIDERCHECK_CONFIG_BODY`, so you can define the configuration in a workflow file.

e.g.

```yaml
- run: terraform version -json | tfprovidercheck
  env:
    TFPROVIDERCHECK_CONFIG_BODY: |
      providers:
        - name: registry.terraform.io/hashicorp/aws
          version: ">= 3.0.0"
```

Then you can prevent configuration from being tampered by `pull_request_target` event.

### JSON Schema

- [tfprovidercheck.json](json-schema/tfprovidercheck.json)
- https://raw.githubusercontent.com/suzuki-shunsuke/tfprovidercheck/refs/heads/main/json-schema/tfprovidercheck.json

If you look for a CLI tool to validate configuration with JSON Schema, [ajv-cli](https://ajv.js.org/packages/ajv-cli.html) is useful.

```sh
ajv --spec=draft2020 -s json-schema/tfprovidercheck.json -d tfprovidercheck.yaml
```

#### Input Complementation by YAML Language Server

[Please see the comment too.](https://github.com/szksh-lab/.github/issues/67#issuecomment-2564960491)

Version: `main`

```yaml
# yaml-language-server: $schema=https://raw.githubusercontent.com/suzuki-shunsuke/tfprovidercheck/main/json-schema/tfprovidercheck.json
```

Or pinning version:

```yaml
# yaml-language-server: $schema=https://raw.githubusercontent.com/suzuki-shunsuke/tfprovidercheck/v1.0.2/json-schema/tfprovidercheck.json
```

## Compared with .terraform.lock.hcl and required_providers block

About [.terraform.lock.hcl](https://developer.hashicorp.com/terraform/language/files/dependency-lock), .terraform.lock.hcl doesn't work as the allow list of providers because `terraform init` adds missing providers automatically.
This means malicious providers not included in .terraform.lock.hcl can be executed in CI.

About [required_providers block](https://developer.hashicorp.com/terraform/language/providers/requirements#requiring-providers), there are several reasons tfprovidercheck is useful compared with required_providers block.

First, it's difficult to validate required_providers in child Terraform Modules.

Second, required_providers block can be tampered in pull requests CI without code review.
On the other hand, you can prevent tfprovidercheck configuration from being tampered by several ways such as [GitHub Actions' `pull_request_target`](#bulb-prevent-configuration-from-being-tampered).

Third, tfprovidercheck enables you to manage the allow list of providers with a single YAML outside of Terraform working directory.
So administrators (SRE, Platform Engineer, DevOps Engineer, etc) can keep the security and governance easily while delegating the management of Terraform configuration to product teams.

If you validate providers with required_providers block, admins need to have the ownership of required_providers block and review changes of them.
In case of Monorepo, the number of required_providers block is proportion to the number of working directories.
GitHub CODEOWNERS manages the ownership per file, so admins may be supposed to review pull requests even for unrelated changes.
For example, if you [update Terraform Providers by Renovate](https://docs.renovatebot.com/modules/manager/terraform/#required_providers-block), admins need to review pull requests every time providers are updated.
In proportion to the number of working directories in Monorepo, the burden of admins gets higher.
And this also makes provider auto update difficult.

On the other hand, if you validate providers with tfprovidercheck, admins don't care about required_blocks providers unless tfprovidercheck fails, so the burden of admins gets lower.

Fourth, the purpose and intention of tfprovidercheck is so simple and clear that it's easy to handle the error of tfprovidercheck and to maintain tfprovidercheck configuration.
tfprovidercheck is a dedicated security tool to manage the allow list of Terraform Providers and prevent disallowed providers from being used.

## Versioning Policy

https://github.com/suzuki-shunsuke/versioning-policy

## LICENSE

[MIT](LICENSE)
