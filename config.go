package main

import (
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Bind        string                       `yaml:"bind"`
	Timezone    string                       `yaml:"timezone"`
	TimeOut     time.Duration                `yaml:"timeout"`
	Vendors     map[string]map[string]string `yaml:"vendors"`
	Juheapikey  string                       `yaml:"juheapikey"`
	ServiceList map[string]*ServiceConfig    `yaml:"servicelist"`
	Errormsg    map[string]string            `yaml:"errormsg"`
}

type ServiceConfig struct {
	Vendor      string   `yaml:"vendor"`
	Group       string   `yaml:"group"`
	Tpl         string   `yaml:"smstpl"`
	Signname    string   `yaml:"signname"`
	Allowcity   []string `yaml:"allowcitys"`
	MaxSendNums uint64   `yaml:"maxsendnums"` //每天最大发送数量
	Callback    string   `yaml:"callback"`    //成功后回调URL
	Mode        byte     `yaml:"mode"`
	Validtime   int64    `yaml:"validtime"`
	Outformat   string   `yaml:"outformat"`
}

func (cfg *Config) ParseConfigData(data []byte) error {
	if err := yaml.Unmarshal([]byte(data), &cfg); err != nil {
		return err
	}
	return nil
}

func (cfg *Config) ParseConfigFile(fileName string) error {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	return cfg.ParseConfigData(data)
}
