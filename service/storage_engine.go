package main

import (
	"baseservices/kvservice/service/rocksdb"
	"fmt"
	"github.com/tecbot/gorocksdb"
	"path"
)

type StorageEngine interface {
	Get(columnFamily string,opts *gorocksdb.ReadOptions, key []byte) ([]byte, error)
	Put(columnFamily string,opts *gorocksdb.WriteOptions, key, value []byte) error
	Del(columnFamily string,opts *gorocksdb.WriteOptions, key []byte) error
	Close()
	GetPartition() int
	IsReplica() bool
}

func NewStorageEngine(parentDir string, partition int, replica bool) (StorageEngine, error) {
	engine := &_StorageEngine{
		partition:    partition,
		replica:      replica,
		writeOptions: gorocksdb.NewDefaultWriteOptions(),
		readOptions:  gorocksdb.NewDefaultReadOptions(),
	}
	opts := gorocksdb.NewDefaultOptions()
	opts.SetCreateIfMissing(true)
	opts.SetKeepLogFileNum(3)
	evn := gorocksdb.NewDefaultEnv()
	opts.SetEnv(evn)
	rocksDBService, err := rocksdb.NewRocksdbService(opts, path.Join(parentDir, fmt.Sprintf("%d", partition)))
	if err != nil {
		return nil, err
	}
	engine.rocksDBService = rocksDBService
	return engine,nil
}

type _StorageEngine struct {
	partition      int
	replica        bool
	rocksDBService rocksdb.RocksDBService
	writeOptions   *gorocksdb.WriteOptions
	readOptions    *gorocksdb.ReadOptions
}

func (this *_StorageEngine) GetPartition() int {
	return this.partition
}

func (this *_StorageEngine) IsReplica() bool {
	return this.replica
}

func (this *_StorageEngine)Get(columnFamily string,opts *gorocksdb.ReadOptions,key []byte) ([]byte, error){
	if opts == nil{
		opts = this.readOptions
	}
	return this.rocksDBService.Get(columnFamily,opts,key)
}
func (this *_StorageEngine)Put(columnFamily string,opts *gorocksdb.WriteOptions,key, value []byte) error{
	if opts == nil{
		opts = this.writeOptions
	}
	return this.rocksDBService.Put(columnFamily,opts,key,value)

}
func (this *_StorageEngine)Del(columnFamily string,opts *gorocksdb.WriteOptions, key []byte) error{
	if opts == nil{
		opts = this.writeOptions
	}
	return this.rocksDBService.Del(columnFamily,opts,key)
}

func (this *_StorageEngine)Close(){
	this.rocksDBService.Close()
}
