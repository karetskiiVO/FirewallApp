package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jessevdk/go-flags"
)

func main() {
	var options struct {
		Args struct {
			Mode       string
			Additional []string
		} `positional-args:"yes" required:"1"`
	}

	parser := flags.NewParser(&options, flags.Default&(^flags.PrintErrors))
	_, err := parser.Parse()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	switch options.Args.Mode {
	case "socket":
		err := socketSpy(options.Args.Additional)
		if err != nil {
			log.Panic(err)
		}
	case "nfqueue":
		err := nfqtablesSpy(options.Args.Additional)
		if err != nil {
			log.Panic(err)
		}
	default:
		log.Panic("wrong mode")
	}
}
