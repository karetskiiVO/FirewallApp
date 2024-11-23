package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"syscall"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/jessevdk/go-flags"
)

const (
	Length = 100000
)

func htons(i uint16) uint16 {
	return (i<<8)&0xff00 | i>>8
}

func main() {
	var options struct {
		Args struct {
			IFName1 string
			IFName2 string
		} `positional-args:"yes" required:"2"`
	}

	parser := flags.NewParser(&options, flags.Default&(^flags.PrintErrors))
	_, err := parser.Parse()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	wait := make(chan struct{})

	if1, err := net.InterfaceByName(options.Args.IFName1)
	if err != nil {
		log.Panic(err)
	}
	if2, err := net.InterfaceByName(options.Args.IFName2)
	if err != nil {
		log.Panic(err)
	}

	infd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, int(htons(syscall.ETH_P_ALL)))
	if err != nil {
		log.Panic(err)
	}
	outfd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, int(htons(syscall.ETH_P_ALL)))
	if err != nil {
		log.Panic(err)
	}

	err = syscall.Bind(infd, &syscall.SockaddrLinklayer{
		Protocol: htons(syscall.ETH_P_ALL),
		Ifindex:  if1.Index,
	})
	if err != nil {
		log.Panic(err)
	}
	err = syscall.Bind(outfd, &syscall.SockaddrLinklayer{
		Protocol: htons(syscall.ETH_P_ALL),
		Ifindex:  if2.Index,
	})
	if err != nil {
		log.Panic(err)
	}

	go spy(infd, outfd)
	go spy(outfd, infd)

	<-wait
}

func spy(infd, outfd int) {
	buffer := make([]byte, Length)

	for {
		length, err := syscall.Read(infd, buffer)
		if err != nil {
			log.Print(err)
			continue
		}

		log.Print(gopacket.NewPacket(buffer[:length], layers.LayerTypeEthernet, gopacket.Default))

		syscall.Write(outfd, buffer[:length])
	}
}
