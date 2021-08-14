package service

import (
	"encoding/json"
	"github.com/axgle/mahonia"
	"io/ioutil"
	"net/http"
)

type ipInfo struct {
	IP          string `json:"ip"`
	Pro         string `json:"pro"`
	ProCode     string `json:"proCode"`
	City        string `json:"city"`
	CityCode    string `json:"cityCode"`
	Region      string `json:"region"`
	RegionCode  string `json:"regionCode"`
	Addr        string `json:"addr"`
	RegionNames string `json:"regionNames"`
	Err         string `json:"err"`
}

func GetIpInfo(ip string) ipInfo {
	info := ipInfo{}
	get, err := http.Get("http://whois.pconline.com.cn/ipJson.jsp?json=true&ip=" + ip)
	if err != nil {
		println(err)
	} else {
		all, _ := ioutil.ReadAll(get.Body)
		json.Unmarshal([]byte(mahonia.NewDecoder("gbk").ConvertString(string(all))), &info)
	}
	return info
}
