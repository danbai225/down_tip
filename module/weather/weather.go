package weather

import (
	"down_tip/core"
	"encoding/json"
	"fmt"
	logs "github.com/danbai225/go-logs"
	"github.com/getlantern/systray"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/skratchdot/open-golang/open"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var weather *core.Module

type weatherConfig struct {
	LatitudeAndLongitude string
}

var config = weatherConfig{}

func ExportModule() *core.Module {
	weather = core.NewModule("weather", "天气", "天气", onReady, exit, router)
	return weather
}
func router(group *ghttp.RouterGroup) {
	group.GET("/", func(r *ghttp.Request) {
		query := r.GetQueryString("latng")
		if query != "" {
			split := strings.Split(query, ",")
			config.LatitudeAndLongitude = split[1] + "," + split[0]
			r.Response.WriteJson(g.Map{
				"msg":  "获取成功",
				"code": 0,
			})
			weather.SaveConfig(config)
		} else {
			r.Response.WriteJson(g.Map{
				"msg":  "获取失败",
				"code": 1,
			})
		}
	})
}
func onReady(item *systray.MenuItem) {
	weather.UnmarshalConfig(&config)
	if config.LatitudeAndLongitude == "" {
		item.SetTitle("点击获取天气信息")
	}
	go func() {
		for {
			if config.LatitudeAndLongitude != "" {
				go weatherUpdate(item)
			}
			time.Sleep(time.Minute * 5)
		}
	}()
	for {
		select {
		case <-item.ClickedCh:
			if config.LatitudeAndLongitude == "" {
				weather.Notify("未设置地理位置请先设置")
				sprintf := fmt.Sprintf(`https://mapapi.qq.com/web/mapComponents/locationPicker/v/index.html?&type=0&backurl=%s&key=WWBBZ-FMDY6-FNDSG-M3IP4-2QLEF-SQBJH&referer=`, weather.GetRootUrl()) + `%E5%B7%A5%E5%85%B7%E9%9B%86`
				open.Run(sprintf)
			} else {
				weatherUpdate(item)
			}
		}
	}
}

var subItem = map[string]*systray.MenuItem{}
var alerts = make([]*systray.MenuItem, 0)
var alertsMsgMap = map[string]struct{}{}

func weatherUpdate(item *systray.MenuItem) {
	url := fmt.Sprintf("https://api.caiyunapp.com/v2.5/ujp0HddE4bY2SwRc/%s/weather.json?unit=metric:v2&alert=true", config.LatitudeAndLongitude)
	resp, err := http.Get(url)
	if err != nil {
		weather.Notify("获取失败:" + err.Error())
	}
	readAll, _ := ioutil.ReadAll(resp.Body)
	res := Weather{}
	err = json.Unmarshal(readAll, &res)
	if err != nil {
		logs.Err(err)
		return
	}
	realtime := res.Result.Realtime
	alert := res.Result.Alert
	minutely := res.Result.Minutely
	hourly := res.Result.Hourly
	//添加子项
	if len(subItem) == 0 {
		subItem["alert"] = item.AddSubMenuItem("", "预警")
		subItem["temperature"] = item.AddSubMenuItem("", "温度和体感温度")
		subItem["precipitation.nearest"] = item.AddSubMenuItem("", "最近的降水带距离和强度")
		subItem["precipitation.local"] = item.AddSubMenuItem("", "本地的降水强度")
		subItem["minutely.description"] = item.AddSubMenuItem("", "天气情况")
		subItem["hourly"] = item.AddSubMenuItem("预警信息", "未来一天预报")
	}
	//预警
	if len(alert.Content) == 0 {
		subItem["alert"].Hide()
	} else {
		subItem["alert"].Show()
		if len(alerts) > 0 {
			for _, menuItem := range alerts {
				menuItem.Hide()
				menuItem.Disable()
			}
		}
		for _, s := range alert.Content {
			subItem["alert"].AddSubMenuItem(s.Title, s.Description)
			if _, has := alertsMsgMap[s.AlertID]; !has {
				weather.Notify(s.Description)
				alertsMsgMap[s.AlertID] = struct{}{}
			}
		}
	}
	//当前天气
	item.SetTitle("天气状况:" + getSkyString(realtime.Skycon))
	subItem["temperature"].SetTitle(fmt.Sprintf("温度:%.1f°C,体感温度:%.1f°C", realtime.Temperature, realtime.ApparentTemperature))
	subItem["precipitation.nearest"].SetTitle(fmt.Sprintf("最近的降水距离%.1fkm，降水量%.fmm/h", realtime.Precipitation.Nearest.Distance, realtime.Precipitation.Nearest.Intensity))
	if realtime.Precipitation.Local.Intensity > 0 {
		subItem["precipitation.local"].SetTitle(fmt.Sprintf("本地的降水量%.fmm/h", realtime.Precipitation.Local.Intensity))
		subItem["precipitation.local"].Show()
	} else {
		subItem["precipitation.local"].Hide()
	}
	//天气情况
	if minutely.Description != "" {
		subItem["minutely.description"].SetTitle(minutely.Description)
	} else {
		subItem["minutely.description"].Hide()
	}
	//未来一天预报
	if minutely.Description != "" {
		subItem["hourly"].SetTitle(hourly.Description)
	} else {
		subItem["hourly"].Hide()
	}
}
func exit() {

}

//https://open.caiyunapp.com/%E5%BD%A9%E4%BA%91%E5%A4%A9%E6%B0%94_API_%E4%B8%80%E8%A7%88%E8%A1%A8
func getSkyString(code string) string {
	switch code {
	case "CLEAR_DAY":
		return "晴🌞"
	case "CLEAR_NIGHT":
		return "晴🌝"
	case "PARTLY_CLOUDY_DAY":
		return "多云🌥"
	case "PARTLY_CLOUDY_NIGHT":
		return "多云🌥"
	case "CLOUDY":
		return "阴🌥"
	case "LIGHT_HAZE":
		return "轻度雾霾"
	case "MODERATE_HAZE":
		return "中度雾霾"
	case "HEAVY_HAZE":
		return "重度雾霾"
	case "LIGHT_RAIN":
		return "小雨🌧"
	case "MODERATE_RAIN":
		return "中雨🌧"
	case "HEAVY_RAIN":
		return "大雨🌧"
	case "STORM_RAIN":
		return "暴雨🌧"
	case "FOG":
		return "雾🌫"
	case "LIGHT_SNOW":
		return "小雪🌨"
	case "MODERATE_SNOW":
		return "中雪🌨"
	case "HEAVY_SNOW":
		return "大雪🌨"
	case "STORM_SNOW":
		return "暴雪🌨"
	case "DUST":
		return "浮尘"
	case "SAND":
		return "沙尘"
	case "WIND":
		return "大风🌬"
	}
	return code
}

//// Realtime 实时
//type Realtime struct {
//	Status     string    `json:"status"`
//	APIVersion string    `json:"api_version"`
//	APIStatus  string    `json:"api_status"`
//	Lang       string    `json:"lang"`
//	Unit       string    `json:"unit"`
//	Tzshift    float64       `json:"tzshift"`
//	Timezone   string    `json:"timezone"`
//	ServerTime float64       `json:"server_time"`
//	Location   []float64 `json:"location"`
//	Result     struct {
//		Realtime struct {
//			Status      string  `json:"status"`
//			Temperature float64 `json:"temperature"` //温度
//			Humidity    float64 `json:"humidity"` //相对湿度
//			Cloudrate   float64     `json:"cloudrate"` //云量
//			Skycon      string  `json:"skycon"` //天气状况 //https://open.caiyunapp.com/%E5%AE%9E%E5%86%B5%E5%A4%A9%E6%B0%94%E6%8E%A5%E5%8F%A3/v2.5#.E5.A4.A9.E6.B0.94.E7.8E.B0.E8.B1.A1.E4.BB.A3.E7.A0.81
//			Visibility  float64 `json:"visibility"` //能见度
//			Dswrf       float64 `json:"dswrf"` //短波辐射
//			Wind        struct {
//				Speed     float64 `json:"speed"` //风速
//				Direction float64     `json:"direction"` //风向
//			} `json:"wind"`
//			Pressure            float64 `json:"pressure"` //气压
//			ApparentTemperature float64 `json:"apparent_temperature"` //体感温度
//			Precipitation       struct {
//				Local struct {
//					Status     string `json:"status"`
//					Datasource string `json:"datasource"`
//					Intensity  float64    `json:"intensity"`
//				} `json:"local"`
//				Nearest struct {
//					Status    string `json:"status"`
//					Distance  float64    `json:"distance"` //最近降水距离
//					Intensity float64    `json:"intensity"` //本地降水强度
//				} `json:"nearest"`
//			} `json:"precipitation"`
//			AirQuality struct {
//				Pm25 float64 `json:"pm25"` //PM25浓度
//				Pm10 float64 `json:"pm10"` //PM10浓度
//				O3   float64 `json:"o3"` //臭氧浓度
//				So2  float64 `json:"so2"` //二氧化氮浓度
//				No2  float64 `json:"no2"` //二氧化硫浓度
//				Co   float64 `json:"co"` //一氧化碳浓度
//				Aqi  struct {
//					Chn float64 `json:"chn"`
//					Usa float64 `json:"usa"`
//				} `json:"aqi"`
//				Description struct {
//					Usa string `json:"usa"`
//					Chn string `json:"chn"`
//				} `json:"description"`
//			} `json:"air_quality"`
//			LifeIndex struct {
//				Ultraviolet struct {
//					Index float64    `json:"index"` //紫外线指数
//					Desc  string `json:"desc"`
//				} `json:"ultraviolet"`
//				Comfort struct {
//					Index float64    `json:"index"` //舒适度指数
//					Desc  string `json:"desc"`
//				} `json:"comfort"`
//			} `json:"life_index"`
//		} `json:"realtime"`
//		Primary float64 `json:"primary"`
//	} `json:"result"`
//}

type Weather struct {
	Status     string    `json:"status"`
	APIVersion string    `json:"api_version"`
	APIStatus  string    `json:"api_status"`
	Lang       string    `json:"lang"`
	Unit       string    `json:"unit"`
	Tzshift    float64   `json:"tzshift"`
	Timezone   string    `json:"timezone"`
	ServerTime float64   `json:"server_time"`
	Location   []float64 `json:"location"`
	Result     struct {
		Alert struct {
			Status  string `json:"status"`
			Content []struct {
				Province      string    `json:"province"`
				Status        string    `json:"status"`
				Code          string    `json:"code"`
				Description   string    `json:"description"`
				Pubtimestamp  float64   `json:"pubtimestamp"`
				City          string    `json:"city"`
				Adcode        string    `json:"adcode"`
				RegionID      string    `json:"regionId"`
				Latlon        []float64 `json:"latlon"`
				County        string    `json:"county"`
				AlertID       string    `json:"alertId"`
				RequestStatus string    `json:"request_status"`
				Source        string    `json:"source"`
				Title         string    `json:"title"`
				Location      string    `json:"location"`
			} `json:"content"`
		} `json:"alert"`
		Realtime struct {
			Status      string  `json:"status"`
			Temperature float64 `json:"temperature"`
			Humidity    float64 `json:"humidity"`
			Cloudrate   float64 `json:"cloudrate"`
			Skycon      string  `json:"skycon"`
			Visibility  float64 `json:"visibility"`
			Dswrf       float64 `json:"dswrf"`
			Wind        struct {
				Speed     float64 `json:"speed"`
				Direction float64 `json:"direction"`
			} `json:"wind"`
			Pressure            float64 `json:"pressure"`
			ApparentTemperature float64 `json:"apparent_temperature"`
			Precipitation       struct {
				Local struct {
					Status     string  `json:"status"`
					Datasource string  `json:"datasource"`
					Intensity  float64 `json:"intensity"`
				} `json:"local"`
				Nearest struct {
					Status    string  `json:"status"`
					Distance  float64 `json:"distance"`
					Intensity float64 `json:"intensity"`
				} `json:"nearest"`
			} `json:"precipitation"`
			AirQuality struct {
				Pm25 float64 `json:"pm25"`
				Pm10 float64 `json:"pm10"`
				O3   float64 `json:"o3"`
				So2  float64 `json:"so2"`
				No2  float64 `json:"no2"`
				Co   float64 `json:"co"`
				Aqi  struct {
					Chn float64 `json:"chn"`
					Usa float64 `json:"usa"`
				} `json:"aqi"`
				Description struct {
					Usa string `json:"usa"`
					Chn string `json:"chn"`
				} `json:"description"`
			} `json:"air_quality"`
			LifeIndex struct {
				Ultraviolet struct {
					Index float64 `json:"index"`
					Desc  string  `json:"desc"`
				} `json:"ultraviolet"`
				Comfort struct {
					Index float64 `json:"index"`
					Desc  string  `json:"desc"`
				} `json:"comfort"`
			} `json:"life_index"`
		} `json:"realtime"`
		Minutely struct {
			Status          string    `json:"status"`
			Datasource      string    `json:"datasource"`
			Precipitation2H []float64 `json:"precipitation_2h"`
			Precipitation   []float64 `json:"precipitation"`
			Probability     []float64 `json:"probability"`
			Description     string    `json:"description"`
		} `json:"minutely"`
		Hourly struct {
			Status        string `json:"status"`
			Description   string `json:"description"`
			Precipitation []struct {
				Datetime string  `json:"datetime"`
				Value    float64 `json:"value"`
			} `json:"precipitation"`
			Temperature []struct {
				Datetime string  `json:"datetime"`
				Value    float64 `json:"value"`
			} `json:"temperature"`
			Wind []struct {
				Datetime  string  `json:"datetime"`
				Speed     float64 `json:"speed"`
				Direction float64 `json:"direction"`
			} `json:"wind"`
			Humidity []struct {
				Datetime string  `json:"datetime"`
				Value    float64 `json:"value"`
			} `json:"humidity"`
			Cloudrate []struct {
				Datetime string  `json:"datetime"`
				Value    float64 `json:"value"`
			} `json:"cloudrate"`
			Skycon []struct {
				Datetime string `json:"datetime"`
				Value    string `json:"value"`
			} `json:"skycon"`
			Pressure []struct {
				Datetime string  `json:"datetime"`
				Value    float64 `json:"value"`
			} `json:"pressure"`
			Visibility []struct {
				Datetime string  `json:"datetime"`
				Value    float64 `json:"value"`
			} `json:"visibility"`
			Dswrf []struct {
				Datetime string  `json:"datetime"`
				Value    float64 `json:"value"`
			} `json:"dswrf"`
			AirQuality struct {
				Aqi []struct {
					Datetime string `json:"datetime"`
					Value    struct {
						Chn float64 `json:"chn"`
						Usa float64 `json:"usa"`
					} `json:"value"`
				} `json:"aqi"`
				Pm25 []struct {
					Datetime string  `json:"datetime"`
					Value    float64 `json:"value"`
				} `json:"pm25"`
			} `json:"air_quality"`
		} `json:"hourly"`
		Daily struct {
			Status string `json:"status"`
			Astro  []struct {
				Date    string `json:"date"`
				Sunrise struct {
					Time string `json:"time"`
				} `json:"sunrise"`
				Sunset struct {
					Time string `json:"time"`
				} `json:"sunset"`
			} `json:"astro"`
			Precipitation []struct {
				Date string  `json:"date"`
				Max  float64 `json:"max"`
				Min  float64 `json:"min"`
				Avg  float64 `json:"avg"`
			} `json:"precipitation"`
			Temperature []struct {
				Date string  `json:"date"`
				Max  float64 `json:"max"`
				Min  float64 `json:"min"`
				Avg  float64 `json:"avg"`
			} `json:"temperature"`
			Wind []struct {
				Date string `json:"date"`
				Max  struct {
					Speed     float64 `json:"speed"`
					Direction float64 `json:"direction"`
				} `json:"max"`
				Min struct {
					Speed     float64 `json:"speed"`
					Direction float64 `json:"direction"`
				} `json:"min"`
				Avg struct {
					Speed     float64 `json:"speed"`
					Direction float64 `json:"direction"`
				} `json:"avg"`
			} `json:"wind"`
			Humidity []struct {
				Date string  `json:"date"`
				Max  float64 `json:"max"`
				Min  float64 `json:"min"`
				Avg  float64 `json:"avg"`
			} `json:"humidity"`
			Cloudrate []struct {
				Date string  `json:"date"`
				Max  float64 `json:"max"`
				Min  float64 `json:"min"`
				Avg  float64 `json:"avg"`
			} `json:"cloudrate"`
			Pressure []struct {
				Date string  `json:"date"`
				Max  float64 `json:"max"`
				Min  float64 `json:"min"`
				Avg  float64 `json:"avg"`
			} `json:"pressure"`
			Visibility []struct {
				Date string  `json:"date"`
				Max  float64 `json:"max"`
				Min  float64 `json:"min"`
				Avg  float64 `json:"avg"`
			} `json:"visibility"`
			Dswrf []struct {
				Date string  `json:"date"`
				Max  float64 `json:"max"`
				Min  float64 `json:"min"`
				Avg  float64 `json:"avg"`
			} `json:"dswrf"`
			AirQuality struct {
				Aqi []struct {
					Date string `json:"date"`
					Max  struct {
						Chn float64 `json:"chn"`
						Usa float64 `json:"usa"`
					} `json:"max"`
					Avg struct {
						Chn float64 `json:"chn"`
						Usa float64 `json:"usa"`
					} `json:"avg"`
					Min struct {
						Chn float64 `json:"chn"`
						Usa float64 `json:"usa"`
					} `json:"min"`
				} `json:"aqi"`
				Pm25 []struct {
					Date string  `json:"date"`
					Max  float64 `json:"max"`
					Avg  float64 `json:"avg"`
					Min  float64 `json:"min"`
				} `json:"pm25"`
			} `json:"air_quality"`
			Skycon []struct {
				Date  string `json:"date"`
				Value string `json:"value"`
			} `json:"skycon"`
			Skycon08H20H []struct {
				Date  string `json:"date"`
				Value string `json:"value"`
			} `json:"skycon_08h_20h"`
			Skycon20H32H []struct {
				Date  string `json:"date"`
				Value string `json:"value"`
			} `json:"skycon_20h_32h"`
			LifeIndex struct {
				Ultraviolet []struct {
					Date  string `json:"date"`
					Index string `json:"index"`
					Desc  string `json:"desc"`
				} `json:"ultraviolet"`
				CarWashing []struct {
					Date  string `json:"date"`
					Index string `json:"index"`
					Desc  string `json:"desc"`
				} `json:"carWashing"`
				Dressing []struct {
					Date  string `json:"date"`
					Index string `json:"index"`
					Desc  string `json:"desc"`
				} `json:"dressing"`
				Comfort []struct {
					Date  string `json:"date"`
					Index string `json:"index"`
					Desc  string `json:"desc"`
				} `json:"comfort"`
				ColdRisk []struct {
					Date  string `json:"date"`
					Index string `json:"index"`
					Desc  string `json:"desc"`
				} `json:"coldRisk"`
			} `json:"life_index"`
		} `json:"daily"`
		Primary          float64 `json:"primary"`
		ForecastKeypoint string  `json:"forecast_keypoint"`
	} `json:"result"`
}
