package main

import (
	"github.com/coffeehc/microserviceboot/consultool"
	"baseservices/kvservice"
)

type Config struct {
	ConsulConfig  *consultool.ConsulConfig `yaml:"consul"`
	StorageConfig *StorageConfig                   `yaml:"storage"`
}

type StorageConfig struct {
	StorageDir string `yaml:"storageDir"`
	NodeId     int    `yaml:"nodeId"`
	KVClusterConfig *kvservice.KVClusterConfig `yaml:"vkClusterConfig"`
}
