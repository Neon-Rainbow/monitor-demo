package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"monitor/config"
	"monitor/domain"
	"monitor/route"
	"sync"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

var (
	service  *Service
	etcdOnce sync.Once
)

// Service 服务
// 用于注册服务时的数据结构
// 包含了服务ID, 服务地址, 服务端口, 租约时间, 以及路由信息
type Service struct {
	domain.ServiceRegistration

	// LeaseTime 租约时间
	LeaseTime int64 `json:"lease_time"`

	// Routes 路由信息
	Routes []domain.Route `json:"routes"`
}

// newService 用于创建服务实例
// 该函数会初始化服务实例, 并将其赋值给 service 变量
// 服务实例包括服务 ID, 服务地址, 服务端口, 租约时间, 以及路由信息
// 该函数只会被调用一次
//
// 参数:
//   - id: 服务 ID
//   - host: 服务地址
//   - port: 服务端口
//   - leaseTime: 租约时间
//   - routes: 路由信息
func newService(id, name, host string, port int, leaseTime int64, routes []domain.Route) {
	service = &Service{
		ServiceRegistration: domain.ServiceRegistration{
			ID:   id,
			Name: name,
			Host: host,
			Port: port,
		},
		LeaseTime: leaseTime,
		Routes:    routes,
	}
}

// GetService 用于获取服务实例
func GetService() *Service {
	etcdOnce.Do(func() {
		// 获取服务配置
		cfg := config.Get().Service
		// 提取路由信息
		routes := route.ExtractRoutes()
		newService(cfg.ID, cfg.Name, cfg.Host, cfg.Port, cfg.LeaseTime, routes)
	})
	return service
}

// Register 注册服务到 ETCD
// 该函数会将服务信息注册到 ETCD 中
// 服务信息包括服务 ID, 服务地址, 服务端口, 租约时间, 以及路由信息
func (s *Service) Register() error {
	// 创建租约
	leaseResp, err := GetClient().Grant(context.TODO(), s.LeaseTime)
	if err != nil {
		zap.L().Error("创建租约失败", zap.Error(err))
		return err
	}

	// 服务信息序列化为 JSON
	serviceKey := fmt.Sprintf("/service/%v/%v", s.Name, s.ID)
	serviceValue, err := json.Marshal(s)
	if err != nil {
		zap.L().Error("序列化服务信息失败", zap.Error(err))
		return err
	}

	// 注册服务信息
	_, err = GetClient().Put(context.TODO(), serviceKey, string(serviceValue), clientv3.WithLease(leaseResp.ID))
	if err != nil {
		zap.L().Error("注册服务失败", zap.Error(err))
		return err
	}

	// 开始续租
	go s.keepAlive(leaseResp.ID)

	zap.L().Info("服务注册成功", zap.String("key", serviceKey))
	return nil
}

// keepAlive 续租
func (s *Service) keepAlive(leaseID clientv3.LeaseID) {
	keepAliveChan, err := GetClient().KeepAlive(context.TODO(), leaseID)
	if err != nil {
		zap.L().Error("续租失败", zap.Error(err))
		return
	}

	for ka := range keepAliveChan {
		if ka == nil {
			zap.L().Error("续租失败")
			return
		}
	}
}

func (s *Service) toPrometheusTarget() *PrometheusTarget {
	return &PrometheusTarget{
		Targets: []string{fmt.Sprintf("%s:%d", s.Host, s.Port)},
		Labels: map[string]string{
			"job": s.Name,
		},
	}
}
