package grade

import (
	"dash/plus"
	"fmt"
	"github.com/obgnail/clash-api/clash"
)

type GradeProxy struct {
	*clash.Proxy
	Level        float64
	Provider     *GradeProvider
	DelayNow     int
	DelayHistory []int
	Mark         []string
	Point        int
}

func (gradeproxy *GradeProxy) SetMark(mark string) {
	containif := func(slice []string, str string) bool {
		for _, s := range slice {
			if s == str {
				return true
			}
		}
		return false
	}
	if !containif(gradeproxy.Mark, mark) {
		gradeproxy.Mark = append(gradeproxy.Mark, mark)
		return
	}
	fmt.Printf("The marker %s already exists", mark)

}
func (gradeproxy *GradeProxy) Update() {
	var delay *clash.ProxyDelay
	delay, _ = plus.GeneralDelayTest(gradeproxy.Name)
	if delay.Delay == 0 {
		//给你机会别不中用
		delay, _ = plus.GeneralDelayTest(gradeproxy.Name)
	}
	gradeproxy.DelayNow = delay.Delay
	gradeproxy.DelayHistory = append(gradeproxy.DelayHistory, gradeproxy.DelayNow)
	if MysqlOn {
		//OneInsertHistory(gradeproxy)
	}
}
func NewGradeProxy(proxy *clash.Proxy, marks []string) *GradeProxy {
	gradeproxy := new(GradeProxy)
	gradeproxy.Proxy = proxy
	gradeproxy.Level = 1
	delay, _ := plus.GeneralDelayTest(gradeproxy.Name)
	gradeproxy.DelayNow = delay.Delay
	gradeproxy.DelayHistory = append(gradeproxy.DelayHistory, gradeproxy.DelayNow)
	if marks == nil {
		MarkProxy(LocationMap, gradeproxy)

	} else {
		gradeproxy.Mark = marks
	}
	return gradeproxy
}

func InitSource(gradeproviders map[string]*GradeProvider) map[string]*GradeProxy {
	source := make(map[string]*GradeProxy)
	for _, v := range gradeproviders {
		for _, gradeproxies := range v.GradeProxies {
			source[gradeproxies.Name] = gradeproxies

		}

	}
	return source
}

func InitGradeGroup(grouplevel map[string]float64, grouplabeldic map[string]map[string]float64, source map[string]*GradeProxy) map[string]*GradeGroup {
	groups, _ := plus.GetGroups()
	gradegroups := make(map[string]*GradeGroup)

	for k := range groups {
		if v, ok := grouplevel[k]; ok {
			gradegroups[k] = NewGradeGroup(k, v, grouplabeldic[k], source)
		} else {
			gradegroups[k] = NewGradeGroup(k, 1, grouplabeldic[k], source)
		}
	}
	return gradegroups
}
