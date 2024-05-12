package plus

import (
	"github.com/obgnail/clash-api/clash"
)

func GeneralDelayTest(name string) (*clash.ProxyDelay, error) {
	return clash.GetProxyDelay(name, "http://www.gstatic.com/generate_204", 1000)
}

func GetALLDelayNow() map[string]int {
	proxies, err := clash.GetProxies()
	if err != nil {
		panic("wrong get proxies")
	}
	delay_dic := make(map[string]int)
	for k, _ := range proxies {
		go func(k string) {
			delay, _ := GeneralDelayTest(k)
			delay_dic[k] = delay.Delay
		}(k)
	}
	return delay_dic
}

//func TestAllDelay(t *testing.T) {
//
//	alldelay := AllDelayTest()
//}
