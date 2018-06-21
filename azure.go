package main

import (
	"os"

	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2017-09-30/containerservice"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-02-01/resources"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
)

type azureClient struct {
	clientId string
	clientSecret string

	resourceGroupsClient resources.GroupsClient
	containerServiceClient containerservice.ManagedClustersClient
}

func buildAzureClient() (*azureClient, error) {
	environmentName := azure.PublicCloud.Name
	tenantId := os.Getenv("ARM_TENANT_ID")
	subscriptionId := os.Getenv("ARM_SUBSCRIPTION_ID")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	token, err := buildAzureServicePrincipalToken(environmentName, tenantId, clientId, clientSecret)
	if err != nil {
		return nil, err
	}

	resourceGroupsClient := resources.NewGroupsClient(subscriptionId)
	resourceGroupsClient.Authorizer = autorest.NewBearerAuthorizer(token)

	containerServiceClient := containerservice.NewManagedClustersClient(subscriptionId)
	containerServiceClient.Authorizer = autorest.NewBearerAuthorizer(token)

	client := azureClient{
		clientId: clientId,
		clientSecret: clientSecret,
		containerServiceClient: containerServiceClient,
		resourceGroupsClient: resourceGroupsClient,
	}

	return &client, nil
}

func buildAzureServicePrincipalToken(environmentName string, tenantId string, clientId string, clientSecret string) (*adal.ServicePrincipalToken, error) {
	env, err := azure.EnvironmentFromName(environmentName)
	if err != nil {
		return nil, err
	}

	oauthConfig, err := adal.NewOAuthConfig(env.ActiveDirectoryEndpoint, tenantId)
	if err != nil {
		return nil, err
	}

	return adal.NewServicePrincipalToken(*oauthConfig, clientId, clientSecret, env.ResourceManagerEndpoint)
}
