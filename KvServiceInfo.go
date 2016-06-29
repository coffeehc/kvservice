package kvservice

type KVServiceInfo struct {
}

//获取 Api 定义的内容
func (this *KVServiceInfo) GetApiDefine() string {
	return ""
}

//获取 Service 名称
func (this *KVServiceInfo) GetServiceName() string {
	return "kv-service"
}

//获取服务版本号
func (this *KVServiceInfo) GetVersion() string {
	return "1.0.0"
}

//获取服务描述
func (this *KVServiceInfo) GetDescriptor() string {
	return "K/V service"
}

//获取 Service tags
func (this *KVServiceInfo) GetServiceTags() []string {
	return nil
}
