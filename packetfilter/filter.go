package packetfilter

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// Filter - implements filtration
type Filter struct {
	rules []Rule

	defaultAccept bool
}

func istechnical (packet gopacket.Packet) bool {
	networkLayer := packet.NetworkLayer()

	if networkLayer == nil {
		fmt.Println("no type")
		return true
	}

	fmt.Println("::::::::::::::::::::::::", networkLayer.LayerType().String())
	
	return false
}

func (f Filter) Accept(rawpacket []byte) bool {
	packet := gopacket.NewPacket(rawpacket, layers.LayerTypeEthernet, gopacket.Lazy)

	if istechnical(packet) {
		return true
	}

	for _, rule := range f.rules {
		if rule.Accept(packet) {
			return !f.defaultAccept
		}
	}

	return f.defaultAccept
}

func NewFilter(cfgfilename string) (*Filter, error) {
	content, err := os.ReadFile(cfgfilename)
	if err != nil {
		return nil, err
	}
	var rawrules ruleContainer
	json.Unmarshal(content, &rawrules)

	defaultAccept := true
	switch rawrules.DefaultAccept {
	case "accept":
		defaultAccept = true
	case "drop":
		defaultAccept = false
	default:
		return nil, fmt.Errorf("Mismatch default format")
	}

	rules := make([]Rule, len(rawrules.Rules))
	for i, rawrule := range rawrules.Rules {
		rules[i], err = NewRuleFromMap(defaultAccept, rawrule)
		if err != nil {
			return nil, err
		}
	}

	return &Filter{
		defaultAccept: defaultAccept,
		rules: rules,
	}, nil
}
