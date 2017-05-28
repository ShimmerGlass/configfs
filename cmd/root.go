package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.PersistentFlags().StringP("cfg", "", "", "Path to configfs config dir. Default to $HOME/.configfs")
}

var RootCmd = &cobra.Command{
	Use: "",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}
