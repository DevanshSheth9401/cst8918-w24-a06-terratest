package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// subscriptionID represents the Azure subscription ID used for testing.
// You typically want to run tests under a separate "Testing" subscription.
// For lab purposes, you will use your assigned subscription under the Cloud Dev/Ops program tenant.
var subscriptionID string = "9e8be2f9-a0b7-4c4d-841e-112767408a75"

// TestAzureLinuxVMCreation verifies if the Azure Linux VM is successfully created.
func TestAzureLinuxVMCreation(t *testing.T) {
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../",
		// Override the default terraform variables
		Vars: map[string]interface{}{
			"labelPrefix": "shet0028",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the value of output variable
	vmName := terraform.Output(t, terraformOptions, "vm_name")
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")

	// Confirm VM exists
	assert.True(t, azure.VirtualMachineExists(t, vmName, resourceGroupName, subscriptionID))
}

// TestAzureLinuxVMUbuntuVersion verifies if the deployed Azure Linux VM is running the specified Ubuntu version.
func TestAzureLinuxVMUbuntuVersion(t *testing.T) {
	// Initialize terraform options with the directory containing Terraform configurations and override variables.
	terraformOptions := &terraform.Options{
		TerraformDir: "../", // Specify the directory path of the Terraform code.
		Vars: map[string]interface{}{
			"labelPrefix": "shet0028", // Custom variable to set a label prefix for resources.
		},
	}

	// Ensure that resources are destroyed after the test execution to avoid resource leakage.
	defer terraform.Destroy(t, terraformOptions)

	// Execute 'terraform init' and 'terraform apply' commands to provision resources, failing the test on errors.
	terraform.InitAndApply(t, terraformOptions)

	// Retrieve output variables from Terraform to use in assertions.
	vmName := terraform.Output(t, terraformOptions, "vm_name")                          // VM name.
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")   // Resource group name.
	expectedVmImageVersion := terraform.Output(t, terraformOptions, "vm_image_version") // Expected VM image version.

	// Fetch the actual VM image version from Azure to verify against the expected version.
	getVirtualMachineImage := azure.GetVirtualMachineImage(t, vmName, resourceGroupName, subscriptionID)

	// Assert that the actual VM image version matches the expected version.
	assert.Equal(t, expectedVmImageVersion, getVirtualMachineImage.Version)
}

// TestAzureNicExistsAndConnectedVm checks if a Network Interface Card (NIC) exists and is properly attached to the specified VM.
func TestAzureNicExistsAndConnectedVm(t *testing.T) {
	// Initialize terraform options with the directory containing Terraform configurations and override variables.
	terraformOptions := &terraform.Options{
		TerraformDir: "../", // Specify the directory path of the Terraform code.
		Vars: map[string]interface{}{
			"labelPrefix": "shet0028", // Custom variable to set a label prefix for resources.
		},
	}

	// Ensure that resources are destroyed after the test execution to avoid resource leakage.
	defer terraform.Destroy(t, terraformOptions)

	// Execute 'terraform init' and 'terraform apply' commands to provision resources, failing the test on errors.
	terraform.InitAndApply(t, terraformOptions)

	// Retrieve output variables from Terraform to use in assertions.
	vmName := terraform.Output(t, terraformOptions, "vm_name")                        // VM name.
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name") // Resource group name.
	nicName := terraform.Output(t, terraformOptions, "nic_name")                      // NIC name.

	// Fetch the list of NICs attached to the VM from Azure.
	listNic := azure.GetVirtualMachineNics(t, vmName, resourceGroupName, subscriptionID)

	// Assert that the specified NIC is indeed attached to the VM.
	assert.Contains(t, listNic, nicName)
}
