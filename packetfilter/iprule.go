package packetfilter

import (
	"fmt"
	"regexp"

	"github.com/google/gopacket"
)

type IPRule struct {
	srcaddr *regexp.Regexp
	dstaddr *regexp.Regexp
}

func NewIPRuleFromMap(content map[string]any) (*IPRule, error) {
	newrule := &IPRule{}

	var err error

	if rawreg, ok := content["source"]; ok {
		reg, ok := rawreg.(string)
		if !ok {
			return nil, fmt.Errorf("can't cast to string: %v", rawreg)
		}

		newrule.srcaddr, err = regexp.Compile(reg)
		if err != nil {
			return nil, err
		}
	}
	if rawreg, ok := content["destination"]; ok {
		reg, ok := rawreg.(string)
		if !ok {
			return nil, fmt.Errorf("can't cast to string: %v", rawreg)
		}

		newrule.dstaddr, err = regexp.Compile(reg)
		if err != nil {
			return nil, err
		}
	}

	return newrule, nil
}

func (rule IPRule) Accept(packet gopacket.Packet) bool {
	networkLayer := packet.NetworkLayer()
	if networkLayer == nil {
		return false
	}

	src := networkLayer.NetworkFlow().Src().String()
	dst := networkLayer.NetworkFlow().Dst().String()

	accepsrc := true
	accepdst := true
	if rule.dstaddr != nil {
		accepdst = rule.dstaddr.MatchString(dst)
	}
	if rule.srcaddr != nil {
		accepsrc = rule.srcaddr.MatchString(src)
	}

	return accepsrc && accepdst
}
