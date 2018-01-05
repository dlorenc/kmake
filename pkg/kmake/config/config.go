package config

import (
	"io/ioutil"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

type KmakeConfig struct {
	Artifacts []Artifact
	Manifests []string
}

type Artifact struct {
	DockerfilePath    string `yaml:"dockerfilePath"`
	DockerContextPath string `yaml:"dockerContextPath"`
	ImageName         string `yaml:"imageName"`
}

func Parse(path string) (*KmakeConfig, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg KmakeConfig
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return nil, err
	}
	baseDir := filepath.Dir(path)

	makeAbs := func(p string) string {
		if filepath.IsAbs(p) {
			return p
		}
		return filepath.Join(baseDir, p)
	}

	for i := range cfg.Artifacts {
		cfg.Artifacts[i].DockerfilePath = makeAbs(cfg.Artifacts[i].DockerfilePath)
		cfg.Artifacts[i].DockerContextPath = makeAbs(cfg.Artifacts[i].DockerContextPath)
	}

	for i := range cfg.Manifests {
		cfg.Manifests[i] = makeAbs(cfg.Manifests[i])
	}

	return &cfg, nil
}
