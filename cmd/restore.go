package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)

}

var versionCmd = &cobra.Command{
	Use:   "restore",
	Short: "force sync the mongodb to mysql",
	Run:   restore,
}

func restore(cmd *cobra.Command, args []string) {
	fmt.Println("Hugo Static Site Generator v0.9 -- HEAD")
}
