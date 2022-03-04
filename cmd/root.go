package cmd

import (
	"os"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "zeus",
	Short: "OHDSI on Azure",
	Long:  "OHDSI on Azure",

	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		verbose := cmd.Flag("verbose").Value.String()

		if verbose == "true" {
			log.SetLevel(log.DebugLevel)
		} else {
			log.SetLevel(log.InfoLevel)
		}

		return err
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Use verbose output logs")
}
