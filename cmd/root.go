package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/oxodao/mqtt2pulseaudio/messages"
	"github.com/oxodao/mqtt2pulseaudio/services"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pa2mqtt",
	Short: "Pulseaudio2mqtt",
	Long:  `Bridge pulseaudio with your nodered / homeassistant setup`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := services.Load(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err := messages.Subscribe()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		services.GET.EventManager.Go("start")

		for {
			time.Sleep(10 * time.Second)
		}
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
