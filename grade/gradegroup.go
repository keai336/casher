package grade

import (
	"dash/plus"
	"fmt"
	lockcheck "github.com/keai336/MediaUnlockTest"
	"github.com/obgnail/clash-api/clash"
	"net/http"
	"time"
)

type GradeGroup struct {
	*plus.Group
	Level    float64
	LabelDic map[string]float64
	Points   map[string]int
	Block    time.Time
	Source   map[string]*GradeProxy
}

func NewGradeGroup(name string, level float64, labledic map[string]float64, source map[string]*GradeProxy) *GradeGroup {
	gradegroup := new(GradeGroup)
	group, err := plus.GetGroupMessage(name)
	if err != nil {
		panic("wrong get group message")
	}
	pdic := make(map[string]int)
	ldic := make(map[string]float64)
	if labledic != nil {
		ldic = labledic
	}
	gradegroup.Group = group
	gradegroup.Level = level
	gradegroup.Points = pdic
	gradegroup.LabelDic = ldic
	gradegroup.Source = source
	gradegroup.LabelDic = labledic
	return gradegroup
}

func (gradegroup *GradeGroup) SetLabelDic(label string, value float64) {
	gradegroup.LabelDic[label] = value
}
func checklock(f func(client http.Client) lockcheck.Result) lockcheck.Result {
	return OneLockTest(f)
}
func checkmark(gradeproxy *GradeProxy, gradegroup *GradeGroup) float64 {
	marks := gradeproxy.Mark
	labeldic := gradegroup.LabelDic
	p := 0.0
	n := 0
	for _, v := range marks {
		if value, ok := labeldic[v]; ok {
			p += value
			n += 1
		}
	}
	if n == 0 {
		return 1
	}
	return p / float64(n)
}
func (gradegroup *GradeGroup) GiveScore() {
	for _, name := range gradegroup.All {
		if rootname := gradegroup.getroot(name); rootname != "" {
			gradeproxy := gradegroup.Source[rootname]
			priority := checkmark(gradeproxy, gradegroup)
			gradegroup.Points[name] = int(float64(gradeproxy.Point) * priority)

		}

	}
}

func (gradegroup *GradeGroup) getroot(name string) string {
	if _, ok := gradegroup.Source[name]; ok {
		return name
	} else {
		groupmsg, _ := plus.GetGroupMessage(name)
		name = groupmsg.Now
		if name == "" {
			return name
		}
		return gradegroup.getroot(name)
	}

}

// func (gradegroup *GradeGroup)CheckWhereLock()[]string{
//
// }
func (gradegroup *GradeGroup) CheckLock() float64 {
	n := 0
	value := 0
	total := 0
	for k, _ := range gradegroup.LabelDic {
		if checkfunc, ok := LockTestDic[k]; ok {
			total++
			rs := checklock(checkfunc)
			if rs.Status == -1 {
				value = 0
			} else {
				value = 1
			}
			n = n + value
		}

	}
	if total == 0 {
		return -1
	}
	return 2 * float64(n) / float64(total)
}
func (gradegroup *GradeGroup) GetlockScore(unlock map[string]int) {

}
func (gradegroup *GradeGroup) FindBest(unlock map[string]int) {
	nowuse := gradegroup.Now
	//fmt.Println(nowuse)
	switch v := gradegroup.CheckLock(); v {
	case 1:
		unlock[nowuse] = gradegroup.Points[nowuse]
		return
	case -1:
		unlock[nowuse] = gradegroup.Points[nowuse]
	default:
		unlock[nowuse] = int(float64(gradegroup.Points[nowuse]) * v * v * v)
	}
	delete(gradegroup.Points, nowuse)
	if len(gradegroup.Points) != 0 {
		name, _ := maxInMap(gradegroup.Points)
		clash.SwitchProxy(gradegroup.Name, name)
		//fmt.Println("切换到", name)
		gradegroup.UpdateStu()
		gradegroup.FindBest(unlock)
	} else {
		name, _ := maxInMap(unlock)
		clash.SwitchProxy(gradegroup.Name, name)
		//fmt.Println(gradegroup.Name, "最优解", name, value)
		gradegroup.UpdateStu()
		return
	}

}

//func (gradegroup *GradeGroup) ChangeIf() {
//	if !gradegroup.Block.After(time.Now()) {
//		//fmt.Println(gradgroup.Name, "可更改")
//		nowuse := gradegroup.Now
//		nowpoint := gradegroup.Points[nowuse]
//		name, value := maxInMap(gradegroup.Points)
//		nowroot := gradegroup.getroot(nowuse)
//		var nowdelay int
//		if now, ok := gradegroup.Source[nowroot]; ok {
//			nowdelay = now.DelayNow
//		}
//		if value > int(float64(nowpoint)*1.3) {
//			err := clash.SwitchProxy(gradegroup.Name, name)
//			if err != nil {
//				fmt.Println("切换失败", err)
//				return
//			}
//			if gradegroup.CheckLock()!=1{
//
//			}
//			//log
//			if MysqlOn {
//				OneInsertChangeHistroy(gradegroup.Name, nowuse, nowdelay, name, gradegroup.Source[gradegroup.getroot(name)].DelayNow)
//			} else {
//				fmt.Printf("%s %s old:%s-延迟%d-分数%d --> new:%s-延迟%d-分数%d\n", time.Now().Format("2006-01-02 15:04:05"), gradegroup.Name, nowuse, nowdelay, nowpoint, name, gradegroup.Source[gradegroup.getroot(name)].DelayNow, value)
//			}
//		}
//
//	} else {
//		//fmt.Println("无法操作")
//	}
//
//}

func (gradegroup *GradeGroup) ChangeIf() {
	if !gradegroup.Block.After(time.Now()) {
		//fmt.Println(gradgroup.Name, "可更改")
		nowuse := gradegroup.Now
		nowpoint := gradegroup.Points[nowuse]
		lockdic := make(map[string]int)
		v := gradegroup.CheckLock()
		switch v {
		case 0:
			lockdic[nowuse] = 1
		case -1:
			lockdic[nowuse] = gradegroup.Points[nowuse]
		default:
			lockdic[nowuse] = int(float64(gradegroup.Points[nowuse]) * v * v * v)

		}
		delete(gradegroup.Points, nowuse)
		if len(gradegroup.Points) == 0 {
			return
		}
		name, _ := maxInMap(gradegroup.Points)
		clash.SwitchProxy(gradegroup.Name, name)
		gradegroup.UpdateStu()
		gradegroup.FindBest(lockdic)
		nowpoint = lockdic[nowuse]
		newuse := gradegroup.Now
		newpoint := lockdic[newuse]
		//1.3 为僭越值,目的是保当前使用
		var jmyt float64
		switch v {
		case -1, 1:
			jmyt = 1.3
		default:
			jmyt = 1
		}

		if newpoint > int(float64(nowpoint)*jmyt) {
			fmt.Println("use", newuse, newpoint)
		} else {
			clash.SwitchProxy(gradegroup.Name, nowuse)
			gradegroup.UpdateStu()
		}

	} else {
		//fmt.Println("无法操作")
	}

}
func (gradegroup *GradeGroup) UpdateStu() {
	gradegroup.Group, _ = plus.GetGroupMessage(gradegroup.Name)
}
func (gradegroup *GradeGroup) Update() {
	gradegroup.UpdateStu()
	gradegroup.GiveScore()
	gradegroup.ChangeIf()
}
