package code2img

import (
	"github.com/danbai225/tipbar/core"
	"github.com/getlantern/systray"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/gogf/gf/net/ghttp"
	"github.com/ncruces/zenity"
	"github.com/skratchdot/open-golang/open"
	"golang.design/x/clipboard"
	"net/url"
	"runtime"
	"strings"
	"time"
)

var code2img *core.Module
var imgData []byte
var language = "auto"

func ExportModule() *core.Module {
	code2img = core.NewModule("code2img", "代码转图片", "把剪贴板中的代码转成图片", onReady, exit, showImg)
	return code2img
}

func showImg(group *ghttp.RouterGroup) {
	group.GET("/", func(r *ghttp.Request) {
		if len(imgData) > 0 {
			r.Response.Write(imgData)
		}

	})
}
func onReady(item *systray.MenuItem) {
	go2 := item.AddSubMenuItem("点击转换", "转换")
	l := item.AddSubMenuItem("设置语言:"+language, "设置语言")
	for {
		select {
		case <-go2.ClickedCh:
			code := string(clipboard.Read(clipboard.FmtText))
			if code != "" {
				img, err := code2Img(code)
				imgData = img
				if err == nil {
					code2img.Notify("转换成功")
					if runtime.GOOS == "windows" {
						open.Run(code2img.GetRootUrl())
					} else {
						clipboard.Write(clipboard.FmtImage, img)
					}
				} else {
					code2img.Notify("转换失败")
				}
			} else {
				code2img.Notify("没有找到剪贴板中的代码")
			}
		case <-l.ClickedCh:
			entry, err := zenity.Entry("请输入语言")
			if err == nil {
				language = entry
				l.SetTitle("设置语言:" + language)
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
	if language != "" {
		carbonOptions["l"] = language
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
	var browser *rod.Browser

	if path, exists := launcher.LookPath(); exists {
		u := launcher.New().Bin(path).Set("--disable-gpu").Headless(true).MustLaunch()
		browser = rod.New().ControlURL(u).MustConnect()
	} else {
		browser = rod.New().MustConnect()
	}
	defer browser.Close()
	urlstr := "https://carbon.supermario.vip/?" + values.Encode() + "&code=t"
	page := browser.MustPage()
	err := rod.Try(func() {
		page.Timeout(10 * time.Second).MustNavigate(urlstr)
	})
	if err != nil {
		return nil, err
	}
	//defer page.Close()
	pt := page.MustElement(".CodeMirror-lines").MustShape().OnePointInside()
	//模拟鼠标键盘
	mouse := page.Mouse
	keyboard := page.Keyboard
	//移动输入代码
	mouse.MustMove(pt.X, pt.Y-10)
	mouse.MustDown("left")
	mouse.MustUp("left")
	keyboard.MustDown('\b')
	keyboard.MustUp('\b')
	split := strings.Split(code, "\n")
	for i, s := range split {
		if i == len(split)-1 {
			keyboard.InsertText(s)
		} else {
			keyboard.InsertText(s + "\n")
		}
	}
	element := page.MustElement("#export-container")
	box := element.MustShape().Box()
	//logs.Info(box.Width, box.Height)

	element.MustEval(`
getxy =function(){
var element=document.getElementById('export-container')
//计算x坐标
  var actualLeft = element.offsetLeft;
  var current = element.offsetParent;
  while (current !== null){
    actualLeft += current.offsetLeft;
    current = current.offsetParent;
  }
  //计算y坐标
  var actualTop = element.offsetTop;
  var current = element.offsetParent;
  while (current !== null){
    actualTop += (current.offsetTop+current.clientTop);
    current = current.offsetParent;
  }
  //返回结果
  return {x: actualLeft, y: actualTop}
}
`)
	vals := page.MustEval("getxy()")
	i := 90
	img, _ := page.Screenshot(true, &proto.PageCaptureScreenshot{
		Format:  proto.PageCaptureScreenshotFormatJpeg,
		Quality: &i,
		Clip: &proto.PageViewport{
			X:      vals.Get("x").Num(),
			Y:      vals.Get("y").Num(),
			Width:  box.Width,
			Height: box.Height,
			Scale:  2,
		},
		FromSurface: true,
	})
	return img, nil
}
