package ip

import (
	"down_tip/core"
	"encoding/json"
	"github.com/axgle/mahonia"
	"github.com/getlantern/systray"
	"golang.design/x/clipboard"
	"io/ioutil"
	"net/http"
)

var keyLog *core.Module

func ExportModule() *core.Module {
	keyLog = core.NewModule("ip", "ip:"+getIpInfo("").IP, "ip", onReady, exit)
	return keyLog
}
func onReady(item *systray.MenuItem) {
	for {
		select {
		case <-item.ClickedCh:
			info := getIpInfo("")
			clipboard.Write(clipboard.FmtText, []byte(info.IP))
		}
	}
}
func exit() {

}

type ipInfo struct {
	IP          string `json:"ip"`
	Pro         string `json:"pro"`
	ProCode     string `json:"proCode"`
	City        string `json:"city"`
	CityCode    string `json:"cityCode"`
	Region      string `json:"region"`
	RegionCode  string `json:"regionCode"`
	Addr        string `json:"addr"`
	RegionNames string `json:"regionNames"`
	Err         string `json:"err"`
}

func getIpInfo(ip string) ipInfo {
	info := ipInfo{}
	get, err := http.Get("http://whois.pconline.com.cn/ipJson.jsp?json=true&ip=" + ip)
	if err != nil {
		println(err)
	} else {
		all, _ := ioutil.ReadAll(get.Body)
		json.Unmarshal([]byte(mahonia.NewDecoder("gbk").ConvertString(string(all))), &info)
	}
	return info
}
