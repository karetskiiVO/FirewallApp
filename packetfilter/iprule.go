package packetfilter

import (
	"fmt"
	"log"
	"regexp"

	"github.com/google/gopacket"
)

type IPRule struct {
	srcaddr *regexp.Regexp
	dstaddr *regexp.Regexp
}

func (rule IPRule) Accept(packet gopacket.Packet) RuleResult {
	linkLayer := packet.NetworkLayer()
	if linkLayer == nil {
		return NotMatchedType
	}

	src := linkLayer.NetworkFlow().Src().String()
	dst := linkLayer.NetworkFlow().Dst().String()

	log.Printf("%v > %v", src, dst)

	accepsrc := true
	accepdst := true
	if rule.dstaddr != nil {
		accepdst = rule.dstaddr.MatchString(dst)
	}
	if rule.srcaddr != nil {
		accepsrc = rule.srcaddr.MatchString(src)
	}

	if accepsrc && accepdst {
		return AcceptRule
	} 
	return NotAcceptRule
}

func NewIPRuleFromMap(content map[string]any) (Rule, error) {
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
