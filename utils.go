package kvservice

import (
	"github.com/coffeehc/microserviceboot/utils"
	"hash/crc32"
)

var _consistentHashing= crc32.ChecksumIEEE

func GetConsistentHash(key []byte,partition int64) int64 {
	index := _consistentHashing(key)
	return utils.JumpConsistentHash(uint64(index),partition)
}
