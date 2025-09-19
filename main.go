package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azcertificates"
	"github.com/fatih/color"
)

const (
	WarningDays  = 30
	CriticalDays = 7
)

type Config struct {
	VaultURI     string
	ClientID     string
	TenantID     string
	ClientSecret string
	CertName     string
}

func main() {
	config := loadConfig()

	if err := validateConfig(config); err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	ctx := context.Background()

	cred, err := getAzureCredential(config)
	if err != nil {
		log.Fatalf("Authentication failed: %v", err)
	}

	client, err := azcertificates.NewClient(config.VaultURI, cred, nil)
	if err != nil {
		log.Fatalf("Failed to create Key Vault client: %v", err)
	}

	if config.CertName != "" {
		checkSingleCertificate(ctx, client, config.CertName)
	} else {
		checkAllCertificates(ctx, client)
	}
}

func loadConfig() Config {
	var config Config

	flag.StringVar(&config.VaultURI, "vault-uri", "", "Azure Key Vault URI")
	flag.StringVar(&config.ClientID, "client-id", "", "Azure Client ID")
	flag.StringVar(&config.TenantID, "tenant-id", "", "Azure Tenant ID")
	flag.StringVar(&config.ClientSecret, "client-secret", "", "Azure Client Secret")
	flag.StringVar(&config.CertName, "cert-name", "", "Specific certificate name to check")
	flag.Parse()

	if config.VaultURI == "" {
		config.VaultURI = os.Getenv("AZURE_KEY_VAULT_URI")
	}
	if config.ClientID == "" {
		config.ClientID = os.Getenv("AZURE_CLIENT_ID")
	}
	if config.TenantID == "" {
		config.TenantID = os.Getenv("AZURE_TENANT_ID")
	}
	if config.ClientSecret == "" {
		config.ClientSecret = os.Getenv("AZURE_CLIENT_SECRET")
	}
	if config.CertName == "" {
		config.CertName = os.Getenv("AZURE_KEYVAULT_CERT_NAME")
	}

	return config
}

func validateConfig(config Config) error {
	if config.VaultURI == "" {
		return fmt.Errorf("AZURE_KEY_VAULT_URI is required")
	}
	if config.ClientID == "" {
		return fmt.Errorf("AZURE_CLIENT_ID is required")
	}
	if config.TenantID == "" {
		return fmt.Errorf("AZURE_TENANT_ID is required")
	}
	if config.ClientSecret == "" {
		return fmt.Errorf("AZURE_CLIENT_SECRET is required")
	}
	return nil
}

func getAzureCredential(config Config) (azcore.TokenCredential, error) {
	return azidentity.NewClientSecretCredential(
		config.TenantID,
		config.ClientID,
		config.ClientSecret,
		nil,
	)
}

func checkSingleCertificate(ctx context.Context, client *azcertificates.Client, certName string) {
	resp, err := client.GetCertificate(ctx, certName, "", nil)
	if err != nil {
		log.Fatalf("Failed to get certificate %s: %v", certName, err)
	}

	checkCertificateExpiry(certName, resp.Attributes)
}

func checkAllCertificates(ctx context.Context, client *azcertificates.Client) {
	pager := client.NewListCertificatePropertiesPager(nil)

	hasIssues := false
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			log.Fatalf("Failed to list certificates: %v", err)
		}

		for _, cert := range page.Value {
			if cert.ID != nil && cert.Attributes != nil {
				certName := getCertNameFromID(string(*cert.ID))
				if checkCertificateExpiry(certName, cert.Attributes) {
					hasIssues = true
				}
			}
		}
	}

	if !hasIssues {
		color.Green("✓ All certificates are valid for more than %d days", WarningDays)
	}
}

func checkCertificateExpiry(certName string, attrs *azcertificates.CertificateAttributes) bool {
	if attrs == nil || attrs.Expires == nil {
		fmt.Printf("Certificate %s: No expiry date found\n", certName)
		return false
	}

	now := time.Now()
	expiryTime := *attrs.Expires
	daysUntilExpiry := int(expiryTime.Sub(now).Hours() / 24)

	if attrs.Enabled != nil && !*attrs.Enabled {
		color.Yellow("⚠ Certificate %s is disabled", certName)
		return false
	}

	switch {
	case daysUntilExpiry < 0:
		color.Red("✗ EXPIRED: Certificate %s expired %d days ago", certName, -daysUntilExpiry)
		return true
	case daysUntilExpiry < CriticalDays:
		color.Red("✗ CRITICAL: Certificate %s expires in %d days (%s)", certName, daysUntilExpiry, expiryTime.Format("2006-01-02"))
		return true
	case daysUntilExpiry < WarningDays:
		color.Yellow("⚠ WARNING: Certificate %s expires in %d days (%s)", certName, daysUntilExpiry, expiryTime.Format("2006-01-02"))
		return true
	default:
		fmt.Printf("✓ Certificate %s expires in %d days (%s)\n", certName, daysUntilExpiry, expiryTime.Format("2006-01-02"))
		return false
	}
}

func getCertNameFromID(id string) string {
	for i := len(id) - 1; i >= 0; i-- {
		if id[i] == '/' {
			return id[i+1:]
		}
	}
	return id
}
