/**
 * Copyright (c) Mainflux
 *
 * Mainflux server is licensed under an Apache license, version 2.0.
 * All rights not explicitly granted in the Apache license, version 2.0 are reserved.
 * See the included LICENSE file for more details.
 */

package api

import (
	"bytes"
	"encoding/json"
	"log"
	"net"

	"github.com/dereulenspiegel/coap-mux"
	"github.com/dustin/go-coap"
)

type (
	Observer struct {
		Conn    *net.UDPConn
		Addr    *net.UDPAddr
		Message *coap.Message
	}
)

// Map of observers
var ObsMap map[string][]Observer

// sendMessage function
func sendMessage(l *net.UDPConn, a *net.UDPAddr, m *coap.Message) *coap.Message {
	log.Printf("Got message in sendMessage: path=%q: %#v from %v", m.Path(), m, a)
	var res *coap.Message = nil
	if m.IsConfirmable() {
		res = &coap.Message{
			Type:      coap.Acknowledgement,
			Code:      coap.Content,
			MessageID: m.MessageID,
			Token:     m.Token,
			Payload:   []byte(""),
		}
		res.SetOption(coap.ContentFormat, coap.AppJSON)
	}

	if len(m.Payload) == 0 {
		if m.IsConfirmable() {
			res.Payload = []byte("{\"res\": \"Error: msg len can not be 0\"}")
		}
		return res
	}

	// Channel ID
	cid := mux.Var(m, "channel_id")

	// Publish message on MQTT via NATS
	n := NatsMsg{}
	n.Channel = cid
	n.Publisher = ""
	n.Protocol = "coap"
	n.Payload = m.Payload

	b, err := json.Marshal(n)
	if err != nil {
		log.Print(err)
	}

	NatsConn.Publish("msg.coap", b)

	if m.IsConfirmable() {
		res.Payload = []byte("{\"res\": \"sent\"}")
	}
	return res
}

func obsTransmit(n NatsMsg) {

	for _, e := range ObsMap[n.Channel] {

		msg := *(e.Message)
		msg.Payload = n.Payload

		log.Printf("ObsMap[cid] = %v", e)
		log.Printf("msg = %v", msg)

		msg.SetOption(coap.ContentFormat, coap.AppJSON)
		msg.SetOption(coap.LocationPath, msg.Path())

		log.Printf("Transmitting %v", msg)
		err := coap.Transmit(e.Conn, e.Addr, msg)
		if err != nil {
			log.Printf("Error on transmitter, stopping: %v", err)
			return
		}
	}

}

// observeMessage function
func observeMessage(l *net.UDPConn, a *net.UDPAddr, m *coap.Message) *coap.Message {
	log.Printf("Got message in observeMessage: path=%q: %#v from %v", m.Path(), m, a)
	var res *coap.Message = nil

	if m.IsConfirmable() {
		res = &coap.Message{
			Type:      coap.Acknowledgement,
			Code:      coap.Content,
			MessageID: m.MessageID,
			Token:     m.Token,
			Payload:   []byte(""),
		}
		res.SetOption(coap.ContentFormat, coap.AppJSON)
	}

	// Channel ID
	cid := mux.Var(m, "channel_id")

	// Observer
	o := Observer{
		Conn:    l,
		Addr:    a,
		Message: m,
	}

	if m.Option(coap.Observe) != nil {
		if value, ok := m.Option(coap.Observe).(uint32); ok {
			if value == 0 {
				// Register
				found := false
				for _, e := range ObsMap[cid] {
					if e.Addr == o.Addr && bytes.Compare(e.Message.Token, o.Message.Token) == 0 {
						found = true
						break
					}
				}
				if !found {
					log.Println("Register " + cid)
					log.Printf("o.Message = %v", o.Message)
					ObsMap[cid] = append(ObsMap[cid], o)
				}
			} else {
				// Deregister
				for i, e := range ObsMap[cid] {
					if bytes.Compare(e.Message.Token, o.Message.Token) == 0 {
						// Observer found, remove it from array
						log.Println("Deregister " + cid)
						arr := ObsMap[cid]
						arr = append(arr[:i], arr[i+1:]...)
					}
				}
			}
		} else {
			log.Printf("%v", value)
		}
	} else {
		// Interop - old deregister was when there was no Observe option provided
		for i, e := range ObsMap[cid] {
			if bytes.Compare(e.Message.Token, o.Message.Token) == 0 {
				// Observer found, remove it from array
				log.Println("Interop - Deregister " + cid)
				ObsMap[cid] = append((ObsMap[cid])[:i], (ObsMap[cid])[i+1:]...)
			}
		}
	}

	if m.IsConfirmable() {
		res.Payload = []byte("{\"res\": \"observing\"}")
	}
	return res
}
