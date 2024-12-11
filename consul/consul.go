package consul

import (
	"fmt"
	"monitor/config"
	"monitor/domain"
	"sync"

	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
)

var (
	service *Service
	once    sync.Once
)

type Service struct {
	domain.ServiceRegistration
	interval                       string
	timeout                        string
	deregisterCriticalServiceAfter string
}

func newService(id, name, host string, port int, interval, timeout, deregisterCriticalServiceAfter string) *Service {
	return &Service{
		ServiceRegistration: domain.ServiceRegistration{
			ID:   id,
			Name: name,
			Host: host,
			Port: port,
		},
		interval:                       interval,
		timeout:                        timeout,
		deregisterCriticalServiceAfter: deregisterCriticalServiceAfter,
	}
}

func GetService() *Service {
	once.Do(
		func() {
			cfg := config.Get()
			service = newService(
				cfg.Service.ID,
				cfg.Service.Name,
				cfg.Service.Host,
				cfg.Service.Port,
				cfg.Service.Interval,
				cfg.Service.Timeout,
				cfg.Service.DeregisterAfter,
			)
		})
	return service
}

func (s *Service) Register() error {
	registeration := &api.AgentServiceRegistration{
		ID:      s.ID,
		Name:    s.Name,
		Address: s.Host,
		Port:    s.Port,
		Check: &api.AgentServiceCheck{
			HTTP:                           fmt.Sprintf("http://%s:%d/health", s.Host, s.Port),
			Interval:                       s.interval,                       // 健康检查间隔
			Timeout:                        s.timeout,                        // 超时时间
			DeregisterCriticalServiceAfter: s.deregisterCriticalServiceAfter, // 多久之后注销服务
		},
	}
	err := GetClient().Agent().ServiceRegister(registeration)
	if err != nil {
		zap.L().Error("注册服务失败", zap.Error(err))
		return err
	}
	return nil
}
