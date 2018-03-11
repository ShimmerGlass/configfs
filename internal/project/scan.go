package project

import (
	"path"
	"path/filepath"
)

func Scan(glob string) ([]*Project, error) {
	files, err := filepath.Glob(glob)
	if err != nil {
		return nil, err
	}

	res := []*Project{}
	for _, m := range files {
		p, err := Load(path.Dir(m))
		if err != nil {
			return nil, err
		}
		res = append(res, p)
	}

	return res, nil
}
