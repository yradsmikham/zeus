package cmd

import (
	"github.com/spf13/cobra"
)

// Initializes the configuration for the given environment
func apps() (err error) {
	// execute az devops pipeline
	return err
}

var appsCmd = &cobra.Command{
	Use:   "apps deploy [--webtools deploy-broadsea-webtools] [--methods deploy-broadsea-methods]",
	Short: "Deploys a OHDSI Applications",
	Long:  `Deploys a OHDSI Applications (WebAPI, Atlas, ETL-Synthea, and Achilles)`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		return apps()
	},
}

func init() {
	appsCmd.PersistentFlags().BoolP("webtools", "w", false, "Deploy Broadsea Webtools to environment")
	appsCmd.PersistentFlags().BoolP("methods", "m", false, "Deploy Broadsea Methods Library to environment")
	rootCmd.AddCommand(appsCmd)
}
