package cmd

import (
	"log"
	"os/user"
	"path"

	"github.com/aestek/configfs/internal/env"
	"github.com/aestek/configfs/internal/fs"
	"github.com/aestek/configfs/internal/project"
	"github.com/aestek/configfs/internal/server"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(startCmd)
	startCmd.Flags().StringP("mount", "m", "/tmp/cfgcfg", "Mount location")
	startCmd.Flags().StringP("listen", "l", ":9567", "Listen addr")
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "",
	Run: func(cmd *cobra.Command, args []string) {
		go func() {
			mp, _ := cmd.Flags().GetString("mount")
			log.Fatal(fs.Run(mp, func(projectPath string) ([]byte, error) {
				usr, err := user.Current()
				if err != nil {
					log.Fatal(err)
				}

				envsPath := path.Join(usr.HomeDir, ".cfgcfg/envs")

				p, err := project.Load(projectPath)
				if err != nil {
					log.Println(err)
					return nil, err
				}

				envs, err := env.Load(envsPath)
				if err != nil {
					log.Println(err)
					return nil, err
				}

				return p.Contents(envs)
			}))
		}()

		listen, _ := cmd.Flags().GetString("listen")
		log.Fatal(server.Start(listen))
	},
}
