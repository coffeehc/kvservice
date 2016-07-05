package kvserviceapi

import (
	"baseservices/kvservice"
	"net/url"

	"github.com/coffeehc/logger"
	"github.com/coffeehc/microserviceboot/base"
	"github.com/coffeehc/microserviceboot/serviceclient"
)

type KVServiceApi interface {
	Get(cf, key string) (string, *base.Error)
	Put(cf, key, value string) *base.Error
	Del(cf, key string) *base.Error
	GetAll(cf, prefix string)([]*kvservice.KVInfo,*base.Error)
}

const (
	GET_VALUE  = "GET_VALUE"
	POST_VALUE = "POST_VALUE"
	DEL_KEY    = "DEL_KEY"
	GET_VALUES  = "GET_VALUES"
)

func NewKVServiceApi(discoveryConfig *serviceclient.ServiceClientConsulConfig) (KVServiceApi, error) {
	serviceClient, err := serviceclient.NewServiceClient(&kvservice.KVServiceInfo{}, discoveryConfig)
	if err != nil {
		return nil, err
	}
	serviceClient.ApiRegister(GET_VALUE, kvservice.GET_VALUE)
	serviceClient.ApiRegister(POST_VALUE, kvservice.POST_VALUE)
	serviceClient.ApiRegister(DEL_KEY, kvservice.DEL_KEY)
	serviceClient.ApiRegister(GET_VALUES, kvservice.GET_VALUES)
	kvServiceApi := &_KVServiceApi{
		serviceClient: serviceClient,
	}
	logger.Info("创建 sequenceServiceApi")
	return kvServiceApi, nil
}

type _KVServiceApi struct {
	serviceClient *serviceclient.ServiceClient
}

func (this *_KVServiceApi) Get(cf, key string) (*kvservice.KVInfo, *base.Error) {
	kvInfo := new(kvservice.KVInfo)
	queryVaules := make(url.Values)
	queryVaules.Set("cf", cf)
	err := this.serviceClient.SyncCallApiExt(GET_VALUE, map[string]string{
		kvservice.PathParam_Key: url.QueryEscape(key),
	}, queryVaules, nil, kvInfo)
	return kvInfo, err
}
func (this *_KVServiceApi) Put(cf, key, value string) *base.Error {
	kvInfo := &kvservice.KVInfo{
		Key:   key,
		Value: value,
		Cf:    cf,
	}
	return this.serviceClient.SyncCallApiExt(POST_VALUE, nil, nil, serviceclient.NewRequestJsonBody(kvInfo), nil)
}
func (this *_KVServiceApi) Del(cf, key string) *base.Error {
	queryVaules := make(url.Values)
	queryVaules.Set("cf", cf)
	return this.serviceClient.SyncCallApiExt(DEL_KEY, map[string]string{
		kvservice.PathParam_Key: url.QueryEscape(key),
	}, queryVaules, nil, nil)
}

func (this *_KVServiceApi) GetAll(cf, prefix string)([]*kvservice.KVInfo,*base.Error){
	var kvInfos []*kvservice.KVInfo
	queryVaules := make(url.Values)
	queryVaules.Set("cf", cf)
	err := this.serviceClient.SyncCallApiExt(GET_VALUES, map[string]string{
		kvservice.PathParam_Prefix: url.QueryEscape(prefix),
	}, queryVaules, nil, &kvInfos)
	return kvInfos, err
}


