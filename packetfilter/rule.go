package packetfilter

import "github.com/google/gopacket"

// Rult - packet rules
type Rule interface {
	Accept (packet gopacket.Packet) bool
}