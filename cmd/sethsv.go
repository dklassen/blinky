package cmd

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/dklassen/blinky/ergodox"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(sethsvCmd)
}

var sethsvCmd = &cobra.Command{
	Use:   "sethsv",
	Short: "Set RGB value on the ErgoDoz EZ LEDs",
	Long:  `Set RGB value on the ErgoDox EZ LEDs`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 3 {
			return errors.New("Needs 3 arguments for hue, saturation, and value")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		keyboard, err := ergodox.SetupErgodoxEZ()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		hue_arg, err := strconv.ParseInt(args[0], 10, 32)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		saturation_arg, err := strconv.ParseInt(args[1], 10, 32)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		value_arg, err := strconv.ParseInt(args[2], 10, 32)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		hue := uint8(hue_arg)
		saturation := uint8(saturation_arg)
		value := uint8(value_arg)

		_, err = keyboard.SetHSV(hue, saturation, value)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
