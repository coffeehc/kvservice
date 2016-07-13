package kvservice

import (
	"hash/crc32"

	"github.com/coffeehc/microserviceboot/utils"
)

var _consistentHashing = crc32.ChecksumIEEE

func GetConsistentHash(key []byte, partition int) int {
	index := _consistentHashing(key)
	return int(utils.JumpConsistentHash(uint64(index), int64(partition)))
}
