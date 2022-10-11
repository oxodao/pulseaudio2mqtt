package cmd

import (
	"fmt"

	"github.com/oxodao/mqtt2pulseaudio/utils"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use: "version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%v %v (Commit %v) by %v\n", utils.SOFTWARE_NAME, utils.VERSION, utils.COMMIT, utils.AUTHOR)
	},
}
