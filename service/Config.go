package main

import (
	"baseservices/kvservice"

	"github.com/coffeehc/microserviceboot/consultool"
)

type Config struct {
	ConsulConfig  *consultool.ConsulConfig `yaml:"consul"`
	StorageConfig *StorageConfig           `yaml:"storage"`
}

type StorageConfig struct {
	StorageDir      string                     `yaml:"storageDir"`
	NodeId          int                        `yaml:"nodeId"`
	KVClusterConfig *kvservice.KVClusterConfig `yaml:"vkClusterConfig"`
}
