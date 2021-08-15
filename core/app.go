package core

import (
	_ "embed"
	logs "github.com/danbai225/go-logs"
	"github.com/getlantern/systray"
)

//go:embed ico.ico
var iconBs []byte

//<editor-fold desc="APP主体结构体">

type App struct {
	title  []*title
	module []*Module
	config config
}

func NewApp() *App {
	return &App{config: config{}}
}
func (a *App) addTitle(module *Module, titleText string) {
	for i := range a.title {
		t := a.title[i]
		if t.module == module {
			t.content = titleText
			return
		}
	}
	a.title = append(a.title, &title{module: module, content: titleText})
}
func (a *App) removeTitle(module *Module) {
	off := 0
	for i := range a.title {
		if a.title[i].module == module {
			off++
			continue
		} else {
			if i+off != i {
				a.title[i+off] = a.title[i]
			}
		}
	}
	if off > 0 {
		a.title = a.title[:len(a.title)-off]
	}
}
func (a *App) Run() error {
	//加载配置
	err := a.config.load()
	if err != nil {
		return err
	}
	//获取模块配置
	for _, m := range a.module {
		if c, has := a.config.Module[m.name]; has {
			m.Config = c
		}
	}

	//运行主体
	systray.Run(a.onReady, a.exit)
	return nil
}
func (a *App) onReady() {
	logs.Info("程序启动")
	//设置程序基本图标等等。。
	systray.SetTemplateIcon(iconBs, iconBs)

	for _, module := range a.module {
		item := systray.AddMenuItem(module.name, module.tooltip)
		go module.onReady(item)
	}
	systray.SetTooltip("关于这个程序。。。")
	quit := systray.AddMenuItem("Quit", "退出这个程序")
	go func() {
		<-quit.ClickedCh
		systray.Quit()
	}()
}
func (a *App) exit() {
	for _, module := range a.module {
		go module.exit()
	}
	logs.Info("程序退出")
}
func (a *App) RegisterModule(module ...*Module) {
	for i := range module {
		m := module[i]
		m.app = a
		a.module = append(a.module, m)
	}
}

//</editor-fold>

type title struct {
	module  *Module
	content string
}

//<editor-fold desc="模块结构体">

type Module struct {
	onReady func(item *systray.MenuItem)
	exit    func()
	app     *App
	name    string
	tooltip string
	Config  interface{}
}

func (m *Module) SetTitle(title string) {
	m.app.addTitle(m, title)
}
func (m *Module) RemoveTitle() {
	m.app.removeTitle(m)
}
func NewModule(name, tooltip string, onReady func(item *systray.MenuItem), exit func()) *Module {
	module := Module{name: name, tooltip: tooltip}
	module.onReady = onReady
	module.exit = exit
	return &module
}

//</editor-fold>
