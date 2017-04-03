# How to ship a new release?

Publish code to github on master and docker hub's automated build service will publish a docker image @ `jhgaylor/k8s-load-demo-api`. The deployment pulls it from docker's public repository.

# How to crank everything up?

* Create the deployment `kubectl create -f k8s-specs/deployment.yaml`

* Create the service `kubectl create -f k8s-specs/service.yaml`

* start the proxy `kubectl proxy`

* load the service through the proxy by browsing to http://localhost:8001/api/v1/proxy/namespaces/default/services/load-demo-api/

