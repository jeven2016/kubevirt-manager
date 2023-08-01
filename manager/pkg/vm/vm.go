package vm

import (
	"errors"
	"github.com/gin-gonic/gin"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/client-go/tools/clientcmd"
	"kubevirt-client/model"
	kv1 "kubevirt.io/client-go/api/v1"
	"kubevirt.io/client-go/kubecli"
	"log"
	"net/http"
	"strings"
)

type VmManager struct {
	virtClient kubecli.KubevirtClient
}

func NewVmManager(clientConf clientcmd.ClientConfig) (*VmManager, error) {
	virtClient, err := kubecli.GetKubevirtClientFromClientConfig(clientConf)
	if err != nil {
		log.Fatalf("cannot obtain KubeVirt client: %v\n", err)
	}

	return &VmManager{
		virtClient: virtClient,
	}, nil
}

func (m *VmManager) List(ctx *gin.Context) {
	var namespace = ctx.DefaultQuery("namespace", "default")

	//在kubevirt client-go中，VirtualMachine和VirtualMachineInstance是两个不同的概念。
	//VirtualMachine是一个自定义资源（Custom Resource）类型，用于定义虚拟机的规格和配置。它类似于Kubernetes中的Pod，但是比Pod更高级，可以定义虚拟机的CPU、内存、存储等资源，以及网络和设备配置。VirtualMachine对象用于创建和管理虚拟机资源，可以通过client-go库的API进行操作。
	//VirtualMachineInstance是VirtualMachine的实例化对象，表示正在运行的虚拟机实例。它类似于Kubernetes中的Pod实例，但是包含了虚拟机的状态和运行时信息。VirtualMachineInstance对象用于监控和管理虚拟机实例，可以通过client-go库的API获取虚拟机实例的状态、日志和指标等信息，以及执行操作如启动、停止、重启虚拟机。
	//总结起来，VirtualMachine是用于定义虚拟机规格和配置的对象，而VirtualMachineInstance是表示正在运行的虚拟机实例的对象。

	// Fetch list of VMs & VMIs
	//vmList, err := m.virtClient.VirtualMachine(namespace).List(&v1.ListOptions{})
	//if err != nil {
	//	log.Fatal(err)
	//	ctx.AbortWithError(500, err)
	//	return
	//}
	vmiList, err := m.virtClient.VirtualMachineInstance(namespace).List(&metav1.ListOptions{})
	if err != nil {
		ctx.AbortWithError(500, err)
		return
	}

	var list model.VmList
	for _, vmi := range vmiList.Items {
		jsonString, err := json.Marshal(vmi)
		if err != nil {
			ctx.AbortWithError(500, err)
			return
		}
		println(string(jsonString))

		//解析CPU架构
		cpuModel := vmi.Spec.Domain.CPU.Model
		var cpuArchitecture string
		if strings.HasPrefix(cpuModel, "Intel") || strings.HasPrefix(cpuModel, "x86") {
			cpuArchitecture = "x86"
		} else if strings.HasPrefix(cpuModel, "ARM") {
			cpuArchitecture = "ARM"
		} else {
			cpuArchitecture = cpuModel
		}

		//解析硬盘
		var disks []model.DiskInfo
		for _, volume := range vmi.Spec.Volumes {
			if volume.VolumeSource.PersistentVolumeClaim != nil &&
				volume.VolumeSource.PersistentVolumeClaim.ClaimName == "kubevirt.io/volume-data" {
				disk := model.DiskInfo{
					Name: volume.Name,
				}
				disks = append(disks, disk)
			}
		}

		list.Vms = append(list.Vms, model.Vm{
			Name:              vmi.Name,
			Namespace:         vmi.Namespace,
			Uid:               string(vmi.UID),
			CreationTimestamp: vmi.CreationTimestamp.Local(),
			Running:           vmi.IsRunning(),
			Memory:            vmi.Spec.Domain.Resources.Requests.Memory().String(),
			Cpu: model.Cpu{
				Cores:   vmi.Spec.Domain.CPU.Cores,
				Sockets: vmi.Spec.Domain.CPU.Sockets,
				Threads: vmi.Spec.Domain.CPU.Threads,
				Model:   vmi.Spec.Domain.CPU.Model,
				Arch:    cpuArchitecture,
			},
			Disks:    disks,
			Networks: vmi.Spec.Networks,
		})

	}

	ctx.JSON(http.StatusOK, list)
}

func (m *VmManager) CreateVm(ctx *gin.Context) {
	var namespace = ctx.DefaultQuery("namespace", "default")

	//var instance = &kv1.VirtualMachineSpec{}
	//if err := ctx.ShouldBindJSON(&instance); err != nil {
	//	ctx.AbortWithError(500, err)
	//	return
	//}

	vm := m.virtClient.VirtualMachine(namespace)

	running := true

	// 创建一个新的VirtualMachine对象
	vmConfig := &kv1.VirtualMachine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "jevenvm",
			Namespace: "default",
			Labels: map[string]string{
				"app": "jevenvm",
			},

			//Labels: map[string]string{
			//	"kubevirt.io/domain": "jevenvm",
			//	"kubevirt.io/size":   "small",
			//},
		},
		Spec: kv1.VirtualMachineSpec{
			Running: &running,
			Template: &kv1.VirtualMachineInstanceTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{},
				Spec: kv1.VirtualMachineInstanceSpec{
					Domain: kv1.DomainSpec{
						Resources: kv1.ResourceRequirements{
							Requests: v1.ResourceList{
								v1.ResourceMemory: resource.MustParse("256M"),
								//v1.ResourceCPU:     resource.MustParse("2G"),
								//v1.ResourceStorage: resource.MustParse("20G"),
							},
						},
						CPU: &kv1.CPU{
							Cores:   2,
							Sockets: 1,
							Threads: 2,
						},
						Machine: kv1.Machine{
							Type: "q35",
						},
						Devices: kv1.Devices{
							Disks: []kv1.Disk{
								{
									Name: "containerdisk",
									DiskDevice: kv1.DiskDevice{
										Disk: &kv1.DiskTarget{
											Bus: "virtio",
										},
									},
								},
							},
							Interfaces: []kv1.Interface{
								{
									Name: "default",
									InterfaceBindingMethod: kv1.InterfaceBindingMethod{
										Masquerade: &kv1.InterfaceMasquerade{},
									},
								},
							},
						},
					},
					Networks: []kv1.Network{
						{
							Name: "default",
							NetworkSource: kv1.NetworkSource{
								Pod: &kv1.PodNetwork{},
							},
						},
					},
					Volumes: []kv1.Volume{
						{
							Name: "containerdisk",
							VolumeSource: kv1.VolumeSource{
								ContainerDisk: &kv1.ContainerDiskSource{
									Image: "quay.io/kubevirt/cirros-container-disk-demo",
								},
							},
						},
						{
							Name: "cloudInitNoCloud",
							VolumeSource: kv1.VolumeSource{
								CloudInitNoCloud: &kv1.CloudInitNoCloudSource{
									UserDataBase64: "SGkuXG4=",
								},
							},
						},
					},
				},
			},
		},
	}

	inst, err := vm.Create(vmConfig)
	if err != nil {
		ctx.AbortWithError(500, err)
		return
	}

	ctx.JSON(http.StatusCreated, inst)
}

func (m *VmManager) Action(ctx *gin.Context) {
	name := ctx.Param("name")
	action := ctx.Query("action")

	vmInterface := m.virtClient.VirtualMachineInstance("default")

	switch action {
	case "start":
		vm, err := vmInterface.Get(name, &metav1.GetOptions{})
		if err != nil {
			ctx.AbortWithError(500, err)
			return
		}
		vm.Status.Phase = kv1.Running
		_, err = vmInterface.Update()
		if err != nil {
			ctx.AbortWithError(500, err)
			return
		}
		ctx.JSON(http.StatusOK, "start")
	default:
		ctx.AbortWithError(http.StatusBadRequest, errors.New("action not support"))

	}
}
