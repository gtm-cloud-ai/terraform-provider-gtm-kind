# Release Process

This document describes how to create a new release for the terraform-provider-gtm-kind.

## Prerequisites

Before creating a release, ensure you have:

1. **GPG Key Setup**: You need a GPG key for signing releases. The key should be configured in your GitHub repository secrets:
   - `GPG_PRIVATE_KEY`: Your GPG private key (export with `gpg --armor --export-secret-key YOUR_KEY_ID`)
   - `PASSPHRASE`: Your GPG key passphrase

2. **Permissions**: You need write access to the repository to create tags and releases.

## Creating a Release

1. **Update the CHANGELOG.md** (if you have one) with the new version and changes.

2. **Create and push a git tag**:
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

3. **GitHub Actions will automatically**:
   - Build binaries for multiple platforms (Linux, macOS, Windows, FreeBSD)
   - Create a GitHub release
   - Upload the signed binaries and checksums
   - Generate the Terraform Registry manifest

## GPG Key Setup

If you haven't set up GPG signing yet:

1. **Generate a GPG key** (if you don't have one):
   ```bash
   gpg --full-generate-key
   ```

2. **Export your public key** and add it to your GitHub account:
   ```bash
   gpg --armor --export YOUR_KEY_ID
   ```

3. **Export your private key** and add it to repository secrets:
   ```bash
   gpg --armor --export-secret-key YOUR_KEY_ID
   ```

4. **Add the following secrets to your GitHub repository**:
   - `GPG_PRIVATE_KEY`: The output from step 3
   - `PASSPHRASE`: Your GPG key passphrase

## Manual Testing

You can test the build process locally:

```bash
# Install goreleaser
go install github.com/goreleaser/goreleaser@latest

# Test the configuration
goreleaser check

# Build snapshot (without releasing)
goreleaser build --snapshot --clean
```

## Troubleshooting

- **"Resource not accessible by integration" error**: This usually means the GITHUB_TOKEN doesn't have sufficient permissions. The workflow uses the default GITHUB_TOKEN which should work for most cases.

- **GPG signing errors**: Ensure your GPG_PRIVATE_KEY and PASSPHRASE secrets are correctly set in the repository settings.

- **Build failures**: Check that all dependencies are properly declared in go.mod and that the code compiles successfully.
