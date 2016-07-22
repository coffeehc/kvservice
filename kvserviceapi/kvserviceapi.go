package kvserviceapi

import (
	"baseservices/kvservice"
	"net/url"

	"baseservices/kvservice/modules"
	"encoding/base64"
	"github.com/coffeehc/logger"
	"github.com/coffeehc/microserviceboot/base"
	"github.com/coffeehc/microserviceboot/serviceclient"
	"strconv"
)

type KVServiceApi interface {
	Get(cf string, key []byte) (*modules.KVInfo, base.Error)
	Put(cf string, key, value []byte) base.Error
	Del(cf string, key []byte) base.Error
	GetAll(cf string, prefix, startKey []byte, order kvservice.Order, limit int) ([]*modules.KVInfo, base.Error)
}

const (
	GET_VALUE  = "GET_VALUE"
	POST_VALUE = "POST_VALUE"
	DEL_KEY    = "DEL_KEY"
	GET_VALUES = "GET_VALUES"
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

func (this *_KVServiceApi) Get(cf string, key []byte) (*modules.KVInfo, base.Error) {
	kvInfo := new(modules.KVInfo)
	queryVaules := make(url.Values)
	queryVaules.Set("cf", cf)
	err := this.serviceClient.SyncCallApiExt(GET_VALUE, map[string]string{
		kvservice.PathParam_Key: base64.RawURLEncoding.EncodeToString(key),
	}, queryVaules, nil, kvInfo)
	return kvInfo, err
}
func (this *_KVServiceApi) Put(cf string, key, value []byte) base.Error {
	kvInfo := &modules.KVInfo{
		Key:   key,
		Value: value,
		Cf:    &cf,
	}
	return this.serviceClient.SyncCallApiExt(POST_VALUE, nil, nil, serviceclient.NewRequestPBBody(kvInfo), nil)
}
func (this *_KVServiceApi) Del(cf string, key []byte) base.Error {
	queryVaules := make(url.Values)
	queryVaules.Set("cf", cf)
	return this.serviceClient.SyncCallApiExt(DEL_KEY, map[string]string{
		kvservice.PathParam_Key: base64.RawURLEncoding.EncodeToString(key),
	}, queryVaules, nil, nil)
}

func (this *_KVServiceApi) GetAll(cf string, prefix, startKey []byte, order kvservice.Order, limit int) ([]*modules.KVInfo, base.Error) {
	queryVaules := make(url.Values)
	queryVaules.Set("cf", cf)
	queryVaules.Set("start", base64.RawURLEncoding.EncodeToString(startKey))
	queryVaules.Set("order", string(order))
	queryVaules.Set("limit", strconv.Itoa(limit))
	kvInfos := &modules.KVInfos{}
	err := this.serviceClient.SyncCallApiExt(GET_VALUES, map[string]string{
		kvservice.PathParam_Prefix: base64.RawURLEncoding.EncodeToString(prefix),
	}, queryVaules, nil, kvInfos)
	return kvInfos.GetKvInfos(), err
}
