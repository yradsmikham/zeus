package cmd

import (
	"zeus/util"

	"github.com/spf13/cobra"
)

var branch string

// Initializes the configuration for the given environment
func Deploy(app string) (err error) {
	if app == "broadsea-webtools" {
		// Execute Azure DevOps build pipeline to build Broadsea Docker images
		if error := util.ExecuteBroadseaBuild(); error != nil {
			return error
		}
		// Execute Azure DevOps build pipeline to deploy broadsea-webtools docker image
		if error := util.ExecuteBroadseaWebToolsRelease(); error != nil {
			return error
		}
	}
	if app == "broadsea-methods" {
		// Execute Azure DevOps build pipeline to build Broadsea Docker images
		if error := util.ExecuteBroadseaBuild(); error != nil {
			return error
		}
		// Execute Azure DevOps build pipeline to deploy broadsea-methods docker image
		if error := util.ExecuteBroadseaMethodsRelease(); error != nil {
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
