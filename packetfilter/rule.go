package packetfilter

import (
	"fmt"

	"github.com/google/gopacket"
)

type RuleResult int

type Rule interface {
	Accept(packet gopacket.Packet) bool
}

func NewRuleFromMap(content map[string]any) (Rule, error) {
	ruletype, ok := content["type"]
	if !ok {
		return nil, fmt.Errorf("the type field is expected")
	}

	switch ruletype {
	case "and":
		return NewAndRule(content)
	case "or":
		return NewOrRule(content)
	case "ip":
		return NewIPRuleFromMap(content)
	case "tcp":
		return NewTCPRuleFromMap(content)
	case "udp":
		return NewUDPRuleFromMap(content)
	}
	return nil, fmt.Errorf("invalid rule type")
}
