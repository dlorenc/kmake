package builder

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/r2d4/kmake/pkg/kmake/util"
	"github.com/sirupsen/logrus"
)

func LocalBuild(imageName, dockerfile, _ string, t Tagger) (string, error) {
	tag, err := t.Tag()
	if err != nil {
		return "", err
	}
	logrus.Infof("Generated tag: %s", tag)
	fqImageName := fmt.Sprintf("%s:%s", imageName, tag)
	logrus.Infof("Building docker image %s from %s", fqImageName, dockerfile)
	cmd := exec.Command("docker", "build", "-q", "--tag", fqImageName, "--file", dockerfile, filepath.Dir(dockerfile))
	stdout, stderr, err := util.RunCommand(cmd, nil)

	if err != nil {
		return "", errors.Wrapf(err, "docker build: %s %s", string(stdout), string(stderr))
	}

	return fqImageName, nil
}

func Push(imageName string) error {
	logrus.Infof("Pushing docker image %s", imageName)
	output, err := exec.Command("docker", "push", imageName).CombinedOutput()
	logrus.Infof("Push output: %s", output)
	return err
}
