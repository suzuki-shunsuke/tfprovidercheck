# tfprovidercheck

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

tfprovidercheck is a command line tool for security, and prevents malicious Terraform Providers from being executed.
You can define the allow list of Terraform Providers and their versions, and check if disallowed providers aren't used.

## Install

tfprovidercheck is a single binary written in [Go](https://go.dev/). So you only need to install an execurable file into `$PATH`.

1. [Homebrew](https://brew.sh/)

```sh
brew install suzuki-shunsuke/tfprovidercheck/tfprovidercheck
```

2. [Scoop](https://scoop.sh/)

```sh
scoop bucket add suzuki-shunsuke https://github.com/suzuki-shunsuke/scoop-bucket
scoop install tfprovidercheck
```

1. [aqua](https://aquaproj.github.io/)

```sh
aqua g -i suzuki-shunsuke/tfprovidercheck
```

1. Download a prebuilt binary from [GitHub Releases](https://github.com/suzuki-shunsuke/tfprovidercheck/releases) and install it into `$PATH`

## Usage

Please run `terraform init` in advance to get Terraform Providers.

```sh
terraform version -json | tfprovidercheck [-c <configuration file path>]
```

```sh
tfprovidercheck --help
```

## Configuration

There are several ways to configure tfprovidercheck.

1. The command line option `-config [-c]`
1. The environment variable `TFPROVIDERCHECK_CONFIG_BODY`
1. The environment variable `TFPROVIDERCHECK_CONFIG`
1. tfprovidercheck looks for a configuration file `.tfprovidercheck.yaml` from the current directory to the root directory

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
If you run this tool on GitHub Actions, `pull_request_target` event is useful to prevent workflows from being tampered.

[Secure GitHub Actions by pull_request_target](https://dev.to/suzukishunsuke/secure-github-actions-by-pullrequesttarget-641)

tfprovidercheck supports configuring with the environment variable `TFPROVIDERCHECK_CONFIG_BODY`, so you can define the configuraiton in a workflow file.

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

## LICENSE

[MIT](LICENSE)
