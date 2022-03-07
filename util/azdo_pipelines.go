package util

import (
	"os/exec"

	log "github.com/sirupsen/logrus"
)

// ExecuteBroadseaBuild function will run azdo pipeline to build Broadsea docker images
func ExecuteBroadseaBuild() (err error) {
	cmd := exec.Command("az", "pipelines", "run", "--name", "Broadsea Build")

	if output, err := cmd.CombinedOutput(); err != nil {
		log.Error("%s: %s", err, output)
		return err
	}
	log.Info("Broadsea build pipeline executed!")
	return err
}

// ExecuteBroadseaRelease function will run azdo pipeline to deploy Broadsea to app services
func ExecuteBroadseaRelease() (err error) {
	cmd := exec.Command("az", "pipelines", "run", "--name", "Broadsea WebTools Release")

	if output, err := cmd.CombinedOutput(); err != nil {
		log.Error("%s: %s", err, output)
		return err
	}
	log.Info("Broadsea release pipeline executed!")
	return err
}
