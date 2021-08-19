package code2img

import (
	"down_tip/core"
	"github.com/getlantern/systray"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"golang.design/x/clipboard"
	"net/url"
	"time"
)

var code2img *core.Module

func ExportModule() *core.Module {
	code2img = core.NewModule("code2img", "代码转图片", "把剪贴板中的代码转成图片", onReady, exit, nil)
	return code2img
}

func onReady(item *systray.MenuItem) {
	for {
		select {
		case <-item.ClickedCh:
			code := string(clipboard.Read(clipboard.FmtText))
			if code != "" {
				img, err := code2Img(code)
				if err == nil {
					code2img.Notify("转换成功")
					clipboard.Write(clipboard.FmtImage, img)
				} else {
					code2img.Notify("转换失败")
				}
			}
		}
	}
}
func exit() {

}

//https://github.com/carbon-app/carbon/blob/b2e251f429d000ad6c9ee85bb9e052d5cf8db746/lib/constants.js#L624

func code2Img(code string, Options ...map[string]string) ([]byte, error) {
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
	if len(Options) > 0 {
		for k, v := range Options[0] {
			carbonOptions[k] = v
		}
	}
	values := url.Values{}
	for k, v := range carbonOptions {
		values.Set(k, v)
	}
	codeparam := url.Values{}
	codeparam.Set("code", url.PathEscape(code))
	var browser *rod.Browser
	if path, exists := launcher.LookPath(); exists {
		u := launcher.New().Bin(path).MustLaunch()
		browser = rod.New().ControlURL(u).MustConnect()
	} else {
		browser = rod.New().MustConnect()
	}
	urlstr := "https://carbon.supermario.vip/?" + values.Encode() + "&" + codeparam.Encode()
	page := browser.MustPage()
	err := rod.Try(func() {
		page.Timeout(10 * time.Second).MustNavigate(urlstr)
	})
	if err != nil {
		return nil, err
	}
	bytes, err := page.MustElement("#export-container").Screenshot(proto.PageCaptureScreenshotFormatPng, 100)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
