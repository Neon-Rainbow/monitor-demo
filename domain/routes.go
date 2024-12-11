package domain

// Route 路由信息
// 包含了路径, 方法, 描述
// 用于 etcd 服务注册时的路由信息
type Route struct {
	Path   string `json:"path"`
	Method string `json:"method"`
}
