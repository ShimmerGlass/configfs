package configfs

import (
	"bytes"
	"fmt"
	"strings"
)

type Generator struct {
	provider Provider
	env      EnvFn
}

func NewGenerator(provider Provider, env EnvFn) *Generator {
	return &Generator{
		provider: provider,
		env:      env,
	}
}

func (g *Generator) Gen(project string, in []byte) ([]byte, error) {
	keys, err := g.provider.List()
	if err != nil {
		return nil, err
	}

	out := in
	for _, k := range keys {
		env, err := g.env(k.Name)
		if err != nil {
			return nil, err
		}

		v, err := g.provider.Value(k.Name, project, env)
		if err != nil {
			return nil, err
		}

		out = bytes.Replace(out, []byte(k.Name), []byte(v), 1)
	}

	return out, nil
}

func (g *Generator) toLocalized(k, project string) string {
	return fmt.Sprintf("%s_%s", k, strings.ToUpper(
		strings.Replace(project, "-", "_", -1),
	))
}
