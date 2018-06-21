package main

import (
	"fmt"
)

func main() {
	err := run("bug2085")
	if err != nil {
		panic(err)
	}
}

func run(prefix string) error {
	client, err := buildAzureClient()
	if err != nil {
		return fmt.Errorf("Error building Azure Client: %+v", err)
	}

	name := fmt.Sprintf("%s-k8s", prefix)
	resourceGroupName := fmt.Sprintf("%s-resources", prefix)
	location := "West Europe"

	err = client.createResourceGroup(resourceGroupName, location)
	if err != nil {
		return err
	}

	err = client.createContainerService(name, resourceGroupName, location)
	if err != nil {
		return err
	}

	defer client.deleteContainerService(name, resourceGroupName)
	defer client.deleteResourceGroup(resourceGroupName)

	return nil
}

