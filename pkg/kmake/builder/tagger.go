package builder

import (
	"os/exec"
	"strconv"
	"strings"
	"time"

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

type TimeStampTagger struct {
}

func (t *TimeStampTagger) Tag() (string, error) {
	n := time.Now()
	return strconv.FormatInt(n.Unix(), 10), nil
}
