package cmd

import (
	"io"

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
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunWatch(out, cmd)
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
	return nil
}
