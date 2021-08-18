package code2img

/*
The MIT License (MIT)

Copyright (c) 2020 skanehira,danbai225

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
import (
	"context"
	"down_tip/core"
	"github.com/chromedp/cdproto/runtime"
	"github.com/getlantern/systray"
	"golang.design/x/clipboard"

	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"math"
	"net/url"
)

var code2img *core.Module

func ExportModule() *core.Module {
	code2img = core.NewModule("code2img", "代码转图片", "把剪贴板中的代码转成图片", onReady, exit)
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
					clipboard.Write(clipboard.FmtImage, img)
				}
			}
		}
	}
}
func exit() {

}

//https://github.com/carbon-app/carbon/blob/b2e251f429d000ad6c9ee85bb9e052d5cf8db746/lib/constants.js#L624

func code2Img(code string, Options ...map[string]string) ([]byte, error) {
	// create context
	ctx, cancel := chromedp.NewContext(
		context.Background(),
	)
	defer cancel()
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
		"es":     "2x",                 // 出口尺寸
		"wm":     "false",              // 水印
	}
	values := url.Values{}
	for k, v := range carbonOptions {
		values.Set(k, v)
	}
	codeparam := url.Values{}
	codeparam.Set("code", url.PathEscape(code))
	urlstr := "https://carbon.now.sh/?" + values.Encode() + "&" + codeparam.Encode()
	// capture screenshot of an element
	var buf []byte
	if err := chromedp.Run(ctx, elementScreenshot(urlstr, "#export-container  .container-bg", &buf)); err != nil {
		return buf, err
	}
	return buf, nil
}

// elementScreenshot takes a screenshot of a specific element.
func elementScreenshot(urlstr, sel string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.EmulateViewport(2560, 1440),
		chromedp.Navigate(urlstr),
		screenshot(sel, res, chromedp.NodeReady, chromedp.ByID),
	}
}
func screenshot(sel interface{}, picbuf *[]byte, opts ...chromedp.QueryOption) chromedp.QueryAction {
	if picbuf == nil {
		panic("picbuf cannot be nil")
	}

	return chromedp.QueryAfter(sel, func(ctx context.Context, time runtime.ExecutionContextID, nodes ...*cdp.Node) error {
		if len(nodes) < 1 {
			return fmt.Errorf("selector %q did not return any nodes", sel)
		}

		// get layout metrics
		_, _, contentSize, _, _, _, err := page.GetLayoutMetrics().Do(ctx)
		if err != nil {
			return err
		}

		width, height := int64(math.Ceil(contentSize.Width)), int64(math.Ceil(contentSize.Height))

		// force viewport emulation
		err = emulation.SetDeviceMetricsOverride(width, height, 1, false).
			WithScreenOrientation(&emulation.ScreenOrientation{
				Type:  emulation.OrientationTypePortraitPrimary,
				Angle: 0,
			}).
			Do(ctx)
		if err != nil {
			return err
		}

		// get box model
		box, err := dom.GetBoxModel().WithNodeID(nodes[0].NodeID).Do(ctx)
		if err != nil {
			return err
		}
		if len(box.Margin) != 8 {
			return chromedp.ErrInvalidBoxModel
		}

		// take screenshot of the box
		buf, err := page.CaptureScreenshot().
			WithFormat(page.CaptureScreenshotFormatPng).
			WithClip(&page.Viewport{
				X:      math.Round(box.Margin[0]),
				Y:      math.Round(box.Margin[1]),
				Width:  math.Round(box.Margin[4] - box.Margin[0]),
				Height: math.Round(box.Margin[5] - box.Margin[1]),
				Scale:  1.0,
			}).Do(ctx)
		if err != nil {
			return err
		}

		*picbuf = buf
		return nil
	}, append(opts, chromedp.NodeVisible)...)
}
