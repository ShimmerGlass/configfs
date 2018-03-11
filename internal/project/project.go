package project

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	"github.com/aestek/configfs/internal/env"

	"github.com/aestek/configfs/internal/template"
)

type Var struct {
	Anchors []template.Var `json:"-"`
	Env     string         `json:"env,omitempty"`
	Value   string         `json:"value,omitempty"`
}

type Project struct {
	Path string          `json:"path"`
	Name string          `json:"name"`
	Vars map[string]*Var `json:"vars,omitempty"`
}

func Load(projectPath string) (*Project, error) {
	project := &Project{
		Path: projectPath,
		Name: path.Base(projectPath),
		Vars: make(map[string]*Var),
	}

	fileVars, err := project.templateVars()
	if err != nil {
		return nil, err
	}

	cfg, err := project.cfg()
	if err != nil {
		return nil, err
	}

	for _, v := range fileVars {
		if _, ok := project.Vars[v.Name]; ok {
			project.Vars[v.Name].Anchors = append(project.Vars[v.Name].Anchors, v)
			continue
		}
		vv := cfg.Vars[v.Name]
		if vv == nil {
			vv = &Var{}
		}
		vv.Anchors = []template.Var{v}
		project.Vars[v.Name] = vv
	}

	return project, nil
}

func (p *Project) Contents(envs env.Envs) ([]byte, error) {
	tmpl, err := p.template()
	if err != nil {
		return nil, err
	}

	fileVars, err := p.templateVars()
	if err != nil {
		return nil, err
	}

	for _, fv := range fileVars {
		var value string
		v, ok := p.Vars[fv.Name]
		if ok {
			if v.Env != "" {
				value = envs.Get(v.Env, fv.Name)
			} else {
				value = v.Value
			}
		}

		tmpl = template.YamlDistReplace(tmpl, fv, value)
	}

	return tmpl, nil
}

func (p *Project) Write(envs env.Envs) error {
	cfg, _ := json.Marshal(p)
	err := ioutil.WriteFile(p.cfgPath(), cfg, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (p *Project) distPath() string {
	return path.Join(p.Path, "local.toml.dist")
}

func (p *Project) tomlPath() string {
	return path.Join(p.Path, "local.toml")
}

func (p *Project) cfgPath() string {
	return path.Join(p.Path, ".cfgcfg")
}

func (p *Project) template() ([]byte, error) {
	return ioutil.ReadFile(p.distPath())
}

func (p *Project) templateVars() ([]template.Var, error) {
	tmpl, err := p.template()
	if err != nil {
		return nil, err
	}
	return template.YamlDistVars(tmpl), nil
}

func (p *Project) cfg() (*Project, error) {
	cfgContent, err := ioutil.ReadFile(p.cfgPath())
	if os.IsNotExist(err) {
		return &Project{}, nil
	}
	if err != nil {
		return nil, err
	}

	cfg := &Project{}
	err = json.Unmarshal(cfgContent, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
