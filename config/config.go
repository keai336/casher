package config

import (
	"bufio"
	"fmt"
	"github.com/obgnail/clash-api/clash"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

// var GradeProviders map[string]*GradeProvider
// var Source map[string]map[string]*map[string]*GradeProxy
// var GradeGroups map[string]*GradeGroup
var Configs *Config

func init() {
	clash.SetURL("http://10.18.18.31:9090")
	clash.SetSecret("D1u5ETt5")
	Configs = LoadConfig("test.yaml")

}

type Config struct {
	ProviderLevel map[string]float64  `yaml:"providerlevel"`
	GroupLabelDic map[string]float64  `yaml:"grouplabeldic"`
	ProxyMark     map[string][]string `yaml:"proxymark"`
}

func NewConfig(provider map[string]float64, group map[string]float64, proxy map[string][]string) *Config {
	config := new(Config)
	config.ProviderLevel = provider
	config.GroupLabelDic = group
	config.ProxyMark = proxy
	return config
}

func LoadConfig(path string) *Config {
	config := new(Config)
	content, err := ioutil.ReadFile(path)

	if err != nil {
		log.Fatal(err)
	}
	yaml.Unmarshal(content, config)
	return config
}
func DumpConfig(config *Config, path string) {
	v, err := yaml.Marshal(&config)
	if err == nil {
		fmt.Println(string(v))
	}
	file, _ := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0664)

	defer file.Close()

	// 获取bufio.Writer实例
	writer := bufio.NewWriter(file)

	// 写入字符串
	writer.Write(v)

	// 清空缓存 确保写入磁盘
	writer.Flush()
}

func TestRef(t *testing.T) {
	providerleveldic := map[string]float64{"mesl": 1, "planc": 2.3}
	grouplabeldic := map[string]float64{"final": 1, "hk": 2.3}
	proxymark := map[string][]string{"final": []string{"hk", "openai"}}
	config := NewConfig(providerleveldic, grouplabeldic, proxymark)
	//DumpConfig(config, "test.yaml")
	config = LoadConfig("test.yaml")
	fmt.Println(config)
}
