const Websocket = {
  setup: (store) => {
    const ws_url = "ws://localhost:8080/ws/pods/load-demo"
    const sock = new WebSocket(ws_url);

    sock.onopen = (event) => {
      console.log(`Websocket successfully connected to '${ws_url}'.`, event)
    }

    sock.onmessage = (event) => {
      const message = JSON.parse(event.data)
      console.log("onmessage", message)
      store.dispatch({ type: message.Type, object: message.Object })
    }

    sock.onerror = (event) => {
      console.info("Todo: Implement proper error handling (recover the connection!)")
      console.error("Websocket Error", event)
    }

    return sock
  }
}
