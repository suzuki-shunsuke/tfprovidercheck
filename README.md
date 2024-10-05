# tfprovidercheck

[Install](#install) | [Usage](#usage) | [Config](#configuration)

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

3. [aqua](https://aquaproj.github.io/)

```sh
aqua g -i suzuki-shunsuke/tfprovidercheck
```

4. Download a prebuilt binary from [GitHub Releases](https://github.com/suzuki-shunsuke/tfprovidercheck/releases) and install it into `$PATH`

<details>
<summary>Verify downloaded assets from GitHub Releases</summary>

You can verify downloaded assets using some tools.

1. [GitHub CLI](https://cli.github.com/)
1. [slsa-verifier](https://github.com/slsa-framework/slsa-verifier)
1. [Cosign](https://github.com/sigstore/cosign)

--

1. GitHub CLI

tfprovidercheck >= v1.0.1

You can install GitHub CLI by aqua.

```sh
aqua g -i cli/cli
```

```sh
gh release download -R suzuki-shunsuke/tfprovidercheck v1.0.1 -p tfprovidercheck_darwin_arm64.tar.gz
gh attestation verify tfprovidercheck_darwin_arm64.tar.gz \
  -R suzuki-shunsuke/tfprovidercheck \
  --signer-workflow suzuki-shunsuke/go-release-workflow/.github/workflows/release.yaml
```

Output:

```
Loaded digest sha256:4e444f43865f52c1d969a9af9691f062f60e0bc64c713ee2e90c2163c8ce0d67 for file://tfprovidercheck_darwin_arm64.tar.gz
Loaded 1 attestation from GitHub API
✓ Verification succeeded!

sha256:4e444f43865f52c1d969a9af9691f062f60e0bc64c713ee2e90c2163c8ce0d67 was attested by:
REPO                                 PREDICATE_TYPE                  WORKFLOW
suzuki-shunsuke/go-release-workflow  https://slsa.dev/provenance/v1  .github/workflows/release.yaml@7f97a226912ee2978126019b1e95311d7d15c97a
```

2. slsa-verifier

You can install slsa-verifier by aqua.

```sh
aqua g -i slsa-framework/slsa-verifier
```

```sh
gh release download -R suzuki-shunsuke/tfprovidercheck v1.0.1 -p tfprovidercheck_darwin_arm64.tar.gz  -p multiple.intoto.jsonl
slsa-verifier verify-artifact tfprovidercheck_darwin_arm64.tar.gz \
  --provenance-path multiple.intoto.jsonl \
  --source-uri github.com/suzuki-shunsuke/tfprovidercheck \
  --source-tag v1.0.1
```

Output:

```
Verified signature against tlog entry index 137013754 at URL: https://rekor.sigstore.dev/api/v1/log/entries/108e9186e8c5677a1d8396570e05ff2dccf0d5060dbf587764f72ac809690392674ad9b57aa6bcf7
Verified build using builder "https://github.com/slsa-framework/slsa-github-generator/.github/workflows/generator_generic_slsa3.yml@refs/tags/v2.0.0" at commit 877e39ddd975a45a467fc4a5bacdf55d374df198
Verifying artifact tfprovidercheck_darwin_arm64.tar.gz: PASSED

PASSED: SLSA verification passed
```

3. Cosign

You can install Cosign by aqua.

```sh
aqua g -i sigstore/cosign
```

```sh
gh release download -R suzuki-shunsuke/tfprovidercheck v1.0.1
cosign verify-blob \
  --signature tfprovidercheck_1.0.1_checksums.txt.sig \
  --certificate tfprovidercheck_1.0.1_checksums.txt.pem \
  --certificate-identity-regexp 'https://github\.com/suzuki-shunsuke/go-release-workflow/\.github/workflows/release\.yaml@.*' \
  --certificate-oidc-issuer "https://token.actions.githubusercontent.com" \
  tfprovidercheck_1.0.1_checksums.txt
```

Output:

```
Verified OK
```

After verifying the checksum, verify the artifact.

```sh
cat tfprovidercheck_1.0.1_checksums.txt | sha256sum -c --ignore-missing
```

</details>

5. go install

```sh
go install github.com/suzuki-shunsuke/tfprovidercheck/cmd/tfprovidercheck@latest
```

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

## Compared with .terraform.lock.hcl and required_providers block

About [.terraform.lock.hcl](https://developer.hashicorp.com/terraform/language/files/dependency-lock), .terraform.lock.hcl doesn't work as the allow list of providers because `terraform init` adds missing providers automatically.
This means malicious providers not included in .terraform.lock.hcl can be executed in CI.

About [required_providers block](https://developer.hashicorp.com/terraform/language/providers/requirements#requiring-providers), there are several reasons tfprovidercheck is useful compared with required_providers block.

First, it's difficult to validate required_providers in child Terraform Modules.

Second, required_providers block can be tampered in pull requests CI without code review.
On the other hand, you can prevent tfprovidercheck configuration from being tampered by several ways such as [GitHub Actions' `pull_request_target`](#bulb-prevent-configuration-from-being-tampered).

Third, tfprovidercheck enables you to manange the allow list of providers with a single YAML outside of Terraform working directory.
So administrators (SRE, Platform Engineer, DevOps Engineer, etc) can keep the security and governance easily while delegating the management of Terraform configuration to product teams.

If you validate providers with required_providers block, admins need to have the ownership of required_providers block and review changes of them.
In case of Monorepo, the number of required_providers block is proportion to the number of working directories.
GitHub CODEOWNERS manages the ownership per file, so admins may be supposed to review pull requests even for unrelated changes.
For example, if you [update Terraform Providers by Renovate](https://docs.renovatebot.com/modules/manager/terraform/#required_providers-block), admins need to review pull requests every time providers are updated.
In proportion to the number of working directories in Monorepo, the burden of admins gets higher.
And this also makes provider auto update difficult.

On the other hand, if you validate providers with tfprovidercheck, admins don't care about required_blocks providers unless tfprovidercheck fails, so the burden of admins gets lower.

Fourth, the purpose and intention of tfprovidercheck is so simple and clear that it's easy to handle the error of tfprovidercheck and to maintain tfprovidercheck configuraiton.
tfprovidercheck is a dedicated security tool to manage the allow list of Terraform Providers and prevent disallowed providers from being used.

## Versioning Policy

https://github.com/suzuki-shunsuke/versioning-policy

## LICENSE

[MIT](LICENSE)
