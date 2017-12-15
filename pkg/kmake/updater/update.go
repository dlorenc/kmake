package updater

import (
	"os/exec"

	"github.com/pkg/errors"
)

func Update(imageName string) error {
	deployments, err := FindDeploymentsWithImage(imageName)
	if err != nil {
		return err
	}
	for _, d := range deployments {
		out, err := exec.Command("kubectl", "apply", d).Output()
		if err != nil {
			return errors.Wrapf(err, "running docker build: %s %s", out, err)
		}
	}
	return nil
}

func FindDeploymentsWithImage(imageName string) ([]string, error) {
	return nil, nil
}
