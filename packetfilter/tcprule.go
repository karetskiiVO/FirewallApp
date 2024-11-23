package packetfilter

import (
	"fmt"
	"log"
	"regexp"

	"github.com/google/gopacket"
)

type TCPRule struct {
	srcaddr *regexp.Regexp
	dstaddr *regexp.Regexp
}

func NewTCPRuleFromMap(content map[string]any) (*TCPRule, error) {
	newrule := &TCPRule{}

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

func (rule TCPRule) Accept(packet gopacket.Packet) bool {
	transportLayer := packet.TransportLayer()
	if transportLayer == nil {
		return false
	}

	if transportLayer.LayerType().String() != "TCP" {
		return false
	}

	src := transportLayer.TransportFlow().Src().String()
	dst := transportLayer.TransportFlow().Dst().String()

	log.Printf("tcp:%v > %v\n", src, dst)

	accepsrc := true
	accepdst := true
	if rule.dstaddr != nil {
		accepdst = rule.dstaddr.MatchString(dst)
	}
	if rule.srcaddr != nil {
		accepsrc = rule.srcaddr.MatchString(src)
	}

	if accepsrc && accepdst {
		return true
	}
	return false
}
