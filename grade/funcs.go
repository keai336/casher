package grade

import (
	"dash/config"
	"fmt"
	"regexp"
	"strings"
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