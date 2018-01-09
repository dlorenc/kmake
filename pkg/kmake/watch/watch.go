package watch

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/r2d4/kmake/pkg/kmake/docker"
	"github.com/rjeczalik/notify"
	"github.com/sirupsen/logrus"
)

func Watch(imageName, dockerfile string) error {
	logrus.Infof("Starting watch on image: %s dockerfile:%s", imageName, dockerfile)
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

	expandedPaths := []string{}
	if strings.Contains(path, "*") {
		paths, err := filepath.Glob(path)
		if err != nil {
			return err
		}
		expandedPaths = append(expandedPaths, paths...)
	} else {
		expandedPaths = append(expandedPaths, path)
	}

	for i := range expandedPaths {
		fi, err := os.Stat(expandedPaths[i])
		if err != nil {
			return err
		}
		if fi.IsDir() {
			expandedPaths[i] = filepath.Join(path, "...")
		}
	}

	for _, wp := range expandedPaths {
		logrus.Infof("Starting watch on %s", wp)
		if err := notify.Watch(wp, c, notify.All); err != nil {
			return errors.Wrap(err, "notify")
		}
	}
	return nil
}
