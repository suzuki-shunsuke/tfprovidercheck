# Install

tfprovidercheck is written in Go. So you only have to install a binary in your `PATH`.

There are some ways to install tfprovidercheck.

1. [Homebrew](#homebrew)
1. [Scoop](#scoop)
1. [aqua](#aqua)
1. [GitHub Releases](#github-releases)
1. [Build an executable binary from source code yourself using Go](#build-an-executable-binary-from-source-code-yourself-using-go)

## Homebrew

You can install tfprovidercheck using [Homebrew](https://brew.sh/).

[Homebrew Core Formula: tfprovidercheck](https://formulae.brew.sh/formula/tfprovidercheck)

```sh
brew install tfprovidercheck
```

Or

```sh
brew install --cask suzuki-shunsuke/tfprovidercheck/tfprovidercheck
```

## Scoop

You can install tfprovidercheck using [Scoop](https://scoop.sh/).

```sh
scoop bucket add suzuki-shunsuke https://github.com/suzuki-shunsuke/scoop-bucket
scoop install tfprovidercheck
```

## aqua

You can install tfprovidercheck using [aqua](https://aquaproj.github.io/).

```sh
aqua g -i suzuki-shunsuke/tfprovidercheck
```

## Build an executable binary from source code yourself using Go

```sh
go install github.com/suzuki-shunsuke/tfprovidercheck/cmd/tfprovidercheck@latest
```

## GitHub Releases

You can download an asset from [GitHub Releases](https://github.com/suzuki-shunsuke/tfprovidercheck/releases).
Please unarchive it and install a pre built binary into `$PATH`. 

### Verify downloaded assets from GitHub Releases

You can verify downloaded assets using some tools.

1. [GitHub CLI](https://cli.github.com/)
1. [Cosign](https://github.com/sigstore/cosign)
1. [Packslip](https://packslip.dev/) `tfprovidercheck >= v1.0.8`
1. [slsa-verifier](https://github.com/slsa-framework/slsa-verifier) `tfprovidercheck <= v1.0.7`

### 1. GitHub CLI

You can install GitHub CLI by aqua.

```sh
aqua g -i cli/cli
```

```sh
version=v1.0.8
asset=tfprovidercheck_darwin_arm64.tar.gz
gh release download -R suzuki-shunsuke/tfprovidercheck "$version" -p "$asset"
gh attestation verify "$asset" \
  -R suzuki-shunsuke/tfprovidercheck \
  --signer-workflow suzuki-shunsuke/go-release-workflow/.github/workflows/release.yaml
```

### 2. Cosign

You can install Cosign by aqua.

```sh
aqua g -i sigstore/cosign
```

From tfprovidercheck v1.0.8, the bundle file is `tfprovidercheck_<version>_checksums.txt.sigstore.json`.
Cosign v2.4.2 or later is required.

```sh
version=v1.0.8
checksum_file="tfprovidercheck_${version#v}_checksums.txt"
asset=tfprovidercheck_darwin_arm64.tar.gz
gh release download "$version" \
  -R suzuki-shunsuke/tfprovidercheck \
  -p "$asset" \
  -p "$checksum_file" \
  -p "${checksum_file}.sigstore.json"
cosign verify-blob \
  "$checksum_file" \
  --bundle "${checksum_file}.sigstore.json" \
  --certificate-identity-regexp 'https://github\.com/suzuki-shunsuke/go-release-workflow/\.github/workflows/release\.yaml@.*' \
  --certificate-oidc-issuer "https://token.actions.githubusercontent.com"
cat "$checksum_file" | sha256sum -c --ignore-missing -
```

For tfprovidercheck v1.0.7 or older, `*.pem` and `*.sig` are available instead.

```sh
version=v1.0.7
checksum_file="tfprovidercheck_${version#v}_checksums.txt"
asset=tfprovidercheck_darwin_arm64.tar.gz
gh release download "$version" \
  -R suzuki-shunsuke/tfprovidercheck \
  -p "$asset" \
  -p "$checksum_file" \
  -p "${checksum_file}.pem" \
  -p "${checksum_file}.sig"
cosign verify-blob \
  --signature "${checksum_file}.sig" \
  --certificate "${checksum_file}.pem" \
  --certificate-identity-regexp 'https://github\.com/suzuki-shunsuke/go-release-workflow/\.github/workflows/release\.yaml@.*' \
  --certificate-oidc-issuer "https://token.actions.githubusercontent.com" \
  "$checksum_file"
cat "$checksum_file" | sha256sum -c --ignore-missing -
```

### 3. Packslip

`tfprovidercheck >= v1.0.8`

You can install packslip by aqua.

```sh
aqua g -i jdx/packslip
```

```sh
version=v1.0.8
asset=tfprovidercheck_darwin_arm64.tar.gz
gh release download -R suzuki-shunsuke/tfprovidercheck "$version" -p "$asset" -p packslip.sigstore.json
packslip verify packslip.sigstore.json \
    --pin ps1_t2lk53pajlz4pt43b3z24dfmne \
    --identity-prefix https://github.com/suzuki-shunsuke/go-release-workflow/.github/workflows/release.yaml@ \
    --artifact "$asset"
```

- `--pin` requires the release to be signed from the repository `suzuki-shunsuke/tfprovidercheck`, identified by its repository ID rather than its name
- `--identity-prefix` requires the release to be signed by the release workflow of [suzuki-shunsuke/go-release-workflow](https://github.com/suzuki-shunsuke/go-release-workflow)

On success, packslip prints a line starting with `ok:`.
Please make sure it shows the project and version you meant, e.g. `github.com/suzuki-shunsuke/tfprovidercheck 1.0.8`.

### 4. slsa-verifier

> [!WARNING]
> As of tfprovidercheck v1.0.8, `multiple.intoto.jsonl` hasn't been published anymore.

You can install slsa-verifier by aqua.

```sh
aqua g -i slsa-framework/slsa-verifier
```

```sh
version=v1.0.7
asset=tfprovidercheck_darwin_arm64.tar.gz
gh release download -R suzuki-shunsuke/tfprovidercheck "$version" -p "$asset" -p multiple.intoto.jsonl
slsa-verifier verify-artifact "$asset" \
  --provenance-path multiple.intoto.jsonl \
  --source-uri github.com/suzuki-shunsuke/tfprovidercheck \
  --source-tag "$version"
```
