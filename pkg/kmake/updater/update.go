package updater

import (
	"os/exec"

	"github.com/pkg/errors"
)

func Update(checksum string) error {
	setCmd := exec.Command("ks", "param", "set", "hello-node", "image", "hello-node:"+checksum)
	setCmd.Dir = "./hello-node"
	if err := setCmd.Run(); err != nil {
		return errors.Wrapf(err, "setting hello-node image tag %s", "hello-node:"+checksum)
	}

	applyCmd := exec.Command("ks", "apply", "default")
	applyCmd.Dir = "./hello-node"
	if err := applyCmd.Run(); err != nil {
		return errors.Wrap(err, "applying new deployment")
	}
	return nil
}
