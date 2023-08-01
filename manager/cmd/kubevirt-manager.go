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

	// kubecli.DefaultClientConfig() prepares config using kubeconfig.
	// typically, you need to set env variable, kubeconfig=<path-to-kubeconfig>/.kubeconfig
	//os.Setenv("kubeconfig", "~/.kube/config")
	//clientConfig := kubecli.DefaultClientConfig(&pflag.FlagSet{})

	//manager, err := vm.NewVmManager(clientConfig)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//gin.SetMode(gin.ReleaseMode)
	//var root = gin.Default()
	//
	//root.GET("/vms", manager.List)
	//root.POST("/vms", manager.CreateVm)
	//root.POST("/vms/:name", manager.Action)
	//
	//// run as a web server
	//log.Println("starting...")
	//if err = root.Run(":9999"); err != nil {
	//	log.Fatal("unable to start web server", err)
	//}
}
