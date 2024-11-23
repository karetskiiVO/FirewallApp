package packetfilter

import (
	"fmt"
	"reflect"

	"github.com/google/gopacket"
)

type AndRule struct {
	rules []Rule
}

func NewAndRule(content map[string]any) (*AndRule, error) {
	rawrulesMap, ok := content["rules"]
	if !ok {
		return nil, fmt.Errorf("can't solve and rule")
	}

	var err error

	fmt.Println(reflect.TypeOf(rawrulesMap), rawrulesMap)

	rawrules, ok := rawrulesMap.([]map[string]any)

	fmt.Println(rawrules)
	if !ok || len(rawrules) == 0 {
		return nil, fmt.Errorf("wrong array of rules format")
	}

	rules := make([]Rule, len(rawrules))
	for i, rawrule := range rawrules {
		rules[i], err = NewRuleFromMap(rawrule)
		if err != nil {
			return nil, err
		}
	}

	return &AndRule{rules}, nil
}

func (rule AndRule) Accept(packet gopacket.Packet) bool {
	for _, rule := range rule.rules {
		if !rule.Accept(packet) {
			return false
		}
	}

	return true
}

type OrRule struct {
	rules []Rule
}

func NewOrRule(content map[string]any) (*OrRule, error) {
	rawrulesMap, ok := content["rules"]
	if !ok {
		return nil, fmt.Errorf("can't solve and rule")
	}

	var err error

	rawrules, ok := rawrulesMap.([]map[string]any)
	if !ok || len(rawrules) == 0 {
		return nil, fmt.Errorf("wrong array of rules format")
	}

	rules := make([]Rule, len(rawrules))
	for i, rawrule := range rawrules {
		rules[i], err = NewRuleFromMap(rawrule)
		if err != nil {
			return nil, err
		}
	}

	return &OrRule{rules}, nil
}

func (rule OrRule) Accept(packet gopacket.Packet) bool {
	for _, rule := range rule.rules {
		if rule.Accept(packet) {
			return true
		}
	}

	return false
}
