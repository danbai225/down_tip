package service

import (
	"down_tip/config"
	"down_tip/keyboard"
	"down_tip/utils"
	"fmt"
	logs "github.com/danbai225/go-logs"
	"log"
	"time"
)

func GetTheRemainingTime() string {
	sub := config.Conf.DownTimeToday.Sub(time.Now())
	h := uint(sub.Hours())
	tm := sub - time.Duration(h)*time.Hour
	m := uint(tm.Minutes())
	s := uint((tm - time.Duration(m)*time.Minute).Seconds())
	return fmt.Sprintf("%s:%s:%s", utils.ComplementZero(h), utils.ComplementZero(m), utils.ComplementZero(s))
}
func ResetTime() {
	config.Conf.StartingTimeToday = time.Now()
	config.CalculateDownTime()
}
func GetStartingTimeToday() string {
	return "今日时间:" + config.Conf.StartingTimeToday.Format("15:04:05")
}
func AddStartingTimeToday(duration time.Duration) string {
	config.Conf.StartingTimeToday = config.Conf.StartingTimeToday.Add(duration)
	log.Println(config.Conf.StartingTimeToday.Format("15:04:05"))
	return GetStartingTimeToday()
}
func MonitorReset() {
	m := make(map[uint]struct{})
	for {
		now := time.Now()
		i := now.Hour() - int(config.Conf.GetStartHour())
		if _, has := m[config.Conf.GetStartHour()]; !has {
			if i == -1 {
				m[config.Conf.GetStartHour()] = struct{}{}
				if keyboard.Activity() {
					logs.Info("处于活动状态，重置时间")
					ResetTime()

				}
			}
		}
	}

}
