package ip

import (
	"down_tip/module/tcp_proxy"
	"encoding/json"
	"fmt"
	"fyne.io/systray"
	"github.com/axgle/mahonia"
	emoji "github.com/danbai225/flag_emoji"
	logs "github.com/danbai225/go-logs"
	"github.com/danbai225/tipbar/core"
	"github.com/ncruces/zenity"
	"golang.design/x/clipboard"
	"io"
	"net"
	"net/http"
)

var ip *core.Module
var info ipInfo
var pItem, vItem, w4Item, w6Item, addItem *systray.MenuItem
var vip = ""

func ExportModule() *core.Module {
	ip = core.NewModule("ip", "IP信息", "", onReady, nil, nil)
	return ip
}

func onReady(item *systray.MenuItem) {
	pItem = item.AddSubMenuItem("公网:", "")
	w4Item = item.AddSubMenuItem("外网v4:", "")
	w6Item = item.AddSubMenuItem("外网v6:", "")
	vItem = item.AddSubMenuItem("内网:", "")
	addItem = item.AddSubMenuItem("地址信息:", "")
	qItem := item.AddSubMenuItem("ip查询", "点击查询ip信息")
	update()
	for {
		select {
		case <-item.ClickedCh:
			update()
		case <-w4Item.ClickedCh:
			clipboard.Write(clipboard.FmtText, []byte(info.Wv4Ip))
		case <-w6Item.ClickedCh:
			clipboard.Write(clipboard.FmtText, []byte(info.Wv6Ip))
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
		all, _ := io.ReadAll(resp.Body)
		str := string(all)
		if len(str) >= 2 {
			pStr += emoji.Gen(string(str[0]), string(str[1]))
		}
	}
	pItem.SetTitle(fmt.Sprintf("ip:%s", pStr))
	addItem.SetTitle(fmt.Sprintf("地址信息:%s", info.Addr))
	req, err := http.NewRequest("GET", "https://api-ipv4.ip.sb/ip", nil)
	if err != nil {
		logs.Err(err)
		return
	}
	req.Header.Add("User-Agent", "Mozilla/5.0")
	req.Header.Add("Accept-Language", "en-US")
	client := &http.Client{}
	resp, err = client.Do(req)
	if err == nil {
		ipv4, _ := io.ReadAll(resp.Body)
		w4Item.SetTitle(fmt.Sprintf("外网v4:%s", ipv4))
		info.Wv4Ip = string(ipv4)
	}
	req, err = http.NewRequest("GET", "https://api-ipv6.ip.sb/ip", nil)
	if err != nil {
		logs.Err(err)
		return
	}
	req.Header.Add("User-Agent", "Mozilla/5.0")
	req.Header.Add("Accept-Language", "en-US")
	resp, err = client.Do(req)
	if err == nil {
		ipv6, _ := io.ReadAll(resp.Body)
		w6Item.SetTitle(fmt.Sprintf("外网v6:%s", ipv6))
		info.Wv6Ip = string(ipv6)
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
	Wv4Ip       string
	Wv6Ip       string
}

func getIpInfo(ip string) ipInfo {
	i := ipInfo{}
	if tcp_proxy.IsDomain(ip) {
		ips := tcp_proxy.GetIP(ip, []string{"223.5.5.5", "114.114.114.114", "117.50.10.10", "119.29.29.29"})
		if len(ips) > 0 {
			ip = ips[0]
		}
	}
	get, err := http.Get("https://whois.pconline.com.cn/ipJson.jsp?json=true&ip=" + ip)
	if err != nil {
		logs.Err(err)
	} else {
		all, _ := io.ReadAll(get.Body)
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
