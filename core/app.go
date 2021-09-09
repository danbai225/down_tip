package core

import (
	_ "embed"
	"fmt"
	logs "github.com/danbai225/go-logs"
	"github.com/getlantern/systray"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/ncruces/zenity"
	"sync"
	"time"
)

//go:embed ico.ico
var iconBs []byte

//<editor-fold desc="APP主体结构体">

type App struct {
	title     []*title
	module    []*Module
	config    config
	tip       chan tip
	titleLock sync.Mutex
	g         *ghttp.Server
	index     func(r *ghttp.Request)
}

func NewApp(configPath ...string) (*App, error) {
	configP := "config.json"
	if len(configPath) > 0 {
		configP = configPath[0]
	}
	app := App{config: config{configName: configP},
		title: make([]*title, 0), module: make([]*Module, 0), tip: make(chan tip, 10),
		index: func(r *ghttp.Request) {
			r.Response.Write("hello downTip")
		},
	}
	//加载配置
	err := app.config.load()
	if err != nil {
		return nil, err
	}
	app.g = g.Server()
	if app.config.HTTPPort == 0 {
		app.config.HTTPPort = 7989
	}
	app.g.SetPort(int(app.config.HTTPPort))
	if app.config.LogsDir != "" {
		logs.SetLogsDir(app.config.LogsDir)
	}
	return &app, nil
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
	//获取模块配置
	for _, m := range a.module {
		if c, has := a.config.Module[m.name]; has {
			m.Config = c.Config
			m.Port = a.config.HTTPPort
			if m.route != nil {
				group := a.g.Group("/" + m.name)
				m.route(group)
			}
		}
	}
	//运行主体
	systray.Run(a.onReady, a.exit)
	return nil
}
func (a *App) doTip() {
	for t := range a.tip {
		a.titleLock.Lock()
		systray.SetTitle(t.content)
		time.Sleep(t.time)
		systray.SetTitle("")
		a.titleLock.Unlock()
	}
}
func (a *App) doTitle() {
	for {
		for _, t := range a.title {
			a.titleLock.Lock()
			systray.SetTitle(t.content)
			a.titleLock.Unlock()
			time.Sleep(5 * time.Second)
			systray.SetTitle("")
		}
		time.Sleep(time.Second)
	}
}
func middlewareCORS(r *ghttp.Request) {
	r.Response.CORSDefault()
	r.Middleware.Next()
}
func (a *App) onReady() {
	logs.Info("程序启动")
	//设置程序基本图标等等。。
	systray.SetTemplateIcon(iconBs, iconBs)
	//运行http
	a.g.BindHandler("/", a.index)
	// 跨域
	a.g.Use(middlewareCORS)
	for _, module := range a.module {
		item := systray.AddMenuItem(module.itemName, module.tooltip)
		if module.onReady != nil {
			go module.onReady(item)
		}
	}
	systray.SetTooltip("关于这个程序。。。")
	quit := systray.AddMenuItem("Quit", "退出这个程序")
	go func() {
		<-quit.ClickedCh
		systray.Quit()
	}()
	go a.doTip()
	go a.doTitle()
	go a.g.Run()

}
func (a *App) exit() {
	for _, module := range a.module {
		if module.exit != nil {
			go module.exit()
		}
	}
	err := a.g.Shutdown()
	if err != nil {
		logs.Err(err)
	}
	logs.Info("程序退出")
}
func (a *App) RegisterModule(module ...*Module) {
	for i := range module {
		m := module[i]
		if mc, has := a.config.Module[m.name]; !has || !mc.Enable {
			continue
		}
		m.app = a
		a.module = append(a.module, m)
	}
}

//</editor-fold>

func NewModule(name, itemName, tooltip string, onReady func(item *systray.MenuItem), exit func(), route func(*ghttp.RouterGroup)) *Module {
	module := Module{name: name, itemName: itemName, tooltip: tooltip}
	module.onReady = onReady
	module.exit = exit
	module.route = route
	return &module
}

type tip struct {
	module  *Module
	content string
	time    time.Duration
}
type title struct {
	module  *Module
	content string
}

//<editor-fold desc="模块结构体">

type Module struct {
	onReady  func(item *systray.MenuItem)
	exit     func()
	app      *App
	name     string
	itemName string
	tooltip  string
	Config   interface{}
	route    func(*ghttp.RouterGroup)
	Port     uint16
}

func (m *Module) UnmarshalConfig(dst interface{}) error {
	return Unmarshal(m.Config, dst)
}

func (m *Module) SetTitle(title string) {
	m.app.addTitle(m, title)
}
func (m *Module) RemoveTitle() {
	m.app.removeTitle(m)
}
func (m *Module) Tip(str string, time time.Duration) {
	go func() {
		m.app.tip <- tip{
			module:  m,
			content: str,
			time:    time,
		}
	}()
}
func (m *Module) Notify(str string) {
	err := zenity.Notify(str)
	if err != nil {
		logs.Err(err)
	}
}
func (m *Module) SaveConfig(c interface{}) {
	m.Config = c
	m.app.config.saveConfig(m, m.Config)
	err := m.app.config.save()
	if err != nil {
		logs.Err()
	}
}
func (m *Module) GetRootUrl() string {
	return fmt.Sprintf("http://localhost:%d/%s", m.Port, m.name)
}

//</editor-fold>
