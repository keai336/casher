package config

import (
	"bufio"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	ProviderLevel map[string]float64            `yaml:"providerlevel"`
	GroupLabelDic map[string]map[string]float64 `yaml:"grouplabeldic"`
	ProxyMark     map[string][]string           `yaml:"proxymark"`
	GroupLevelDic map[string]float64            `yaml:"groupleveldic"`
	Interval      int                           `yaml:"interval"`
	TimeOut       int                           `yaml:"timeout"`
}

func NewConfig(provider map[string]float64, group map[string]map[string]float64, proxy map[string][]string, grouplevel map[string]float64) *Config {
	config := new(Config)
	config.ProviderLevel = provider
	config.GroupLabelDic = group
	config.ProxyMark = proxy
	config.GroupLevelDic = grouplevel
	return config
}

func LoadConfig(path string) *Config {
	config := new(Config)
	content, err := ioutil.ReadFile(path)

	if err != nil {
		log.Fatal(err)
	}
	if err := yaml.Unmarshal(content, config); err != nil {
		log.Fatal("load err", err)
	}
	if config.Interval == 0 {
		config.Interval = 180
		fmt.Println("no config for interval,use default 180")
	}
	if config.TimeOut == 0 {
		config.TimeOut = 1500
		fmt.Println("no config for timeout,use default 1500")
	}
	return config
}

func DumpConfig(config *Config, path string) {
	v, err := yaml.Marshal(&config)
	if err != nil {
		log.Fatal("marshal err", err)
	}
	file, _ := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0664)

	defer file.Close()

	// 获取bufio.Writer实例
	writer := bufio.NewWriter(file)

	// 写入字符串
	if _, err = writer.Write(v); err != nil {
		log.Fatal("write err", err)
	}

	// 清空缓存 确保写入磁盘
	if err = writer.Flush(); err != nil {
		log.Fatal("wrong flush", err)
	}
}
