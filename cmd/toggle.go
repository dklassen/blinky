package cmd

import (
	"fmt"
	"os"

	"github.com/dklassen/blinky/ergodox"
	hid "github.com/dklassen/blinky/hidapi"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(toggleCmd)
}

var toggleCmd = &cobra.Command{
	Use:   "toggle",
	Short: "Toggle the ErgoDoz EZ LED state",
	Long: `Toggle the LED state on the ErgoDox EZ
						from ON to OFF and vice versa`,
	Run: func(cmd *cobra.Command, args []string) {
		keyboard, err := ergodox.SetupErgodoxEZ()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		val, err := keyboard.Write([]byte{hid.RGBLightToggle})
		fmt.Println(val)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
