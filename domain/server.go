package domain

// ServiceRegistration 基础的服务注册信息
type ServiceRegistration struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Host string `json:"host"`
	Port int    `json:"port"`
}
