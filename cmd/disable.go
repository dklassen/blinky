package cmd

import (
	"fmt"
	"os"

	hid "github.com/dklassen/blinky/hidapi"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(disableCmd)
}

var disableCmd = &cobra.Command{
	Use:   "disable",
	Short: "Turn off ErgoDoz EZ LEDs",
	Long:  `Turn off the ErgoDox EZ LEDs`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		var dev *hid.DeviceInfo
		var keyboard *hid.Device

		devices, err := hid.Enumerate(0x0, 0x0)
		if err != nil {
			fmt.Println("Something wrong with devices:", err)
			os.Exit(1)
		}

		for _, device := range devices {
			if device.ProductString == "ErgoDox EZ" {
				dev = &device
				break
			}
		}

		if dev == nil {
			fmt.Println("Unable to find ErgoDox EZ board")
			os.Exit(1)
		}

		keyboard, err = dev.Open()
		if err != nil {
			fmt.Println("Unable to open connection to ErgoDox EZ board")
			os.Exit(1)
		}

		keyboard.Write([]byte{0, hid.RGBLightDisable})
	},
}
