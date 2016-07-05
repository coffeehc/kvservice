package rocksdb

import (
	"fmt"
	"os"
	"path"
	"sync"

	"github.com/coffeehc/logger"
	"github.com/tecbot/gorocksdb"
)

const Defalut_FC = "default"

type RocksDBService interface {
	InitColumnFamily(columnFamily string, opts *gorocksdb.Options) (*gorocksdb.ColumnFamilyHandle, error)
	Get(columnFamily string, opts *gorocksdb.ReadOptions, key []byte) ([]byte, error)
	Put(columnFamily string, opts *gorocksdb.WriteOptions, key, value []byte) error
	Del(columnFamily string, opts *gorocksdb.WriteOptions, key []byte) error
	GetAll(columnFamily string, opts *gorocksdb.ReadOptions, prefixKey []byte) (*gorocksdb.Iterator, error)
	Close()
}

func NewRocksdbService(opts *gorocksdb.Options, dbPath string) (RocksDBService, error) {
	_, err := os.Stat(path.Join(dbPath, "CURRENT"))
	cfs := make([]string, 0)
	if err == nil {
		cfs, err = gorocksdb.ListColumnFamilies(opts, dbPath)
		if err != nil {
			logger.Error("获取 CF 错误:%s", err)
			return nil, err
		}
	} else {
		cfs = []string{Defalut_FC}
	}
	var db *gorocksdb.DB
	handlers := make(map[string]*gorocksdb.ColumnFamilyHandle, 0)
	cfOpts := make([]*gorocksdb.Options, len(cfs))
	for i, _ := range cfs {
		cfOpts[i] = gorocksdb.NewDefaultOptions()
	}
	opts.SetCreateIfMissingColumnFamilies(true)
	_db, hs, err := gorocksdb.OpenDbColumnFamilies(opts, dbPath, cfs, cfOpts)
	if err != nil {
		return nil, fmt.Errorf("打开rocksdb数据文件出错:%s", err)
	}
	db = _db
	for i, cf := range cfs {
		handlers[cf] = hs[i]
	}
	return &_RocksDBService{
		db:             db,
		columnFamilies: handlers,
		rwMuext:        new(sync.RWMutex),
	}, nil
}

type _RocksDBService struct {
	db             *gorocksdb.DB
	columnFamilies map[string]*gorocksdb.ColumnFamilyHandle
	rwMuext        *sync.RWMutex
}

func (this *_RocksDBService) InitColumnFamily(columnFamily string, opts *gorocksdb.Options) (*gorocksdb.ColumnFamilyHandle, error) {
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
	cfHandler, err := this.getCFHandler(columnFamily, true)
	if err != nil {
		return nil, err
	}
	s, err := this.db.GetCF(opts, cfHandler, key)
	defer s.Free()
	if err != nil {
		return nil, err
	}
	if s.Size() == 0 {
		return nil, nil
	}
	data := make([]byte, s.Size())
	copy(data, s.Data())
	return data, err

}
func (this *_RocksDBService) Put(columnFamily string, opts *gorocksdb.WriteOptions, key, value []byte) error {
	cfHandler, err := this.getCFHandler(columnFamily, true)
	if err != nil {
		return err
	}
	return this.db.PutCF(opts, cfHandler, key, value)

}
func (this *_RocksDBService) Del(columnFamily string, opts *gorocksdb.WriteOptions, key []byte) error {
	cfHandler, err := this.getCFHandler(columnFamily, true)
	if err != nil {
		return err
	}
	return this.db.DeleteCF(opts, cfHandler, key)
}

func (this *_RocksDBService) GetAll(columnFamily string, opts *gorocksdb.ReadOptions, prefixKey []byte) (*gorocksdb.Iterator, error) {
	cfHandler, err := this.getCFHandler(columnFamily, true)
	if err != nil {
		return nil, err
	}
	return this.db.NewIteratorCF(opts, cfHandler), nil
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

func (this *_RocksDBService) getCFHandler(cfName string, missIsCreate bool) (*gorocksdb.ColumnFamilyHandle, error) {
	if cfName == "" {
		cfName = Defalut_FC
	}
	if handler, ok := this.columnFamilies[cfName]; ok {
		return handler, nil
	}
	if missIsCreate {
		return this.InitColumnFamily(cfName, gorocksdb.NewDefaultOptions())
	}
	return nil, NO_COLUMNFAMILY

}
