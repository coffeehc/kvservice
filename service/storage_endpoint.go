package main

import (
	"baseservices/kvservice"
	"fmt"
	"github.com/coffeehc/microserviceboot/base"
	"github.com/coffeehc/microserviceboot/serviceboot"
	"github.com/coffeehc/web"
	"net/http"
)

func (this *KVService) get_value(request *http.Request, pathFragments map[string]string, reply web.Reply) {
	defer serviceboot.ErrorRecover(reply)
	key, ok := pathFragments[kvservice.PathParam_Key]
	if !ok {
		panic(base.BuildBizErr("没有指定 key 值"))
	}
	cf := request.FormValue("cf")
	v, err := this.storageService.Get(cf,[]byte(key))
	serviceboot.PanicErr(err)
	if v == nil{
		panic(base.BuildBizErr("no value",400,404))
	}
	reply.With(&kvservice.KVInfo{
		Cf:cf,
		Value: fmt.Sprintf("%s", v),
		Key:   key,
	})
}

func (this *KVService) post_value(request *http.Request, pathFragments map[string]string, reply web.Reply) {
	defer serviceboot.ErrorRecover(reply)
	kvinfo := new(kvservice.KVInfo)
	serviceboot.UnmarshalWhitJson(request, kvinfo)
	err := this.storageService.Put(kvinfo.Cf,[]byte(kvinfo.Key), []byte(kvinfo.Value))
	serviceboot.PanicErr(err)
	reply.With(kvinfo)
}

func (this *KVService) del_key(request *http.Request, pathFragments map[string]string, reply web.Reply) {
	defer serviceboot.ErrorRecover(reply)
	key, ok := pathFragments[kvservice.PathParam_Key]
	if !ok {
		panic(base.BuildBizErr("没有指定 key 值"))
	}
	cf := request.FormValue("cf")
	err := this.storageService.Del(cf,[]byte(key))
	serviceboot.PanicErr(err)
}
