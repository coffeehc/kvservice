package main

import (
	"baseservices/kvservice"
	"fmt"
	"os"

	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/tecbot/gorocksdb"
)

const DB_PARTITION = 1024

type DataService interface {
	Get(columnFamily string, key []byte) ([]byte, error)
	Put(columnFamily string, key, value []byte) error
	Del(columnFamily string, key []byte) error
	GetAll(columnFamily string, opts *gorocksdb.ReadOptions, prefixKey []byte) Iterator
	GetEngine(key []byte) (StorageEngine, error)
	Close()
}

type StorageService interface {
	DataService
}

func NewStorageService(config *StorageConfig) (StorageService, error) {
	if config.NodeId == 0 {
		return nil, errors.New("nodeId不能为0")
	}
	clusterConfig := config.KVClusterConfig
	clusterConfig.Init()
	err := os.MkdirAll(config.StorageDir, 0760)
	if err != nil {
		return nil, fmt.Errorf("创建存储目录[%s]出错:%s", config.StorageDir, err)
	}
	engines := make(map[int64]StorageEngine, DB_PARTITION)
	for i := 0; i < DB_PARTITION; i++ {
		if i%config.NodeId == 0 {
			engine, err := NewStorageEngine(config.StorageDir, i, false)
			if err != nil {
				return nil, err
			}
			engines[int64(i)] = engine
		}
	}
	storageService := &_StorageService{
		engins: engines,
	}
	return storageService, nil
}

type _StorageService struct {
	engins map[int64]StorageEngine
}

var NO_ENGINE = errors.New("key没有存储在该节点")

func (this *_StorageService) GetEngine(key []byte) (StorageEngine, error) {
	index := kvservice.GetConsistentHash(key, DB_PARTITION)
	engine, ok := this.engins[index]
	if !ok {
		return engine, NO_ENGINE
	}
	return engine, nil
}

func (this *_StorageService) Get(columnFamily string, key []byte) ([]byte, error) {
	engine, err := this.GetEngine(key)
	if err != nil {
		return nil, err
	}
	return engine.Get(columnFamily, nil, key)
}

func (this *_StorageService) Put(columnFamily string, key, value []byte) error {
	engine, err := this.GetEngine(key)
	if err != nil {
		return err
	}
	return engine.Put(columnFamily, nil, key, value)
}

func (this *_StorageService) Del(columnFamily string, key []byte) error {
	engine, err := this.GetEngine(key)
	if err != nil {
		return err
	}
	return engine.Del(columnFamily, nil, key)
}

func (this *_StorageService) GetAll(columnFamily string, opts *gorocksdb.ReadOptions, prefixKey []byte) Iterator {
	iterator := newIterator(1024, prefixKey)
	for _, engine := range this.engins {
		iter, err := engine.GetAll(columnFamily, opts, prefixKey)
		if err == nil {
			iterator.add(iter)
		}
	}
	iterator.wait()
	return iterator
}

func (this *_StorageService) Close() {
	for index, engine := range this.engins {
		engine.Close()
		delete(this.engins, index)
	}
}
