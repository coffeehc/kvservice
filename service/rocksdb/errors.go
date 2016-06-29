package rocksdb

import "github.com/syndtr/goleveldb/leveldb/errors"

var (
	NO_COLUMNFAMILY = errors.New("columnFamily 不存在")
)
