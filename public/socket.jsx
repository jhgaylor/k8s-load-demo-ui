const Websocket = {
  setup: (store) => {
    const ws_url = "ws://localhost:8080/ws/pods/load-demo"
    const sock = new WebSocket(ws_url);

    const metricsIntervals = {}

    const getMetrics = (pod) => {
      const name = pod.metadata.name;
      const role = pod.metadata.labels.role;
      if (role !== 'api') {
        store.dispatch({ type: 'GET_METRICS', metrics: null, pod: pod })
        return
      }
      metricsIntervals[name] = setInterval(() => {
        const url = `http://localhost:8001/api/v1/proxy/namespaces/load-demo/pods/${name}/metrics`
        $.get(url, function (data) {
          const lines = data.split('\n').filter((line) => (line[0] !== '#')).map((line) => (line.split(" ")))
          const metrics = {}
          lines.forEach((line) => {
            metrics[line[0]] = line[1]
          })
          store.dispatch({ type: 'GET_METRICS', metrics: metrics, pod: pod })
        })
      }, 3000)
    }

    sock.onopen = (event) => {
      console.log(`Websocket successfully connected to '${ws_url}'.`, event)
    }

    sock.onmessage = (event) => {
      const message = JSON.parse(event.data)
      console.log("onmessage", message)
      const name = message.Object.metadata.name;

      if (message.Type === 'ADDED') {
        // make sure there is something polling the metrics for this pod
        // this mechanism should dispatch an action to set the metrics values on the pod
        getMetrics(message.Object)
      } else if (message.Type === 'DELETED') {
        // clean up the polling mechanism
        clearInterval(metricsIntervals[name])
        delete metricsIntervals[name]
      }
      store.dispatch({ type: message.Type, object: message.Object })
    }

    sock.onerror = (event) => {
      console.info("Todo: Implement proper error handling (recover the connection!)")
      console.error("Websocket Error", event)
    }

    return sock
  }
}
