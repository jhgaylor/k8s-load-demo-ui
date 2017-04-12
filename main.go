package main

import (
  "encoding/json"
  "flag"
  "fmt"
  "github.com/gorilla/websocket"
  "gopkg.in/gin-gonic/gin.v1"
  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
  "k8s.io/client-go/kubernetes"
  _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
  "k8s.io/client-go/rest"
  "k8s.io/client-go/tools/clientcmd"
  "net/http"
  "os"
)

var upgrader = websocket.Upgrader{}

func bindWsHandlers(k8s *kubernetes.Clientset, router *gin.Engine) {
  router.GET("/ws/pods/:namespace", func(c *gin.Context) {
    namespace := c.Param("namespace")
    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
      fmt.Println("Failed to set websocket upgrade: %+v", err)
      return
    }
    watch, err := k8s.CoreV1().Pods(namespace).Watch(metav1.ListOptions{})
    if err != nil {
      fmt.Println("Failed to get a watch: %+v", err)
      return
    }
    ch := watch.ResultChan()
    for {
      msg := <-ch
      msgJson, err := json.Marshal(msg)
      if err != nil {
        fmt.Println("Failed converting event to json: %+v", err)
        return
      }
      conn.WriteMessage(websocket.TextMessage, msgJson)
    }
  })
}

func bindHttpHandlers(k8s *kubernetes.Clientset, router *gin.Engine) {
  router.GET("/ping", func(c *gin.Context) {
    c.String(http.StatusOK, "Pong")
  })
  router.GET("/pods/:namespace", func(c *gin.Context) {
    namespace := c.Param("namespace")
    pods, err := k8s.CoreV1().Pods(namespace).List(metav1.ListOptions{})
    if err != nil {
      c.String(http.StatusInternalServerError, err.Error())
    }
    c.JSON(http.StatusOK, pods)
  })
}

func getKubernetesClient(external bool, kubeconfig string) (*kubernetes.Clientset, error) {
  var err error
  var config *rest.Config

  if external {
    // creates a config from current context of kubeconfig
    config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
  } else {
    // creates the in-cluster config
    config, err = rest.InClusterConfig()
  }
  if err != nil {
    return nil, err
  }

  // creates the clientset
  clientset, err := kubernetes.NewForConfig(config)
  if err != nil {
    return nil, err
  }

  return clientset, nil
}

func main() {
  external := flag.Bool("external", false, "Are we outside of a cluster?")
  kubeconfig := flag.String("kubeconfig", "/Users/jake/.kube/config", "absolute path to the kubeconfig file")
  flag.Parse()

  port := "8080"
  if _port := os.Getenv("PORT"); _port != "" {
    port = _port
  }

  k8sClient, err := getKubernetesClient(*external, *kubeconfig)
  if err != nil {
    panic(err.Error())
  }

  router := gin.Default()

  bindHttpHandlers(k8sClient, router)
  bindWsHandlers(k8sClient, router)
  router.Run(fmt.Sprintf(":%v", port))
}
