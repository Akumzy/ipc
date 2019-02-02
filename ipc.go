package ipc

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

// IPC channel
type IPC struct {
	SendChannel           chan payload
	ReceiveChannel        chan payloadReceive
	ReceiveListerners     map[string][]func(data interface{})
	ReceiveSendListerners map[string][]func(emitName string, data string)
}

var (
	rLock  sync.Mutex
	rRLock sync.Mutex
)

// Payload this is the payload structure
type payload struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
	Error interface{} `json:"error"`
	SR    bool        `json:"SR"` //send and receive
}

// PayloadReceive this is the payload structure
type payloadReceive struct {
	Event string `json:"event"`
	Data  string `json:"data"`
	Error string `json:"error"`
	SR    bool   `json:"SR"` //send and receive
}

func (ipc IPC) receiveHandler() {
	for {
		// Grab the next message from the broadcast channel
		payload := <-ipc.ReceiveChannel
		// Send it out to every client that is currently connected
		if payload.SR {
			for _, handler := range ipc.ReceiveSendListerners[payload.Event] {
				replyChannel := payload.Event + "___RC___"
				handler(replyChannel, payload.Data)
			}
		} else {
			for _, handler := range ipc.ReceiveListerners[payload.Event] {
				handler(payload.Data)
			}
		}
	}
}

// Send socket
func (ipc IPC) Send(event string, data interface{}) {
	ipc.SendChannel <- payload{Event: event, Data: data}
}

// Reply back to sender
func (ipc IPC) Reply(event string, data, err interface{}) {
	ipc.SendChannel <- payload{Event: event, Data: data, SR: true, Error: err}
}

// On socket listener
func (ipc IPC) On(event string, handler func(data interface{})) {
	rLock.Lock()
	defer rLock.Unlock()
	h := ipc.ReceiveListerners[event]
	h = append(h, handler)
	ipc.ReceiveListerners[event] = h
}

// OnReceiveAndReply receive and send back
func (ipc IPC) OnReceiveAndReply(event string, handler func(replyChannel string, data string)) {
	rRLock.Lock()
	defer rRLock.Unlock()
	h := ipc.ReceiveSendListerners[event]
	h = append(h, handler)
	ipc.ReceiveSendListerners[event] = h

}

// Start create ipc
func (ipc IPC) Start() {

	go func() {
		busy := false
		for {
			if !busy {
				busy = true
				msg := <-ipc.SendChannel
				data, err := Marshal(msg)
				if err != nil {
					log.Println(err)
				} else {
					fmt.Print(data + "\\n")
				}
				busy = false
			}

		}
	}()
	for {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		if text != "" {
			var payload payloadReceive
			text = strings.Replace(text, "\n", "", -1)
			if text != "" {
				if err := json.Unmarshal([]byte(text), &payload); err != nil {
					log.Println(err)
					continue
				}
				ipc.ReceiveChannel <- payload
			}

		}
	}
}

// Marshal to json
func Marshal(v interface{}) (string, error) {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(&v); err != nil {
		return "", err
	}
	return strings.TrimSpace(buf.String()), nil
}

// New return now ipc
func New() *IPC {
	ipc := &IPC{}
	ipc.SendChannel = make(chan payload)
	ipc.ReceiveChannel = make(chan payloadReceive)
	ipc.ReceiveListerners = make(map[string][]func(data interface{}))
	ipc.ReceiveSendListerners = make(map[string][]func(emitName string, data string))
	return ipc
}
