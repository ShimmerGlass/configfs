package main

import (
	"fmt"
	"strings"
)

type Generator struct {
	provider Provider
	env      EnvFn
	configFn ConfigFn
}

func NewGenerator(provider Provider, env EnvFn, configFn ConfigFn) *Generator {
	return &Generator{
		provider: provider,
		env:      env,
		configFn: configFn,
	}
}

func (g *Generator) Gen(project, in string) (string, error) {
	keys, err := g.provider.List()
	if err != nil {
		return "", err
	}

	config, err := g.configFn()
	if err != nil {
		return "", err
	}

	out := in
	for _, k := range keys {
		env, err := g.env(k)
		if err != nil {
			return "", err
		}

		outK := k
		for _, l := range config.Localized {
			if g.toLocalized(l, project) == k {
				outK = l
				break
			}
		}

		v, err := g.provider.Value(k, env)
		if err != nil {
			return "", err
		}

		out = strings.Replace(out, outK, v, 1)
	}

	return out, nil
}

func (g *Generator) toLocalized(k, project string) string {
	return fmt.Sprintf("%s_%s", k, strings.ToUpper(
		strings.Replace(project, "-", "_", -1),
	))
}
