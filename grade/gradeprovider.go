package grade

import (
	"bufio"
	"dash/plus"
	"fmt"
	"github.com/obgnail/clash-api/clash"
	"log"
	"os"
	"regexp"
	"sync"
)

type GradeProvider struct {
	*plus.Provider
	Level        float64
	GradeProxies map[string]*GradeProxy
}

func NewGradeProvider(name string, level float64, proxymarksdic map[string][]string) *GradeProvider {
	gradeprovider := new(GradeProvider)
	provider, err := plus.GetProviderMessage(name)
	if err != nil {
		panic("wrong getprovidermessage")
	}
	gradeprovider.Provider = provider
	gradeprovider.Level = level
	gradeprovider.InitProxies(proxymarksdic)
	return gradeprovider

}

func GetAllGradeProviders(providerleveldic map[string]float64, proxymarksdic map[string][]string) map[string]*GradeProvider {
	providers := make(map[string]*GradeProvider)
	if providerss, err := plus.GetProviders(); err != nil {
		panic(err)
	} else {
		for k := range providerss {
			if v, ok := providerleveldic[k]; ok {
				providers[k] = NewGradeProvider(k, v, proxymarksdic)
			} else {
				providers[k] = NewGradeProvider(k, 1, proxymarksdic)

			}

		}
	}
	return providers
}

func (gradeprovider *GradeProvider) InitProxies(proxiesmarksdic map[string][]string) {
	proxies := gradeprovider.Proxies
	gradeproxies := make(map[string]*GradeProxy)
	var lock sync.Mutex
	var wait sync.WaitGroup
	getblacklist := func() []*regexp.Regexp {
		var lines []*regexp.Regexp
		file, err := os.Open("blacklist")
		if err != nil {
			fmt.Printf("无法打开文件: %v", err)
			defer file.Close()
			return lines
		}

		// 使用 bufio.Scanner 逐行读取文件
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lines = append(lines, regexp.MustCompile(scanner.Text()))
		}

		// 检查是否在读取过程中遇到错误
		if err := scanner.Err(); err != nil {
			log.Fatalf("读取文件时出错: %v", err)
		}
		return lines
	}
	blacklist := getblacklist()
	for _, v := range proxies {
		wait.Add(1)
		go func(v *clash.Proxy) {
			defer wait.Done()
			for _, i := range blacklist {
				if i.MatchString(v.Name) {
					fmt.Println(v.Name, "is excluded")
					return
				}
			}
			gradeproxy := NewGradeProxy(v, proxiesmarksdic[v.Name])
			gradeproxy.Provider = gradeprovider
			lock.Lock()
			defer lock.Unlock()
			gradeproxies[v.Name] = gradeproxy

		}(v)
	}
	wait.Wait()
	gradeprovider.GradeProxies = gradeproxies
}

type Range struct {
	minDelay int
	maxDelay int
	minScore int
	maxScore int
}

func getScore(delay int, max int) int {
	// 特殊处理延迟为0的情况
	if delay == 0 {
		return 0
	}

	// 定义延迟范围和对应的评分范围
	ranges := []Range{
		{0, 50, 100, 100},
		{50, 100, 90, 100},
		{100, 200, 80, 90},
		{200, 300, 70, 80},
		{300, 500, 50, 70},
		{500, 800, 30, 50},
		{800, max, 0, 30},
		{max, -1, 0, 0}, // max 到 n 的情况，这里 -1 表示无上限
	}

	for _, r := range ranges {
		if delay >= r.minDelay && (delay < r.maxDelay || r.maxDelay == -1) {
			if r.minScore == r.maxScore {
				return r.minScore
			}
			// 线性插值计算评分
			return r.minScore + (r.maxScore-r.minScore)*(delay-r.minDelay)/(r.maxDelay-r.minDelay)
		}
	}

	return 0 // 默认返回 0 分
}

func (gradeprovider *GradeProvider) GiveScore(gradeproxy *GradeProxy) {
	relativeVariance := func(data []int) (float64, float64) {
		n := len(data)
		if n == 0 {
			return 0, 0
		}

		// 计算平均值
		sum := 0
		for _, x := range data {
			sum += x
		}
		mean := float64(sum) / float64(n)

		// 计算方差
		sumOfSquares := 0.0
		for _, x := range data {
			diff := float64(x) - mean
			sumOfSquares += diff * diff
		}
		variance := sumOfSquares / float64(n)

		// 计算相对方差
		relativeVariance := variance / (mean * mean)
		return mean, relativeVariance
	}
	var delaypoint float64
	var avpoint float64
	//if gradeproxy.DelayNow != 0 {
	//	delaypoint = 10 * float64(plus.TimeOut) / float64(gradeproxy.DelayNow) * 0.6
	//	av, rv := relativeVariance(gradeproxy.DelayHistory)
	//	avpoint = 10 * float64(plus.TimeOut) / av * 0.4
	//	//fmt.Println(rv)
	//	if 1-3*rv <= 0 {
	//		delaypoint = 0
	//	}
	//	delaypoint = delaypoint * (1 - 2*rv)
	//}
	//mark := gradeprovider.Level * gradeproxy.Level * (delaypoint + avpoint)
	//gradeproxy.Point = int(mark)
	delaypoint = float64(getScore(gradeproxy.DelayNow, plus.TimeOut))
	av, rv := relativeVariance(gradeproxy.DelayHistory)
	avpoint = float64(getScore(int(av), plus.TimeOut))
	if wfdkxk := 1 - 2*rv; wfdkxk <= 0 {
		delaypoint = 0
	} else {
		delaypoint = delaypoint * wfdkxk
	}
	mark := gradeprovider.Level * gradeproxy.Level * (0.4*delaypoint + 0.6*avpoint)
	gradeproxy.Point = int(mark)

}
func (gradeprovider *GradeProvider) Update() {
	provider, err := plus.GetProviderMessage(gradeprovider.Name)
	if err != nil {
		panic("wrong getprovider message")
	}
	gradeprovider.Provider = provider
	var wait sync.WaitGroup

	for _, v := range gradeprovider.GradeProxies {
		wait.Add(1)
		go func(v *GradeProxy) {
			defer wait.Done()
			v.Update()
			gradeprovider.GiveScore(v)
		}(v)
	}
	wait.Wait()
}
