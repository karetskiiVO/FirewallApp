package packetfilter

import (
	"fmt"

	"github.com/google/gopacket"
)

type RuleResult int

type Rule interface {
	Accept(packet gopacket.Packet) bool
}

func NewRuleFromMap(defaultAccept bool, content map[string]any) (Rule, error) {
	ruletype, ok := content["type"]
	if !ok {
		return nil, fmt.Errorf("the type field is expected")
	}

	switch ruletype {
	case "ip":
		return NewIPRuleFromMap(content)
	case "tcp":
		return NewTCPRuleFromMap(content)
	case "udp":
		return NewUDPRuleFromMap(content)
	}
	return nil, fmt.Errorf("invalid rule type")
}
