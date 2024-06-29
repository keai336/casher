package plus

import (
	"fmt"
	lockcheck "github.com/keai336/MediaUnlockTest"
	"github.com/obgnail/clash-api/clash"
	"testing"
)

func TestGroups(t *testing.T) {
	clash.SetURL("http://10.18.18.31:9090")
	clash.SetSecret("D1u5ETt5")
	Groups, err := GetGroups()
	fmt.Println(Groups)
	if err != nil {
		fmt.Println(err)

	} else {

		for k, _ := range Groups {
			if k == "美国" {

				group, _ := GetGroupMessage(k)
				for _, v := range group.All {
					connections, _ := GetSpconnections("Spotify")
					Killspconnections(connections)
					clash.SwitchProxy("美国", v)
					rs := lockcheck.Spotify(lockcheck.AutoHttpClient)
					//rs := lockcheck.ChatGPT(lockcheck.AutoHttpClient)
					fmt.Println(v, rs.Status)
				}
			}
		}
	}
	if message, err := GetGroupMessage("EMBY"); err != nil {
		t.Log(11)
	} else {
		t.Log(message.All)
	}
}
