/*
Copyright © 2022 BoiseITGuru.find @Emerald City DAO

*/
package cmdcli

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/bjartek/overflow/overflow"
	"github.com/spf13/cobra"
)

// getConfigCmd represents the getConfig command
var getConfigCmd = &cobra.Command{
	Use:   "getConfig",
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

		configJSON, err := json.MarshalIndent(o.State.Config(), "", " ")
		if err != nil {
			log.Fatal(err.Error())
		}

		fmt.Println(string(configJSON))
	},
}

func init() {
	emulatorCmd.AddCommand(getConfigCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getConfigCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getConfigCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
