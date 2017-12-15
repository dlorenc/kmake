package builder

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func Build(imageName, dockerfile string) (string, error) {
	logrus.Infof("Building docker image %s from %s", imageName, dockerfile)
	digest, err := exec.Command("docker", "build", "-q", "--tag", imageName, "--file", dockerfile, filepath.Dir(dockerfile)).Output()
	if err != nil {
		return "", errors.Wrapf(err, "running docker build: %s %s", digest, err)
	}
	digestStr := strings.TrimSpace(string(digest))
	d := strings.Split(digestStr, ":")
	checksum := d[1]

	logrus.Infof("Docker build of %s complete: %s", imageName, checksum)
	out, err := exec.Command("docker", "tag", imageName, fmt.Sprintf("%s:%s", imageName, checksum)).CombinedOutput()
	if err != nil {
		return "", errors.Wrapf(err, "tagging image: %s", out)
	}

	return checksum, nil
}
