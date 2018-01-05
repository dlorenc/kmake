package cmd

import (
	"io"

	"github.com/pkg/errors"
	"github.com/r2d4/kmake/pkg/kmake/builder"
	"github.com/r2d4/kmake/pkg/kmake/updater"
	"github.com/r2d4/kmake/pkg/kmake/watch"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	imageName      string
	dockerfilePath string
	projectId      string
	remote         bool
)

type WatchSrv struct {
	watch  func(string, string) error
	build  func(string, string, string, builder.Tagger) (string, error)
	update func(string) error
}

//TODO(@r2d4): make these interfaces and configurable
var defaultWatcher = WatchSrv{
	watch:  watch.Watch,
	build:  builder.LocalBuild,
	update: updater.KsonnetUpdater,
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
	cmd.Flags().StringVar(&imageName, "image-name", "", "the name of the image to build")
	cmd.Flags().StringVar(&dockerfilePath, "dockerfile", "Dockerfile", "the relative path to the dockerfile")
	cmd.Flags().BoolVar(&remote, "remote", false, "local or remote builds")
	cmd.Flags().StringVar(&projectId, "project-id", "", "project id to use for cloud builds")
	return cmd
}

func RunWatch(out io.Writer, cmd *cobra.Command) error {
	if remote {
		defaultWatcher.build = builder.RemoteBuild
	}

	tagger := &builder.CommitTagger{}

	for {
		if err := defaultWatcher.watch(imageName, dockerfilePath); err != nil {
			return errors.Wrap(err, "watch")
		}

		digest, err := defaultWatcher.build(imageName, dockerfilePath, projectId, tagger)
		if err != nil {
			return errors.Wrap(err, "build")
		}

		if err := defaultWatcher.update(digest); err != nil {
			return errors.Wrap(err, "update")
		}
	}

	return nil
}
