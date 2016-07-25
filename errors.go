package kvservice

const ERROR_CODE_KVSERVICE_SCOPE  = 0x00010000

var (
	ERROR_CODE_KVSERVICE_CREATE_CF_ERROR int64 = ERROR_CODE_KVSERVICE_SCOPE | 0x1 //创建CF失败
	ERROR_CODE_KVSERVICE_GET_CF_ERROR int64 = ERROR_CODE_KVSERVICE_SCOPE | 0x2 //获取 CF失败
	ERROR_CODE_KVSERVICE_GET_DATA_ERROR int64 = ERROR_CODE_KVSERVICE_SCOPE | 0x3 //获取数据失败
	ERROR_CODE_KVSERVICE_PUT_DATA_ERROR int64 = ERROR_CODE_KVSERVICE_SCOPE | 0x4 //Put 数据失败
	ERROR_CODE_KVSERVICE_DELETE_DATA_ERROR int64 = ERROR_CODE_KVSERVICE_SCOPE | 0x5 //删除数据失败
	ERROR_CODE_KVSERVICE_NOFIND_PARTITION int64 = ERROR_CODE_KVSERVICE_SCOPE | 0x6 // 没有找到对应的分片数据
)
