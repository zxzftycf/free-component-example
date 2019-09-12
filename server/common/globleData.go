package common

var Logger LogI
var Locker RedisLockI
var Router RouteI
var Pusher PushI

//组件配置名和组件配置值的映射
type OneComponentConfig map[string]string

//组件名和组件配置的映射map
type OneServerConfig map[string]OneComponentConfig

//服务器名和服务器组件map的映射
type AllServerConfig map[string]OneServerConfig

var ServerName string
var ServerConfig AllServerConfig
var AllComponentMap map[string]interface{}
var ComponentMap map[string]interface{}
var ServerIndex string
var ClientInterfaceMap map[string]bool

func init() {
	AllComponentMap = make(map[string]interface{})
	ComponentMap = make(map[string]interface{})
	ClientInterfaceMap = make(map[string]bool)

	ClientInterfaceMap["Login.Login"] = true
	ClientInterfaceMap["Register.Register"] = true
}
