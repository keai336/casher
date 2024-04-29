package main

import (
	"fmt"
	"github.com/juju/errors"
	"github.com/obgnail/clash-api/clash"
	"testing"
)

type Info struct {
	Up     int `json:"Upload"`
	Down   int `json:"Download"`
	Total  int `json:"Total"`
	Expire int `json:"Expire"`
}
type Provider struct {
	Proxies      []*clash.Proxy `json:"proxies"`
	Name         string         `json:"name"`
	Info         `json:"subscriptionInfo"`
	TotalUsed    int
	Left         int
	LastTestTime string `json:"updatedAt"`
}

func (provider *Provider) Show(i int) string {
	units := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}

	var unitIndex int
	size := float64(i)
	for size >= 1024 && unitIndex < len(units)-1 {
		size /= 1024
		unitIndex++
	}

	return fmt.Sprintf("%.2f%s", size, units[unitIndex])
}

func GetProviders() (map[string]*Provider, error) {
	container := struct {
		providers map[string]*Provider `json:"providers"`
	}{}
	err := clash.UnmarshalRequest("get", "/providers/proxies", nil, nil, &container)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return container.providers, nil
}

func GetProviderMessage(name string) (*Provider, error) {
	provider := &Provider{}
	route := "/providers/proxies/" + name
	err := clash.UnmarshalRequest("get", route, nil, nil, &provider)
	if err != nil {
		return nil, errors.Trace(err)
	}
	provider.TotalUsed = provider.Up + provider.Down
	provider.Left = provider.Total - provider.TotalUsed
	return provider, nil
}

func TestProvider(t *testing.T) {
	clash.SetURL("http://10.18.18.31:9090")
	clash.SetSecret("D1u5ETt5")
	Groups, err := GetProviderMessage("planc")
	if err != nil {
		t.Log(err)
	} else {
		t.Log(Groups.Name)
		t.Log(Groups.Info)
		t.Log(Groups.LastTestTime)
		//for _, v := range Groups.Proxies {
		//	fmt.Println(v.Name)
		//}
		t.Log(Groups.Show(Groups.Left))
		t.Log(Groups.Show(Groups.TotalUsed))
		t.Log(Groups.Show(Groups.Total))

		providers, err := GetProviders()
		if err != nil {
			t.Log(err)

		} else {
			t.Log(providers)
			for k, v := range providers {
				t.Logf("%s: 1%s", k, v.LastTestTime)

			}
		}
	}
}
