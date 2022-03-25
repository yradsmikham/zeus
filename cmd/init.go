package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/microsoft/azure-devops-go-api/azuredevops"
	"github.com/microsoft/azure-devops-go-api/azuredevops/pipelinepermissions"
	"github.com/microsoft/azure-devops-go-api/azuredevops/pipelines"
	"github.com/microsoft/azure-devops-go-api/azuredevops/pipelineschecks"
	"github.com/microsoft/azure-devops-go-api/azuredevops/serviceendpoint"
	"github.com/microsoft/azure-devops-go-api/azuredevops/webapi"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var environment string
var organization string
var project string
var service_connection_id string
var AZURE_DEVOPS_EXT_GITHUB_PAT string
var AZURE_DEVOPS_EXT_PAT string

var pipeline_mapping = map[string]string{
	"Broadsea Build":            "broadsea_build_pipeline",
	"Broadsea Methods Release":  "broadsea_methods_release_pipeline",
	"Broadsea WebTools Release": "broadsea_webtools_release_pipeline",
	"Vocabulary Build":          "vocabulary_build_pipeline",
	"Vocabulary Release":        "vocabulary_release_pipeline",
	"Yvonne Test Pipeline":      "broadsea_build_pipeline",
}

// Init function initializes the configuration for a given environment
func Init(environment string, organization string, project string) (err error) {
	requiredSystemTools := []string{"git", "terraform", "az"}
	for _, tool := range requiredSystemTools {
		path, err := exec.LookPath(tool)
		if err != nil {
			return err
		}
		log.Info("Using: " + tool + " from " + path)
	}

	if error := installAzDevOpsExt(); error != nil {
		return error

	}
	log.Info("Azure DevOps Extension is installed.")

	/*if error := azLogin(organization, project); error != nil {
		return error
	}*/

	// Validate OHDSI Application pipelines exists
	if error := checkPipelines(); error != nil {
		return error
	}

	return err
}

// Verify that environment variables are set locally
func verifyEnvVariables(variable string) (err error) {
	_, ok := os.LookupEnv(variable)
	if !ok {
		log.Error(variable, " not set\n")
		os.Exit(1)
	}
	return err
}

// Install Azure Devops Extension to Azure CLI if it is not already installed
func installAzDevOpsExt() (err error) {
	addCmd := exec.Command("az", "extension", "add", "--name", "azure-devops")
	if output, err := addCmd.CombinedOutput(); err != nil {
		log.Error("%s: %s", err, output)
		return err
	}
	setCmd := exec.Command("az", "config", "set", "extension.use_dynamic_install=yes_without_prompt")
	if output, err := setCmd.CombinedOutput(); err != nil {
		log.Error("%s: %s", err, output)
		return err
	}
	return err
}

// TODO: Check to see if a service connection exists for Github
// Challenging because will require parsing JSON output
/*func checkServiceConnection() (err error) {
	cmd := exec.Command("az", "devops", "service-endpoint", "list")
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Error("Unable to list service connections.")
		log.Error("%s: %s", err, output)
		return err
	}
	return err
}

// TO-DO: Create Service Connection between AzureDevops and GitHub
// Return Service Connection ID
func createServiceConnectionGithub() (id string, err error) {
	log.Info("Checking for Github Service Connection...")
	cmd := exec.Command("az", "devops", "service-endpoint", "github", "create", "--name", "github_service_connection_zeus", "--github-url", "https://github.com")
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Error("%s: %s", err, string(output))
		return "Unable to create Service Connection. Please try again.", err
	} else {
		log.Info(string(output))
	}
	// Needs to return the service connection ID
	return "", err
} */

func createServiceEndpoint(ctx context.Context, connection *azuredevops.Connection, project string) (svc_con_id string, err error) {
	client, err := serviceendpoint.NewClient(ctx, connection)
	if err != nil {
		log.Error("Unable to connect to Azure DevOps.")
		log.Fatal(err)
	} else {
		log.Info("Successfully Connected to Azure DevOps!")
	}

	endpoint := serviceendpoint.ServiceEndpoint{
		AdministratorsGroup: &webapi.IdentityRef{},
		Authorization:       &serviceendpoint.EndpointAuthorization{},
		CreatedBy:           &webapi.IdentityRef{},
		Data:                &map[string]string{},
		Description:         new(string),
		IsReady:             new(bool),
		IsShared:            new(bool),
		Name:                new(string),
		OperationStatus:     nil,
		Owner:               new(string),
		ReadersGroup:        &webapi.IdentityRef{},
		Type:                new(string),
		Url:                 new(string),
	}
	*endpoint.Name = "zeus-svc-con"
	*endpoint.Type = "Github"
	*endpoint.Url = "https://github.com"
	auth := serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{},
		Scheme:     new(string),
	}

	auth.Parameters = &map[string]string{
		"apitoken": os.Getenv(AZURE_DEVOPS_EXT_PAT),
	}
	*auth.Scheme = "Token"
	*endpoint.Authorization = auth
	*endpoint.IsShared = false
	*endpoint.IsReady = true

	serviceEndpointCreateArgs := serviceendpoint.CreateServiceEndpointArgs{
		Endpoint: &endpoint,
		Project:  &project,
	}

	responseValue, err := client.CreateServiceEndpoint(ctx, serviceEndpointCreateArgs)

	if err != nil {
		log.Info("There was a problem creating service endpoint.")
		log.Fatal(err)
	} else {
		log.Info("Service Endpoint for Github '", string(*responseValue.Name), "' was created. Service Endpoint ID: ", *responseValue.Id)
	}
	uuid := *responseValue.Id
	return uuid.String(), err
}

/* // Update Service Connection to enable for pipelines
func updateServiceConnection(id string) (err error) {
	log.Info("Updating Service Connection...")
	cmd := exec.Command("az", "devops", "service-endpoint", "update", "--id", id, "--enable-for-all", "true")
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Error("%s: %s", err, string(output))
		return err
	}
	return err
} */

func updateServiceEndpoint(id string, ctx context.Context, connection *azuredevops.Connection, project string) (err error) {
	log.Info("Attempting to update Service Endpoint: ", id)
	client, err := pipelinepermissions.NewClient(ctx, connection)
	if err != nil {
		log.Error("Unable to connect to Azure DevOps.")
		log.Fatal(err)
	} else {
		log.Info("Successfully Connected to Azure DevOps!")
	}

	permission := pipelinepermissions.Permission{
		Authorized:   new(bool),
		AuthorizedBy: &webapi.IdentityRef{},
		AuthorizedOn: &azuredevops.Time{},
	}
	*permission.Authorized = true

	resource := pipelineschecks.Resource{
		Id:   &id,
		Type: new(string),
	}
	*resource.Type = "endpoint"

	pipelinePermission := pipelinepermissions.PipelinePermission{
		Authorized:   new(bool),
		AuthorizedBy: &webapi.IdentityRef{},
		AuthorizedOn: &azuredevops.Time{},
		Id:           new(int),
	}
	*pipelinePermission.Authorized = true

	resourcePipelinePermissions := pipelinepermissions.ResourcePipelinePermissions{
		AllPipelines: &permission,
		Resource:     &resource,
	}

	updatePipelinepPermissionsForResourceArgs := pipelinepermissions.UpdatePipelinePermisionsForResourceArgs{
		ResourceAuthorization: &resourcePipelinePermissions,
		Project:               &project,
		ResourceType:          new(string),
		ResourceId:            &id,
	}
	*updatePipelinepPermissionsForResourceArgs.ResourceType = "endpoint"

	responseValue, err := client.UpdatePipelinePermisionsForResource(ctx, updatePipelinepPermissionsForResourceArgs)
	if err != nil {
		log.Info("There was a problem updating service endpoint.")
		log.Fatal(err)
	} else {
		log.Info("Service Endpoint updated successfully!")
		log.Info("Grant access permissions to all pipelines: ", *responseValue.AllPipelines.Authorized)
	}

	return err

}

/*// Configure organization and project for Azure DevOps
func azLogin(organization string, project string) (err error) {
	cmd := exec.Command("az", "devops", "configure", "--defaults", "organization="+organization, "project="+project)
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Error("%s: %s", err, output)
		return err
	}
	log.Info("Azure DevOps defaults configured.")
	log.Info("Successfully logged into Azure DevOps.")
	return err
} */

// Verify if required pipelines already exist in Azure DevOps sub
func checkPipelines() (err error) {
	//pipelines := []string{"Broadsea Build", "Broadsea Methods Release", "Broadsea WebTools Release", "Vocabulary Build", "Vocabulary Release", "Yvonne Test Pipeline"}
	pipelines := []string{"Broadsea Build Pipeline", "Yvonne Test Pipeline"}
	for _, pipeline := range pipelines {
		log.Info("Verifying pipelines: ", pipeline)
		output, err := exec.Command("az", "pipelines", "list", "--name", pipeline).CombinedOutput()
		if len(output) == 3 {
			fmt.Println(pipeline + " does not exist. Importing new pipeline.")
			if error := importPipelines(pipeline, pipeline_mapping[pipeline]); error != nil {
				return error
			}
		} else {
			log.Info(pipeline, " validated")
		}
		if err != nil {
			log.Error("There was an error validating pipeline " + pipeline)
			log.Error(" %s: %s", err, output)
		}
	}
	return err
}

func verifyPipelines(ctx context.Context, connection *azuredevops.Connection, project string) (err error) {
	//pipelinesList := []string{"Broadsea Build Pipeline", "Yvonne Test Pipeline"}
	log.Info("Verifying OHDSI pipelines")
	client := pipelines.NewClient(ctx, connection)
	if err != nil {
		log.Error("Unable to connect to Azure DevOps.")
		log.Fatal(err)
	} else {
		log.Info("Successfully Connected to Azure DevOps!")
	}

	listPipelinesArgs := pipelines.ListPipelinesArgs{
		Project: &project,
		//OrderBy: new(string),
	}

	responseValue, err := client.ListPipelines(ctx, listPipelinesArgs)
	if err != nil {
		log.Fatal(err)
	}

	index := 0
	for responseValue != nil {
		// Log the page of team project names
		for _, pipes := range (*responseValue).Value {
			log.Printf("Name[%v] = %v", index, *pipes.Name)
			index++
		}
	}

	return err
}

// TO-DO: Import required pipelines if they do not exists
func importPipelines(pipeline string, yaml_name string) (err error) {
	yaml_path := "pipelines/" + yaml_name + ".yaml"
	log.Info("Importing pipeline yaml: ", yaml_path)
	cmd := exec.Command("az", "pipelines", "create", "--name", pipeline, "--description", "This pipeline was created by Zeus", "--repository", "https://github.com/yradsmikham/OHDSIonAzure", "--branch", "master", "--yaml-path", yaml_path, "--repository-type", "github", "--service-connection", "8b78bfd2-ad49-449f-b220-b329efd6f601")
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Error("%s: %s", err, output)
		return err
	}
	return err
}

var initCmd = &cobra.Command{
	Use:   "init [--env environment]",
	Short: "Initializes an OHDSI CDM instance",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if cmd.Flags().Changed("service-connection-id") {
			// Go straight to checking pipelines
			log.Info("Proceeding to use provided service connection.")
		} else {
			// Verify that environment variables are set
			verifyEnvVariables("AZURE_DEVOPS_EXT_PAT")
			verifyEnvVariables("AZURE_DEVOPS_EXT_GITHUB_PAT")

			// Define variables to log in Azure DevOps
			organizationUrl := organization
			personalAccessToken := os.Getenv("AZURE_DEVOPS_EXT_PAT") // todo: replace value with your PAT

			// Create a connection to your organization
			connection := azuredevops.NewPatConnection(organizationUrl, personalAccessToken)
			ctx := context.Background()

			// Create Service Connection
			if id, error := createServiceEndpoint(ctx, connection, project); error != nil {
				return error
			} else {
				updateServiceEndpoint(id, ctx, connection, project)
			}

			if error := verifyPipelines(ctx, connection, project); error != nil {
				return error
			}
		}
		return Init(environment, organization, project)
	},
}

func init() {
	initCmd.Flags().StringVar(&environment, "env", "", "OHDSI environment")
	initCmd.Flags().StringVar(&organization, "org", "", "Azure DevOps organization (URL)")
	initCmd.Flags().StringVar(&project, "proj", "", "Azure DevOps project name")
	initCmd.Flags().StringVar(&service_connection_id, "service-connection-id", "", "Service Connection for Github")
	if error := initCmd.MarkFlagRequired("env"); error != nil {
		return
	}
	rootCmd.AddCommand(initCmd)
}
