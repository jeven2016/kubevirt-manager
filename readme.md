
## Kubevirt.io
该工程基于kubevirt/api源码工程生成， 且对应kuberntes v1.27.1版本的client工程，官方提供的client-go
依赖包存在缺陷无法在go1.8+以上版本运行，存在bug且社区不活跃一直没有人修复。 故需要以api工程
重新构建对应k8s版本及go版本的客户端。