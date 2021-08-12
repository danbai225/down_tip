package keyboard

import (
	"fmt"
	hook "github.com/robotn/gohook"
	"time"
)

func MonitorStr(keys []string) bool {
	chekI := 0
	EvChan := hook.Start()
	defer hook.End()
	for ev := range EvChan {
		if ev.Kind == hook.KeyDown {
			fmt.Println("hook: ", ev)
			if string(ev.Keychar) == keys[chekI] {
				chekI++
				if chekI == len(keys) {
					return true
				}
			} else {
				chekI = 0
			}
		}
	}
	return false
}
func Activity() bool {
	EvChan := hook.Start()
	defer hook.End()
	t := time.Now()
	k := false
	m := false
	for ev := range EvChan {
		seconds := time.Now().Sub(t).Seconds()
		if seconds > 1 && seconds < 10 {
			if ev.Kind == hook.KeyDown {
				k = true
			} else if ev.Kind == hook.MouseMove {
				m = true
			}
		} else if seconds > 10 {
			t = time.Now()
			k = false
			m = false
		}
		if k && m {
			return true
		}
	}
	return false
}
