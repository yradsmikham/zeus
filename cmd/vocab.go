package cmd

/*

import (
	"os/exec"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)


var account_name string
var account_key string
var container string
var path string

// Uploads vocabulary files to Storage Account
func Vocab(operation string) (err error) {
	log.Info("Attempting to", operation, "vocabulary files")
	if operation == "upload" {
		if error := uploadVocab(pipeline); error != nil {
			return error
		}
	}
	if operation == "import" {
		log.Info("Queuing vocabulary build and release pipeline...")
		// Execute vocab build and release pipeline
	}
	return err
}

func uploadVocab(account_name string, account_key string, container string, vocab_path string) (err error) {
	cmd := exec.Command("az", "storage", "blob", "upload-batch", "--account-name", account_name, "--account-key", account_key, "--destination", container, "--source", vocab_path, "--pattern", "*.csv")
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Error("%s: %s", err, output)
		return err
	}
	log.Info("Uploading vocabulary files to ", container)

	return err
}

func importVocab(env string) (err error) {
	// TODO: Not using "env"
	// TODO: Need to execute Vocabulary Release pipeline as well
	if error := util.ExecPipeline("vocab-release"); error != nil {
		return error
	}

	return err
}

var vocabCmd = &cobra.Command{
	Use:   "vocab <operation>",
	Short: "Performs different operations for vocabulary files",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 || args[0] != "import" || args[0] != "upload" {
			log.Error("Must specify operation (i.e. 'import' or 'upload') ")
		}
	},
}

var uploadCmd = &cobra.Command{
	Use:   "upload --account-name account-name --account-key account-key --container container-name --path path-to-vocab-files",
	Short: "Uploads Vocabulary Files to Storage Account Blob",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		accountName, _ := cmd.Flags().GetString("account-name")
		accountKey, _ := cmd.Flags().GetString("account-key")
		container, _ := cmd.Flags().GetString("container")
		path, _ := cmd.Flags().GetString("path")

		return uploadVocab(accountName, accountKey, container, path)
	},
}

var importCmd = &cobra.Command{
	Use:   "import --environment environment-name",
	Short: "Uploads Vocabulary Files to Storage Account Blob",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		environment, _ := cmd.Flags().GetString("environment")
		return importVocab(environment)
	},
}

func init() {
	vocabCmd.AddCommand(uploadCmd)
	vocabCmd.AddCommand(importCmd)

	uploadCmd.PersistentFlags().StringP("container", "c", "", "Name of Azure Blob Container")
	uploadCmd.MarkPersistentFlagRequired("container")
	uploadCmd.PersistentFlags().StringP("path", "p", "", "Path to Vocabulary Files")
	uploadCmd.MarkPersistentFlagRequired("path")
	uploadCmd.PersistentFlags().StringP("account-name", "n", "", "Azure Storage Account Name")
	uploadCmd.MarkPersistentFlagRequired("account-name")
	uploadCmd.PersistentFlags().StringP("account-key", "k", "", "Azure Storage Account Key")
	uploadCmd.MarkPersistentFlagRequired("account-key")

	// TO-DO: This needs more thought. Should "environment" correspond to a git branch? Or possibly the name of the Azure resource like SQL Database name?
	importCmd.PersistentFlags().StringP("environment", "e", "", "OHDSI CDM environment")
	importCmd.MarkPersistentFlagRequired("environment")

	rootCmd.AddCommand(vocabCmd)
}

*/
