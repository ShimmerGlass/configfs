package cmd

import (
	"log"
	"os/user"
	"path/filepath"
)

func configDir(provided string) string {
	if provided != "" {
		return provided
	}

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	return filepath.Join(usr.HomeDir, ".configfs")
}
