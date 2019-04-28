package main

import (
	"encoding/json"
	"log"

	"github.com/Akumzy/ipc"
)

var ipcIO *ipc.IPC

type Who struct {
	Name string `json:"name,omitempty"`
}

func main() {
	ipcIO = ipc.New()
	go func() {
		// Me trying to write spanish
		ipcIO.SendAndReceive("hola", "Hola amigo, coma este nombre?", func(payload interface{}) {
			log.Println(payload)
		})

		ipcIO.On("who", func(data interface{}) {
			var who Who
			text := data.(string)
			if err := json.Unmarshal([]byte(text), &who); err != nil {
				log.Println(err)
				return
			}
			log.Println(who.Name)
		})

		ipcIO.OnReceiveAndReply("yoo", func(reply string, d interface{}) {
			log.Println(d)
			ipcIO.Reply(reply, "Sup Node", nil)
		})

	}()
	// You either start the IPC in it's routine
	// or start your own code in a go routine
	ipcIO.Start()
}
