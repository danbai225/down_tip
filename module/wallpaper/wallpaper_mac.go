// Package wallpaper
//go:build darwin

package wallpaper

import (
	"fmt"
	"os/exec"
)

func SetWallpaper(filename string) error {
	command := exec.Command(`osascript`, `-e`, fmt.Sprintf("tell application \"System Events\" to tell every desktop to set picture to \"%s\"", filename))
	return command.Run()
}
