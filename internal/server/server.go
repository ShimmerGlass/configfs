package server

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/user"
	"path"

	"github.com/aestek/configfs/internal/env"
	"github.com/aestek/configfs/internal/project"
	_ "github.com/aestek/configfs/internal/server/statik"
	"github.com/rakyll/statik/fs"
	"github.com/rs/cors"
)

func Start(listen string) error {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	gopath := os.Getenv("GOPATH")

	envsPath := path.Join(usr.HomeDir, ".cfgcfg/envs")

	mux := http.NewServeMux()

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	mux.Handle("/public/", http.StripPrefix("/public/", http.FileServer(statikFS)))

	mux.HandleFunc("/envs", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			envs := env.Envs{}
			err := json.NewDecoder(r.Body).Decode(&envs)
			if err != nil {
				http.Error(rw, err.Error(), 400)
				return
			}

			err = envs.Write(envsPath)
			if err != nil {
				http.Error(rw, err.Error(), 500)
				return
			}
		}

		envs, err := env.Load(envsPath)
		if err != nil {
			http.Error(rw, err.Error(), 500)
			return
		}

		json.NewEncoder(rw).Encode(envs)
	})

	mux.HandleFunc("/projects", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			project := &project.Project{}
			err := json.NewDecoder(r.Body).Decode(&project)
			if err != nil {
				http.Error(rw, err.Error(), 400)
				return
			}

			envs, err := env.Load(envsPath)
			if err != nil {
				http.Error(rw, err.Error(), 500)
				return
			}

			err = project.Write(envs)
			if err != nil {
				http.Error(rw, err.Error(), 500)
				return
			}
		}

		pattern := path.Join(gopath, "src", "*", "*", "*", "local.toml.dist")

		projects, err := project.Scan(pattern)
		if err != nil {
			http.Error(rw, err.Error(), 500)
			return
		}

		json.NewEncoder(rw).Encode(projects)
	})

	handler := cors.Default().Handler(mux)
	return http.ListenAndServe(listen, handler)
}
