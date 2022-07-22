package ip

import (
	"encoding/json"
	"fmt"
	"fyne.io/systray"
	"github.com/axgle/mahonia"
	logs "github.com/danbai225/go-logs"
	"github.com/danbai225/tipbar/core"
	"github.com/ncruces/zenity"
	"golang.design/x/clipboard"
	"io/ioutil"
	"net"
	"net/http"
)

var ip *core.Module
var info ipInfo
var pItem, vItem, addItem *systray.MenuItem
var vip = ""

func ExportModule() *core.Module {
	ip = core.NewModule("ip", "IP信息", "ip", onReady, nil, nil)
	return ip
}

func onReady(item *systray.MenuItem) {
	pItem = item.AddSubMenuItem("公网:", "公网地址")
	vItem = item.AddSubMenuItem("内网:", "内网地址")
	addItem = item.AddSubMenuItem("地址信息:", "地址信息")
	s, _ := LocalIPv4s()
	for i, s2 := range s {
		if i == 0 {
			vip = s2
			vItem.SetTitle(fmt.Sprintf("内网:%s", s2))
		}
		ii := vItem.AddSubMenuItem(s2, "")
		go cpItem(s2, ii)
	}
	qItem := item.AddSubMenuItem("ip查询", "ip查询")
	update()
	for {
		select {
		case <-item.ClickedCh:
			update()
		case <-pItem.ClickedCh:
			clipboard.Write(clipboard.FmtText, []byte(info.IP))
		case <-vItem.ClickedCh:
			clipboard.Write(clipboard.FmtText, []byte(vip))
		case <-addItem.ClickedCh:
			clipboard.Write(clipboard.FmtText, []byte(info.Addr))
		case <-qItem.ClickedCh:
			entry, err := zenity.Entry("请输入需要查询的ip")
			if err != nil {
				continue
			}
			i := getIpInfo(entry)
			_ = zenity.Info(fmt.Sprintf("ip:%s\naddr:%s", i.IP, i.Addr))
		}
	}
}
func cpItem(val string, item *systray.MenuItem) {
	for {
		_, ok := <-item.ClickedCh
		if !ok {
			break
		}
		clipboard.Write(clipboard.FmtText, []byte(val))
	}
}
func update() {
	info = getIpInfo("")
	pItem.SetTitle(fmt.Sprintf("ip:%s", info.IP))
	addItem.SetTitle(fmt.Sprintf("地址信息:%s", info.Addr))
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
	i := ipInfo{}
	get, err := http.Get("https://whois.pconline.com.cn/ipJson.jsp?json=true&ip=" + ip)
	if err != nil {
		logs.Err(err)
	} else {
		all, _ := ioutil.ReadAll(get.Body)
		err = json.Unmarshal([]byte(mahonia.NewDecoder("gbk").ConvertString(string(all))), &i)
		if err != nil {
			logs.Err(err)
		}
	}
	return i
}
func LocalIPv4s() ([]string, error) {
	var ips []string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ips, err
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			ips = append(ips, ipnet.IP.String())
		}
	}

	return ips, nil
}
