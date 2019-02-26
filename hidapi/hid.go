package hid

/*
#cgo pkg-config: hidapi

#include <stddef.h>
#include <stdlib.h>
#include <stdio.h>
#include <hidapi/hidapi.h>
#include "hidapi.h"
*/
import "C"
import (
	"errors"
	"fmt"
	"sync"
	"unsafe"

	wchar "github.com/vitaminwater/cgo.wchar"
)

const (
	Version = byte(iota + 1)
	RGBLightEnable
	RGBLightDisable
	RGBLightToggle
	RGBLightSetMode
	RGBLightGetMode
	RGBLightSetHSV
)

var enumerateLock sync.Mutex

type DeviceInfo struct {
	Path          string
	ProductString string
	VendorID      uint16
	ProductID     uint16
}

type Device struct {
	DeviceInfo
	hid_device *C.hid_device
	lock       sync.Mutex
}

func Enumerate(vendorID uint16, productID uint16) ([]DeviceInfo, error) {
	var err error
	var infos []DeviceInfo

	enumerateLock.Lock()
	defer enumerateLock.Unlock()

	first := C.hid_enumerate(C.ushort(vendorID), C.ushort(productID))
	if first == nil {
		return nil, fmt.Errorf("hidapi: no devices found")
	}

	defer C.hid_free_enumeration(first)

	for next := first; next != nil; next = next.next {
		deviceInfo := DeviceInfo{
			Path:      C.GoString(next.path),
			VendorID:  uint16(next.vendor_id),
			ProductID: uint16(next.product_id),
		}

		deviceInfo.ProductString, err = wchar.WcharStringPtrToGoString(unsafe.Pointer(next.product_string))
		if err != nil {
			return nil, fmt.Errorf("Could not convert *C.wchar_t product_string from hid_device_info to go string. Error: %s\n", err)
		}

		infos = append(infos, deviceInfo)
	}

	return infos, err
}

func (info DeviceInfo) Open() (*Device, error) {
	enumerateLock.Lock()
	defer enumerateLock.Unlock()

	path := C.CString(info.Path)
	defer C.free(unsafe.Pointer(path))

	device := C.hid_open_path(path)
	if device == nil {
		return nil, errors.New("hidapi: failed to open device")
	}
	return &Device{
		DeviceInfo: info,
		hid_device: device,
	}, nil
}

func (dev *Device) Write(b []byte) (int, error) {
	if len(b) == 0 {
		return 0, nil
	}

	dev.lock.Lock()
	device := dev.hid_device
	dev.lock.Unlock()

	if device == nil {
		return 0, errors.New("hidapi: device closed")
	}

	written := int(C.hid_write(device, (*C.uchar)(&b[0]), C.size_t(len(b))))
	if written == -1 {
		dev.lock.Lock()
		device = dev.hid_device
		dev.lock.Unlock()

		if device == nil {
			return 0, errors.New("hidapi: device closed")
		}

		message := C.hid_error(device)
		if message == nil {
			return 0, errors.New("hidapi: unknown failure")
		}

		failure, _ := wchar.WcharStringPtrToGoString(unsafe.Pointer(message))
		return 0, errors.New("hidapi: " + failure)
	}
	return written, nil
}

func (dev *Device) ReadTimeout(b []byte, timeout int) (err error) {
	if len(b) == 0 {
		return nil
	}

	res := C.hid_read_timeout(dev.hid_device, (*C.uchar)(&b[0]), C.size_t(len(b)), C.int(timeout))
	resInt := int(res)
	if resInt == -1 {
		return dev.lastError()
	}

	return nil
}

func (dev *Device) lastError() error {
	return errors.New(dev.lastErrorString())
}

func (dev *Device) lastErrorString() string {
	wcharPtr := C.hid_error(dev.hid_device)
	str, err := wchar.WcharStringPtrToGoString(unsafe.Pointer(wcharPtr))
	if err != nil {
		return fmt.Sprintf("Error retrieving error string: %s", err)
	}
	return str
}
