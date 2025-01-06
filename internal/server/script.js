const ws = new WebSocket("ws://@addr/ws")
const handler = null

ws.onopen = () => {
    console.log("[Kronos] Development server connected.")
    isConnected = true;
}

ws.onmessage = (event) => {
    if (handler !== null) {
        clearTimeout(handler)
    }

    handler = setTimeout(() => {
        console.log(event)
        ws.close()
        location.reload()
    }, 500)
}

ws.onclose = () => {
    console.error("[Kronos] Development server disconnected.")
}

