package main

import (
	"fmt"
	"github.com/juju/errors"
	"github.com/obgnail/clash-api/clash"
	"testing"
)

type Group struct {
	History []*clash.History `json:"history"`
	Name    string           `json:"name"`
	Type    string           `json:"type"`
	UDP     bool             `json:"udp"`
	Now     string           `json:"now"`
	All     []string         `json:"all"`
	Alive   bool             `json:"alive"`
}

func GetGroups() (map[string]*Group, error) {
	container := struct {
		Groups []*Group `json:"proxies"`
	}{}
	err := clash.UnmarshalRequest("get", "/group", nil, nil, &container)
	if err != nil {
		return nil, errors.Trace(err)
	}
	fi_map := make(map[string]*Group)
	for _, v := range container.Groups {
		fi_map[v.Name] = v
	}

	return fi_map, nil
}

func GetGroupMessage(name string) (*Group, error) {
	group := &Group{}
	route := "/group/" + name
	err := clash.UnmarshalRequest("get", route, nil, nil, &group)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return group, nil
}

func TestGroups(t *testing.T) {
	clash.SetURL("http://10.18.18.31:9090")
	clash.SetSecret("D1u5ETt5")
	Groups, err := GetGroups()
	if err != nil {
		fmt.Println(err)

	} else {

		for k, _ := range Groups {
			if k == "Final" {

				group, _ := GetGroupMessage(k)
				for _, v := range group.History {
					println(v.Time, v.Delay, v.MeanDelay)
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

//func GetProxyMessage(name string) (*Proxy, error) {
//	proxy := &Proxy{}
//	route := "/proxies/" + name
//	err := UnmarshalRequest("get", route, nil, nil, &proxy)
//	if err != nil {
//		return nil, errors.Trace(err)
//	}
//	return proxy, nil
//}
