const log = msg => {
    document.getElementById('log').textContent += msg + '\n'
}

const ws = new WebSocket("ws://localhost:8080/ws")
ws.onopen    = ()  => log("🟢 connected")
ws.onmessage = e => log("← " + e.data)
ws.onclose   = ()  => log("🔴 disconnected")
ws.onerror   = e  => log("❌ " + e)

function send() {
    const author  = document.getElementById('author').value
    const content = document.getElementById('message').value
    const msg = JSON.stringify({ author, content })
    ws.send(msg)
    log("→ " + msg)
}