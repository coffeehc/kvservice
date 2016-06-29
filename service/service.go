package main

import (
	"baseservices/kvservice"
	"github.com/coffeehc/microserviceboot/base"
	"github.com/coffeehc/microserviceboot/consultool"
	"github.com/coffeehc/microserviceboot/serviceboot"
	"github.com/coffeehc/web"
)

func main() {
	serviceboot.ServiceLauncher(&KVService{})
}

type KVService struct {
	config                   *Config
	serviceDiscoveryRegister base.ServiceDiscoveryRegister
	storageService           StorageService
}

func (this *KVService) Init(configPath string, server *web.Server) error {
	serviceConfig := new(Config)
	err := base.LoadConfig(base.GetDefaultConfigPath(configPath), serviceConfig)
	if err != nil {
		return err
	}
	this.config = serviceConfig
	serviceDiscoveryRegister, err := consultool.NewConsulServiceRegister(consultool.WarpConsulConfig(serviceConfig.ConsulConfig))
	if err != nil {
		return err
	}
	this.serviceDiscoveryRegister = serviceDiscoveryRegister
	this.storageService, err = NewStorageService(serviceConfig.StorageConfig)
	if err != nil {
		return err
	}
	return nil
}

func (this *KVService) Run() error {
	return nil

}
func (this *KVService) Stop() error {
	this.storageService.Close()
	return nil
}
func (this *KVService) GetServiceInfo() base.ServiceInfo {
	return &kvservice.KVServiceInfo{}
}
func (this *KVService) GetEndPoints() []base.EndPoint {
	return []base.EndPoint{
		base.EndPoint{
			Metadata:    kvservice.GET_VALUE,
			HandlerFunc: this.get_value,
		},
		base.EndPoint{
			Metadata:    kvservice.POST_VALUE,
			HandlerFunc: this.post_value,
		},
		base.EndPoint{
			Metadata:    kvservice.DEL_KEY,
			HandlerFunc: this.del_key,
		},
	}
}

func (this *KVService) GetServiceDiscoveryRegister() base.ServiceDiscoveryRegister {
	return this.serviceDiscoveryRegister
}
