/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/fitan/genx/gen"
	"github.com/fitan/genx/plugs/log"
	"github.com/fitan/genx/plugs/trace"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
	"os"
)

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "g",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		x, err := gen.NewX("./")
		if err != nil {
			slog.Error("new x error", err)
			os.Exit(1)
		}

		x.RegImpl(&log.Plug{})
		x.RegImpl(&trace.Plug{})
		x.Gen()
	},
}

func init() {
	rootCmd.AddCommand(genCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// genCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// genCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
