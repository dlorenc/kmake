package updater

import (
	"os/exec"

	"github.com/pkg/errors"
)

func Update(checksum string) error {
	envSetCmd := exec.Command("ks", "env", "set", "default")
	envSetCmd.Dir = "./hello-node"
	if err := envSetCmd.Run(); err != nil {
		return errors.Wrap(err, "setting default env")
	}

	applyCmd := exec.Command("ks", "apply", "default", "--ext-str", "image=hello-node:"+checksum)
	applyCmd.Dir = "./hello-node"
	out, err := applyCmd.Output()
	if err != nil {
		return errors.Wrapf(err, "applying new deployment: %s", out)
	}
	return nil
}
