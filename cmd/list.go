package cmd

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"sort"

	"github.com/aestek/configfs/pkg"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var forceEnv string

func init() {
	RootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVar(&forceEnv, "env", "", "Force env")
}

var listCmd = &cobra.Command{
	Use: "list",
	Run: func(cmd *cobra.Command, args []string) {
		cfgDir, err := cmd.Flags().GetString("cfg")
		if err != nil {
			log.Fatal(err)
		}
		cfgDir = configDir(cfgDir)

		config := configfs.NewConfigManager(filepath.Join(cfgDir, ".config")).Load
		envManager := configfs.NewEnv(config).Env
		provider := configfs.NewTomlProvider(configDir(cfgDir))
		entries, err := provider.List()
		if err != nil {
			log.Fatal(err)
		}

		sort.Slice(entries, func(i, j int) bool {
			return strings.Compare(entries[i].Name, entries[j].Name) < 0
		})

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "Project", "Env", "Value"})

		for _, entry := range entries {
			env := forceEnv
			if env == "" {
				env, err = envManager(entry.Name)
				if err != nil {
					log.Fatal(err)
				}
			}

			value, err := provider.Value(entry.Name, entry.Project, env)
			if err != nil {
				log.Fatal(err)
			}

			table.Append([]string{
				entry.Name,
				entry.Project,
				env,
				value,
			})
		}

		table.Render()
	},
}
