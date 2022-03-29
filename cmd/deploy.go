package cmd

import (
	"zeus/util"

	"github.com/spf13/cobra"
)

var branch string

var pipelineNames = map[string]string{
	"broadsea-build":   "Broadsea Build Pipeline",
	"broadsea-release": "Broadsea Release Pipeline",
	"vocab-build":      "Vocabulary Build Pipeline",
	"vocab-release":    "Vocabulary Release Pipeline",
	"achilles-build":   "Achilles Build",
	"achilles-release": "Achilles Release",
}

// Initializes the configuration for the given environment
func Deploy(app string) (err error) {
	if app == "broadsea" {
		// Execute Azure DevOps build pipeline to build Broadsea Docker images
		if error := util.ExecPipeline(pipelineNames["broadsea-build"]); error != nil {
			return error
		}
		// Execute Azure DevOps build pipeline to deploy Broadsea docker image
		if error := util.ExecPipeline(pipelineNames["broadsea-release"]); error != nil {
			return error
		}
	} else if app == "vocab" {
		// Execute Azure DevOps build pipeline to build vocabulary dacpac
		if error := util.ExecPipeline(pipelineNames["vocab-build"]); error != nil {
			return error
		}
		// Execute Azure DevOps build pipeline to deploy vocabulary dacpac
		if error := util.ExecPipeline(pipelineNames["vocab-release"]); error != nil {
			return error
		}
	} else if app == "achilles" {
		// Execute Azure DevOps build pipeline for Achilles-Synthea-ETL Docker image
		if error := util.ExecPipeline(pipelineNames["achilles-build"]); error != nil {
			return error
		}
		// Execute Azure DevOps build pipeline to deploy Achilles-Synthea-ETL Docker image
		if error := util.ExecPipeline(pipelineNames["achilles-release"]); error != nil {
			return error
		}
	} else if app == "broadsea-build" {
		// Execute Azure DevOps build pipeline to build Broadsea Docker image
		if error := util.ExecPipeline(pipelineNames["broadsea-build"]); error != nil {
			return error
		}
	} else if app == "broadsea-release" {
		// Execute Azure DevOps build pipeline to deploy Broadsea Docker image
		if error := util.ExecPipeline(pipelineNames["broadsea-release"]); error != nil {
			return error
		}
	} else if app == "vocab-build" {
		// Execute Azure DevOps build pipeline to build vocabulary dacpac
		if error := util.ExecPipeline(pipelineNames["vocab-build"]); error != nil {
			return error
		}
	} else if app == "vocab-release" {
		// Execute Azure DevOps build pipeline to deploy vocabulary dacpac
		if error := util.ExecPipeline(pipelineNames["vocab-release"]); error != nil {
			return error
		}
	} else if app == "achilles-build" {
		// Execute Azure DevOps build pipeline to build Achilles-Synthea-ETL Docker image
		if error := util.ExecPipeline(pipelineNames["achilles-build"]); error != nil {
			return error
		}
	} else if app == "achilles-release" {
		// Execute Azure DevOps build pipeline to deploy Achilles-Synthea-ETL Docker image
		if error := util.ExecPipeline(pipelineNames["achilles-release"]); error != nil {
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
		return Deploy(app)
	},
}

func init() {
	deployCmd.Flags().StringVar(&branch, "b", "main", "Environment branch name for pipeline to deploy on")
	rootCmd.AddCommand(deployCmd)
}
