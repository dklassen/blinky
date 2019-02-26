package ergodox

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"time"

	hid "github.com/dklassen/blinky/hidapi"
)

const VendorID = 0xfeed
const ProductID = 0x1307

/*
|-----------------|-----------------------------------|
| old mode number | new mode name                     |
|-----------------|-----------------------------------|
|        1        | RGBLIGHT_MODE_STATIC_LIGHT        |
|        2        | RGBLIGHT_MODE_BREATHING           |
|        3        | RGBLIGHT_MODE_BREATHING + 1       |
|        4        | RGBLIGHT_MODE_BREATHING + 2       |
|        5        | RGBLIGHT_MODE_BREATHING + 3       |
|        6        | RGBLIGHT_MODE_RAINBOW_MOOD        |
|        7        | RGBLIGHT_MODE_RAINBOW_MOOD + 1    |
|        8        | RGBLIGHT_MODE_RAINBOW_MOOD + 2    |
|        9        | RGBLIGHT_MODE_RAINBOW_SWIRL       |
|       10        | RGBLIGHT_MODE_RAINBOW_SWIRL + 1   |
|       11        | RGBLIGHT_MODE_RAINBOW_SWIRL + 2   |
|       12        | RGBLIGHT_MODE_RAINBOW_SWIRL + 3   |
|       13        | RGBLIGHT_MODE_RAINBOW_SWIRL + 4   |
|       14        | RGBLIGHT_MODE_RAINBOW_SWIRL + 5   |
|       15        | RGBLIGHT_MODE_SNAKE               |
|       16        | RGBLIGHT_MODE_SNAKE + 1           |
|       17        | RGBLIGHT_MODE_SNAKE + 2           |
|       18        | RGBLIGHT_MODE_SNAKE + 3           |
|       19        | RGBLIGHT_MODE_SNAKE + 4           |
|       20        | RGBLIGHT_MODE_SNAKE + 5           |
|       21        | RGBLIGHT_MODE_KNIGHT              |
|       22        | RGBLIGHT_MODE_KNIGHT + 1          |
|       23        | RGBLIGHT_MODE_KNIGHT + 2          |
|       24        | RGBLIGHT_MODE_CHRISTMAS           |
|       25        | RGBLIGHT_MODE_STATIC_GRADIENT     |
|       26        | RGBLIGHT_MODE_STATIC_GRADIENT + 1 |
|       27        | RGBLIGHT_MODE_STATIC_GRADIENT + 2 |
|       28        | RGBLIGHT_MODE_STATIC_GRADIENT + 3 |
|       29        | RGBLIGHT_MODE_STATIC_GRADIENT + 4 |
|       30        | RGBLIGHT_MODE_STATIC_GRADIENT + 5 |
|       31        | RGBLIGHT_MODE_STATIC_GRADIENT + 6 |
|       32        | RGBLIGHT_MODE_STATIC_GRADIENT + 7 |
|       33        | RGBLIGHT_MODE_STATIC_GRADIENT + 8 |
|       34        | RGBLIGHT_MODE_STATIC_GRADIENT + 9 |
|       35        | RGBLIGHT_MODE_RGB_TEST            |
|       36        | RGBLIGHT_MODE_ALTERNATING         |
|-----------------|-----------------------------------|

*/

var RGBLIGHT_MODE = map[string]byte{
	"static":        1,
	"breathing":     2,
	"rainbow":       6,
	"rainbow_swirl": 9,
	"snake":         15,
	"knight":        21,
	"christmas":     24,
	"gradient":      25,
}

type stop struct {
	error
}

type ErgodoxEZ struct {
	device  *hid.Device
	version int
}

func (keyboard *ErgodoxEZ) SetHSV(hue uint8, saturation uint8, value uint8) ([]byte, error) {
	// hue is represented from 0-360. We have to convert fit within 0xff.
	hue_0_255 := (hue >> 8) & 0xff
	hue_256_360 := hue & 0xff

	data := []byte{
		hid.RGBLightSetHSV,
		hue_0_255,
		hue_256_360,
		saturation,
		value,
	}
	written, err := keyboard.Write(data)
	result := make([]byte, written)
	keyboard.Read(result)
	return result, err
}

func (keyboard *ErgodoxEZ) SetMode(mode string) ([]byte, error) {
	if !IsValidRGBMode(mode) {
		return nil, fmt.Errorf("Invalid RGB mode: %s", mode)
	}
	written, err := keyboard.Write([]byte{hid.RGBLightSetMode, RGBLIGHT_MODE[mode]})
	if err != nil {
		return nil, err
	}

	result := make([]byte, written)
	keyboard.Read(result)
	return result, nil
}

func (keyboard *ErgodoxEZ) Write(data []byte) (int, error) {
	data = append([]byte{0}, data...)
	val, err := keyboard.device.Write(data)
	return val, err
}

func (keyboard *ErgodoxEZ) Read(bytes []byte) error {
	return keyboard.device.ReadTimeout(bytes, 300)
}

func IsValidRGBMode(mode string) bool {
	_, ok := RGBLIGHT_MODE[mode]
	return ok
}

func Find(vendorID uint16, productID uint16) (dev *hid.DeviceInfo, err error) {
	devices, err := hid.Enumerate(vendorID, productID)
	if err != nil {
		return nil, fmt.Errorf("Something went wrong with device lookup:", err)
	}

	for _, device := range devices {
		if device.ProductString == "ErgoDox EZ" {
			dev = &device
			break
		}
	}

	if dev == nil {
		return nil, errors.New("Unable to find ErgodoxEZ board")
	}
	return dev, nil
}

func Open(dev *hid.DeviceInfo, attempts int, sleep time.Duration) (keyboard *hid.Device, err error) {
	if keyboard, err = dev.Open(); err != nil {
		if s, ok := err.(stop); ok {
			return keyboard, s.error
		}

		if attempts--; attempts > 0 {
			time.Sleep(sleep)
			return Open(dev, attempts, 2*sleep)
		}
		return keyboard, err
	}
	return keyboard, nil
}

func SetupErgodoxEZ() (*ErgodoxEZ, error) {
	dev, err := Find(VendorID, ProductID)
	if err != nil {
		return nil, err
	}

	keyboard, err := Open(dev, 10, 0)
	if err != nil {
		return nil, err
	}

	written, err := keyboard.Write([]byte{0, hid.Version})
	read := make([]byte, written)
	keyboard.ReadTimeout(read, 300)
	buf := bytes.NewBuffer(read[:len(read)-1])
	version, _ := binary.ReadVarint(buf)

	return &ErgodoxEZ{keyboard, int(version)}, err
}
