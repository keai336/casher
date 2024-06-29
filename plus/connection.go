package plus

import (
	"fmt"
	"github.com/juju/errors"
	"github.com/obgnail/clash-api/clash"
	"strings"
	"sync"
)

func GetConnecions() ([]*Connection, error) {
	container := struct {
		Connections []*Connection `json:"connections"`
	}{}
	err := clash.UnmarshalRequest("get", "/connections", nil, nil, &container)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return container.Connections, nil
}

type Connection struct {
	Id   string `json:"id"`
	Rule string `json:"rulePayload"`
}

func GetSpconnections(key string) ([]*Connection, error) {
	if connections, err := GetConnecions(); err != nil {
		return nil, err
	} else {
		fi_ls := []*Connection{}
		for _, v := range connections {
			if strings.Contains(v.Rule, key) {
				fi_ls = append(fi_ls, v)

			}
		}
		return fi_ls, nil
	}

}

func Killspconnections(connections []*Connection) {
	wai := sync.WaitGroup{}
	for _, v := range connections {
		wai.Add(1)
		go func() {
			defer wai.Done()
			if err := KillSpconnection(v.Id); err != nil {
				fmt.Println(err)
				return
			}
		}()

	}
	wai.Wait()
}

func KillSpconnection(id string) error {
	route := "/connections/" + id
	code, _, err := clash.EasyRequest("delete", route, nil, nil)
	if code != 200 || err != nil {
		return err
	}
	fmt.Println("成功", id)
	return nil
}
func KillConnection() error {
	code, _, err := clash.EasyRequest("delete", "/connections", nil, nil)
	if code != 200 || err != nil {
		return err
	}
	return nil
}
