package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/dklassen/blinky/ergodox"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(modeCmd)
}

var modeCmd = &cobra.Command{
	Use:   "mode",
	Short: "Set RGB mode on the ErgoDoz EZ LEDs",
	Long:  `Set RGB mode on the ErgoDox EZ LEDs`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a specified mode")
		}
		if ergodox.IsValidRGBMode(args[0]) {
			return nil
		}
		return fmt.Errorf("invalid mode specified: %s", args[0])
	},
	Run: func(cmd *cobra.Command, args []string) {
		keyboard, err := ergodox.SetupErgodoxEZ()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		_, err = keyboard.SetMode(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
