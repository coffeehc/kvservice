package main

import (
	"github.com/coffeehc/kvservice"
	"fmt"
	"os"

	"github.com/coffeehc/logger"
	"github.com/tecbot/gorocksdb"
	"sync"
	"github.com/coffeehc/microserviceboot/base"
)

const DB_PARTITION = 1024

type DataService interface {
	Get(columnFamily string, key []byte) ([]byte, base.Error)
	Put(columnFamily string, key, value []byte) base.Error
	Del(columnFamily string, key []byte) base.Error
	GetAll(columnFamily string, opts *gorocksdb.ReadOptions, prefixKey, startKey []byte, order string) Iterator
	GetEngine(key []byte) (StorageEngine, base.Error)
	Close()
}

type StorageService interface {
	DataService
}

func NewStorageService(config *StorageConfig) (StorageService, base.Error) {
	if config.NodeId == 0 {
		return nil, base.NewError(base.ERROR_CODE_BASE_CONFIG_ERROR,"nodeId不能为0")
	}
	clusterConfig := config.KVClusterConfig
	clusterConfig.Init()
	err := os.MkdirAll(config.StorageDir, 0760)
	if err != nil {
		return nil, base.NewError(base.ERROR_CODE_BASE_INIT_ERROR,fmt.Sprintf("创建存储目录[%s]出错:%s", config.StorageDir, err))
	}
	engines := make(map[int]StorageEngine, DB_PARTITION)
	for i := 0; i < DB_PARTITION; i++ {
		if i%config.NodeId == 0 {
			engine, err := NewStorageEngine(config.StorageDir, i, false)
			if err != nil {
				return nil, err
			}
			engines[i] = engine
		}
	}
	logger.Info("启动存储引擎[%d]个", len(engines))
	storageService := &_StorageService{
		engins:    engines,
		partition: DB_PARTITION,
	}
	return storageService, nil
}

type _StorageService struct {
	engins    map[int]StorageEngine
	partition int
}

func (this *_StorageService) GetEngine(key []byte) (StorageEngine, base.Error) {
	index := kvservice.GetConsistentHash(key, DB_PARTITION)
	engine, ok := this.engins[index]
	if !ok {
		return engine, base.NewError(kvservice.ERROR_CODE_KVSERVICE_NOFIND_PARTITION,"key没有存储在该节点")
	}
	return engine, nil
}

func (this *_StorageService) Get(columnFamily string, key []byte) ([]byte, base.Error) {
	engine, err := this.GetEngine(key)
	if err != nil {
		return nil, err
	}
	return engine.Get(columnFamily, nil, key)
}

func (this *_StorageService) Put(columnFamily string, key, value []byte) base.Error {
	engine, err := this.GetEngine(key)
	if err != nil {
		return err
	}
	return engine.Put(columnFamily, nil, key, value)
}

func (this *_StorageService) Del(columnFamily string, key []byte) base.Error {
	engine, err := this.GetEngine(key)
	if err != nil {
		return err
	}
	return engine.Del(columnFamily, nil, key)
}

func (this *_StorageService) GetAll(columnFamily string, opts *gorocksdb.ReadOptions, prefixKey, startKey []byte, order string) Iterator {
	iterator := newIterator(this.partition, prefixKey, startKey, order)
	wait := new(sync.WaitGroup)
	for partition, engine := range this.engins {
		wait.Add(1)
		func(_partition int, _engine StorageEngine) {
			defer wait.Done()
			iter, err := _engine.GetAll(columnFamily, opts)
			if err == nil {
				iterator.add(_partition, iter)
			}
		}(partition, engine)
	}
	wait.Wait()
	iterator.addEnd()
	return iterator
}

func (this *_StorageService) Close() {
	for index, engine := range this.engins {
		engine.Close()
		delete(this.engins, index)
	}
}
