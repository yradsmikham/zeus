package cmd

import (
	"fmt"
	"os/exec"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var environment string
var organization string
var project string

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

	if error := azLogin(organization, project); error != nil {
		return error
	}

	if error := checkPipelines(); error != nil {
		return error
	}

	return err
}

// Install Azure Devops Extension to Azure CLI if it is not already installed
func installAzDevOpsExt() (err error) {
	cmd := exec.Command("az", "extension", "add", "--name", "azure-devops")
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Error("%s: %s", err, output)
		return err
	}
	return err
}

// Configure organization and project for Azure DevOps
func azLogin(organization string, project string) (err error) {
	cmd := exec.Command("az", "devops", "configure", "--defaults", "organization="+organization, "project="+project)
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Error("%s: %s", err, output)
		return err
	}
	log.Info("Azure DevOps defaults configured.")
	log.Info("Successfully logged into Azure DevOps.")
	return err
}

// Verify if required pipelines already exist in Azure DevOps sub
func checkPipelines() (err error) {
	pipelines := []string{"Broadsea Build", "Broadsea WebTools Release", "Broadsea Methods Release"}
	for i, pipeline := range pipelines {
		num := i + 1
		log.Info("Verifying pipelines: ", num, " ", pipeline)
		output, err := exec.Command("az", "pipelines", "list", "--name", pipeline).CombinedOutput()
		if len(output) == 0 {
			fmt.Println(pipeline + " does not exist. Importing new pipeline.")
			if error := importPipelines(pipeline); error != nil {
				return error
			}
		}
		if err != nil {
			log.Error("There was an error validating pipeline " + pipeline)
			log.Error(" %s: %s", err, output)
		}
	}
	return err
}

// Import required pipelines if they do not exists
func importPipelines(pipeline string) (err error) {
	cmd := exec.Command("az", "pipelines", "create", "--name", pipeline, "--description", "pipeline created by Zeus", "--repository", "https://dev.azure.com/US-HLS-AppInnovations/_git/OHDSIonAzure", "--branch", "main", "--yml-path", "azure-pipelines.yml", "--repository-type", "tfsgit")
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
		return Init(environment, organization, project)
	},
}

func init() {
	initCmd.Flags().StringVar(&environment, "env", "", "OHDSI environment")
	initCmd.Flags().StringVar(&organization, "org", "", "Azure DevOps organization (URL)")
	initCmd.Flags().StringVar(&project, "proj", "", "Azure DevOps project name")
	if error := initCmd.MarkFlagRequired("env"); error != nil {
		return
	}
	rootCmd.AddCommand(initCmd)
}
