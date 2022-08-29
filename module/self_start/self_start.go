package self_start

import (
	"errors"
	"fmt"
	"fyne.io/systray"
	logs "github.com/danbai225/go-logs"
	"github.com/danbai225/tipbar/core"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
)

var selfStart *core.Module
var conf = selfStartConfig{}

type selfStartConfig struct {
	Enable bool `json:"enable"`
}

func ExportModule() *core.Module {
	selfStart = core.NewModule("self_start", "自启动", "", onReady, nil, nil)
	return selfStart
}

func onReady(item *systray.MenuItem) {
	err := selfStart.UnmarshalConfig(&conf)
	if err != nil {
		logs.Err(err)
	}
	if conf.Enable {
		item.Check()
	}
	for {
		select {
		case <-item.ClickedCh:
			if item.Checked() {
				item.Uncheck()
				conf.Enable = false
			} else {
				item.Check()
				conf.Enable = true
			}
			selfStart.SaveConfig(conf)
			err := start(item.Checked())
			if err != nil {
				logs.Err(err)
			}
		}
	}
}

const (
	macListFile = `
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>Label</key>
	<string>cn.p00q.tipbar</string>
	<key>ProgramArguments</key>
        <array>
        	<string>/Applications/DownTip.app/Contents/MacOS/down_tip</string>
        </array>
	<key>RunAtLoad</key>
	<true/>
	<key>WorkingDirectory</key>
	<string>/Applications/DownTip.app/Contents/MacOS</string>
	<key>StandardErrorPath</key>
	<string>/tmp/down_tip.err</string>
	<key>StandardOutPath</key>
	<string>/tmp/down_tip.out</string>
</dict>
</plist>


`
	winBat = `start %s`
)

func start(on bool) error {
	var err error
	var path, content string
	current, err := user.Current()
	if err != nil {
		return err
	}
	switch runtime.GOOS {
	case "darwin":
		path = fmt.Sprintf("%s/Library/LaunchAgents/cn.p00q.tipbar.plist", current.HomeDir)
		content = macListFile
	//case "linux":
	case "windows":
		path = fmt.Sprintf("%s\\AppData\\Roaming\\Microsoft\\Windows\\Start Menu\\Programs\\Startup\\down_tip.bat", current.HomeDir)
		abs, _ := filepath.Abs(os.Args[0])
		content = fmt.Sprintf(winBat, abs)
	default:
		return errors.New("不支持的系统")
	}
	return writer(on, path, content)
}
func writer(on bool, path, content string) error {
	var err error
	if on {
		stat, _ := os.Stat(path)
		if stat == nil {
			err = os.WriteFile(path, []byte(content), os.ModePerm)
		}
	} else {
		err = os.Remove(path)
	}
	return err
}
