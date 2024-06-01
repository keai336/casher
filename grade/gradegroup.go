package grade

import (
	"dash/plus"
	"fmt"
	lockcheck "github.com/keai336/MediaUnlockTest"
	"github.com/obgnail/clash-api/clash"
	"net/http"
	"sync"
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
	var n, total int
	var mu sync.Mutex
	var wg sync.WaitGroup

	for k := range gradegroup.LabelDic {
		if checkfunc, ok := LockTestDic[k]; ok {
			wg.Add(1)
			go func(cf func(r http.Client) lockcheck.Result) {
				defer wg.Done()
				rs := checklock(cf)
				value := 0
				if rs.Status != -1 {
					value = 1
				}
				mu.Lock()
				n += value
				mu.Unlock()
			}(checkfunc)
			total++
		}
	}

	wg.Wait()

	if total == 0 {
		return -1
	}
	return float64(n) / float64(total)
}
func (gradegroup *GradeGroup) GetlockScore(unlock map[string]int) {

}
func (group *GradeGroup) getlockchecknum() int {
	count := 0
	for key := range group.LabelDic {
		if _, exists := LockTestDic[key]; exists {
			count++
		}
	}
	return count
}

func (proxy *GradeProxy) getlockchecknum() int {
	count := 0
	for _, key := range proxy.Mark {
		if _, exists := LockTestDic[key]; exists {
			count++
		}
	}
	return count
}
func (gradegroup *GradeGroup) FindBest(unlock map[string]int) {
	switch num := gradegroup.getlockchecknum(); num {
	case 0:
		name, _ := maxInMap(gradegroup.Points)
		clash.SwitchProxy(gradegroup.Name, name)
		gradegroup.UpdateStu()
		return
	default:
		nowuse := gradegroup.Now
		switch v := gradegroup.CheckLock(); v {
		case 1:
			unlock[nowuse] = gradegroup.Points[nowuse]
			return
		case -1:
			unlock[nowuse] = 0
		default:
			unlock[nowuse] = int(float64(gradegroup.Points[nowuse]) * v * v * v)
		}
		delete(gradegroup.Points, nowuse)
		if len(gradegroup.Points) != 0 {
			var next string
			var mark int = 1
			for range gradegroup.Points {
				next, _ = maxInMap(gradegroup.Points)
				if pcn := gradegroup.Source[next].getlockchecknum(); pcn != 0 {
					mark = mark * 0
					break
				} else {
					unlock[next] = gradegroup.Points[next]
					delete(gradegroup.Points, next)
				}

			}
			switch mark {
			case 0:
				clash.SwitchProxy(gradegroup.Name, next)
				gradegroup.UpdateStu()
				gradegroup.FindBest(unlock)
			case 1:
				name, _ := maxInMap(unlock)
				clash.SwitchProxy(gradegroup.Name, name)
				gradegroup.UpdateStu()
				return

			}
		} else {
			name, _ := maxInMap(unlock)
			clash.SwitchProxy(gradegroup.Name, name)
			gradegroup.UpdateStu()
			return
		}

		//	if len(gradegroup.Points) != 0 {
		//		name, _ := maxInMap(gradegroup.Points)
		//		switch gradegroup.Source[name].getlockchecknum() {
		//		case 0:
		//
		//		default:
		//			clash.SwitchProxy(gradegroup.Name, name)
		//			gradegroup.UpdateStu()
		//			gradegroup.FindBest(unlock)
		//		}
		//	} else {
		//		name, _ := maxInMap(unlock)
		//		clash.SwitchProxy(gradegroup.Name, name)
		//		gradegroup.UpdateStu()
		//		return
		//	}
		//
		//}
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
		var jmyt float64
		nowuse := gradegroup.Now
		nowpoint := gradegroup.Points[nowuse]
		var newuse string
		var newpoint int
		switch gradegroup.getlockchecknum() {
		case 0:
			name, _ := maxInMap(gradegroup.Points)
			clash.SwitchProxy(gradegroup.Name, name)
			gradegroup.UpdateStu()
			newuse = gradegroup.Now
			newpoint = gradegroup.Points[newuse]
			jmyt = 1.3

		default:
			lockdic := make(map[string]int)
			v := gradegroup.CheckLock()
			switch v {
			case 0, -1:
				lockdic[nowuse] = 1
			default:
				lockdic[nowuse] = int(float64(gradegroup.Points[nowuse]) * v * v * v)

			}
			delete(gradegroup.Points, nowuse)
			if len(gradegroup.Points) == 0 {
				return
			}
			var next string
			for range gradegroup.Points {
				next, _ = maxInMap(gradegroup.Points)
				if pcn := gradegroup.Source[next].getlockchecknum(); pcn != 0 {
					break
				} else {
					lockdic[next] = gradegroup.Points[next]
					delete(gradegroup.Points, next)
				}

			}
			clash.SwitchProxy(gradegroup.Name, next)
			gradegroup.UpdateStu()
			gradegroup.FindBest(lockdic)
			nowpoint = lockdic[nowuse]
			newuse = gradegroup.Now
			newpoint = lockdic[newuse]
			switch v {
			case -1, 1:
				jmyt = 1.3
			default:
				jmyt = 1
			}

		}
		//1.3 为僭越值,目的是保当前使用

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
