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
			manifest = strings.Replace(manifest, replaceableName, fqImageName, 0)
		}

		cmd := exec.Command("kubectl", "apply", "-f", "-")
		r := strings.NewReader(manifest)
		stdout, stderr, err := util.RunCommand(cmd, r)
		logrus.Infof("Applying manifest: %s, %s, %s", mp, stdout, stderr)
		if err != nil {
			return err
		}
	}
	return nil
}
