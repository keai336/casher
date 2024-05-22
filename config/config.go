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
	Url           string                        `yaml:"url"`
	Secret        string                        `yaml:"secret"`
}

func NewConfig(provider map[string]float64, group map[string]map[string]float64, proxy map[string][]string, grouplevel map[string]float64) *Config {
	config := new(Config)
	config.ProviderLevel = provider
	config.GroupLabelDic = group
	config.ProxyMark = proxy
	config.GroupLevelDic = grouplevel
	return config
}
func NewEmputyConfig() *Config {
	config := new(Config)
	config.ProviderLevel = make(map[string]float64)
	config.GroupLabelDic = make(map[string]map[string]float64)
	config.ProxyMark = make(map[string][]string)
	config.GroupLevelDic = make(map[string]float64)
	return config
}
func LoadConfig(path string) *Config {
	config := new(Config)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		if path != "config.yaml" {
			log.Fatal(err)
		} else {
			config = NewEmputyConfig()
			fmt.Println("初始化配置文件")
			fmt.Println("配置url")
			fmt.Scanln(&config.Url)
			fmt.Println("配置secret")
			fmt.Scanln(&config.Secret)
			fileName := "config.yaml"

			// 创建文件
			file, err := os.Create(fileName)
			if err != nil {
				//fmt.Printf("Error creating file %s: %v\n", fileName, err)
				panic(err)
			}
			defer file.Close()

			fmt.Printf("File %s created successfully.\n", fileName)
		}
	} else {
		if err := yaml.Unmarshal(content, config); err != nil {
			log.Fatal("load err", err)
		}
	}
	if config.Interval == 0 {
		config.Interval = 180
		fmt.Println("no config for interval,use default 180")
	}
	if config.TimeOut == 0 {
		config.TimeOut = 1500
		fmt.Println("no config for timeout,use default 1500")
	}
	if config.Url == "" {
		config.Url = "http://127.0.0.1:9090"
		fmt.Println("no config for Url,use default 127.0.0.1:9090")
	}

	return config
}

func DumpConfig(config *Config, path string) {
	v, err := yaml.Marshal(&config)
	if err != nil {
		log.Fatal("marshal err", err)
	}
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0664)
	if err != nil {
		panic(err)
	}

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
