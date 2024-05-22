package plus

import (
	"fmt"
	"github.com/obgnail/clash-api/clash"
	"testing"
)

func TestGeneralDelayTest(t *testing.T) {
	clash.SetURL("http://10.18.18.31:9090")
	//err := clash.GetTraffic(func(traffic *clash.Traffic) (stop bool) {
	//	time.Sleep(5 * time.Second)
	//	return true
	clash.SetSecret("D1u5ETt5")
	delay, err := GeneralDelayTest("111")
	fmt.Println(err)
	println(delay.Delay)
}
func TestAllDelay(t *testing.T) {
	clash.SetURL("http://127.0.0.1:9090")
	clash.SetSecret("adfxf")
	config, err := clash.GetVersion()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(config.Version)
}
