package cmd

import (
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/aestek/configfs/pkg"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(mountCmd)
	mountCmd.SetArgs([]string{"source", "destination"})
	mountCmd.Flags().StringP("project", "p", "", "Your project name")
}

var mountCmd = &cobra.Command{
	Use: "mount",
	Run: func(cmd *cobra.Command, args []string) {
		cfgDir, err := cmd.Flags().GetString("cfg")
		if err != nil {
			log.Fatal(err)
		}
		cfgDir = configDir(cfgDir)

		if len(args) != 2 {
			log.Fatal("You must provide a source (your template file) and a destination (where the config file is mounted)")
		}

		project, err := cmd.Flags().GetString("project")
		if err != nil {
			log.Fatal(err)
		}

		config := configfs.NewConfigManager(filepath.Join(cfgDir, ".config")).Load
		configProvider := configfs.NewTomlProvider(cfgDir)
		envManager := configfs.NewEnv(config).Env
		generator := configfs.NewGenerator(configProvider, envManager)

		closeFs, errs := configfs.MountFS(args[1], func() ([]byte, error) {
			tmpl, err := ioutil.ReadFile(args[0])
			if err != nil {
				return nil, err
			}

			out, err := generator.Gen(project, tmpl)
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
	},
}
