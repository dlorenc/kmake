package cmd

import (
	"io"

	"github.com/Sirupsen/logrus"
	"github.com/r2d4/kmake/pkg/kmake/builder"
	"github.com/r2d4/kmake/pkg/kmake/updater"
	"github.com/r2d4/kmake/pkg/kmake/watch"
	"github.com/spf13/cobra"
)

var (
	imageName      string
	dockerfilePath string
)

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
	return cmd
}

func RunWatch(out io.Writer, cmd *cobra.Command) error {
	if err := watch.Watch(imageName, dockerfilePath); err != nil {
		return err
	}

	if err := builder.Build(imageName, dockerfilePath); err != nil {
		return err
	}

	if err := updater.Update(imageName); err != nil {
		return err
	}
	return nil
}
