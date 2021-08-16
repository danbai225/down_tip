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
	Module     map[string]interface{} `json:"module"`
	configName string
}

func (c *config) getConfig(module *Module) interface{} {
	return c.Module[module.name]
}
func (c *config) saveConfig(module *Module, conf interface{}) {
	c.Module[module.name] = conf
}
func (c *config) load() error {
	if c.configName == "" {
		return errors.New("配置文件名为空")
	}
	c.Module = make(map[string]interface{})
	file, err := ioutil.ReadFile(getConfigPath(c.configName))
	if err != nil {
		return err
	}
	return json.Unmarshal(file, c)
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
