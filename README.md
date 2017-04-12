# Running locally

To run locally, you need to indicate to the go binary that we are external to a cluster.

`go run main.go --external --kubeconfig /Users/jake/.kube/config`

Test the WS api

`curl -i -N -H "Connection: Upgrade" -H "Upgrade: websocket" -H "Sec-Websocket-Version: 13" -H "Sec-Websocket-Key: jake" http://localhost:8080/ws/pods/load-demo`


# How to ship a new release?

Publish code to github on master and docker hub's automated build service will publish a docker image @ `jhgaylor/k8s-load-demo-api`. The deployment pulls it from docker's public repository.

# How to crank everything up?

* Make sure that the k8s-load-demo-charts/ui chart is installed in your kubernetes cluster

* start the proxy `kubectl proxy`

* load the service through the proxy by browsing to http://localhost:8001/api/v1/proxy/namespaces/default/services/load-demo-ui/

