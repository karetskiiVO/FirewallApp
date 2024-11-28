package packetfilter

import (
	"fmt"
	"regexp"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type DNSRule struct {
	depricatedQuestions []*regexp.Regexp

	requiresAA *bool
	requiresTC *bool
	requiresRD *bool
	requiresRA *bool
}

func NewDNSRuleFromMap(content map[string]any) (*DNSRule, error) {
	var depricated []*regexp.Regexp
	var requiresAA *bool
	var requiresTC *bool
	var requiresRD *bool
	var requiresRA *bool

	if depricatedRaw, ok := content["depricated"]; ok {
		regsRaw, ok := depricatedRaw.([]any)

		if !ok {
			return nil, fmt.Errorf("wrong depricated format")
		}

		depricated = make([]*regexp.Regexp, len(regsRaw))
		for i, regAny := range regsRaw {
			reg, ok := regAny.(string)
			if !ok {
				return nil, fmt.Errorf("can't cast to string raw regular expresstion")
			}

			var err error
			depricated[i], err = regexp.Compile(reg)
			if err != nil {
				return nil, err
			}
		}
	}

	if requiresAAraw, ok := content["AA"]; ok {
		requiresAA = new(bool)
		*requiresAA, ok = requiresAAraw.(bool)

		if !ok {
			return nil, fmt.Errorf("wrong AA format")
		}
	}

	if requiresTCraw, ok := content["TC"]; ok {
		requiresTC = new(bool)
		*requiresTC, ok = requiresTCraw.(bool)

		if !ok {
			return nil, fmt.Errorf("wrong TC format")
		}
	}

	if requiresRDraw, ok := content["RD"]; ok {
		requiresRD = new(bool)
		*requiresRD, ok = requiresRDraw.(bool)

		if !ok {
			return nil, fmt.Errorf("wrong RD format")
		}
	}

	if requiresRAraw, ok := content["RA"]; ok {
		requiresRA = new(bool)
		*requiresRA, ok = requiresRAraw.(bool)

		if !ok {
			return nil, fmt.Errorf("wrong RA format")
		}
	}

	return &DNSRule{
		depricatedQuestions: depricated,
		requiresAA: requiresAA,
		requiresTC: requiresTC,
		requiresRD: requiresRD,
		requiresRA: requiresRA,
	}, nil
}

func (rule DNSRule) Accept(packet gopacket.Packet) bool {
	application := packet.ApplicationLayer()

	if application == nil {
		return false
	}

	if application.LayerType().String() != "DNS" {
		return false
	}
	dns := application.(*layers.DNS)

	acceptQuestion := true
	acceptAA := true
	acceptTC := true
	acceptRD := true
	acceptRA := true

	if rule.depricatedQuestions != nil && len(rule.depricatedQuestions) != 0 {
		for _, question := range dns.Questions {
			for _, dropAddr := range rule.depricatedQuestions {
				if dropAddr.MatchString(string(question.Name)) {
					acceptQuestion = false
				}
			}
		}
	}

	if rule.requiresAA != nil {
		acceptAA = (*rule.requiresAA == dns.AA)
	}
	if rule.requiresTC != nil {
		acceptTC = (*rule.requiresTC == dns.TC)
	}
	if rule.requiresRD != nil {
		acceptRD = (*rule.requiresRD == dns.RD)
	}
	if rule.requiresRA != nil {
		acceptRA = (*rule.requiresRA == dns.RA)
	}

	return acceptQuestion && acceptAA && acceptTC && acceptRD && acceptRA
}
