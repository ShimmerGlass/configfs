package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"os/user"

	"path/filepath"

	"io/ioutil"

	_ "bazil.org/fuse/fs/fstestutil"
)

func main() {
	source := flag.String("source", "local.toml.dist", "Config template")
	dest := flag.String("dest", "local.toml", "Config file")
	project := flag.String("project", "", "Project name")
	flag.Parse()

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	path := filepath.Join(usr.HomeDir, ".configfs")

	config := NewConfigManager(filepath.Join(path, ".config")).Load
	configProvider := NewTomlProvider(filepath.Join(usr.HomeDir, ".configfs"))
	envManager := NewEnv(config).Env
	generator := NewGenerator(configProvider, envManager, config)

	closeFs, errs := MountFS(*dest, func() ([]byte, error) {
		tmpl, err := ioutil.ReadFile(*source)
		if err != nil {
			return nil, err
		}

		out, err := generator.Gen(*project, tmpl)
		if err != nil {
			return nil, err
		}

		return out, err
	})

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	go func() {
		<-sigChan
		closeFs()
		os.Exit(0)
	}()

	for err := range errs {
		log.Fatal(err)
	}
}
