/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/aceaura/aries/util/settings"
	"github.com/spf13/cobra"
)

// startSchedulerCmd represents the startScheduler command
var startSchedulerCmd = &cobra.Command{
	Use:   "startScheduler",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := settings.Flags.BindPFlags(cmd.Flags()); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(startSchedulerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startSchedulerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startSchedulerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
