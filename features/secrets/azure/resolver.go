package azure

import (
	"context"
	"os"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
	core "github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald"
)

// the resolver for Azure Key Vault secrets
type azureSecretsResolver struct {
	kvClient *azsecrets.Client
}

// making sure it complies with the interface defined by Goald
var _ goald.ISecretResolver = (*azureSecretsResolver)(nil)

var (
	resolver                       *azureSecretsResolver
	resolverOnce, resolverInitOnce sync.Once
)

// SecretsResolver init
func SecretsResolver() *azureSecretsResolver {
	// making sure the resolver is created only once
	resolverOnce.Do(func() {
		resolver = &azureSecretsResolver{}
	})

	return resolver
}

// ResolveSecret implements [goald.ISecretResolver]; lazily builds the Key Vault client, once
func (thisResolver *azureSecretsResolver) ResolveSecret(name string) string {
	// ensuring the resolver is initialized
	resolverInitOnce.Do(func() {
		vaultURL := os.Getenv("KEYVAULT_URL")
		core.PanicMsgIf(vaultURL == "", "KEYVAULT_URL is not set, can't resolve secrets")

		cred, errCred := azidentity.NewDefaultAzureCredential(nil)
		core.PanicMsgIfErr(errCred, "No Azure credential available")

		var errClient error
		thisResolver.kvClient, errClient = azsecrets.NewClient(vaultURL, cred, nil)
		core.PanicMsgIfErr(errClient, "Could not create the Key Vault client")
	})

	// using the initialized resolver to get the secret
	resp, err := thisResolver.kvClient.GetSecret(context.Background(), name, "", nil)
	core.PanicMsgIfErr(err, "Could not get secret '%s' from Key Vault", name)

	return *resp.Value
}
