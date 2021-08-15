package core

import (
	"encoding/json"
	"io/ioutil"
)

var configPath = "config.json"

type config struct {
	Module map[string]interface{} `json:"module"`
}

func (c *config) getConfig(module *Module) interface{} {
	return c.Module[module.name]
}
func (c *config) saveConfig(module *Module, conf interface{}) {
	c.Module[module.name] = conf
}
func (c *config) load() error {
	c.Module = make(map[string]interface{})
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}
	return json.Unmarshal(file, c)
}
func (c *config) save() error {
	marshal, err2 := json.Marshal(c)
	if err2 != nil {
		return err2
	}
	return ioutil.WriteFile(configPath, marshal, 0644)
}
func Unmarshal(src, dst interface{}) error {
	marshal, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(marshal, dst)
}
