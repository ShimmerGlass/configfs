package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"path/filepath"
)

type TomlProvider struct {
	root string
}

func NewTomlProvider(root string) *TomlProvider {
	return &TomlProvider{root: root}
}

func (p *TomlProvider) List() ([]string, error) {
	files, err := filepath.Glob(filepath.Join(p.root, "/*.toml"))
	if err != nil {
		return nil, err
	}

	var res []string

	for _, f := range files {
		entries := map[string]interface{}{}
		_, err := toml.DecodeFile(f, &entries)
		if err != nil {
			return nil, err
		}

		for k := range entries {
			res = append(res, k)
		}
	}

	return res, nil
}

func (p *TomlProvider) Value(k, project, env string) (string, error) {
	entries := map[string]interface{}{}
	_, err := toml.DecodeFile(filepath.Join(p.root, fmt.Sprintf("%s.toml", env)), &entries)
	if err != nil {
		return "", err
	}

	projectConfig, ok := entries[project]
	if ok {
		projectConfig, ok := projectConfig.(map[string]interface{})
		if ok {
			v, ok := projectConfig[k].(string)
			if ok {
				return v, nil
			}
		}
	}

	v, _ := entries[k].(string)
	return v, nil
}
