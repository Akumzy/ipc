const IPC = require('ipc-node-go').default
const ipc = new IPC('./example').init()
ipc.on('log', console.log)
let count = 0

// setInterval(() => {
//   count++
//   ipc.send('count', count)
//   ipc.send('count-object', { num: count })
// }, 2000)
ipc.on('hello', d => {
  console.log(d)
})
ipc.sendAndReceive('yoo', '/home/akumzy/Documents/Petrobase_Drive', (err, d) => {
  console.log(err)
  console.log(d)
})
process.on('exit', () => {
  ipc.kill()
})
// { "event": "app:addRecursive", "data": "/home/akumzy/Documents/Petrobase_Drive", "SR": true }
