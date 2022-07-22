package socks5proxy

import (
	"fyne.io/systray"
	logs "github.com/danbai225/go-logs"
	"github.com/danbai225/tcpproxy"
	"github.com/danbai225/tipbar/core"
	"github.com/gogf/gf/container/gset"
	"github.com/gogf/gf/net/ghttp"
	"github.com/miekg/dns"
	"github.com/ncruces/zenity"
	"regexp"
	"sync"
	"time"
)

//https://github.com/danbai225/tcpproxy客户端套壳socks5
var socks5 *core.Module

type socks5Config struct {
	Host     string
	Port     string
	Password string
	LPort    string `json:"l_port"`
}

var config = socks5Config{}

func ExportModule() *core.Module {
	socks5 = core.NewModule("socks5proxy", "socks5proxy", "代理模块", onReady, exit, router)
	return socks5
}
func router(group *ghttp.RouterGroup) {
	group.GET("/", func(r *ghttp.Request) {

	})
}

var connflag = false
var rootItem *systray.MenuItem
var client *tcpproxy.Client

func onReady(item *systray.MenuItem) {
	_ = socks5.UnmarshalConfig(&config)
	item.SetTitle("点击运行客户端")
	rootItem = item
	if config.Host != "" && config.Port != "" {
		conn()
		item.SetTitle("点击断开")
	}
	for {
		select {
		case <-item.ClickedCh:
			if !connflag {
				if config.Host == "" {
					host, _ := zenity.Entry("Host",
						zenity.Title("请输入Host"))
					config.Host = host
					port, _ := zenity.Entry("Port",
						zenity.Title("请输入Port"))
					config.Port = port
					pass, _ := zenity.Entry("Password",
						zenity.Title("请输入Password"))
					config.Password = pass
					if host != "" && pass != "" && port != "" {
						socks5.SaveConfig(config)
					} else {
						continue
					}
				}
				go conn()
				item.SetTitle("点击断开")
			} else {
				connflag = false
				rootItem.SetTitle("点击运行客户端")
				_ = client.Stop()
				logs.Info("客户端已经关闭")
			}
		}
	}
}
func conn() {
	if connflag {
		return
	}
	go func() {
		defer func() {
			rootItem.SetTitle("点击运行客户端")
		}()
		if isDomain(config.Host) {
			ips := getIP(config.Host, []string{"223.5.5.5", "114.114.114.114", "8.8.8.8", "1.1.1.1"})
			for _, ip := range ips {
				client = tcpproxy.Client{}.New(config.Password, ip+":"+config.Port, ":"+config.LPort)
				err := client.Start()
				connflag = false
				if err != nil {
					logs.Err(err)
				} else {
					logs.Info("连接成功：", ip)
					break
				}
			}
		}
	}()
	connflag = true
}
func exit() {

}
func getIP(domain string, dnsList []string) []string {
	set := gset.New(true)
	var dst []string
	group := &sync.WaitGroup{}
	for _, s := range dnsList {
		group.Add(1)
		go func(ip string) {
			defer group.Done()
			c := dns.Client{
				Timeout: 5 * time.Second,
			}
			m := dns.Msg{}
			m.SetQuestion(dns.Fqdn(domain), dns.TypeA)
			r, _, err := c.Exchange(&m, ip+":53")
			if err != nil {
				logs.Err(err)
				return
			}
			for _, ans := range r.Answer {
				record, isType := ans.(*dns.A)
				if isType {
					set.Add(record.A.String())
				}
			}
		}(s)
	}
	group.Wait()
	slice := set.Slice()
	for _, i := range slice {
		dst = append(dst, i.(string))
	}
	return dst
}
func isDomain(text string) bool {
	compile := regexp.MustCompile(".*[a-zA-Z].*")
	return compile.MatchString(text)
}
