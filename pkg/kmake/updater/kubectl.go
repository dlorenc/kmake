package updater

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"

	"github.com/r2d4/kmake/pkg/kmake/config"
	"github.com/r2d4/kmake/pkg/kmake/util"
	"github.com/sirupsen/logrus"
)

func KubectlUpdater(tag string, manifests []string, artifacts []config.Artifact) error {

	for _, mp := range manifests {
		manifestContents, err := ioutil.ReadFile(mp)
		if err != nil {
			return err
		}
		manifest := string(manifestContents)
		for _, a := range artifacts {
			replaceableName := fmt.Sprintf("%s:replaceme", a.ImageName)
			fqImageName := fmt.Sprintf("%s:%s", a.ImageName, tag)
			// logrus.Infof("Replacing %s with %s", replaceableName, fqImageName)
			manifest = strings.Replace(manifest, replaceableName, fqImageName, -1)
		}

		cmd := exec.Command("kubectl", "apply", "-f", "-")
		// logrus.Infof("Manifest: %s", manifest)
		r := strings.NewReader(manifest)
		stdout, stderr, err := util.RunCommand(cmd, r)
		logrus.Infof("Applying manifest: %s, %s, %s", mp, stdout, stderr)
		if err != nil {
			return err
		}
	}
	return nil
}
