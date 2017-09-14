/**
 * Copyright (c) 2016 Mainflux
 *
 * Mainflux server is licensed under an Apache license, version 2.0.
 * All rights not explicitly granted in the Apache license, version 2.0 are reserved.
 * See the included LICENSE file for more details.
 */

package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/dustin/go-coap"
	"github.com/fatih/color"
	"github.com/mainflux/coap-adapter/api"

	"github.com/cenkalti/backoff"
	"github.com/nats-io/go-nats"
)

const (
	help string = `
Usage: mainflux-influxdb [options]
Options:
	-a, --addr	CoAP server host
	-p, --port	CoAP server port
	-n, --nhost	NATS host
	-m, --nport	NATS port
	-h, --help	Show help
`
)

type (
	Opts struct {
		COAPHost string
		COAPPort string

		NatsHost string
		NatsPort string

		Help bool
	}

	NatsMsg struct {
		Channel   string `json:"channel"`
		Publisher string `json:"publisher"`
		Protocol  string `json:"protocol"`
		Payload   []byte `json:"payload"`
	}
)

var opts Opts

func tryNatsConnect() error {
	var err error

	log.Print("Connecting to NATS... ")
	api.NatsConn, err = nats.Connect("nats://" + opts.NatsHost + ":" + opts.NatsPort)
	return err
}

func main() {
	flag.StringVar(&opts.COAPHost, "a", "localhost", "CoAP host.")
	flag.StringVar(&opts.COAPPort, "p", "5683", "CoAP port.")
	flag.StringVar(&opts.NatsHost, "n", "localhost", "NATS broker address.")
	flag.StringVar(&opts.NatsPort, "m", "4222", "NATS broker port.")
	flag.BoolVar(&opts.Help, "h", false, "Show help.")
	flag.BoolVar(&opts.Help, "help", false, "Show help.")

	flag.Parse()

	if opts.Help {
		fmt.Printf("%s\n", help)
		os.Exit(0)
	}

	// Connect to NATS broker
	if err := backoff.Retry(tryNatsConnect, backoff.NewExponentialBackOff()); err != nil {
		log.Fatalf("NATS: Can't connect: %v\n", err)
	} else {
		log.Println("OK")
	}

	// Initialize map of Observers
	api.ObsMap = make(map[string][]api.Observer)

	// Subscribe to NATS
	api.NatsConn.Subscribe("msg.http", api.MsgHandler)
	api.NatsConn.Subscribe("msg.mqtt", api.MsgHandler)

	// Print banner
	color.Cyan(banner)
	color.Cyan(fmt.Sprintf("Magic happens on port %s", opts.COAPPort))

	// Serve CoAP
	coapAddr := fmt.Sprintf("%s:%s", opts.COAPHost, opts.COAPPort)
	coap.ListenAndServe("udp", coapAddr, api.COAPServer())
}

var banner = `
╔╦╗┌─┐┬┌┐┌┌─┐┬  ┬ ┬─┐ ┬   ╔═╗┌─┐╔═╗╔═╗
║║║├─┤││││├┤ │  │ │┌┴┬┘───║  │ │╠═╣╠═╝
╩ ╩┴ ┴┴┘└┘└  ┴─┘└─┘┴ └─   ╚═╝└─┘╩ ╩╩  

     == Industrial IoT System ==
       
     Made with <3 by Mainflux Team
[w] http://mainflux.io
[t] @mainflux

`
