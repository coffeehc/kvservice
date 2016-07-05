package kvservice

type KVClusterConfig struct {
	Nodes    []string `yaml:"nodes" json:"nodes"`
	Replicas int      `yaml:"replicas" json:"replicas"`
}

func (this *KVClusterConfig) Init() {
	if len(this.Nodes) < this.Replicas {
		this.Replicas = len(this.Nodes)
		//TODO 暂时不考虑副本的问题,主要是处理太复杂
	}
}
