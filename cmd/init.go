package cmd

import (
	"context"
	"flag"
	"math/rand"
	"os"
	"os/exec"
	"reflect"
	"time"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/kyokomi/emoji"
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
var AZURE_DEVOPS_ORGANIZATION string
var AZURE_DEVOPS_PROJECT string

var pipeline_mapping = map[string]string{
	"Broadsea Build Pipeline":   "broadsea_build_pipeline",
	"Broadsea Release Pipeline": "broadsea_webtools_release_pipeline",
	//"Broadsea Methods Release": "broadsea_methods_release_pipeline",
	"Vocabulary Build Pipeline":   "vocabulary_build_pipeline",
	"Vocabulary Release Pipeline": "vocabulary_release_pipeline",
	"Yvonne Test Pipeline":        "broadsea_build_pipeline",
}

var (
	words     = flag.Int("words", 2, "The number of words in the pet name")
	separator = flag.String("separator", "-", "The separator between words in the pet name")
)

// Init function initializes the configuration for a given environment
func Init(environment string, organization string, project string) (err error) {
	log.Info("Confirming prerequisites are installed: ")
	requiredSystemTools := []string{"git", "terraform", "az"}
	for _, tool := range requiredSystemTools {
		path, err := exec.LookPath(tool)
		if err != nil {
			return err
		}
		log.Info(emoji.Sprint("Using: " + tool + " from " + path + " :white_check_mark:"))
	}

	if error := installAzDevOpsExt(); error != nil {
		return error

	}
	log.Info(emoji.Sprint("Azure DevOps Extension is installed :white_check_mark:"))

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

func createServiceEndpoint(ctx context.Context, connection *azuredevops.Connection, project string) (svc_con_id string, err error) {
	rand.Seed(time.Now().UTC().UnixNano())
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
	rand := petname.Generate(*words, *separator)

	client, err := serviceendpoint.NewClient(ctx, connection)
	if err != nil {
		log.Error("Unable to connect to Azure DevOps.")
		log.Fatal(err)
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

	*endpoint.Name = "zeus-svc-con-" + rand
	*endpoint.Type = "Github"
	*endpoint.Url = "https://github.com"
	auth := serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{},
		Scheme:     new(string),
	}

	auth.Parameters = &map[string]string{
		"apitoken": os.Getenv(AZURE_DEVOPS_EXT_GITHUB_PAT),
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

func updateServiceEndpoint(id string, ctx context.Context, connection *azuredevops.Connection, project string) (err error) {
	log.Info("Attempting to update Service Endpoint: ", id)
	client, err := pipelinepermissions.NewClient(ctx, connection)
	if err != nil {
		log.Error("Unable to connect to Azure DevOps.")
		log.Fatal(err)
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

func PipelineExists(arrayType interface{}, item interface{}) bool {
	arr := reflect.ValueOf(arrayType)

	for i := 0; i < arr.Len(); i++ {
		if arr.Index(i).Interface() == item {
			return true
		}
	}

	return false
}

func verifyPipelines(ctx context.Context, connection *azuredevops.Connection, project string) (err error) {
	ohdsiPipelinesList := []string{"Broadsea Build Pipeline", "Broadsea Release Pipeline", "Vocabulary Build Pipeline", "Vocabulary Release Pipeline", "Yvonne Test Pipeline"}
	pipelinesList := []string{}
	log.Info("Verifying OHDSI Pipelines:")
	client := pipelines.NewClient(ctx, connection)
	if err != nil {
		log.Error("Unable to connect to Azure DevOps.")
		log.Fatal(err)
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
		// List all existing pipelines
		for _, pipes := range (*responseValue).Value {
			pipelinesList = append(pipelinesList, *pipes.Name)
			index++
		}
		// if continuationToken has a value, then there is at least one more page of pipelines to get
		if responseValue.ContinuationToken != "" {
			// Get next page of pipelines
			pipelineArgs := pipelines.ListPipelinesArgs{
				ContinuationToken: &responseValue.ContinuationToken,
			}
			responseValue, err = client.ListPipelines(ctx, pipelineArgs)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			responseValue = nil
		}
	}

	for _, pipes := range ohdsiPipelinesList {
		if PipelineExists(pipelinesList, pipes) {
			log.Info(emoji.Sprintf("%v :white_check_mark:", pipes))
			return
		} else {
			log.Warning(emoji.Sprintf("%v :x:", pipes))
			log.Info("Attempting to import pipeline...")
			createPipeline(pipes, pipeline_mapping[pipes])
		}
	}

	return err
}

// Import required pipelines if they do not exists
func createPipeline(pipeline string, yaml_name string) (err error) {
	yaml_path := "pipelines/" + yaml_name + ".yaml"
	log.Info("Importing pipeline yaml: ", yaml_path)
	cmd := exec.Command("az", "pipelines", "create", "--name", pipeline, "--description", "This pipeline was created by Zeus", "--repository", "https://github.com/yradsmikham/OHDSIonAzure", "--branch", "master", "--yaml-path", yaml_path, "--repository-type", "github", "--service-connection", "8b78bfd2-ad49-449f-b220-b329efd6f601")
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Error("%s: %s", err, output)
		return err
	}
	return err
}

// CURRENTLY UNSUPPORTED: https://github.com/microsoft/azure-devops-go-api/issues/79
/*func createPipeline(ctx context.Context, connection *azuredevops.Connection, project string, pipelineName string) (err error) {
	client := pipelines.NewClient(ctx, connection)
	if err != nil {
		log.Error("Unable to connect to Azure DevOps.")
		log.Fatal(err)
	}

	configurationType := pipelines.ConfigurationType("yaml")
	createPipelineConfigurationParameters := pipelines.CreatePipelineConfigurationParameters{
		Type: &configurationType,
		Path: ,
		Repository: ,

	}

	createPipelineParameters := pipelines.CreatePipelineParameters{
		Configuration: &createPipelineConfigurationParameters,
		//Folder:        new(string),
		Name: new(string),
	}
	*createPipelineParameters.Name = pipelineName

	createPipelinesArgs := pipelines.CreatePipelineArgs{
		InputParameters: &createPipelineParameters,
		Project:         &project,
	}

	responseValue, err := client.CreatePipeline(ctx, createPipelinesArgs)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Info(*responseValue.Name + " created successfully")
	}
	return err
} */

var initCmd = &cobra.Command{
	Use:   "init [--env environment]",
	Short: "Initializes an OHDSI CDM instance",
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		var organizationUrl string
		var projectName string

		// Verify that environment variables are set
		verifyEnvVariables("AZURE_DEVOPS_EXT_PAT")
		verifyEnvVariables("AZURE_DEVOPS_EXT_GITHUB_PAT")
		personalAccessToken := os.Getenv("AZURE_DEVOPS_EXT_PAT")

		if cmd.Flags().Changed("org") {
			organizationUrl = organization
		} else {
			verifyEnvVariables("AZURE_DEVOPS_ORGANIZATION")
			organizationUrl = os.Getenv("AZURE_DEVOPS_ORGANIZATION")
		}

		// Create a connection to your organization
		connection := azuredevops.NewPatConnection(organizationUrl, personalAccessToken)
		ctx := context.Background()

		if cmd.Flags().Changed("proj") {
			projectName = project
		} else {
			verifyEnvVariables("AZURE_DEVOPS_PROJECT")
			projectName = os.Getenv("AZURE_DEVOPS_PROJECT")
		}

		if cmd.Flags().Changed("service-connection-id") {
			// Go straight to checking pipelines
			log.Info("Proceeding to use provided service connection.")
		} else {
			// Create Service Connection
			if id, error := createServiceEndpoint(ctx, connection, projectName); error != nil {
				return error
			} else {
				updateServiceEndpoint(id, ctx, connection, projectName)
			}

		}
		// Verify Pipelines
		if error := verifyPipelines(ctx, connection, projectName); error != nil {
			return error
		}
		return Init(environment, organization, projectName)
	},
}

func init() {
	initCmd.Flags().StringVar(&environment, "env", "", "OHDSI environment")
	initCmd.Flags().StringVar(&organization, "org", "", "Azure DevOps organization (URL)")
	initCmd.Flags().StringVar(&project, "proj", "", "Azure DevOps project name")
	initCmd.Flags().StringVar(&service_connection_id, "service-connection-id", "", "Service Connection for Github")
	/*if error := initCmd.MarkFlagRequired("env"); error != nil {
		return
	}*/
	rootCmd.AddCommand(initCmd)
}
