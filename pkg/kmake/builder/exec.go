package builder

import (
	"os/exec"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
)

func Build(imageName, dockerfile string) error {
	logrus.Infof("Building docker image %s from %s", imageName, dockerfile)
	out, err := exec.Command("docker", "build", "--file", dockerfile, filepath.Dir(dockerfile), "--tag", imageName).CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "running docker build: %s %s", out, err)
	}
	logrus.Infof("Docker build of %s complete", imageName)
	return nil
}
