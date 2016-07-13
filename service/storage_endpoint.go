package main

import (
	"baseservices/kvservice"
	"net/http"

	"baseservices/kvservice/modules"
	"encoding/base64"
	"github.com/coffeehc/microserviceboot/base"
	"github.com/coffeehc/microserviceboot/serviceboot"
	"github.com/coffeehc/web"
	"github.com/coffeehc/web/protobuf"
	"strconv"
)

func (this *KVService) get_value(request *http.Request, pathFragments map[string]string, reply web.Reply) {
	defer serviceboot.ErrorRecover(reply)
	keyStr, ok := pathFragments[kvservice.PathParam_Key]
	if !ok {
		panic(base.BuildBizErr("没有指定 key 值"))
	}
	key, err := base64.RawURLEncoding.DecodeString(keyStr)
	if err != nil {
		panic(base.BuildBizErr("无法解析Key"))
	}
	cf := request.FormValue("cf")
	v, err := this.storageService.Get(cf, key)
	serviceboot.PanicErr(err)
	if v == nil {
		panic(base.BuildBizErr("no value", 400, 404))
	}
	reply.With(&modules.KVInfo{
		Cf:    &cf,
		Value: v,
		Key:   key,
	}).As(protobuf.Transport_PB)
}

func (this *KVService) post_value(request *http.Request, pathFragments map[string]string, reply web.Reply) {
	defer serviceboot.ErrorRecover(reply)
	kvinfo := &modules.KVInfo{}
	serviceboot.UnmarshalWhitProtobuf(request, kvinfo)
	err := this.storageService.Put(kvinfo.GetCf(), kvinfo.GetKey(), kvinfo.GetValue())
	serviceboot.PanicErr(err)
	reply.With(kvinfo)
}

func (this *KVService) del_key(request *http.Request, pathFragments map[string]string, reply web.Reply) {
	defer serviceboot.ErrorRecover(reply)
	keyStr, ok := pathFragments[kvservice.PathParam_Key]
	if !ok {
		panic(base.BuildBizErr("没有指定 key 值"))
	}
	key, err := base64.RawURLEncoding.DecodeString(keyStr)
	if err != nil {
		panic(base.BuildBizErr("无法解析Key"))
	}
	cf := request.FormValue("cf")
	err = this.storageService.Del(cf, key)
	serviceboot.PanicErr(err)
}

func (this *KVService) get_vaules(request *http.Request, pathFragments map[string]string, reply web.Reply) {
	defer serviceboot.ErrorRecover(reply)
	prefix := serviceboot.ParsePathParamToBinary(pathFragments, kvservice.PathParam_Prefix)
	start, _ := base64.RawURLEncoding.DecodeString(request.FormValue("start"))
	cf := request.FormValue("cf")
	iterator := this.storageService.GetAll(cf, nil, prefix, start, request.FormValue("order"))
	limit, err := strconv.Atoi(request.FormValue("limit"))
	if err != nil {
		limit = 100
	}
	defer iterator.Close()
	kvs := make([]*modules.KVInfo, 0)
	count := 0
	if (limit < 0 || count < limit) && iterator.Next() {
		value := iterator.Value()
		kvInfo := &modules.KVInfo{
			Cf:    &cf,
			Key:   value.Key,
			Value: value.Value,
		}
		kvs = append(kvs, kvInfo)
		count++
	}
	reply.With(&modules.KVInfos{
		KvInfo: kvs,
	}).As(protobuf.Transport_PB)
}
