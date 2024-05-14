package config

import (
	"fmt"
	"testing"
)

func TestRef(t *testing.T) {
	providerleveldic := map[string]float64{"mesl": 1, "planc": 2.3}
	grouplabeldic := map[string]map[string]float64{"Final": map[string]float64{"final": 1, "hk": 2.3}}
	proxymark := map[string][]string{"final": []string{"hk", "openai"}}
	grouplevl := map[string]float64{"Final": 1}
	config := NewConfig(providerleveldic, grouplabeldic, proxymark, grouplevl)
	DumpConfig(config, "test.yaml")
	config = LoadConfig("test.yaml")
	fmt.Println(config)
}
