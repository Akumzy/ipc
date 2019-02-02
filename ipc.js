import { EventEmitter } from 'events'
import { spawn } from 'child_process'
class IPC extends EventEmitter {
  constructor(binPath, arg) {
    super()
    // Path to Golang binary
    this.binPath = binPath
    this.arg = arg || []
    this.go = null
    this.closed = false
  }
  init() {
    this.closed = false
    const self = this
    const go = spawn(this.binPath, [...this.arg])
    this.go = go
    go.stderr.setEncoding('utf8')
    go.stdout.setEncoding('utf8')
    ;['close', 'error', 'end', 'data'].forEach(event => go.stderr.on(event, e => self.emit('log', e)))
    let outBuffer = ''
    go.stdout.on('data', s => {
      if (isJSON(s)) {
        let payload = parseJSON(s)
        if (typeof payload === 'object' && payload !== null) {
          self.emit('data', payload)
          let { error, data, event } = payload
          self.emit(event, data, error)
        }

        return
      }
      outBuffer += s
      if (s.endsWith('\\n')) {
        let d = outBuffer.split('\\n')
        let payload = parseJSON(d[0])
        if (typeof payload === 'object' && payload !== null) {
          self.emit('data', payload)
          let { error, data, event } = payload
          self.emit(event, data, error)
        }
        outBuffer = ''
      }
    })
  }

  kill() {
    this.closed = true
    if (this.go) this.go.kill()
  }
  send(eventType, data) {
    this._send(eventType, data, false)
  }
  _send(eventType, data, SR) {
    if (!this.go || this.closed) return
    if (this.go.killed) return
    if (this.go && this.go.stdin) {
      let payload
      if (typeof data === 'object' || Array.isArray(data)) payload = JSON.stringify(data)
      else payload = data
      let d = JSON.stringify({
        event: eventType,
        data: payload,
        SR: !!SR
      })
      if (this.go.stdin) {
        this.go.stdin.write(d + '\n')
      }
    }
  }
  sendAndReceive(eventName, data, cb) {
    this.send(eventName, eventName, data, true)
    let rc = eventName + '___RC___'
    this.on(rc, (data, error) => {
      if (typeof cb === 'function') cb(data, error)
    })
  }
}
function parseJSON(s) {
  try {
    let data = s.replace(/\n/g, ',').replace(/\'\'/, "'")
    if (data.endsWith(',')) {
      data = data.slice(0, -1)
    }
    let payload = JSON.parse(`[${data}]`)
    return payload[0]
  } catch (error) {
    return null
  }
}
function isJSON(s) {
  try {
    JSON.parse(s)
    return true
  } catch (error) {
    return false
  }
}

export default IPC
