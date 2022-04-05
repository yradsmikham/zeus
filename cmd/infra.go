package cmd

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// TO-DO: validate environment resources in Azure Portal

func Infra(operation string, envName string) (err error) {
	if operation == "create" {
		log.Info("Creating directory ", envName)
		os.MkdirAll(envName+"/"+envName+"-bootstrap", 0755)
		os.MkdirAll(envName+"/"+envName+"-cdm", 0755)
	}
	return (err)
}

var infraCmd = &cobra.Command{
	Use:   "infra <operation> --env <environment-name>",
	Short: "Performs actions on OHDSI Infrastructure",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var operation string
		if len(args) > 0 {
			log.Info("operation create")
			operation = args[0]
		} else {
			log.Fatal("`zeus infra` command takes in a command `create` or `validate`. Please specify one.")
		}
		return Infra(operation, environment)
	},
}

func init() {
	infraCmd.Flags().StringVar(&environment, "env", "", "OHDSI CDM environment")
	infraCmd.MarkFlagRequired("env")
	rootCmd.AddCommand(infraCmd)
}
