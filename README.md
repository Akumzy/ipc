# IPC

`ipc` provides an event listeners and event emitter methods or functions for ipc(Inter-process communication) using the process `stdin` as it's meduim.

## What's the motive behind this package creation?

I was having one or more issues with any of `Nodejs` based `fs-watcher` module I could find,
So I extended my search outside Node and found a package written in `Go` that addressed all those issues for me but the biggest challenge I was facing then was communicating with the two process (electron (Nodjs) and Go process) I tried using websocket but it didn't work out well until I found a post on how to read/scan the `stdin` using `bufio.NewReader(os.Stdin)` which drived the idea from to create this package.

## Note: To use this package you must not log/print stuffs with `fmt` package because `fmt` writes to process `stdout`

This package was made with `Nodejs` as the parent process in mind.

## Usage

To use this package you will surely need it's `Nodejs` version [ipc-node-go](https://github.com/Akumzy/ipc-node)

```shell
go get github.com/Akumzy/ipc
```

```go
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

		ipcIO.On("who", func(payload interface{}) {
			var who Who
			text := payload.(string)
			if err := json.Unmarshal([]byte(text), &who); err != nil {
				log.Println(err)
				return
			}
			log.Println(who.Name)
		})

		ipcIO.OnReceiveAndReply("yoo", func(reply string, payload interface{}) {
			log.Println(payload)
			ipcIO.Reply(reply, "Sup Node", nil)
		})

	}()
	// You either start the IPC in it's routine
	// or start your own code in a go routine
	ipcIO.Start()
}
```