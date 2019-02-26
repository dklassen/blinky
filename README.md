## Blinky

> I did it all for the keebs.


Thanks to Tenderlove for doing all the work -> https://github.com/tenderlove/qmk_firmware/commit/a6be80b03ecafa1e26db073a08239ff03d2889c9. Full credit goes to them I just ported to GO.

Blinky is a cmd relay using the Human Interface Device (HID) protocol to communicate with an ErgoDox EZ keyboard and control the LED's. To get this up and running you'll have to fiddle with your [firmware](github.com/dklassen/qmk_firmware) and install a library.

## Firmware Changes

The ErgoDox EZ stock firmware doesn't come enabled to listen for HID commands. We will need to modify the firmware to process any commands from the relay.

The file responsible for that is [here](https://github.com/dklassen/qmk_firmware/blob/8968344241f6c0e248d4b2c35d767aaba0accac2/keyboards/ergodox_ez/keymaps/dklassen/rgb_hid_protocol.c). To enable this we need to make some adjustments to the `keyboards/ergodox_ez/rules.mk` file:

```c++
RGB_HID_ENABLE = yes

ifeq ($(strip $(RGB_HID_ENABLE)), yes)
	RAW_ENABLE = yes
  SRC += rgb_hid_protocol.c
endif
```

Here we enable the the HID stuff by setting `RGB_HID_ENABLE` and `RAW_ENABLE` to `yes`. This has the wonderful effect of enabling some code called `raw_hid_task` and `raw_hid_receive` check out [here](https://github.com/qmk/qmk_firmware/blob/ae79b60e6bfc03c7fa84076e508f0f2241f087c2/tmk_core/protocol/lufa/lufa.c) for the details.

`raw_hid_receive` is exactly what we end up overriding! We
ve written our own function which on receiving some data messes around with the rgb functionality of the keyboard.

## System Libraries

This package requires that the `hidapi` package be installed. I have only really made this work for Mac where installation is acomplished by:

```
$ brew install hidapi
```
