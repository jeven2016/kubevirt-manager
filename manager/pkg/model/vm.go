package model

import (
	"time"
)

type Cpu struct {
	Cores   uint32 `json:"cores"`
	Sockets uint32 `json:"sockets"`
	Threads uint32 `json:"threads"`
	Model   string `json:"model"`
	Arch    string `json:"arch"`
}

type DiskInfo struct {
	Name string `json:"name"`
}

type Vm struct {
	Name              string     `json:"name"`
	Namespace         string     `json:"namespace"`
	Uid               string     `json:"uid"`
	CreationTimestamp time.Time  `json:"creationTimestamp"`
	Running           bool       `json:"running"`
	Memory            string     `json:"memory"`
	Cpu               Cpu        `json:"cpu"`
	Disks             []DiskInfo `json:"disks"`
	//Networks          []v1.Network `json:"network"`
}

type VmList struct {
	Vms []Vm `json:"vms"`
}
