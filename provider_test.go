package main

import (
	"fmt"
	"github.com/juju/errors"
	"github.com/obgnail/clash-api/clash"
	"testing"
	"time"
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
	VehicleType  string `json:"vehicleType"`
}

func (provider *Provider) ShowFlow(i int) string {
	units := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}

	var unitIndex int
	size := float64(i)
	for size >= 1024 && unitIndex < len(units)-1 {
		size /= 1024
		unitIndex++
	}

	return fmt.Sprintf("%.2f%s", size, units[unitIndex])
}

func (provider *Provider) ShowDeadline() string {
	timestamp := int64(provider.Expire) // 例如，你可以将其替换为任何时间戳

	// 将时间戳转换为时间对象
	t := time.Unix(timestamp, 0)

	// 将时间对象格式化为字符串
	formattedTime := t.Format("2006-01-02 15:04:05") // 这里的日期和时间格式是示例，你可以根据需要更改

	// 输出格式化后的时间字符串
	//fmt.Println("Formatted Time:", formattedTime)
	return formattedTime
}

func (provider *Provider) LeftTime() string {
	if provider.Expire == 0 {
		return "expired"
	}
	timestamp1 := int64(provider.Expire) // 第一个时间戳，你可以替换为任何时间戳

	// 获取当前时间的时间戳
	timestamp2 := time.Now().Unix()

	// 计算时间差
	duration := time.Duration(timestamp1-timestamp2) * time.Second

	// 将时间差格式化为天、小时、分钟和秒
	days := int(duration.Hours() / 24)
	hours := int(duration.Hours()) % 24
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60

	// 输出格式化后的时间差
	return fmt.Sprintf("时间差：%d天 %d小时 %d分钟 %d秒\n", days, hours, minutes, seconds)
}
func GetProviders() (map[string]*Provider, error) {
	container := struct {
		Providers map[string]*Provider `json:"providers"`
	}{}
	err := clash.UnmarshalRequest("get", "/providers/proxies", nil, nil, &container)
	if err != nil {
		return nil, errors.Trace(err)
	}
	fi := make(map[string]*Provider)
	for k, v := range container.Providers {
		if v.VehicleType != "Compatible" {
			fi[k] = v

		}
	}

	return fi, nil
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
		t.Log(Groups.ShowFlow(Groups.Left))
		t.Log(Groups.ShowFlow(Groups.TotalUsed))
		t.Log(Groups.ShowFlow(Groups.Total))

		providers, err := GetProviders()
		if err != nil {
			t.Log(err)

		} else {
			t.Log(providers)
			for k, v := range providers {
				t.Logf("%s: %s,%s", k, v.ShowDeadline(), v.LeftTime())

			}
		}
	}
}
