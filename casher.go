package main

import (
	"dash/config"
	"dash/grade"
	"fmt"
	"github.com/obgnail/clash-api/clash"
)

type Casher struct {
	Providers map[string]*grade.GradeProvider
	Groups    map[string]*grade.GradeGroup
	Source    map[string]map[string]*grade.GradeProxy
	*config.Config
	ConfigPath string
}

func (casher *Casher) Init() {
	casher.Providers = grade.GetAllGradeProviders(casher.ProviderLevel, casher.ProxyMark)
	casher.Source = grade.InitSource(casher.Providers)
	casher.Groups = grade.InitGradeGroup(casher.GroupLevelDic, casher.GroupLabelDic, casher.Source)
}
func NewOneCasher(configpt string) *Casher {
	casher := new(Casher)
	clash.SetURL("http://10.18.18.31:9090")
	clash.SetSecret("D1u5ETt5")
	casher.ConfigPath = configpt
	casher.Config = config.LoadConfig(configpt)
	casher.Init()
	return casher

}
func (casher *Casher) Update() {
	grade.AllUpdate(casher.Providers)
	grade.AllUpdate(casher.Groups)
}

func (casher *Casher) OffDuty() {
	removeEmptyValues := func(m interface{}) {
		switch v := m.(type) {
		case map[string]float64:
			for k, _ := range v {
				if v[k] == 0.0 {
					delete(v, k)
				}
			}
		case map[string][]string:
			for k, _ := range v {
				if len(v[k]) == 0 {
					delete(v, k)
				}

			}
		case map[string]map[string]float64:
			for k, _ := range v {
				if len(v[k]) == 0 {
					delete(v, k)
				}
			}

		}
	}

	removeDuplicates := func(slice []string) []string {
		// 创建一个 map 用于记录已经出现过的元素
		seen := make(map[string]bool)

		// 创建一个新的切片来保存去重后的结果
		result := []string{}

		// 遍历切片中的每个元素
		for _, value := range slice {
			// 如果元素在 map 中不存在，说明是第一次出现，则将其加入结果切片中，并在 map 中标记为已出现
			if !seen[value] {
				result = append(result, value)
				seen[value] = true
			}
		}

		return result
	}

	for _, k := range casher.Providers {
		casher.ProviderLevel[k.Name] = k.Level
	}
	removeEmptyValues(casher.ProviderLevel)
	for _, k := range casher.Groups {
		casher.GroupLevelDic[k.Name] = k.Level
		casher.GroupLabelDic[k.Name] = k.LabelDic
	}
	removeEmptyValues(casher.GroupLevelDic)
	removeEmptyValues(casher.GroupLabelDic)
	for _, k := range casher.Source {
		for _, proxy := range k {
			casher.ProxyMark[proxy.Name] = removeDuplicates(proxy.Mark)
		}
	}
	removeEmptyValues(casher.ProxyMark)
	config.DumpConfig(casher.Config, casher.ConfigPath)
}

func (casher *Casher) SetProxyMark(proxyname string, mark string) {
	for _, proxiesmap := range casher.Source {
		if proxy, ok := proxiesmap[proxyname]; ok {
			proxy.SetMark(mark)
			return

		}
	}
	fmt.Println("no proxy named", proxyname)
}

func (casher *Casher) DelProxyMark(proxyname string, mark string) {
	removeValue := func(slice []string, value string) []string {
		// 从切片末尾开始查找指定值元素的索引
		for i := len(slice) - 1; i >= 0; i-- {
			if slice[i] == value {
				// 将找到的指定值元素从切片中删除
				copy(slice[i:], slice[i+1:])
				// 调整切片的长度
				slice = slice[:len(slice)-1]
				return slice
			}
		}
		fmt.Printf("no mark named %s", mark)
		return slice
	}

	for _, proxiesmap := range casher.Source {
		if proxy, ok := proxiesmap[proxyname]; ok {
			proxy.Mark = removeValue(proxy.Mark, mark)
		}
		return

	}
	fmt.Println("no proxy named", proxyname)
}
