package cmd

import (
	"io"

	"github.com/r2d4/kmake/pkg/kmake/builder"
	"github.com/r2d4/kmake/pkg/kmake/config"
	"github.com/r2d4/kmake/pkg/kmake/updater"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	configFile string
)

func NewCmdBuild(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build",
		Short: "builds a set of dockerfiles",
		Run: func(cmd *cobra.Command, args []string) {
			if err := RunBuild(out, cmd); err != nil {
				logrus.Errorf("build: %s", err)
			}
		},
	}
	cmd.Flags().StringVar(&configFile, "config-file", "kmake.yaml", "the path of the config file to build")
	cmd.Flags().BoolVar(&remote, "remote", false, "local or remote builds")
	cmd.Flags().StringVar(&projectId, "project-id", "", "project id to use for cloud builds")
	return cmd
}

func RunBuild(out io.Writer, cmd *cobra.Command) error {
	tagger := &builder.CommitTagger{}
	tag, err := tagger.Tag()
	if err != nil {
		return err
	}
	cfg, err := config.Parse(configFile)
	if err != nil {
		return err
	}

	updater := updater.KubectlUpdater

	for _, a := range cfg.Artifacts {
		img, err := builder.LocalBuild(a.ImageName, a.DockerfilePath, "", tagger)
		if err != nil {
			return err
		}
		if err := builder.Push(img); err != nil {
			return err
		}
		if err := updater(tag, cfg.Manifests, cfg.Artifacts); err != nil {
			return err
		}
	}

	return nil
}
