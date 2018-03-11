package env

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Envs map[string]map[string]string

func (e Envs) Get(env, k string) string {
	a, ok := e[env]
	if !ok {
		return ""
	}

	return a[k]
}

func (e Envs) Write(path string) error {
	raw, _ := json.Marshal(e)
	err := ioutil.WriteFile(path, raw, 0644)
	if err != nil {
		return err
	}

	return nil
}

func Load(path string) (Envs, error) {
	ct, err := ioutil.ReadFile(path)
	if os.IsNotExist(err) {
		return make(Envs), nil
	}
	if err != nil {
		return nil, err
	}

	envs := Envs{}
	err = json.Unmarshal(ct, &envs)
	if err != nil {
		return nil, err
	}
	return envs, nil
}
