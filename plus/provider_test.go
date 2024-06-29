package plus

import (
	"github.com/obgnail/clash-api/clash"
	"testing"
)

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
