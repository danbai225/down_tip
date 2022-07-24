// Package wallpaper
//go:build windows
// +build windows

package wallpaper

import (
	"syscall"
	"unsafe"
)

// UI Params for windows
const (
	spiSetdeskwallpaper         = 0x0014
	uiParam                     = 0x0000
	updateAndSetToRegistryParam = 0x0003
)

// user32.dll and its proc
var (
	user32                = syscall.NewLazyDLL("user32.dll")
	systemParametersInfoW = user32.NewProc("SystemParametersInfoW")
)

func SetWallpaper(filename string) error {
	filenameUTF16Ptr, err := syscall.UTF16PtrFromString(filename)
	if err != nil {
		return err
	}

	systemParametersInfoW.Call(
		uintptr(spiSetdeskwallpaper),              // DLL Message
		uintptr(uiParam),                          // UI Param
		uintptr(unsafe.Pointer(filenameUTF16Ptr)), // User argument e.g. file name
		uintptr(updateAndSetToRegistryParam),      // Param to update the user profile and set this change into registry
	)

	return nil
}
