package util

import (
	"os/exec"

	log "github.com/sirupsen/logrus"
)

// ExecuteBroadseaBuild function will run azdo pipeline to build Broadsea docker images
func ExecuteBroadseaBuild() (err error) {
	log.Info("Queuing run for Broadsea Build pipeline")
	cmd := exec.Command("az", "pipelines", "run", "--name", "Broadsea Build")
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Error("%s: %s", err, output)
		return err
	}
	log.Info("Broadsea build pipeline executed!")
	return err
}

// ExecuteBroadseaRelease function will run azdo pipeline to deploy Broadsea to app services
func ExecuteBroadseaWebToolsRelease() (err error) {
	log.Info("Queuing run for Broadsea WebTools Release pipeline")
	cmd := exec.Command("az", "pipelines", "run", "--name", "Broadsea WebTools Release")

	if output, err := cmd.CombinedOutput(); err != nil {
		log.Error("%s: %s", err, output)
		return err
	}
	log.Info("Broadsea Webtools release pipeline executed!")

	return err
}

// ExecuteBroadseaRelease function will run azdo pipeline to deploy Broadsea to app services
func ExecuteBroadseaMethodsRelease() (err error) {

	log.Info("Queuing run for Broadsea Methods Release pipeline")
	cmd := exec.Command("az", "pipelines", "run", "--name", "Broadsea Methods Release")

	if output, err := cmd.CombinedOutput(); err != nil {
		log.Error("%s: %s", err, output)
		return err
	}
	log.Info("Broadsea Methods release pipeline executed!")

	return err
}

// Execute function will run azdo pipeline to deploy Broadsea to app services
func ExecuteVocabBuild() (err error) {

	log.Info("Queuing run for Vocabulary Build pipeline")
	cmd := exec.Command("az", "pipelines", "run", "--name", "Vocabulary Build")

	if output, err := cmd.CombinedOutput(); err != nil {
		log.Error("%s: %s", err, output)
		return err
	}
	log.Info("Vocabulary build pipeline executed!")

	return err
}

func ExecuteVocabRelease() (err error) {

	log.Info("Queuing run for Vocabulary Release pipeline")
	cmd := exec.Command("az", "pipelines", "run", "--name", "Vocabulary Release")

	if output, err := cmd.CombinedOutput(); err != nil {
		log.Error("%s: %s", err, output)
		return err
	}
	log.Info("Vocabulary release pipeline executed!")

	return err
}
