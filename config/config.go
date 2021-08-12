package config

import (
	"down_tip/utils"
	"encoding/json"
	"errors"
	"fmt"
	logs "github.com/danbai225/go-logs"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

//全局配置对象
var Conf Config

//本地时区
var loc, _ = time.LoadLocation("Local")

//默认配置文件位置
var configFile = "config.json"

// Config 配置文件对象
type Config struct {
	StartTime         string `json:"start_time"`
	startHour         uint
	startMinute       uint
	EndTime           string `json:"end_time"`
	endHour           uint
	endMinute         uint
	DownTimeToday     time.Time
	StartingTimeToday time.Time
	Elasticity        bool `json:"elasticity"`
}

func (c *Config) GetStartHour() uint {
	return c.startHour
}
func (c *Config) GetStartMinute() uint {
	return c.startMinute
}

// InitConfig 启动程序时初始化配置文件
func InitConfig() {
	err := loadConfiguration()
	if err != nil {
		//加载失败退出程序
		logs.Err(err)
		os.Exit(1)
	}
}

//加载配置
func loadConfiguration() error {
	bs, err := ioutil.ReadFile(configFile)
	if err != nil {
		return errors.New("加载配置文件失败")
	}
	//解析配置文件json
	json.Unmarshal(bs, &Conf)
	//根据配置文件中时间解析时分
	h, m, err := analysisTime(Conf.StartTime)
	if err != nil {
		return err
	}
	Conf.startHour = h
	Conf.startMinute = m
	//根据配置文件中时间解析时分
	h, m, err = analysisTime(Conf.EndTime)
	if err != nil {
		return err
	}
	Conf.endHour = h
	Conf.endMinute = m
	//生成今日开始时间
	parse, err := time.ParseInLocation("2006-01-02 15:04:05", fmt.Sprintf("%s %s:%s:00", time.Now().Format("2006-01-02 "), utils.ComplementZero(Conf.startHour), utils.ComplementZero(Conf.startMinute)), loc)
	if err != nil {
		return err
	}
	Conf.StartingTimeToday = parse
	CalculateDownTime()
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
func CalculateDownTime() {
	//根据开始时间和配置时间生成 结束时间
	stt, _ := time.ParseInLocation("2006-01-02 15:04:05", fmt.Sprintf("2001-10-22 %s:%s:00", utils.ComplementZero(Conf.startHour), utils.ComplementZero(Conf.startMinute)), loc)
	ett, _ := time.ParseInLocation("2006-01-02 15:04:05", fmt.Sprintf("2001-10-22 %s:%s:00", utils.ComplementZero(Conf.endHour), utils.ComplementZero(Conf.endMinute)), loc)
	sub := ett.Sub(stt)
	Conf.DownTimeToday = Conf.StartingTimeToday.Add(sub)
}
