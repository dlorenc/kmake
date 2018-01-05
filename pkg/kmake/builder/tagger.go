package builder

import (
	"os/exec"
	"strings"

	"github.com/pkg/errors"
	"github.com/r2d4/kmake/pkg/kmake/util"
)

type Tagger interface {
	Tag() (string, error)
}

type CommitTagger struct {
}

func (c *CommitTagger) Tag() (string, error) {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	stdout, stderr, err := util.RunCommand(cmd, nil)
	if err != nil {
		return "", errors.Wrapf(err, "tag: %s %s", stdout, stderr)
	}

	return strings.TrimSuffix(string(stdout), "\n"), nil
}
