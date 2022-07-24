/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmdcli

import (
	"github.com/bjartek/overflow/overflow"
	"github.com/spf13/cobra"
)

// saveConfigCmd represents the saveConfig command
var saveConfigCmd = &cobra.Command{
	Use:   "saveConfig",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var overflowConfig *overflow.OverflowBuilder

		if tempFlowConfig != "" {
			overflowConfig = overflow.NewOverflowBuilder("emulator", false, 0).Config(tempFlowConfig)
		} else {
			overflowConfig = overflow.NewOverflowBuilder("emulator", false, 0)
		}

		overflowConfig.InitializeAccounts = false
		overflowConfig.DeployContracts = false

		o := overflowConfig.Start()

		o.State.SaveEdited([]string{"flow.json"})
	},
}

func init() {
	emulatorCmd.AddCommand(saveConfigCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// saveConfigCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// saveConfigCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
