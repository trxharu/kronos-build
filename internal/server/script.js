let ws = null
let isConnected = false

function init(addr) {
    ws = new WebSocket(`ws://${addr}/ws`)

    ws.onopen = () => {
        console.log("[Kronos] Development server connected.")
        isConnected = true;
    }

    ws.onmessage = (event) => {
        console.log(event.data)
        location.reload()
    }
}

init("@addr")

if (isConnected) {
    ws.onclose = () => {
        console.error("[Kronos] Development server disconnected.")
        isConnected = false
        setInterval(() => init("@addr"), 2000) 
    }
}
