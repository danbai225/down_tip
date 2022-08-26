package ip

import (
	"down_tip/module/socks5proxy"
	"encoding/json"
	"fmt"
	"fyne.io/systray"
	"github.com/axgle/mahonia"
	emoji "github.com/danbai225/flag_emoji"
	logs "github.com/danbai225/go-logs"
	"github.com/danbai225/tipbar/core"
	"github.com/ncruces/zenity"
	"golang.design/x/clipboard"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

var ip *core.Module
var info ipInfo
var pItem, vItem, wItem, addItem *systray.MenuItem
var vip = ""

func ExportModule() *core.Module {
	ip = core.NewModule("ip", "IP信息", "", onReady, nil, nil)
	return ip
}

func onReady(item *systray.MenuItem) {
	pItem = item.AddSubMenuItem("公网:", "")
	wItem = item.AddSubMenuItem("外网:", "")
	vItem = item.AddSubMenuItem("内网:", "")
	addItem = item.AddSubMenuItem("地址信息:", "")
	qItem := item.AddSubMenuItem("ip查询", "点击查询ip信息")
	update()
	for {
		select {
		case <-item.ClickedCh:
			update()
		case <-wItem.ClickedCh:
			clipboard.Write(clipboard.FmtText, []byte(info.WIp))
		case <-pItem.ClickedCh:
			clipboard.Write(clipboard.FmtText, []byte(info.IP))
		case <-vItem.ClickedCh:
			clipboard.Write(clipboard.FmtText, []byte(vip))
		case <-addItem.ClickedCh:
			clipboard.Write(clipboard.FmtText, []byte(info.Addr))
		case <-qItem.ClickedCh:
			entry, err := zenity.Entry("请输入需要查询的ip/域名")
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
	s, _ := LocalIPv4s()
	for i, s2 := range s {
		if i == 0 {
			vip = s2
			vItem.SetTitle(fmt.Sprintf("内网:%s", s2))
		}
		ii := vItem.AddSubMenuItem(s2, "")
		go cpItem(s2, ii)
	}
	info = getIpInfo("")
	pStr := info.IP
	resp, err := http.Get("https://ip.gs/country-iso?ip=" + info.IP)
	if err == nil {
		all, _ := ioutil.ReadAll(resp.Body)
		str := string(all)
		if len(str) >= 2 {
			pStr += emoji.Gen(string(str[0]), string(str[1]))
		}
	}
	pItem.SetTitle(fmt.Sprintf("ip:%s", pStr))
	addItem.SetTitle(fmt.Sprintf("地址信息:%s", info.Addr))
	get, err := http.Get("https://ip.gs")
	if err == nil {
		all, _ := ioutil.ReadAll(get.Body)
		info.WIp = string(all)
		wStr := string(all)
		resp, err = http.Get("https://ip.gs/country-iso?ip=" + strings.ReplaceAll(info.WIp, "\n", ""))
		if err == nil {
			all, _ = ioutil.ReadAll(resp.Body)
			str := string(all)
			if len(str) >= 2 {
				wStr += emoji.Gen(string(str[0]), string(str[1]))
			}
		}
		wItem.SetTitle(fmt.Sprintf("外网:%s", wStr))
	}
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
	WIp         string
}

func getIpInfo(ip string) ipInfo {
	i := ipInfo{}
	if socks5proxy.IsDomain(ip) {
		ips := socks5proxy.GetIP(ip, []string{"223.5.5.5", "114.114.114.114", "117.50.10.10", "119.29.29.29"})
		if len(ips) > 0 {
			ip = ips[0]
		}
	}
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
