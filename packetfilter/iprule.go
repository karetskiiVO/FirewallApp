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

func (rule IPRule) Accept(packet gopacket.Packet) bool {
	linkLayer := packet.LinkLayer()
	if linkLayer != nil {
		return false
	}

	src := linkLayer.LinkFlow().Src().String()
	dst := linkLayer.LinkFlow().Dst().String()

	log.Printf("%v > %v", src, dst)

	accepsrc := true
	accepdst := true
	if rule.srcaddr != nil {
		accepdst = rule.srcaddr.MatchString(dst)
	}
	if rule.dstaddr != nil {
		accepsrc = rule.srcaddr.MatchString(src)
	}

	return accepdst && accepsrc
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
