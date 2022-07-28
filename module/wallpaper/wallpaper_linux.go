//go:build linux

package wallpaper

import (
	"fmt"
	"os/exec"
)

//Ubuntu
func SetWallpaper(filename string) error {
	return exec.Command("gsettings", "set", "org.gnome.desktop.background", "picture-uri", fmt.Sprintf("file:%s", filename)).Run()
}
