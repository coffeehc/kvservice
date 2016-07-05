package kvservice

import (
	"fmt"

	"github.com/coffeehc/microserviceboot/base"
	"github.com/coffeehc/web"
)

var (
	PathParam_Key    = "key"
	PathParam_Prefix = "prefix"
)

var (
	GET_VALUE  = base.EndPointMeta{Path: fmt.Sprintf("/api/v1/value/{%s}", PathParam_Key), Method: web.GET, Description: "获取Key 对应的值"}
	POST_VALUE = base.EndPointMeta{Path: "/api/v1/values", Method: web.POST, Description: "新增一个值"}
	DEL_KEY    = base.EndPointMeta{Path: fmt.Sprintf("/api/v1/value/{%s}", PathParam_Key), Method: web.DELETE, Description: "删除一个值"}
	GET_VALUES = base.EndPointMeta{Path: fmt.Sprintf("/api/v1/values/{%s}", PathParam_Prefix), Method: web.GET, Description: "获取对应前缀的值"}
)

type KVInfo struct {
	Cf    string `json:"cf"`
	Key   string `json:"k"`
	Value string `json:"v"`
}
