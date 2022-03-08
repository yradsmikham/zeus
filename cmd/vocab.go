package cmd

import (
	"os/exec"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var container string
var vocab_path string
var account_name string
var account_key string

// Uploads vocabulary files to Storage Account
func Vocab(container string, vocab_path string, account_name string, account_key string) (err error) {
	cmd := exec.Command("az", "storage", "blob", "upload-batch", "-d", container, "--source", vocab_path, "--account-name", account_name, "--account-key", account_key)
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Error("%s: %s", err, output)
		return err
	}
	log.Info("Uploading vocabulary files to ", container)
	return err
}

var vocabCmd = &cobra.Command{
	Use:   "init [--env environment]",
	Short: "Initializes an OHDSI CDM instance",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		return Init(environment, organization, project)
	},
}

func init() {
	vocabCmd.Flags().StringVar(&container, "c", "", "Name of Azure Blob Container")
	vocabCmd.Flags().StringVar(&vocab_path, "p", "", "Path to Vocabulary Files")
	vocabCmd.Flags().StringVar(&account_name, "n", "", "Azure Storage Account Name")
	vocabCmd.Flags().StringVar(&account_key, "k", "", "Azure Storage Account Key")
	rootCmd.AddCommand(vocabCmd)
}
