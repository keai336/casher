package plus

import (
	"fmt"
	"github.com/obgnail/clash-api/clash"
	"testing"
)

func TestGetconnection(t *testing.T) {
	clash.SetURL("http://10.18.18.31:9090")
	clash.SetSecret("D1u5ETt5")
	config, err := clash.GetVersion()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(config.Version)
	container, err := GetSpconnections("Spotify")
	for k, v := range container {
		fmt.Println(k, v.Rule)
	}
	Killspconnections(container)
}
