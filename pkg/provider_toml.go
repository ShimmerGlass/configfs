package configfs

import (
	"fmt"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type TomlProvider struct {
	root string
}

func NewTomlProvider(root string) *TomlProvider {
	return &TomlProvider{root: root}
}

func (p *TomlProvider) List() ([]ConfigEntry, error) {
	files, err := filepath.Glob(filepath.Join(p.root, "/*.toml"))
	if err != nil {
		return nil, err
	}

	var res []ConfigEntry

	for _, f := range files {
		entries := map[string]interface{}{}
		_, err := toml.DecodeFile(f, &entries)
		if err != nil {
			return nil, err
		}

		for k, v := range entries {
			pk, ok := v.(map[string]interface{})
			if ok {
				for sk := range pk {
					res = append(res, ConfigEntry{
						Name:    sk,
						Project: k,
					})
				}
			} else {
				res = append(res, ConfigEntry{
					Name: k,
				})
			}
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
			v, ok := projectConfig[k]
			if ok {
				return p.v(v), nil
			}
		}
	}

	v, _ := entries[k]
	return p.v(v), nil
}

func (p *TomlProvider) v(v interface{}) string {
	return fmt.Sprintf("%v", v)
}
