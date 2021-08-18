package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

var appPath = ""

type config struct {
	Module     map[string]ModuleConfig `json:"module"`
	configName string
	HTTPPort   uint16 `json:"http_port"`
}

func (c *config) getConfig(module *Module) interface{} {
	if moduleConf, has := c.Module[module.name]; has {
		return moduleConf.Config
	}
	return nil
}
func (c *config) saveConfig(module *Module, conf interface{}) {
	if moduleConf, has := c.Module[module.name]; has {
		moduleConf.Config = conf
		c.Module[module.name] = moduleConf
	}
}
func (c *config) load() error {
	if c.configName == "" {
		return errors.New("配置文件名为空")
	}
	c.Module = make(map[string]ModuleConfig)
	file, err := ioutil.ReadFile(getConfigPath(c.configName))
	if err != nil {
		return err
	}
	return json.Unmarshal(file, c)
}

type ModuleConfig struct {
	Enable bool        `json:"enable"`
	Config interface{} `json:"config"`
}

func ExecPathDir() string {
	if appPath == "" {
		file, err := exec.LookPath(os.Args[0])
		if err != nil {
			return ""
		}
		appPath, _ = filepath.Abs(file)
	}
	return path.Dir(appPath)
}
func (c *config) save() error {
	marshal, err2 := json.Marshal(c)
	if err2 != nil {
		return err2
	}
	return ioutil.WriteFile(getConfigPath(c.configName), marshal, 0644)
}
func Unmarshal(src, dst interface{}) error {
	marshal, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(marshal, dst)
}
func getConfigPath(configFileName string) string {
	_, err := os.Stat(configFileName)
	if err == nil {
		return configFileName
	}
	exp := fmt.Sprintf("%s%c%s", ExecPathDir(), os.PathSeparator, configFileName)
	_, err = os.Stat(exp)
	if err == nil {
		os.Chdir(ExecPathDir())
		return exp
	}
	return configFileName
}
