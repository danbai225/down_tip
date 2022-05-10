package code2img

import (
	logs "github.com/danbai225/go-logs"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"net/url"
	"os"
	"testing"
	"time"
)

func TestName(t *testing.T) {
	var carbonOptions = map[string]string{
		"bg":     "rgba(74,144,226,1)", // 背景颜色
		"t":      "VSCode",             // 主题
		"wt":     "none",               // 窗口主题
		"l":      "auto",               // language
		"ds":     "true",               // 阴影
		"dsyoff": "15px",               // dropShadowOffsetY
		"dsblur": "25px",               // dropShadowBlurRadius
		"wc":     "true",               // 窗口控件
		"wa":     "true",               // 宽度调整
		"pv":     "30px",               // 填充垂直
		"ph":     "50px",               // 填充水平
		"ln":     "true",               // 行号
		"fl":     "1",                  // 第一个行号
		"fm":     "Source Code Pro",    // 字体系列
		"fs":     "13.5px",             // 字体大小
		"lh":     "152%",               // 行高
		"si":     "false",              //平方图像
		"es":     "1x",                 // 出口尺寸
		"wm":     "false",              // 水印
	}
	values := url.Values{}
	for k, v := range carbonOptions {
		values.Set(k, v)
	}
	codeparam := url.Values{}
	codeparam.Set("code", url.PathEscape(`package main

import (
	"down_tip/core"
	"down_tip/module/code2img"
	"down_tip/module/down"
	"down_tip/module/ip"
	"down_tip/module/keylog"
	"down_tip/module/qrcode"
	"down_tip/module/weather"
	logs "github.com/danbai225/go-logs"
)

func main() {
	a, err := core.NewApp()
	if err != nil {
		logs.Err(err)
		return
	}
	//注册模块
	a.RegisterModule(down.ExportModule(), keylog.ExportModule(), ip.ExportModule(), qrcode.ExportModule(), code2img.ExportModule(), weather.ExportModule())
	err = a.Run()
	if err != nil {
		logs.Err(err)
		return
	}
}
`))

	var browser *rod.Browser
	if path, exists := launcher.LookPath(); exists {
		u := launcher.New().Bin(path).MustLaunch()
		browser = rod.New().ControlURL(u).MustConnect()
	} else {

	}
	urlstr := "https://carbon.supermario.vip/?" + values.Encode() + "&" + codeparam.Encode()
	page := browser.MustPage()
	err := rod.Try(func() {
		page.Timeout(10 * time.Second).MustNavigate(urlstr)
	})
	screenshot, err := page.MustElement("#export-container").Screenshot(proto.PageCaptureScreenshotFormatPng, 1)
	if err != nil {
		logs.Err(err)
	}
	os.WriteFile("test.png", screenshot, os.ModePerm)
}
