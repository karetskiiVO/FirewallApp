package packetfilter

import (
	"encoding/json"
	"os"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// Filter - implements filtration
type Filter struct {
	rules []Rule

	defaultAccept bool
}

func (f Filter) Accept(rawpacket []byte) bool {
	packet := gopacket.NewPacket(rawpacket, layers.LayerTypeEthernet, gopacket.Lazy)

	for _, rule := range f.rules {
		switch rule.Accept(packet) {
		case NotMatchedType:
			continue
		case AcceptRule:
			return !f.defaultAccept
		case NotAcceptRule:
			continue
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

	rules := make([]Rule, len(rawrules.Rules))
	for i, rawrule := range rawrules.Rules {
		rules[i], err = NewRuleFromMap(rawrules.DefaultAccept, rawrule)
		if err != nil {
			return nil, err
		}
	}

	return &Filter{
		defaultAccept: rawrules.DefaultAccept,
		rules: rules,
	}, nil
}
