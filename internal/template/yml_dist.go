package template

import (
	"bytes"
	"regexp"
)

func YamlDistVars(in []byte) (vars []Var) {
	reg := regexp.MustCompile(`\{\{var\.(\w+)\}\}`)
	idx := reg.FindAllIndex(in, -1)
	for _, i := range idx {
		vars = append(vars, Var{
			StartPos: i[0],
			EndPos:   i[1],
			Name:     string(in[i[0]+6 : i[1]-2]),
		})
	}

	return
}

func YamlDistReplace(in []byte, v Var, value string) []byte {
	return bytes.Replace(in, []byte("{{var."+v.Name+"}}"), []byte(value), 1)
}
