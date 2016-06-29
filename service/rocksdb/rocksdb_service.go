package rocksdb

import (
	"fmt"
	"github.com/tecbot/gorocksdb"
	"github.com/coffeehc/logger"
	"sync"
)

type RocksDBService interface {
	InitColumnFamily(columnFamily string, opts *gorocksdb.Options) (*gorocksdb.ColumnFamilyHandle, error)
	Get(columnFamily string, opts *gorocksdb.ReadOptions, key []byte) ([]byte, error)
	Put(columnFamily string, opts *gorocksdb.WriteOptions, key, value []byte) error
	Del(columnFamily string, opts *gorocksdb.WriteOptions, key []byte) error
	Close()
}

func NewRocksdbService(opts *gorocksdb.Options, dbPath string) (RocksDBService, error) {
	cfs, err := gorocksdb.ListColumnFamilies(opts, dbPath)
	if err != nil {
		logger.Error("获取 CF 错误:%s", err)
		return nil,err
	}
	var db *gorocksdb.DB
	handlers :=make(map[string]*gorocksdb.ColumnFamilyHandle, 0)
	if len(cfs)>0{
		cfOpts := make([]*gorocksdb.Options,len(cfs))
		for i,_:=range cfs{
			cfOpts[i] = gorocksdb.NewDefaultOptions()
		}
		_db,hs,err:=gorocksdb.OpenDbColumnFamilies(opts,dbPath,cfs,cfOpts)
		if err != nil {
			return nil, fmt.Errorf("打开rocksdb数据文件出错:%s", err)
		}
		db=_db
		for i,cf:=range cfs{
			handlers[cf] = hs[i]
		}
	}else{
		db, err = gorocksdb.OpenDb(opts, dbPath)
		if err != nil {
			return nil, fmt.Errorf("打开rocksdb数据文件出错:%s", err)
		}
	}
	return &_RocksDBService{
		db: db,
		columnFamilies:handlers,
		rwMuext:new(sync.RWMutex),
	}, nil
}

type _RocksDBService struct {
	db             *gorocksdb.DB
	columnFamilies map[string]*gorocksdb.ColumnFamilyHandle
	rwMuext        *sync.RWMutex
}

func (this *_RocksDBService)InitColumnFamily(columnFamily string, opts *gorocksdb.Options) (*gorocksdb.ColumnFamilyHandle, error) {
	if handler, ok := this.columnFamilies[columnFamily]; ok {
		logger.Warn("columnFamily:[%s]已经存在", columnFamily)
		return handler, nil
	}
	handler, err := this.db.CreateColumnFamily(opts, columnFamily)
	if err != nil {
		return nil, err
	}
	this.columnFamilies[columnFamily] = handler
	return handler, nil

}

func (this *_RocksDBService) Get(columnFamily string, opts *gorocksdb.ReadOptions, key []byte) ([]byte, error) {
	if columnFamily == "" {
		return this.db.GetBytes(opts, key)
	}
	cf, ok := this.columnFamilies[columnFamily]
	if !ok {
		return nil, NO_COLUMNFAMILY
	}
	s, err := this.db.GetCF(opts, cf, key)
	defer s.Free()
	if err!=nil{
		return nil,err
	}
	if s.Size() == 0{
		return nil,nil
	}
	data := make([]byte, s.Size())
	copy(data, s.Data())
	return data, err

}
func (this *_RocksDBService) Put(columnFamily string, opts *gorocksdb.WriteOptions, key, value []byte) error {
	if columnFamily == "" {
		return this.db.Put(opts, key, value)
	}
	cf, ok := this.columnFamilies[columnFamily]
	if !ok {
		handler, err := this.InitColumnFamily(columnFamily, gorocksdb.NewDefaultOptions())
		if err != nil {
			return err
		}
		cf = handler
	}
	return this.db.PutCF(opts, cf, key, value)

}
func (this *_RocksDBService) Del(columnFamily string, opts *gorocksdb.WriteOptions, key []byte) error {
	if columnFamily == "" {
		return this.db.Delete(opts, key)
	}
	cf, ok := this.columnFamilies[columnFamily]
	if !ok {
		return NO_COLUMNFAMILY
	}
	return this.db.DeleteCF(opts, cf, key)
}

func (this *_RocksDBService) Close() {
	for _, handler := range this.columnFamilies {
		handler.Destroy()
	}
	this.db.Close()
}

type DataInfo struct {
	Key   []byte
	Value []byte
}
