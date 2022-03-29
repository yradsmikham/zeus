package util

import (
	"os/exec"

	log "github.com/sirupsen/logrus"
)

// Execute pipeline
func ExecPipeline(name string) (err error) {
	log.Info("Queuing run for %v Build pipeline", name)
	cmd := exec.Command("az", "pipelines", "run", "--name", name)
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Error("%s: %s", err, output)
		return err
	}
	log.Info("%v executed!", name)
	return err
}
