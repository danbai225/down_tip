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
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var wallpaper *core.Module

type wallpaperConfig struct {
	Resolution string //分辨率
	Auto       bool   //自动
	Interval   int    //间隔，分钟
	Categories string //类别
	Q          string //关键字
}

var config = wallpaperConfig{}
var timer *time.Ticker

func ExportModule() *core.Module {
	wallpaper = core.NewModule("wallpaper", "壁纸", "wallpaper", onReady, nil, nil)
	return wallpaper
}
func getCategoriesName(c string) string {
	switch c {
	case "100":
		return "一般"
	case "010":
		return "动漫"
	case "001":
		return "美图"
	}
	return ""
}
func onReady(item *systray.MenuItem) {
	wallpaper.UnmarshalConfig(&config)
	if config.Resolution == "" {
		config.Resolution = "16x9"
	}
	if config.Categories == "" {
		config.Categories = "100"
	}
	change := item.AddSubMenuItem("更换", "")
	resolutionRatio := item.AddSubMenuItem(fmt.Sprintf("设置分辨率/比例 当前:%s", config.Resolution), "")
	autoStr := "自动更换"
	if config.Auto {
		autoStr = fmt.Sprintf("自动更换(%d分钟)", config.Interval)
	}
	auto := item.AddSubMenuItemCheckbox(autoStr, "", config.Auto)
	categoriesItem := item.AddSubMenuItem(fmt.Sprintf("类型(%s)", getCategoriesName(config.Categories)), "")
	categories1 := categoriesItem.AddSubMenuItem("一般", "")
	categories2 := categoriesItem.AddSubMenuItem("动漫", "")
	categories3 := categoriesItem.AddSubMenuItem("美图", "")
	qitemStr := "关键字"
	if config.Q != "" {
		qitemStr = fmt.Sprintf("关键字(%s)", config.Q)
	}
	qItem := item.AddSubMenuItem(qitemStr, "")
	//100/101/111*
	updateTimer()
	go timeOut()
	for {
		select {
		case <-change.ClickedCh:
			changeWallpaper()
		case <-resolutionRatio.ClickedCh:
			entry, err := zenity.Entry("请输入壁纸分辨率例如 1920x1080或者比例 16x9")
			if err == nil {
				config.Resolution = entry
				resolutionRatio.SetTitle(fmt.Sprintf("设置分辨率/比例 当前:%s", config.Resolution))
			}
			wallpaper.SaveConfig(config)
		case <-auto.ClickedCh:
			if auto.Checked() {
				auto.Uncheck()
				auto.SetTitle("自动更换")
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
		case <-categories1.ClickedCh:
			config.Categories = "100"
			categoriesItem.SetTitle(fmt.Sprintf("类型(%s)", getCategoriesName(config.Categories)))
			wallpaper.SaveConfig(config)
		case <-categories2.ClickedCh:
			config.Categories = "010"
			categoriesItem.SetTitle(fmt.Sprintf("类型(%s)", getCategoriesName(config.Categories)))
			wallpaper.SaveConfig(config)
		case <-categories3.ClickedCh:
			config.Categories = "001"
			categoriesItem.SetTitle(fmt.Sprintf("类型(%s)", getCategoriesName(config.Categories)))
			wallpaper.SaveConfig(config)
		case <-qItem.ClickedCh:
			entry, err := zenity.Entry("请输入关键字(英文)：")
			if err == nil {
				config.Q = entry
				qItem.SetTitle(fmt.Sprintf("关键字(%s)", entry))
			}
			wallpaper.SaveConfig(config)
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
	url := fmt.Sprintf("https://wallhaven.cc/api/v1/search?ratios=%s&sorting=random&categories=%s&q=%s", config.Resolution, config.Categories, config.Q)
	all, err := httpGet(url)
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
	split := strings.Split(path, ".")
	suffix := split[len(split)-1]
	data, err := httpGet(path)
	if err != nil {
		logs.Err(err)
		return
	}
	_ = os.MkdirAll("./wallpaper", os.ModePerm)

	filename := fmt.Sprintf("./wallpaper/%s.%s", rdata.Data[0].Id, suffix)
	err = ioutil.WriteFile(filename, data, os.ModePerm)
	if err != nil {
		logs.Err(err)
		return
	}
	abs, _ := filepath.Abs(filename)
	err = SetWallpaper(abs)
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
