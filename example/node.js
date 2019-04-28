const IPC = require('ipc-node-go')
const ipc = new IPC('./example')
ipc.init()

ipc.on('log', console.log)
// Log all error from stderr
ipc.on('error', console.error)
ipc.on('ping', ()=> {
  console.log(new Date())
})

//
ipc.send('who', {name:'Akuma Nodejs'})
// send and receive and an acknowledgement
ipc.sendAndReceive('yoo', 'Hello, how are you doing?', (err, d) => {
  err ? console.error(err) : console.log(d)
})
// listen for event and reply to the channel
ipc.onReceiveAnSend('hola', (channel,data)=>{
  console.log(data)
  ipc.send(channel, 'cool thanks')
})

