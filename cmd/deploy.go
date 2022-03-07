package cmd

import (
	"fmt"

	"../util"
	"github.com/spf13/cobra"
)

var broadsea_app string

// Initializes the configuration for the given environment
func Deploy(broadsea_app string) (err error) {
	fmt.Println("Hello, World!")
	if broadsea_app == "broadsea-webtools" {
		// Execute Azure DevOps build pipeline to build broadsea-methods docker image
		if error := util.ExecuteBroadseaBuild(); error != nil {
			return error
		}
		// Execute Azure DevOps build pipeline to deploy broadsea-methods docker image
		if error := util.ExecuteBroadseaRelease(); error != nil {
			return error
		}
	}
	return err
}

var deployCmd = &cobra.Command{
	Use:   "deploy <broadsea-application-name> --env <environment-name>",
	Short: "Deploys a OHDSI Applications",
	Long:  `Deploys a OHDSI Applications (WebAPI, Atlas, ETL-Synthea, and Achilles)`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		return Deploy(broadsea_app)
	},
}

func init() {
	deployCmd.PersistentFlags().BoolP("webtools", "w", false, "Deploy Broadsea Webtools to environment")
	deployCmd.PersistentFlags().BoolP("methods", "m", false, "Deploy Broadsea Methods Library to environment")
	rootCmd.AddCommand(deployCmd)
}
