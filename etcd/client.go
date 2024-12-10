package etcd

import (
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"monitor/config"
	"sync"
	"time"
)

var (
	client *clientv3.Client
	once   sync.Once
)

// newClient 用于创建 etcd client
func newClient() {
	var err error
	// 创建 etcd client
	client, err = clientv3.New(clientv3.Config{
		Endpoints:   config.Get().Etcd.Endpoints,
		DialTimeout: time.Duration(config.Get().Etcd.Timeout) * time.Second,
	})
	if err != nil {
		zap.L().Fatal("创建 etcd client 失败", zap.Error(err))
		return
	}
	zap.L().Info("创建 etcd client 成功")
	return
}

// GetClient 用于获取 etcd client
// 通过 sync.Once 来保证只初始化一次
//
// 返回值:
//   - *clientv3.Client: etcd client 指针
func GetClient() *clientv3.Client {
	once.Do(func() {
		newClient()
	})
	return client
}
