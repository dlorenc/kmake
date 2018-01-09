package cmd

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
	"github.com/r2d4/kmake/pkg/kmake/builder"
	"github.com/r2d4/kmake/pkg/kmake/config"
	"github.com/r2d4/kmake/pkg/kmake/updater"
	"github.com/r2d4/kmake/pkg/kmake/watch"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	projectId string
	remote    bool
)

type WatchSrv struct {
	watch  func(string, string) error
	build  func(string, string, string, builder.Tagger) (string, error)
	update func(string, []string, []config.Artifact) error
}

//TODO(@r2d4): make these interfaces and configurable
var defaultWatcher = WatchSrv{
	watch:  watch.Watch,
	build:  builder.LocalBuild,
	update: updater.KubectlUpdater,
}

func NewCmdWatch(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "watch",
		Short: "watches an dockerfile and its dependencies",
		Run: func(cmd *cobra.Command, args []string) {
			if err := RunWatch(out, cmd); err != nil {
				logrus.Errorf("watch: %s", err)
			}
		},
	}
	cmd.Flags().StringVar(&configFile, "config-file", "kmake.yaml", "the path of the config file to build")
	cmd.Flags().BoolVar(&remote, "remote", false, "local or remote builds")
	cmd.Flags().StringVar(&projectId, "project-id", "", "project id to use for cloud builds")
	return cmd
}

func RunWatch(out io.Writer, cmd *cobra.Command) error {
	if remote {
		defaultWatcher.build = builder.RemoteBuild
	}

	cfg, err := config.Parse(configFile)
	fmt.Println(cfg)
	if err != nil {
		return err
	}

	tagger := &builder.TimeStampTagger{}
	done := make(chan bool)

	for _, artifact := range cfg.Artifacts {
		go func(a config.Artifact) error {
			for {
				if err := defaultWatcher.watch(a.ImageName, a.DockerfilePath); err != nil {
					return errors.Wrap(err, "watch")
				}

				digest, err := defaultWatcher.build(a.ImageName, a.DockerfilePath, projectId, tagger)
				if err != nil {
					return errors.Wrap(err, "build")
				}

				if err := defaultWatcher.update(digest, cfg.Manifests, []config.Artifact{a}); err != nil {
					return errors.Wrap(err, "update")
				}
			}
		}(artifact)
	}

	<-done

	return nil
}
