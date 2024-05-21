package main

import (
	"fmt"
	"github.com/obgnail/clash-api/clash"
	"sync"
	"testing"
)

func AllDelayTest(proxies map[string]*clash.Proxies) map[string]int {
	fidic := make(map[string]int)
	var wait sync.WaitGroup
	for k, _ := range proxies {
		wait.Add(1)
		go func(k string) {
			if delay, err := clash.GetProxyDelay(k, "http://www.gstatic.com/generate_204", 1000); err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("%s : %dms\n", k, delay.Delay)
				fidic[k] = delay.Delay
			}
			wait.Done()
		}(k)
	}
	wait.Wait()
	return fidic

}

func TestA(t *testing.T) {
	clash.SetURL("http://10.18.18.31:9090")
	//err := clash.GetTraffic(func(traffic *clash.Traffic) (stop bool) {
	//	time.Sleep(5 * time.Second)
	//	return true
	clash.SetSecret("D1u5ETt5")
	//})
	//全部代理测试
	proxies, _ := clash.GetProxies()
	//var wait sync.WaitGroup
	for _, k := range proxies {
		fmt.Printf("%s : now: %s  type:%s\n ", k.Name, k.Now, k.Type)
	}
	//fi := AllDelayTest(proxies)
	//fmt.Println(fi)

	proxy, _ := clash.GetProxyMessage("Final")
	fmt.Println(proxy.Name)
	fmt.Println(proxy.Type)
	fmt.Println(proxy.History)
	fmt.Println(proxy.UDP)
	delay, _ := clash.GetProxyDelay("Final", "http://www.gstatic.com/generate_204", 1000)
	println(delay.Message)
	rule, _ := clash.GetRules()
	for _, v := range rule {
		fmt.Println(v, "")
	}
	configs, _ := clash.GetConfigs()
	for k, v := range configs {
		fmt.Println(k, v)
	}
	//fmt.Println(configs)
}
