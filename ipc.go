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
	sendChannel           chan payload
	receiveListerners     map[string][]Handler
	receiveSendListerners map[string][]HandlerWithReply
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
	RS    bool        `json:"RC"` // receive and send
}

/*Handler When the underline type of data is being
  access through `type assertion` if the data has a
  literal value the underlining type will be return
  else a `JSON` representive of the data will be return
*/
type Handler func(data interface{})

/*HandlerWithReply  When the underline type of data is being
  access through `type assertion` if the data has a literal
  value the underlining type will be return else a `JSON` representive of
  the data will be return.
  `replyChannel` is the event name you'll pass to `ipc.Reply` method to respond
   to the sender
*/
type HandlerWithReply func(replyChannel string, data interface{})

// PayloadReceive this is the payload structure
type payloadReceive struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
	SR    bool        `json:"SR"` //send and receive
}

// Send data to parent process
func (ipc IPC) Send(event string, data interface{}) {
	ipc.sendChannel <- payload{Event: event, Data: data}
}

// Reply back to sender
func (ipc IPC) Reply(event string, data, err interface{}) {
	ipc.sendChannel <- payload{Event: event, Data: data, SR: true, Error: err}
}

// On listens for events from parent process
func (ipc IPC) On(event string, handler Handler) {
	rLock.Lock()
	defer rLock.Unlock()
	h := ipc.receiveListerners[event]
	h = append(h, handler)
	ipc.receiveListerners[event] = h
}

/*OnReceiveAndReply listen for an events and as well reply back to
 * the same sender with the help of `ipc.Reply` method
 */
func (ipc IPC) OnReceiveAndReply(event string, handler HandlerWithReply) {
	rRLock.Lock()
	defer rRLock.Unlock()
	h := ipc.receiveSendListerners[event]
	h = append(h, handler)
	ipc.receiveSendListerners[event] = h

}

/*SendAndReceive send and listen for reply event
 */
func (ipc IPC) SendAndReceive(event string, data interface{}, handler Handler) {
	ipc.sendChannel <- payload{Event: event, Data: data, RS: true}
	channel := event + "___RS___"
	ipc.On(channel, handler)
}

/*Start `ipc`
* the `ipc.Start` method will blocks executions
* so is either you put in a seperate `Go routine` or put you own code in
* a different `Go routine`
 */
func (ipc IPC) Start() {

	go func() {
		for {
			msg := <-ipc.sendChannel
			data, err := Marshal(msg)
			if err != nil {
				log.Println(err)
			} else {
				fmt.Print(data + "\\n")
			}

		}
	}()
	for {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		if text != "" {
			var payload payloadReceive
			text = strings.Replace(text, "\n", "", -1)
			// check if the text is not empty string

			if text != "" {
				if err := json.Unmarshal([]byte(text), &payload); err != nil {
					log.Println(err)
					continue
				}
				if payload.SR {
					for _, handler := range ipc.receiveSendListerners[payload.Event] {
						replyChannel := payload.Event + "___RC___"
						handler(replyChannel, payload.Data)
					}
				} else {
					for _, handler := range ipc.receiveListerners[payload.Event] {

						handler(payload.Data)
					}
				}
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
	ipc.sendChannel = make(chan payload)
	ipc.receiveListerners = make(map[string][]Handler)
	ipc.receiveSendListerners = make(map[string][]HandlerWithReply)
	return ipc
}
