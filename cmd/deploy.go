package cmd

import (
	"context"
	"os"
	"zeus/util"

	"github.com/microsoft/azure-devops-go-api/azuredevops"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var branch string

var pipelineNames = map[string]string{
	"broadsea":         "",
	"vocab":            "",
	"achilles":         "",
	"broadsea-build":   "Broadsea Build Pipeline",
	"broadsea-release": "Broadsea Release Pipeline",
	"vocab-build":      "Vocabulary Build Pipeline",
	"vocab-release":    "Vocabulary Release Pipeline",
	"achilles-build":   "Achilles Build",
	"achilles-release": "Achilles Release",
}

// Initializes the configuration for the given environment
func Deploy(app string, ctx context.Context, connection *azuredevops.Connection, id []int, project string) (err error) {
	if app == "broadsea-build" {
		// Execute Azure DevOps build pipeline to build Broadsea Docker image
		if error := util.ExecPipeline(ctx, connection, pipelineNames["broadsea-build"], id[0], project); error != nil {
			return error
		}
	} else if app == "broadsea-release" {
		// Execute Azure DevOps build pipeline to deploy Broadsea Docker image
		if error := util.ExecPipeline(ctx, connection, pipelineNames["broadsea-release"], id[0], project); error != nil {
			return error
		}
	} else if app == "vocab-build" {
		// Execute Azure DevOps build pipeline to build vocabulary dacpac
		if error := util.ExecPipeline(ctx, connection, pipelineNames["vocab-build"], id[0], project); error != nil {
			return error
		}
	} else if app == "vocab-release" {
		// Execute Azure DevOps build pipeline to deploy vocabulary dacpac
		if error := util.ExecPipeline(ctx, connection, pipelineNames["vocab-release"], id[0], project); error != nil {
			return error
		}
	} else if app == "achilles-build" {
		// Execute Azure DevOps build pipeline to build Achilles-Synthea-ETL Docker image
		if error := util.ExecPipeline(ctx, connection, pipelineNames["achilles-build"], id[0], project); error != nil {
			return error
		}
	} else if app == "achilles-release" {
		// Execute Azure DevOps build pipeline to deploy Achilles-Synthea-ETL Docker image
		if error := util.ExecPipeline(ctx, connection, pipelineNames["achilles-release"], id[0], project); error != nil {
			return error
		}
	}
	return err
}

var deployCmd = &cobra.Command{
	Use:   "deploy <broadsea-application-name> --branch <environment-branch>",
	Short: "Deploys a OHDSI Applications",
	Long:  `Deploys a OHDSI Applications (WebAPI, Atlas, ETL-Synthea, and Achilles)`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var app = "broadsea-app"

		if len(args) > 0 {
			app = args[0]
		}

		_, found := pipelineNames[app]

		if found {
			// Define variables to log in Azure DevOps
			organizationUrl := organization
			personalAccessToken := os.Getenv("AZURE_DEVOPS_EXT_PAT") // todo: replace value with your PAT

			// Create a connection to your organization
			connection := azuredevops.NewPatConnection(organizationUrl, personalAccessToken)
			ctx := context.Background()

			if app == "broadsea" {
				broadseaBuildId, err := util.ReturnBuildId(pipelineNames["broadsea-build"], project, ctx, connection)
				if err != nil {
					log.Error("%s: %s", err, broadseaBuildId)
					return err
				}
				Deploy("broadsea-build", ctx, connection, broadseaBuildId, project)
				broadseaReleaseId, err := util.ReturnBuildId(pipelineNames["broadsea-release"], project, ctx, connection)
				if err != nil {
					log.Error("%s: %s", err, broadseaReleaseId)
					return err
				}
				Deploy("broadsea-release", ctx, connection, broadseaReleaseId, project)
			} else if app == "vocab" {
				vocabBuildId, err := util.ReturnBuildId(pipelineNames["vocab-build"], project, ctx, connection)
				if err != nil {
					log.Error("%s: %s", err, vocabBuildId)
					return err
				}
				Deploy("vocab-build", ctx, connection, vocabBuildId, project)
				vocabReleaseId, err := util.ReturnBuildId(pipelineNames["vocab-release"], project, ctx, connection)
				if err != nil {
					log.Error("%s: %s", err, vocabReleaseId)
					return err
				}
				Deploy("vocab-release", ctx, connection, vocabReleaseId, project)
			} else if app == "achilles" {
				achillesBuildId, err := util.ReturnBuildId(pipelineNames["achilles-build"], project, ctx, connection)
				if err != nil {
					log.Error("%s: %s", err, achillesBuildId)
					return err
				}
				Deploy("achilles-build", ctx, connection, achillesBuildId, project)
				achillesReleaseId, err := util.ReturnBuildId(pipelineNames["achilles-release"], project, ctx, connection)
				if err != nil {
					log.Error("%s: %s", err, achillesReleaseId)
					return err
				}
				Deploy("achilles-release", ctx, connection, achillesReleaseId, project)
			} else {
				id, err := util.ReturnBuildId(pipelineNames[app], project, ctx, connection)
				if err != nil {
					log.Error("%s: %s", err, id)
					return err
				}
				return Deploy(app, ctx, connection, id, project)
			}
		} else {
			log.Fatal(app, " is not a valid option.")
		}
		return err
	},
}

func init() {
	deployCmd.Flags().StringVar(&branch, "b", "main", "Environment branch name for pipeline to deploy on")
	deployCmd.Flags().StringVar(&organization, "org", "", "Azure DevOps organization (URL)")
	deployCmd.Flags().StringVar(&project, "proj", "", "Azure DevOps project name")
	rootCmd.AddCommand(deployCmd)
}
