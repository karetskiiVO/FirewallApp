package main

import (
	"fmt"

	"github.com/jessevdk/go-flags"
	"github.com/juturnas/go-netfilter-queue"
	"github.com/karetskiiVO/FirewallApp/packetfilter"
)

func nfqtablesSpy(args []string) error {
	var options struct {
		Args struct {
			QueueNum   uint16
			ConfigFile string
		} `positional-args:"yes" required:"3"`

		QueueLen uint32 `short:"l" long:"length" default:"128"`
	}

	parser := flags.NewParser(&options, flags.Default&(^flags.PrintErrors))
	_, err := parser.ParseArgs(args)
	if err != nil {
		return err
	}

	// err = syscall.Exec(
	// 	"iptables",
	// 	[]string{
	// 		"-t", "mangle",
	// 		"-A", "FORWARD",
	// 		"-j", "NFQUEUE",
	// 		"--queue-num", fmt.Sprint(options.Args.QueueNum)},
	// 	[]string{},
	// )
	// if err != nil {
	// 	return err
	// }

	nfq, err := netfilter.NewNFQueue(options.Args.QueueNum, options.QueueLen, netfilter.NF_DEFAULT_PACKET_SIZE)
	if err != nil {
		return err
	}
	defer nfq.Close()

	packets := nfq.GetPackets()
	filter, err := packetfilter.NewFilter(options.Args.ConfigFile)
	if err != nil {
		return err
	}

	for true {
		select {
		case packet := <-packets:
			fmt.Println(packet.Packet)

			if filter.Accept(packet.Packet.Data()) {
				packet.SetVerdict(netfilter.Verdict(netfilter.NF_ACCEPT))
			} else {
				packet.SetVerdict(netfilter.Verdict(netfilter.NF_DROP))
			}
		}
	}

	return nil
}
