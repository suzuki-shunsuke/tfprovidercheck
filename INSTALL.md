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
brew install suzuki-shunsuke/tfprovidercheck/tfprovidercheck
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
1. [slsa-verifier](https://github.com/slsa-framework/slsa-verifier)
1. [Cosign](https://github.com/sigstore/cosign)

### 1. GitHub CLI

You can install GitHub CLI by aqua.

```sh
aqua g -i cli/cli
```

```sh
version=v1.0.1
asset=tfprovidercheck_darwin_arm64.tar.gz
gh release download -R suzuki-shunsuke/tfprovidercheck "$version" -p "$asset"
gh attestation verify "$asset" \
  -R suzuki-shunsuke/tfprovidercheck \
  --signer-workflow suzuki-shunsuke/go-release-workflow/.github/workflows/release.yaml
```

### 2. slsa-verifier

You can install slsa-verifier by aqua.

```sh
aqua g -i slsa-framework/slsa-verifier
```

```sh
version=v1.0.1
asset=tfprovidercheck_darwin_arm64.tar.gz
gh release download -R suzuki-shunsuke/tfprovidercheck "$version" -p "$asset" -p multiple.intoto.jsonl
slsa-verifier verify-artifact "$asset" \
  --provenance-path multiple.intoto.jsonl \
  --source-uri github.com/suzuki-shunsuke/tfprovidercheck \
  --source-tag "$version"
```

### 3. Cosign

You can install Cosign by aqua.

```sh
aqua g -i sigstore/cosign
```

```sh
version=v1.0.1
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
cat "$checksum_file" | sha256sum -c --ignore-missing
```
