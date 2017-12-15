package watch

import (
	"path"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/r2d4/kmake/pkg/kmake/docker"
	"github.com/rjeczalik/notify"
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

	// TODO(@r2d4): make context configurable
	dockerCtx := filepath.Dir(dockerfile)

	c := make(chan notify.EventInfo, 1)
	defer notify.Stop(c)
	for _, dep := range deps {
		dep := dep
		go watchFile(path.Join(dockerCtx, dep), c)
	}
	ei := <-c
	logrus.Infof("event: %s", ei)

	return nil
}

func watchFile(path string, c chan notify.EventInfo) error {
	logrus.Infof("Starting watch on %s", path)
	if err := notify.Watch(path, c, notify.All); err != nil {
		return errors.Wrap(err, "notify")
	}
	return nil
}
