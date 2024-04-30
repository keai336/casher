package main

import (
	"github.com/obgnail/clash-api/clash"
	"testing"
)

type GradeProvider struct {
	*Provider
	Level int
}

func NewGradeProvider(name string, level int) *GradeProvider {
	gradeprovider := new(GradeProvider)
	provider, err := GetProviderMessage(name)
	if err != nil {
		panic("wrong getprovidermessage")
	}
	gradeprovider.Provider = provider
	gradeprovider.Level = level
	return gradeprovider

}

func (gradeprovider *GradeProvider) InitProxies() map[string]*GradeProxy {
	proxies := gradeprovider.Proxies
	gradeproxies := make(map[string]*GradeProxy)
	for _, v := range proxies {
		gradeproxy := NewGradeProxy(v)
		gradeproxy.Provider = gradeprovider
		gradeproxies[v.Name] = gradeproxy

	}
	return gradeproxies
}

type GradeGroup struct {
	*Group
	Level    int
	LabelDic map[string]int
}

type GradeProxy struct {
	*clash.Proxy
	Level    int
	Provider *GradeProvider
}

func NewGradeProxy(proxy *clash.Proxy) *GradeProxy {
	gradeproxy := new(GradeProxy)
	gradeproxy.Proxy = proxy
	gradeproxy.Level = 1
	return gradeproxy
}

func TestGradeprovider(t *testing.T) {
	clash.SetURL("http://10.18.18.31:9090")
	clash.SetSecret("D1u5ETt5")
	gradeprovider := NewGradeProvider("mesl", 1)

	for k, v := range gradeprovider.InitProxies() {

		t.Log(k, v.Provider.Name, v.Level, v.Provider.ShowFlow(v.Provider.Left))

	}
}
