package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/aestek/configfs/internal/fs"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(fsPathCmd)
}

var fsPathCmd = &cobra.Command{
	Use:   "fspath",
	Short: "",
	Run: func(cmd *cobra.Command, args []string) {
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(fs.Name(dir))
	},
}
