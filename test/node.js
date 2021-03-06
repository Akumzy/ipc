const IPC = require('ipc-node-go')
const ipc = new IPC('./test')
ipc.init()

ipc.on('log', console.log)
// Log all error from stderr
ipc.on('error', console.error)


//
ipc.send('who', { name: 'Akuma Nodejs' })
// send and receive and an acknowledgement

// listen for event and reply to the channel
ipc.onReceiveAnSend('hola', (channel, data) => {
  console.log(data)
  ipc.send(channel, 'cool thanks')
})
setInterval(() => {
  ipc.sendAndReceive('yoo', 'Hello, how are you doing?', (err, d) => {
    err ? console.error(err) : console.log(d)
  })
}, 20000)
