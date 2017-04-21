package main

import (
  "flag"
  "fmt"
  "github.com/gin-contrib/cors"
  "github.com/gin-gonic/gin"
  "github.com/juju/loggo"
  "os"
)

func main() {
  external := flag.Bool("external", false, "Are we outside of a cluster?")
  kubeconfig := flag.String("kubeconfig", "/Users/jake/.kube/config", "absolute path to the kubeconfig file")
  debug := flag.Bool("debug", false, "Are we in debug mode?")
  flag.Parse()

  ConfigureLogger(*debug)
  logger := loggo.GetLogger("app.main")

  port := "8080"
  if _port := os.Getenv("PORT"); _port != "" {
    port = _port
  }

  k8sClient, err := GetKubernetesClient(*external, *kubeconfig)
  if err != nil {
    panic(err.Error())
  }
  logger.Debugf("Connected to Kubernetes.")

  router := gin.Default()
  demoRouter := DemoRouter{Engine: router, k8s: k8sClient}

  demoRouter.BindHTTPHandlers()
  demoRouter.BindWSHandlers()
  demoRouter.Use(cors.Default())
  demoRouter.Run(fmt.Sprintf(":%v", port))
}
