package cmd

import (
	"fmt"
	"log"
	"path/filepath"

	"os"

	"github.com/aestek/configfs/pkg"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(setEnvCmd)
}

var setEnvCmd = &cobra.Command{
	Use: "set-env",
	Run: func(cmd *cobra.Command, args []string) {
		cfgDir, err := cmd.Flags().GetString("cfg")
		if err != nil {
			log.Fatal(err)
		}
		cfgDir = configDir(cfgDir)

		if len(args) != 1 {
			fmt.Println("Please provide an env to set.")
			os.Exit(1)
		}

		configManager := configfs.NewConfigManager(filepath.Join(cfgDir, ".config"))
		config, err := configManager.Load()
		if err != nil {
			log.Fatal(err)
		}

		config.Env = args[0]
		err = configManager.Save(config)
		if err != nil {
			log.Fatal(err)
		}
	},
}
