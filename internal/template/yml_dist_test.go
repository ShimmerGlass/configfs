package template

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestYamlDistVars(t *testing.T) {
	tmpl := `[http]
port = {{var.GLOBAL_SERVER_LISTEN_PORT}}
host = "{{var.GLOBAL_SERVER_LISTEN_ADDR}}"`

	assert.Equal(t,
		[]Var{
			{"GLOBAL_SERVER_LISTEN_PORT", 14, 47},
			{"GLOBAL_SERVER_LISTEN_ADDR", 56, 89},
		},
		YamlDistVars([]byte(tmpl)),
	)
}

func TestYamlDistReplace(t *testing.T) {
	tmpl := `[http]
port = {{var.GLOBAL_SERVER_LISTEN_PORT}}
host = "{{var.GLOBAL_SERVER_LISTEN_ADDR}}"`

	cmpl := YamlDistReplace([]byte(tmpl), Var{"GLOBAL_SERVER_LISTEN_PORT", 14, 47}, "3000")

	assert.Equal(t, []byte(`[http]
port = 3000
host = "{{var.GLOBAL_SERVER_LISTEN_ADDR}}"`), cmpl)

}
