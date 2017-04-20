package main

import (
  "k8s.io/client-go/kubernetes"
  _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
  "k8s.io/client-go/rest"
  "k8s.io/client-go/tools/clientcmd"
)

// type KubernetesService interface {
//   Initialize() (*kubernetes.Clientset, error)
// }

// type KubeClient struct {
//   client *kubernetes.Clientset
// }

// TODO: export a type and interface . this will expose both client mgmt and helpers to get data

func GetKubernetesClient(external bool, kubeconfig string) (*kubernetes.Clientset, error) {
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
