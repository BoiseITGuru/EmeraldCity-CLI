/*
Copyright © 2022 BoiseITGuru.find @Emerald City DAO

*/
package cmdcli

import (
	"fmt"

	"github.com/bjartek/overflow/overflow"
	"github.com/spf13/cobra"
)

// listAccountsCmd represents the listAccounts command
var listAccountsCmd = &cobra.Command{
	Use:   "listAccounts",
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

		o := overflowConfig.Start()

		fmt.Println(o.State.Accounts())
	},
}

func init() {
	emulatorCmd.AddCommand(listAccountsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listAccountsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listAccountsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
