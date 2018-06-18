package fs

import (
	"encoding/base64"
)

func Name(path string) string {
	return base64.URLEncoding.EncodeToString([]byte(path))
}
