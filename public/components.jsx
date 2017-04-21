// See for help: https://github.com/reactjs/redux/blob/master/examples/counter/src/index.js

// Note: you need to be running the proxy to get metrics because the UI is running in the browser.
// Todo: have the api return this data directly so the proxy doesn't need to be running on the demo machine
// Note: you also have to have this in chrome https://chrome.google.com/webstore/detail/allow-control-allow-origi/nlfbmbojpeacfghkpbjhddihlkkiljbi/related?hl=en
//       because CORS.
// http://localhost:8001/api/v1/proxy/namespaces/load-demo/pods/${pod.metadata.name}/metrics
// http://localhost:8001/api/v1/proxy/namespaces/load-demo/pods/load-demo-api-3070421063-f7xt9/metrics

class Pod extends React.Component {
  static propTypes = {
    pod: React.PropTypes.object.isRequired
  }

  render() {
    const {pod} = this.props
    var pod_classes = `pod ${pod.metadata.labels.role}`
    const status = pod.status
    if (status.phase === 'Pending') {
      pod_classes += " pending"
      const conditions = status.conditions || []
      conditions.forEach((condition) => {
        if (condition.reason === 'Unschedulable') {
          pod_classes += " unschedulable"
        }
      })
    }
    const request_total = pod.metrics && pod.metrics['http_requests_total{code="200",handler="PasswordHasher",method="get"}'] || null
    const uptime = moment().diff(pod.metadata.creationTimestamp, 'seconds')
    const restarts = (pod.status.containerStatuses || []).reduce((acc, status) => {
      return status.restartCount + acc
    }, 0)
    var requests_per_second = '--'
    if (request_total || request_total === 0) {
      requests_per_second = (request_total / uptime).toFixed(1)
    }
    // console.log("details", request_total, uptime, requests_per_second)
    return (
      <div className={pod_classes}>
        <p>{pod.metadata.labels.role}</p>
        <p>R/s {requests_per_second}</p>
        <p>Restarts {restarts}</p>
      </div>
    )
  }
}

class PodsContainer extends React.Component {
  static propTypes = {
    data: React.PropTypes.object.isRequired
  }

  getPods = () => {
    const uuids = Object.keys(this.props.data.pods)
    const pods = uuids.map((id) => this.props.data.pods[id])
    return pods
  }

  getPodsGroupedByNode = () => {
    return this.getPods().reduce((memo, pod) => {
      memo[pod.spec.nodeName] = memo[pod.spec.nodeName] || []
      memo[pod.spec.nodeName].push(pod)
      return memo
    }, {})
  }

  getApiPods = () => {
    return this.getPods().filter((pod) => (pod.metadata.labels.role === 'api'))
  }

  getUiPods = () => {
    return this.getPods().filter((pod) => (pod.metadata.labels.role === 'ui'))
  }

  //
  //   { this.renderPods() }
  // </div>
  render() {
    const { data } = this.props
    return (
      <div className="container">
        <p>Pods: {Object.keys(data.pods).length}</p>
        { this.renderPodsGroupedByNodes() }
      </div>
    )
  }

  renderPodsGroupedByNodes() {
    const podsGroupedByNode = this.getPodsGroupedByNode()
    const keys = Object.keys(podsGroupedByNode)
    const nodes = keys.map((nodeName) => (podsGroupedByNode[nodeName]))
    return nodes.map((pods, index) => {
      return (
        <div className="node" key={keys[index]} >
          <div className="pods-container">
            { pods.map(this.renderPod.bind(this))}
          </div>
        </div>
      )
    })
  }

  renderPods() {
    const pods = this.getPods()
    return pods.map(this.renderPod.bind(this))
  }

  renderPod(pod) {
    return (<Pod pod={pod} key={pod.metadata.uid} />)
  }
}
