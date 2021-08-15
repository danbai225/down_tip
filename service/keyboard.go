package service

import (
	logs "github.com/danbai225/go-logs"
	hook "github.com/robotn/gohook"
)

//键码
//http://www.atoolbox.net/Tool.php?Id=815

var keyLog map[byte]uint64

func monitorInput() {
	keyLog = make(map[byte]uint64)
	EvChan := hook.Start()
	defer hook.End()
	for ev := range EvChan {
		if ev.Kind == hook.KeyDown {
			if _, has := keyLog[byte(ev.Keychar)]; !has {
				logs.Info(ev.String())
				keyLog[byte(ev.Keychar)] = 1
			}
			keyLog[byte(ev.Keychar)]++
		}
	}
}
func GetKeyLog() interface{} {
	type Key struct {
		KeyCode byte
		Val     uint64
	}
	keys := make([]Key, 0)
	for b, u := range keyLog {
		keys = append(keys, Key{KeyCode: b, Val: u})
	}
	return keys
}
