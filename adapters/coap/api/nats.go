/**
 * Copyright (c) 2017 Mainflux
 *
 * Mainflux server is licensed under an Apache license, version 2.0.
 * All rights not explicitly granted in the Apache license, version 2.0 are reserved.
 * See the included LICENSE file for more details.
 */

package api

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/go-nats"
	"log"
)

type (
	NatsMsg struct {
		Channel   string `json:"channel"`
		Publisher string `json:"publisher"`
		Protocol  string `json:"protocol"`
		Payload   []byte `json:"payload"`
	}
)

var (
	NatsConn *nats.Conn
)

func MsgHandler(nm *nats.Msg) {
	fmt.Printf("Received a message: %s\n", string(nm.Data))

	// And write it into the database
	m := NatsMsg{}
	if len(nm.Data) > 0 {
		if err := json.Unmarshal(nm.Data, &m); err != nil {
			println("Can not decode NATS msg")
			return
		}
	}

	println("Calling obsTransmit()")
	fmt.Println(m.Publisher, m.Protocol, m.Channel, m.Payload)
	obsTransmit(m)
}

func NatsInit(host string, port string) error {
	/** Connect to NATS broker */
	var err error

	log.Print("Connecting to NATS... ")
	NatsConn, err = nats.Connect("nats://" + host + ":" + port)

	return err
}
