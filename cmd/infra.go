package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// TO-DO: validate environment resources in Azure Portal
type Response []struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Sha         string `json:"sha"`
	Size        int    `json:"size"`
	URL         string `json:"url"`
	HTMLURL     string `json:"html_url"`
	GitURL      string `json:"git_url"`
	DownloadURL string `json:"download_url"`
	Type        string `json:"type"`
	Links       struct {
		Self string `json:"self"`
		Git  string `json:"git"`
		HTML string `json:"html"`
	} `json:"_links"`
}

// Example of struct
/*{
    "name": "main.tf",
    "path": "infra/terraform/bootstrap/main.tf",
    "sha": "fed44b33803f34c317dfe49c2e547b4b8e85ff4b",
    "size": 34597,
    "url": "https://api.github.com/repos/yradsmikham/OHDSIonAzure/contents/infra/terraform/bootstrap/main.tf?ref=master",
    "html_url": "https://github.com/yradsmikham/OHDSIonAzure/blob/master/infra/terraform/bootstrap/main.tf",
    "git_url": "https://api.github.com/repos/yradsmikham/OHDSIonAzure/git/blobs/fed44b33803f34c317dfe49c2e547b4b8e85ff4b",
    "download_url": "https://raw.githubusercontent.com/yradsmikham/OHDSIonAzure/master/infra/terraform/bootstrap/main.tf",
    "type": "file",
    "_links": {
      "self": "https://api.github.com/repos/yradsmikham/OHDSIonAzure/contents/infra/terraform/bootstrap/main.tf?ref=master",
      "git": "https://api.github.com/repos/yradsmikham/OHDSIonAzure/git/blobs/fed44b33803f34c317dfe49c2e547b4b8e85ff4b",
      "html": "https://github.com/yradsmikham/OHDSIonAzure/blob/master/infra/terraform/bootstrap/main.tf"
    }
  },
*/

func Infra(operation string, envName string) (err error) {
	bootstrapGitApi := "https://api.github.com/repos/yradsmikham/OHDSIonAzure/contents/infra/terraform/bootstrap"
	bootstrapGitRaw := "https://raw.githubusercontent.com/yradsmikham/OHDSIonAzure/master/infra/terraform/bootstrap/"
	bootstrapDir := envName + "/" + envName + "-bootstrap"
	cdmGitApi := "https://api.github.com/repos/yradsmikham/OHDSIonAzure/contents/infra/terraform/omop"
	cdmGitRaw := "https://raw.githubusercontent.com/yradsmikham/OHDSIonAzure/master/infra/terraform/omop/"
	cdmDir := envName + "/" + envName + "-cdm"

	// Create Bootstrap directory and OHDSI CDM env directory
	if operation == "create" {
		log.Info("Creating directory ", envName)
		os.MkdirAll(bootstrapDir, 0755)
		os.MkdirAll(cdmDir, 0755)
	}
	curlTfFiles(bootstrapGitApi, bootstrapGitRaw, bootstrapDir)
	curlTfFiles(cdmGitApi, cdmGitRaw, cdmDir)
	return (err)
}

func curlTfFiles(apiUrl string, rawUrl string, dir string) (err error) {
	// Copy TF files over to each directory
	// curl -X GET https://api.github.com/repos/yradsmikham/OHDSIonAzure/contents/infra/terraform/bootstrap
	resp, err := http.Get(apiUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var result Response
	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to the go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}

	// Loop through the data node for the FirstName
	for _, tfFile := range result {
		fmt.Println(tfFile.Name)
		log.Info("Curling File from Repo: " + rawUrl + tfFile.Name)
		cmd := exec.Command("curl", "--output-dir", dir, "-OL", rawUrl+tfFile.Name)
		if output, err := cmd.CombinedOutput(); err != nil {
			log.Error("%s: %s", err, output)
			return err
		}
	}
	return err
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
