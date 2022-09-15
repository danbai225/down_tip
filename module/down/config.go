package down

import (
	"errors"
	"fmt"
	"github.com/danbai225/tipbar/core"
	"strconv"
	"strings"
	"time"
)

// 配置对象
var conf config

// 本地时区
var loc, _ = time.LoadLocation("Local")

// Config 配置文件对象
type config struct {
	StartTime         string `json:"start_time"`
	startHour         uint
	startMinute       uint
	EndTime           string `json:"end_time"`
	endHour           uint
	endMinute         uint
	DownTimeToday     time.Time
	StartingTimeToday time.Time
	Elasticity        bool   `json:"elasticity"`
	OrderMeal         string `json:"order_meal"`
}

func (c *config) getStartHour() uint {
	return c.startHour
}
func (c *config) getStartMinute() uint {
	return c.startMinute
}

// 加载配置
func loadConfiguration() error {
	//解析配置
	core.Unmarshal(down.Config, &conf)
	//根据配置文件中时间解析时分
	h, m, err := analysisTime(conf.StartTime)
	if err != nil {
		return err
	}
	conf.startHour = h
	conf.startMinute = m
	//根据配置文件中时间解析时分
	h, m, err = analysisTime(conf.EndTime)
	if err != nil {
		return err
	}
	conf.endHour = h
	conf.endMinute = m
	//生成今日开始时间
	parse, err := time.ParseInLocation("2006-01-02 15:04:05", fmt.Sprintf("%s %s:%s:00", time.Now().Format("2006-01-02 "), complementZero(conf.startHour), complementZero(conf.startMinute)), loc)
	if err != nil {
		return err
	}
	conf.StartingTimeToday = parse
	calculateDownTime()
	return nil
}
func analysisTime(time string) (uint, uint, error) {
	split := strings.Split(time, ":")
	if len(split) != 2 {
		return 0, 0, errors.New("时间格式出错 例 08:30")
	}
	//解析开始时间
	var h, m uint
	parseUint, _ := strconv.ParseUint(split[0], 10, 16)
	h = uint(parseUint)
	parseUint, _ = strconv.ParseUint(split[1], 10, 16)
	m = uint(parseUint)
	return h, m, nil
}
func calculateDownTime() {
	//根据开始时间和配置时间生成 结束时间
	stt, _ := time.ParseInLocation("2006-01-02 15:04:05", fmt.Sprintf("2001-10-22 %s:%s:00", complementZero(conf.startHour), complementZero(conf.startMinute)), loc)
	ett, _ := time.ParseInLocation("2006-01-02 15:04:05", fmt.Sprintf("2001-10-22 %s:%s:00", complementZero(conf.endHour), complementZero(conf.endMinute)), loc)
	sub := ett.Sub(stt)
	conf.DownTimeToday = conf.StartingTimeToday.Add(sub)
}
