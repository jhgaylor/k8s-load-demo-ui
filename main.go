package main

import (
  "flag"
  "fmt"
  "html"
  "log"
  "net/http"
  // "time"
  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
  "k8s.io/client-go/kubernetes"
  _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
  "k8s.io/client-go/rest"
  "k8s.io/client-go/tools/clientcmd"
)

func main() {
  external := flag.Bool("external", false, "are we inside of a cluster")
  kubeconfig := flag.String("kubeconfig", "/Users/jake/.kube/config", "absolute path to the kubeconfig file")
  flag.Parse()

  var err error
  var config *rest.Config

  if *external {
    // creates a config from current context of kubeconfig
    config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
  } else {
    // creates the in-cluster config
    config, err = rest.InClusterConfig()
  }
  if err != nil {
    panic(err.Error())
  }
  // creates the clientset
  clientset, err := kubernetes.NewForConfig(config)
  if err != nil {
    panic(err.Error())
  }
  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
    if err != nil {
      fmt.Fprintf(w, "Error, %v", err.Error())
      panic(err.Error())
    }
    fmt.Fprintf(w, "There are %d pods in the cluster\n", len(pods.Items))
    fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
  })

  log.Fatal(http.ListenAndServe(":8080", nil))

}
