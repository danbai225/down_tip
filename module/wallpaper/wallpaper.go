package wallpaper

import (
	"encoding/json"
	"fmt"
	"fyne.io/systray"
	logs "github.com/danbai225/go-logs"
	"github.com/danbai225/tipbar/core"
	"github.com/ncruces/zenity"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var wallpaper *core.Module

type wallpaperConfig struct {
	Resolution string //分辨率
	Auto       bool   //自动
	Interval   int    //间隔，分钟
}

var config = wallpaperConfig{}
var timer *time.Ticker

func ExportModule() *core.Module {
	wallpaper = core.NewModule("wallpaper", "壁纸", "wallpaper", onReady, nil, nil)
	return wallpaper
}

func onReady(item *systray.MenuItem) {
	wallpaper.UnmarshalConfig(&config)
	if config.Resolution == "" {
		config.Resolution = "16:9"
	}
	change := item.AddSubMenuItem("更换", "")
	resolutionRatio := item.AddSubMenuItem(fmt.Sprintf("设置分辨率/比例 当前:%s", config.Resolution), "")
	auto := item.AddSubMenuItemCheckbox("自动更换", "", config.Auto)
	updateTimer()
	go timeOut()
	for {
		select {
		case <-change.ClickedCh:
			changeWallpaper()
		case <-resolutionRatio.ClickedCh:
			entry, err := zenity.Entry("请输入壁纸分辨率例如 1920x1080或者比例 16:9")
			if err == nil {
				config.Resolution = entry
				resolutionRatio.SetTitle(fmt.Sprintf("设置分辨率/比例 当前:%s", config.Resolution))
			}
			wallpaper.SaveConfig(config)
		case <-auto.ClickedCh:
			if auto.Checked() {
				auto.Uncheck()
			} else {
				auto.Check()
				entry, err := zenity.Entry("请输入更换间隔分钟数")
				if err == nil {
					m, err2 := strconv.Atoi(entry)
					if err2 == nil {
						config.Interval = m
					}
				}
				if config.Interval <= 0 {
					config.Interval = 1
				}
				auto.SetTitle(fmt.Sprintf("自动更换(%d分钟)", config.Interval))
			}
			config.Auto = auto.Checked()
			wallpaper.SaveConfig(config)
			updateTimer()
		}
	}
}
func updateTimer() {
	if timer == nil {
		timer = time.NewTicker(time.Minute * time.Duration(config.Interval))
	} else {
		timer.Reset(time.Minute * time.Duration(config.Interval))
	}
}
func timeOut() {
	for {
		<-timer.C
		if config.Auto {
			changeWallpaper()
		}
	}
}
func changeWallpaper() {
	all, err := httpGet(fmt.Sprintf("https://wallhaven.cc/api/v1/search?ratios=%s&sorting=random", config.Resolution))
	rdata := Rdata{}
	err = json.Unmarshal(all, &rdata)
	if err != nil {
		logs.Err(err)
		return
	}
	if len(rdata.Data) == 0 {
		return
	}
	path := rdata.Data[0].Path
	suffix := strings.Split(path, ".")[1]
	data, err := httpGet(path)
	if err != nil {
		logs.Err(err)
		return
	}
	filename := fmt.Sprintf("%s/%s.%s", os.TempDir(), rdata.Data[0].Id, suffix)
	err = ioutil.WriteFile(filename, data, os.ModePerm)
	if err != nil {
		logs.Err(err)
		return
	}
	err = SetWallpaper(filename)
	if err != nil {
		logs.Err(err)
		return
	}
}

type Rdata struct {
	Data []struct {
		Id         string   `json:"id"`
		Url        string   `json:"url"`
		ShortUrl   string   `json:"short_url"`
		Views      int      `json:"views"`
		Favorites  int      `json:"favorites"`
		Source     string   `json:"source"`
		Purity     string   `json:"purity"`
		Category   string   `json:"category"`
		DimensionX int      `json:"dimension_x"`
		DimensionY int      `json:"dimension_y"`
		Resolution string   `json:"resolution"`
		Ratio      string   `json:"ratio"`
		FileSize   int      `json:"file_size"`
		FileType   string   `json:"file_type"`
		CreatedAt  string   `json:"created_at"`
		Colors     []string `json:"colors"`
		Path       string   `json:"path"`
		Thumbs     struct {
			Large    string `json:"large"`
			Original string `json:"original"`
			Small    string `json:"small"`
		} `json:"thumbs"`
	} `json:"data"`
	Meta struct {
		CurrentPage int         `json:"current_page"`
		LastPage    int         `json:"last_page"`
		PerPage     int         `json:"per_page"`
		Total       int         `json:"total"`
		Query       interface{} `json:"query"`
		Seed        string      `json:"seed"`
	} `json:"meta"`
}

func httpGet(url string) (data []byte, err error) {
	client := new(http.Client) //初始化一个http客户端结构体
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(resp.Body)
}
