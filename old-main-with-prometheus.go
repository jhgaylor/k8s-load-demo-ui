package main

// import (
//   "encoding/json"
//   "flag"
//   "fmt"
//   "github.com/prometheus/client_golang/prometheus"
//   "github.com/prometheus/client_golang/prometheus/promhttp"
//   "html"
//   metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
//   "k8s.io/client-go/kubernetes"
//   _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
//   "k8s.io/client-go/rest"
//   "k8s.io/client-go/tools/clientcmd"

//   "log"
//   "net/http"
// )

// func bindHttpHandlers(k8s *kubernetes.Clientset) {
//   http.Handle("/metrics", promhttp.Handler())
//   http.HandleFunc("/", prometheus.InstrumentHandlerFunc("GetPodsInfo", func(w http.ResponseWriter, r *http.Request) {
//     pods, err := k8s.CoreV1().Pods(*namespace).List(metav1.ListOptions{})
//     if err != nil {
//       fmt.Fprintf(w, "Api error, %v", err.Error())
//       panic(err.Error())
//     }
//     podsJson, err := json.Marshal(pods.Items)
//     if err != nil {
//       fmt.Fprintf(w, "JSON encoding error, %v", err.Error())
//     }
//     fmt.Fprintf(w, "There are %d pods in the cluster\n", len(pods.Items))
//     fmt.Fprintf(w, "\n\n%v\n\n", string(podsJson))
//     fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
//   }))
// }

// func main() {
//   external := flag.Bool("external", false, "are we inside of a cluster")
//   kubeconfig := flag.String("kubeconfig", "/Users/jake/.kube/config", "absolute path to the kubeconfig file")
//   namespace := flag.String("namespace", "load-demo", "The namespace to observe.")
//   flag.Parse()

//   var err error
//   var config *rest.Config

//   if *external {
//     // creates a config from current context of kubeconfig
//     config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
//   } else {
//     // creates the in-cluster config
//     config, err = rest.InClusterConfig()
//   }
//   if err != nil {
//     panic(err.Error())
//   }
//   // creates the clientset
//   clientset, err := kubernetes.NewForConfig(config)
//   if err != nil {
//     panic(err.Error())
//   }

//   bindHttpHandlers(clientset)
//   log.Fatal(http.ListenAndServe(":8080", nil))

// }
