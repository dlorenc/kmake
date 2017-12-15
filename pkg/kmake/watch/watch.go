package watch

import (
	"github.com/pkg/errors"
	"github.com/r2d4/kmake/pkg/kmake/docker"
	"github.com/sirupsen/logrus"
)

func Watch(imageName, dockerfile string) error {
	logrus.Infof("Starting watch on image:%s dockerfile:%s", imageName, dockerfile)
	res, err := docker.ParseDockerfile(dockerfile)
	if err != nil {
		return errors.Wrapf(err, "parsing %s", dockerfile)
	}
	deps := docker.GetDependencies(res)
	logrus.Infof("Found dockerfile dependencies: %s", deps)

	return nil
}
