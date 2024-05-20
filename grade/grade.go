package grade

import (
	"dash/config"
	"dash/plus"
	"fmt"
	"github.com/obgnail/clash-api/clash"
	"regexp"
	"strings"
	"sync"
	"time"
)

func GetLocationMarkReMap() map[string]*regexp.Regexp {
	s := "美国: US|United ?States|USA|美国 @ 英国: UK|United ?Kingdom|英国 @ 台湾: TW|台湾|Taiwan @ 香港: HK|Hong ?Kong|香港 @ 日本: JP|Japan|日本 @ 新加坡: SG|Singapore|新加坡 @ 韩国: KR|Korea|South ?Korea|Republic ?of ?Korea|韩国 @ 法国: FR|France|法国 @ 德国: DE|Germany|德国 @ 意大利: IT|Italy|意大利 @ 西班牙: ES|Spain|西班牙 @ 澳大利亚: AU|Australia|澳大利亚 @ 巴西: BR|Brazil|巴西 @ 加拿大: CA|Canada|加拿大 @ 俄罗斯: RU|Russia|俄罗斯 @ 印度: IN|India|印度 @ 墨西哥: MX|Mexico|墨西哥 @ 荷兰: NL|Netherlands|荷兰 @ 瑞士: CH|Switzerland|瑞士 @ 瑞典: SE|Sweden|瑞典 @ 挪威: NO|Norway|挪威 @ 丹麦: DK|Denmark|丹麦 @ 芬兰: FI|Finland|芬兰 @ 新西兰: NZ|New Zealand|新西兰 @ 阿根廷: AR|Argentina|阿根廷 @ 南非: ZA|South Africa|南非 @ 希腊: GR|Greece|希腊 @ 土耳其: TR|Turkey|土耳其 @ 泰国: TH|Thailand|泰国 @ 马来西亚: MY|Malaysia|马来西亚 @ 印度尼西亚: ID|Indonesia|印度尼西亚 @ 菲律宾: PH|Philippines|菲律宾 @ 以色列: IL|Israel|以色列 @ 沙特阿拉伯: SA|Saudi Arabia|沙特阿拉伯 @ 阿联酋: AE|United Arab Emirates|阿联酋 @ 埃及: EG|Egypt|埃及 @ 尼日利亚: NG|Nigeria|尼日利亚 @ 肯尼亚: KE|Kenya|肯尼亚 @ 摩洛哥: MA|Morocco|摩洛哥 @ 越南: VN|Vietnam|越南 @ 智利: CL|Chile|智利 @ 秘鲁: PE|Peru|秘鲁 @ 哥伦比亚: CO|Colombia|哥伦比亚 @ 委内瑞拉: VE|Venezuela|委内瑞拉 @ 玻利维亚: BO|Bolivia|玻利维亚 @ 厄瓜多尔: EC|Ecuador|厄瓜多尔 @ 巴拿马: PA|Panama|巴拿马 @ 哥斯达黎加: CR|Costa Rica|哥斯达黎加 @ 牙买加: JM|Jamaica|牙买加 @ 古巴: CU|Cuba|古巴 @ 海地: HT|Haiti|海地 @ 多米尼加: DO|Dominican Republic|多米尼加 @ 波多黎各: PR|Puerto Rico|波多黎各 @ 巴哈马: BS|Bahamas|巴哈马 @ 特立尼达和多巴哥: TT|Trinidad and Tobago|特立尼达和多巴哥 @ 巴巴多斯: BB|Barbados|巴巴多斯 @ 圣卢西亚: LC|Saint Lucia|圣卢西亚 @ 圣文森特和格林纳丁斯: VC|Saint Vincent and the Grenadines|圣文森特和格林纳丁斯 @ 格林纳达: GD|Grenada|格林纳达 @ 安提瓜和巴布达: AG|Antigua and Barbuda|安提瓜和巴布达 @ 多米尼克: DM|Dominica|多米尼克 @ 圣基茨和尼维斯: KN|Saint Kitts and Nevis|圣基茨和尼维斯 @ 马尔代夫: MV|Maldives|马尔代夫 @ 斐济: FJ|Fiji|斐济 @ 萨摩亚: WS|Samoa|萨摩亚 @ 汤加: TO|Tonga|汤加 @ 瓦努阿图: VU|Vanuatu|瓦努阿图 @ 所罗门群岛: SB|Solomon Islands|所罗门群岛 @ 帕劳: PW|Palau|帕劳 @ 密克罗尼西亚: FM|Micronesia|密克罗尼西亚 @ 马绍尔群岛: MH|Marshall Islands|马绍尔群岛 @ 基里巴斯: KI|Kiribati|基里巴斯 @ 图瓦卢: TV|Tuvalu|图瓦卢 @ 瑙鲁: NR|Nauru|瑙鲁"
	ss := strings.Split(s, "@")
	finalmap := make(map[string]*regexp.Regexp)
	for _, v := range ss {
		kv := strings.Split(v, ":")
		k := strings.TrimSpace(kv[0])
		r := strings.TrimSpace(kv[1])
		//fmt.Println(k, r)
		finalmap[k] = regexp.MustCompile(r)
	}
	return finalmap
}

var LocationMap = GetLocationMarkReMap()

func MarkProxy(ma map[string]*regexp.Regexp, gradeproxy *GradeProxy) {
	for k, v := range ma {
		if v.MatchString(gradeproxy.Name) {
			gradeproxy.Mark = append(gradeproxy.Mark, k)
			return
		}
	}
}

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
	for _, v := range proxies {
		wait.Add(1)
		go func(v *clash.Proxy) {
			defer wait.Done()
			defer lock.Unlock()
			gradeproxy := NewGradeProxy(v, proxiesmarksdic[v.Name])
			gradeproxy.Provider = gradeprovider
			lock.Lock()
			gradeproxies[v.Name] = gradeproxy

		}(v)
	}
	wait.Wait()
	gradeprovider.GradeProxies = gradeproxies
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
	if gradeproxy.DelayNow != 0 {
		delaypoint = 10 * float64(plus.TimeOut) / float64(gradeproxy.DelayNow) * 0.6
		av, rv := relativeVariance(gradeproxy.DelayHistory)
		avpoint = 10 * float64(plus.TimeOut) / av * 0.4
		//fmt.Println(rv)
		if 1-3*rv <= 0 {
			delaypoint = 0
		}
		delaypoint = delaypoint * (1 - 2*rv)
	}
	mark := gradeprovider.Level * gradeproxy.Level * (delaypoint + avpoint)
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

		gradeproxy, ok := gradegroup.Source[name]
		if ok {
			priority := checkmark(gradeproxy, gradegroup)
			gradegroup.Points[name] = int(float64(gradeproxy.Point) * priority)

		}

	}
}
func (gradegroup *GradeGroup) ChangeIf() {
	if !gradegroup.Block.After(time.Now()) {
		//fmt.Println(gradgroup.Name, "可更改")
		nowuse := gradegroup.Now
		nowpoint := gradegroup.Points[nowuse]
		name, value := maxInMap(gradegroup.Points)
		if value > int(float64(nowpoint)*1.3) {
			err := clash.SwitchProxy(gradegroup.Name, name)
			if err != nil {
				fmt.Println("切换失败", err)
				return
			}
			fmt.Printf("%s old:%s-延迟%d-分数%d --> new:%s-延迟%d-分数%d\n", gradegroup.Name, nowuse, gradegroup.Source[nowuse].DelayNow, nowpoint, name, gradegroup.Source[name].DelayNow, value)
		}

	}

}
func (gradegroup *GradeGroup) Update() {
	gradegroup.Group, _ = plus.GetGroupMessage(gradegroup.Name)
	gradegroup.GiveScore()
	gradegroup.ChangeIf()
}

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
	delay, _ := plus.GeneralDelayTest(gradeproxy.Name)
	gradeproxy.DelayNow = delay.Delay
	gradeproxy.DelayHistory = append(gradeproxy.DelayHistory, gradeproxy.DelayNow)
	//OneInsertHistory(gradeproxy)
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

func maxInMap(m map[string]int) (string, int) {
	var maxKey string
	var maxValue int
	first := true
	for key, value := range m {
		if first || value > maxValue {
			maxKey = key
			maxValue = value
			first = false
		}
	}
	return maxKey, maxValue
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

type Updater interface {
	ToUpdater()
	Update()
}

func UpdateOne(one Updater) {
	one.Update()

}
func AllUpdate(all interface{}) {
	switch v := all.(type) {
	case map[string]*GradeProvider:
		for _, provider := range v {
			provider.Update()
		}
	case map[string]*GradeGroup:
		for _, proxy := range v {
			proxy.Update()
		}
	case config.Config:

	default:
		fmt.Println("Unsupported type")
	}
}
