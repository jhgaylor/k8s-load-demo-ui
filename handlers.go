package main

import (
  "encoding/json"
  "github.com/gorilla/websocket"
  "github.com/juju/loggo"
  "gopkg.in/gin-gonic/gin.v1"
  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
  "k8s.io/client-go/kubernetes"
  "net/http"
)

type JakeRouter interface {
  BindHTTPHandlers()
  BindWSHandlers()
}

type DemoRouter struct {
  *gin.Engine
  k8s *kubernetes.Clientset
}

func (r *DemoRouter) BindHTTPHandlers() {
  healthLogger := loggo.GetLogger("http.health")
  readinessLogger := loggo.GetLogger("http.readiness")
  podsLogger := loggo.GetLogger("http.pods")
  r.GET("/healthz", makeHealthHandler(&healthLogger))
  r.GET("/readiness", makeReadinessHandler(&readinessLogger))
  r.GET("/pods/:namespace", makePodsHandler(r.k8s, &podsLogger))
}

func (r *DemoRouter) BindWSHandlers() {

  // podsLogger := loggo.GetLogger("ws.pods")
  // router.GET("/ws/pods/:namespace", makePodsWatchHandler(k8sClient, podsLogger))
}

func makeHealthHandler(logger *loggo.Logger) func(*gin.Context) {
  return func(c *gin.Context) {
    logger.Debugf("OK")
    c.String(http.StatusOK, "OK")
  }
}

// TODO: this should be aware of if the kubernetes client is throwing errors.
func makeReadinessHandler(logger *loggo.Logger) func(*gin.Context) {
  return func(c *gin.Context) {
    logger.Debugf("OK")
    c.String(http.StatusOK, "OK")
  }
}

func makePodsHandler(k8s *kubernetes.Clientset, logger *loggo.Logger) func(*gin.Context) {
  return func(c *gin.Context) {
    namespace := c.Param("namespace")
    // pods, err := kubeclient.GetPods(namespace)
    logger.Debugf("Attempting to get pods for namespace %s", namespace)
    pods, err := k8s.CoreV1().Pods(namespace).List(metav1.ListOptions{})
    if err != nil {
      logger.Errorf("Failed to get pods. %v", err.Error())
      c.String(http.StatusInternalServerError, err.Error())
      return
    }
    logger.Debugf("Successfully got pods for namespace %s", namespace)
    c.JSON(http.StatusOK, pods)
  }
}

func makePodsWatchHandler(k8s *kubernetes.Clientset, logger *loggo.Logger) func(*gin.Context) {
  var upgrader = websocket.Upgrader{}
  return func(c *gin.Context) {
    namespace := c.Param("namespace")
    logger.Debugf("Attempting to upgrade the http connection to ws.")
    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
      logger.Errorf("Failed to set websocket upgrade: %+v", err)
      c.String(http.StatusInternalServerError, err.Error())
      return
    }
    logger.Debugf("Successfully upgraded the connection to ws.")
    logger.Debugf("Attempting to get a watch on pods for namespace %s", namespace)
    watch, err := k8s.CoreV1().Pods(namespace).Watch(metav1.ListOptions{})
    if err != nil {
      logger.Errorf("Failed to get a watch: %+v", err)
      c.String(http.StatusInternalServerError, err.Error())
      return
    }
    logger.Debugf("Watching pods for namespace %s", namespace)
    ch := watch.ResultChan()
    for {
      msg := <-ch
      msgJson, err := json.Marshal(msg)
      logger.Debugf("Received an event for namespace %s. Event: %v", namespace, string(msgJson))
      if err != nil {
        logger.Errorf("Failed converting event to json: %+v", err)
        c.String(http.StatusInternalServerError, err.Error())
        return
      }
      conn.WriteMessage(websocket.TextMessage, msgJson)
    }
  }
}
