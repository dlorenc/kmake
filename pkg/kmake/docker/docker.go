package docker

import (
	"log"
	"os"

	"github.com/moby/moby/builder/dockerfile/parser"
	"github.com/pkg/errors"
)

const (
	add  = "add"
	copy = "copy"
)

func ParseDockerfile(path string) (*parser.Result, error) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	res, err := parser.Parse(f)
	if err != nil {
		return nil, errors.Wrap(err, "parsing dockerfile")
	}
	return res, nil
}

func GetDependencies(res *parser.Result) []string {
	deps := []string{}
	for _, value := range res.AST.Children {
		switch value.Value {
		case add, copy:
			src := value.Next.Value
			deps = append(deps, src)
		}
	}
	return deps
}
