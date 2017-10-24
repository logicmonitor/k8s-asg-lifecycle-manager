package kubectl

import (
	"bytes"
	"fmt"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

// Kubectl kubectl
type Kubectl struct{}

// Exec run a kubectl command
func (k Kubectl) Exec(args []string) error {
	var outbuf, errbuf bytes.Buffer
	cmd := exec.Command("kubectl", args...)
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	log.Infof("Running command kubectl %v", args)
	err := cmd.Run()
	if err != nil {
		log.Error(errbuf.String())
		return fmt.Errorf(errbuf.String())
	}
	log.Info(outbuf.String())
	return nil
}
