package main

import (
	"context"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"kubevirt.io/client-go/versioned"
)

func main() {
	cfg, err := clientcmd.BuildConfigFromFlags("", "/home/cloud/.kube/config")
	if err != nil {
		panic(err)
	}
	client := versioned.NewForConfigOrDie(cfg)
	list, err := client.KubevirtV1().VirtualMachineInstances(v1.NamespaceAll).List(context.Background(), v1.ListOptions{})
	if err != nil {
		panic(err)
	}
	println(len(list.Items))

}
