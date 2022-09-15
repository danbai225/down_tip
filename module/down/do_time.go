package down

import (
	"fmt"
	"github.com/ncruces/zenity"
	"strconv"
	"time"
)

var orderMeal = false

func getTheRemainingTime() string {
	sub := conf.DownTimeToday.Sub(time.Now())
	h := uint(sub.Hours())
	tm := sub - time.Duration(h)*time.Hour
	m := uint(tm.Minutes())
	s := uint((tm - time.Duration(m)*time.Minute).Seconds())
	if s < 0 || s > 60 {
		return " no Time"
	}
	h1, m1, err := analysisTime(conf.OrderMeal)
	if err == nil {
		now := time.Now()
		if now.Hour() == int(h1) && now.Minute() == int(m1) {
			if !orderMeal {
				orderMeal = true
				_ = zenity.Notify("该点饭啦！！！")
			}
		} else {
			orderMeal = false
		}
	}
	return fmt.Sprintf("%s:%s:%s", complementZero(h), complementZero(m), complementZero(s))
}
func resetTime() {
	conf.StartingTimeToday = time.Now()
	calculateDownTime()
}
func getStartingTimeToday() string {
	return "今日时间:" + conf.StartingTimeToday.Format("15:04:05")
}
func addStartingTimeToday(duration time.Duration) string {
	conf.StartingTimeToday = conf.StartingTimeToday.Add(duration)
	calculateDownTime()
	return getStartingTimeToday()
}
func complementZero(i uint) string {
	itoa := strconv.Itoa(int(i))
	if i < 10 {
		return "0" + itoa
	}
	return itoa
}
