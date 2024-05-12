package main

import (
	config2 "dash/config"
	"dash/grade"
	"fmt"
	"sync"
	"time"
)

func main() {
	config := config2.Configs
	providerdic := config.ProviderLevel
	gradeproviders := grade.GetAllGradeProviders(providerdic)
	source := grade.InitSource(gradeproviders)
	groupdic := config.GroupLabelDic
	gradegroups := grade.InitGradeGroup(groupdic, source)
	one := new(grade.GradeProxy)
	var wait sync.WaitGroup
	wait.Add(1)
	go func() {
		defer wait.Done()
		for {
			grade.AllUpdate(gradeproviders)
			grade.AllUpdate(gradegroups)
			one = (*source["mesl"])["🇨🇳 TW1 Hinet [0.2X]"]
			fmt.Println(one.Name, one.Mark, one.Point, one.DelayHistory)
			//gradeprovider.Update()
			//for _, v := range gradegroups {
			//	v.SetLabelDic("香港", 2.5)
			//	v.Update()
			//	name, point := maxInMap(v.Points)
			//	fmt.Println(v.Name, name, point)
			//	//fmt.Println(v.Source)
			time.Sleep(5 * time.Second)
		}

	}()
	marks := []string{"1", "2", "3"}

	for _, v := range marks {
		one.SetMark(v)
		time.Sleep(5 * time.Second)

	}
	//for k, v := range gradeproxies {
	//	fmt.Println(k, v.Point)
	wait.Wait()
}
