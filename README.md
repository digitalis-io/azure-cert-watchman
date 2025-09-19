# Azure Certificates Watchman

A command-line tool that monitors SSL certificate expiration in Azure Key Vault. Provides early warnings for certificates approaching expiration.

## Features

- Connect to Azure Key Vault using service principal authentication
- List and check all SSL certificates or a specific certificate
- Color-coded alerts:
  - 🔴 **CRITICAL**: Less than 7 days until expiration
  - 🟡 **WARNING**: Less than 30 days until expiration
  - 🟢 **OK**: More than 30 days until expiration
- Single binary distribution for easy deployment

## Installation

### Build from source

```bash
# Clone the repository
git clone <repository-url>
cd azure-cert-watchman

# Build for current platform
go build -o azure-cert-watchman main.go

# Or use the build script for cross-platform binaries
./build.sh
```

## Configuration

The tool accepts configuration through environment variables or command-line flags.

### Environment Variables

```bash
export AZURE_KEY_VAULT_URI="https://your-vault.vault.azure.net/"
export AZURE_CLIENT_ID="your-client-id"
export AZURE_TENANT_ID="your-tenant-id"
export AZURE_CLIENT_SECRET="your-client-secret"
export AZURE_KEYVAULT_CERT_NAME="specific-cert-name"  # Optional
```

### Command-line Flags

```bash
./azure-cert-watchman \
  --vault-uri="https://your-vault.vault.azure.net/" \
  --client-id="your-client-id" \
  --tenant-id="your-tenant-id" \
  --client-secret="your-client-secret" \
  --cert-name="specific-cert-name"  # Optional
```

## Usage

### Check all certificates in a vault

```bash
./azure-cert-watchman
```

### Check a specific certificate

```bash
./azure-cert-watchman --cert-name="my-certificate"
# or
export AZURE_KEYVAULT_CERT_NAME="my-certificate"
./azure-cert-watchman
```

## Example Output

```
✓ Certificate web-ssl expires in 180 days (2024-06-15)
⚠ WARNING: Certificate api-cert expires in 25 days (2024-02-10)
✗ CRITICAL: Certificate app-cert expires in 5 days (2024-01-20)
✗ EXPIRED: Certificate old-cert expired 10 days ago
```

## Azure Permissions

The service principal needs the following Key Vault permissions:
- **Certificates**: Get, List

## Integration with CI/CD

This tool can be easily integrated into CI/CD pipelines:

```bash
# Example: GitHub Actions
- name: Check Certificate Expiration
  run: |
    ./azure-cert-watchman
  env:
    AZURE_KEY_VAULT_URI: ${{ secrets.AZURE_KEY_VAULT_URI }}
    AZURE_CLIENT_ID: ${{ secrets.AZURE_CLIENT_ID }}
    AZURE_TENANT_ID: ${{ secrets.AZURE_TENANT_ID }}
    AZURE_CLIENT_SECRET: ${{ secrets.AZURE_CLIENT_SECRET }}
```

## Exit Codes

- `0`: All certificates are valid
- `1`: Configuration error or authentication failure
- `1`: One or more certificates are expired or expiring soon

## License

MIT