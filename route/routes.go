package route

// Info 路由信息
// 包含了路径, 方法, 描述
// 用于 etcd 服务注册时的路由信息
type Info struct {
	Path   string `json:"path"`
	Method string `json:"method"`
}
