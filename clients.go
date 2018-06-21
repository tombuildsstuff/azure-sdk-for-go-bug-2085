package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2017-09-30/containerservice"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-02-01/resources"
)

func pointerToInt32(input int32) *int32 {
	return &input
}

func pointerToString(input string) *string {
	return &input
}

func (c *azureClient) createContainerService(name, resourceGroupName, location string) error {
	log.Printf("Creating Managed Kubernetes Service..")
	ctx := context.TODO()
	service := containerservice.ManagedCluster{
		Location: &location,
		ManagedClusterProperties: &containerservice.ManagedClusterProperties{
			DNSPrefix: pointerToString(name),
			LinuxProfile: &containerservice.LinuxProfile{
				AdminUsername: pointerToString("exampleuser"),
				SSH: &containerservice.SSHConfiguration{
					PublicKeys: &[]containerservice.SSHPublicKey{
						{
							KeyData: pointerToString("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC3ydA75xi8jonETRItwFox1cVw0CctGgUXjPJ91yZoqMzT7AO8clK3traYH8bF6SHKy+Ia+9FK4bA3FIynjCTNUCya/fCXgLoa2U9n3HL+2c6KKB7reCJgll5GfRJMkBPJVJAlgMj6rCbra414bJfphIYFqlItD4I6VEqs8cJ8eThIvq+nsjRkIMqUiaMMiFaRqFZY1/8App0XyXuCyWnep+U6TKzbpKuWxF9fZTdmgYcCqW4cVH8sJCvx18qtIcQJuVt9oxjqvYBFA/tbMTibfAcEvobFBWCMA0u3E84Tdjp9QIlgUxDjPKrVlOOFm2Kv+wU8pFwYePFWtOKlAfad"),
						},
					},
				},
			},
			ServicePrincipalProfile: &containerservice.ServicePrincipalProfile{
				ClientID: pointerToString(c.clientId),
				Secret: pointerToString(c.clientSecret),
			},
			AgentPoolProfiles: &[]containerservice.AgentPoolProfile{
				{
					Name: pointerToString("default"),
					Count: pointerToInt32(int32(3)),
					OsType: containerservice.Linux,
					VMSize: containerservice.VMSizeTypes("NotValid"),
				},
			},
		},
	}

	future, err := c.containerServiceClient.CreateOrUpdate(ctx, resourceGroupName, name, service)
	if err != nil {
		return err
	}

	err = future.WaitForCompletion(ctx, c.containerServiceClient.Client)
	if err != nil {
		return err
	}

	return nil
}

func (c *azureClient) deleteContainerService(name, resourceGroupName string) error {
	ctx := context.TODO()
	log.Printf("Deleting Container Service..")
	future, err := c.containerServiceClient.Delete(ctx, resourceGroupName, name)
	if err != nil {
		return err
	}

	log.Printf("Waiting for deletion of Container Service to complete..")
	err = future.WaitForCompletion(ctx, c.resourceGroupsClient.Client)
	if err != nil {
		return err
	}

	return nil
}

func (c *azureClient) createResourceGroup(name, location string) error {
	ctx := context.TODO()
	group := resources.Group{
		Location: &location,
	}

	log.Printf("Creating Resource Group..")
	_, err := c.resourceGroupsClient.CreateOrUpdate(ctx, name, group)
	if err != nil {
		return fmt.Errorf("Error creating Resource Group %q: %+v", name, err)
	}

	return nil
}

func (c *azureClient) deleteResourceGroup(name string) error {
	ctx := context.TODO()
	log.Printf("Deleting Resource Group..")
	future, err := c.resourceGroupsClient.Delete(ctx, name)
	if err != nil {
		return err
	}

	log.Printf("Waiting for deletion of Resource Group to complete..")
	err = future.WaitForCompletion(ctx, c.resourceGroupsClient.Client)
	if err != nil {
		return err
	}

	return nil
}