package main

import (
	"regexp"
)

type EnvFn func(k string) (string, error)

type Env struct {
	configFn ConfigFn
}

func NewEnv(configFn ConfigFn) *Env {
	return &Env{
		configFn: configFn,
	}
}

func (m *Env) Env(k string) (string, error) {
	config, err := m.configFn()
	if err != nil {
		return "", err
	}

	for pattern, env := range config.EnvPatterns {
		r, err := regexp.Compile(pattern)
		if err != nil {
			return "", err
		}

		if r.MatchString(k) {
			return env, nil
		}
	}

	return config.Env, nil
}
