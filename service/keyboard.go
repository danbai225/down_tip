package service

import (
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
				keyLog[byte(ev.Keychar)] = 1
			}
			keyLog[byte(ev.Keychar)]++
		}
	}
}
func GetKeyLog() map[byte]uint64 {
	return keyLog
}
