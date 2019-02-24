package main

import (
	"encoding/json"
	"log"

	"github.com/Akumzy/ipc"
)

var ipcIO *ipc.IPC

// Count numbers
type Count struct {
	Num int `json:"num"`
}

func main() {
	ipcIO := ipc.New()
	go func() {
		name := map[string]string{}
		name["name"] = "Golang"
		ipcIO.Send("hello", name)

		ipcIO.On("count", func(data interface{}) {
			count := data.(float64)
			log.Println(count)
		})

		ipcIO.On("count-object", func(data interface{}) {
			var count Count
			text := data.(string)
			if err := json.Unmarshal([]byte(text), &count); err != nil {
				log.Println(err)
				return
			}
			log.Println(count)
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
