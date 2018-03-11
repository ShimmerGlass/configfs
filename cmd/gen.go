package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(genCmd)
	genCmd.Flags().StringP("project", "p", "", "Your project name")
}

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate a config file without mounting it",
	Run: func(cmd *cobra.Command, args []string) {
		/*cfgDir, err := cmd.Flags().GetString("cfg")
		if err != nil {
			log.Fatal(err)
		}
		cfgDir = configDir(cfgDir)

		if len(args) != 1 {
			log.Fatal("You must provide a source (your template file)")
		}

		project, err := cmd.Flags().GetString("project")
		if err != nil {
			log.Fatal(err)
		}

		config := configfs.NewConfigManager(filepath.Join(cfgDir, ".config")).Load
		configProvider := configfs.NewTomlProvider(cfgDir)
		envManager := configfs.NewEnv(config).Env
		generator := configfs.NewGenerator(configProvider, envManager)

		tmpl, err := ioutil.ReadFile(args[0])
		if err != nil {
			log.Fatal(err)
		}

		out, err := generator.Gen(project, tmpl)
		if err != nil {
			log.Fatal(err)
		}

		os.Stdout.Write(out)*/
	},
}
