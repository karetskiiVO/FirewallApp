package main

import (
	"fmt"
	"log"
	"net"
	"syscall"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/jessevdk/go-flags"
	"github.com/karetskiiVO/FirewallApp/packetfilter"
)

func htons(i uint16) uint16 {
	return (i<<8)&0xff00 | i>>8
}

func spy(infd, outfd int, filter *packetfilter.Filter, mtu int) {
	buffer := make([]byte, mtu)

	for {
		length, err := syscall.Read(infd, buffer)
		if err != nil {
			log.Print(err)
			continue
		}

		fmt.Println("========================")
		fmt.Println(gopacket.NewPacket(buffer[:length], layers.LayerTypeEthernet, gopacket.Default))
		if filter.Accept(buffer[:length]) {
			syscall.Write(outfd, buffer[:length])
		} else {
			fmt.Println(">>>>>>DROPPED<<<<<<")
		}
		fmt.Println("========================\n")
	}
}

func socketSpy(args []string) error {
	var options struct {
		Args struct {
			IFName1    string
			IFName2    string
			ConfigFile string
		} `positional-args:"yes" required:"3"`
	}

	parser := flags.NewParser(&options, flags.Default&(^flags.PrintErrors))
	_, err := parser.ParseArgs(args)
	if err != nil {
		return err
	}

	wait := make(chan struct{})

	if1, err := net.InterfaceByName(options.Args.IFName1)
	if err != nil {
		return err
	}
	if2, err := net.InterfaceByName(options.Args.IFName2)
	if err != nil {
		return err
	}

	infd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, int(htons(syscall.ETH_P_ALL)))
	if err != nil {
		return err
	}
	outfd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, int(htons(syscall.ETH_P_ALL)))
	if err != nil {
		return err
	}

	err = syscall.Bind(infd, &syscall.SockaddrLinklayer{
		Protocol: htons(syscall.ETH_P_ALL),
		Ifindex:  if1.Index,
	})
	if err != nil {
		return err
	}
	defer syscall.Close(infd)
	err = syscall.Bind(outfd, &syscall.SockaddrLinklayer{
		Protocol: htons(syscall.ETH_P_ALL),
		Ifindex:  if2.Index,
	})
	if err != nil {
		return err
	}
	defer syscall.Close(outfd)

	filter, err := packetfilter.NewFilter(options.Args.ConfigFile)
	if err != nil {
		return err
	}

	mtu := if1.MTU
	if mtu > if2.MTU {
		mtu = if2.MTU
	}
	go spy(infd, outfd, filter, mtu)
	go spy(outfd, infd, filter, mtu)

	<-wait

	return nil
}
