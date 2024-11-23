package packetfilter

import (
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
		if rule.Accept(packet) {
			return !f.defaultAccept
		}
	}

	return f.defaultAccept
}

func NewFilter() *Filter {
	return &Filter{}
}
