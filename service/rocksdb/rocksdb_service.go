package rocksdb

import (
	"fmt"
	"os"
	"path"

	"github.com/coffeehc/logger"
	"github.com/tecbot/gorocksdb"
	"github.com/coffeehc/microserviceboot/base"
	"github.com/coffeehc/kvservice"
)

const Defalut_FC = "default"

type RocksDBService interface {
	InitColumnFamily(columnFamily string, opts *gorocksdb.Options) (*gorocksdb.ColumnFamilyHandle, base.Error)
	Get(columnFamily string, opts *gorocksdb.ReadOptions, key []byte) ([]byte, base.Error)
	Put(columnFamily string, opts *gorocksdb.WriteOptions, key, value []byte) base.Error
	Del(columnFamily string, opts *gorocksdb.WriteOptions, key []byte) base.Error
	GetAll(columnFamily string, opts *gorocksdb.ReadOptions) (*gorocksdb.Iterator, base.Error)
	Close()
}

func NewRocksdbService(opts *gorocksdb.Options, dbPath string) (RocksDBService, base.Error) {
	_, err := os.Stat(path.Join(dbPath, "CURRENT"))
	cfs := make([]string, 0)
	if err == nil {
		cfs, err = gorocksdb.ListColumnFamilies(opts, dbPath)
		if err != nil {
			logger.Error("获取 CF 错误:%s", err)
			return nil, base.NewError(base.ERROR_CODE_BASE_INIT_ERROR, err.Error())
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
		return nil, base.NewError(base.ERROR_CODE_BASE_INIT_ERROR, fmt.Sprintf("打开rocksdb数据文件出错:%s", err))
	}
	db = _db
	for i, cf := range cfs {
		handlers[cf] = hs[i]
	}
	return &_RocksDBService{
		db:             db,
		columnFamilies: handlers,
	}, nil
}

type _RocksDBService struct {
	db             *gorocksdb.DB
	columnFamilies map[string]*gorocksdb.ColumnFamilyHandle
}

func (this *_RocksDBService) InitColumnFamily(columnFamily string, opts *gorocksdb.Options) (*gorocksdb.ColumnFamilyHandle, base.Error) {
	if handler, ok := this.columnFamilies[columnFamily]; ok {
		logger.Warn("columnFamily:[%s]已经存在", columnFamily)
		return handler, nil
	}
	handler, err := this.db.CreateColumnFamily(opts, columnFamily)
	if err != nil {
		return nil, base.NewError(kvservice.ERROR_CODE_KVSERVICE_CREATE_CF_ERROR, err.Error())
	}
	this.columnFamilies[columnFamily] = handler
	return handler, nil

}

func (this *_RocksDBService) Get(columnFamily string, opts *gorocksdb.ReadOptions, key []byte) ([]byte, base.Error) {
	cfHandler, err := this.getCFHandler(columnFamily, true)
	if err != nil {
		return nil, err
	}
	s, err1 := this.db.GetCF(opts, cfHandler, key)
	if err1 != nil {
		return nil, base.NewError(kvservice.ERROR_CODE_KVSERVICE_GET_DATA_ERROR, err1.Error())
	}
	defer s.Free()
	if s.Size() == 0 {
		return nil, nil
	}
	data := make([]byte, s.Size())
	copy(data, s.Data())
	return data, nil

}
func (this *_RocksDBService) Put(columnFamily string, opts *gorocksdb.WriteOptions, key, value []byte) base.Error {
	cfHandler, err := this.getCFHandler(columnFamily, true)
	if err != nil {
		return err
	}
	err1 := this.db.PutCF(opts, cfHandler, key, value)
	if err1 != nil {
		return base.NewError(kvservice.ERROR_CODE_KVSERVICE_PUT_DATA_ERROR, err1.Error())
	}
	return nil
}
func (this *_RocksDBService) Del(columnFamily string, opts *gorocksdb.WriteOptions, key []byte) base.Error {
	cfHandler, err := this.getCFHandler(columnFamily, true)
	if err != nil {
		return err
	}
	err1 := this.db.DeleteCF(opts, cfHandler, key)
	if err1 != nil {
		return base.NewError(kvservice.ERROR_CODE_KVSERVICE_DELETE_DATA_ERROR, err1.Error())
	}
	return nil
}

func (this *_RocksDBService) GetAll(columnFamily string, opts *gorocksdb.ReadOptions) (*gorocksdb.Iterator, base.Error) {
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

func (this *_RocksDBService) getCFHandler(cfName string, missIsCreate bool) (*gorocksdb.ColumnFamilyHandle, base.Error) {
	if cfName == "" {
		cfName = Defalut_FC
	}
	if handler, ok := this.columnFamilies[cfName]; ok {
		return handler, nil
	}
	if missIsCreate {
		return this.InitColumnFamily(cfName, gorocksdb.NewDefaultOptions())
	}
	return nil, base.NewError(kvservice.ERROR_CODE_KVSERVICE_GET_CF_ERROR,fmt.Sprintf("没有[%s]的 CF",cfName))

}
