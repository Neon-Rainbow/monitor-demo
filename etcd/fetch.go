package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

type PrometheusTarget struct {
	Targets []string          `json:"targets"`
	Labels  map[string]string `json:"labels"`
}

// fetchServicesFromEtcd 用于从 etcd 中获取服务列表
//
// 参数:
//   - ctx: 上下文
//   - prefix: etcd 中服务列表的前缀
//
// 返回值:
//   - []*PrometheusTarget: PrometheusTarget 列表
//   - error: 错误信息
func fetchServicesFromEtcd(ctx context.Context, prefix string) ([]*PrometheusTarget, error) {
	cli := GetClient()
	resp, err := cli.Get(ctx, prefix, clientv3.WithPrefix())

	var targets []*PrometheusTarget

	if err != nil {
		return nil, fmt.Errorf("failed to get services from ETCD: %w", err)
	}
	for _, kv := range resp.Kvs {
		var s Service
		// 反序列化 ETCD 值为 Service 类型
		if err := json.Unmarshal(kv.Value, &s); err != nil {
			zap.L().Error("无法将 ETCD 值反序列化为 Service 类型", zap.Error(err))
			return nil, fmt.Errorf("无法将 ETCD 值反序列化为 Service 类型: %w", err)
		}
		targets = append(targets, s.toPrometheusTarget())
	}
	return targets, nil
}

// generatePrometheusFileSd 生成Prometheus file_sd_configs使用的目标文件 (generate the target file for Prometheus file_sd_configs)
//
// 参数:
//   - targets: PrometheusTarget 列表
//   - outputFile: 输出文件
//
// 返回值:
//   - error: 错误信息
func generatePrometheusFileSd(targets []*PrometheusTarget, outputFile string) error {
	// 序列化为 JSON
	data, err := json.MarshalIndent(targets, "", "  ")
	if err != nil {
		return fmt.Errorf("无法 target 序列化为 json: %w", err)
	}

	// 写入文件
	if err := os.WriteFile(outputFile, data, 0644); err != nil {
		return fmt.Errorf("无法写入目标文件: %w", err)
	}

	zap.L().Info("更新 Prometheus file_sd 文件成功", zap.String("file", outputFile))
	return nil
}

// AutoFetchServices 定时从ETCD中拉取服务列表并生成Prometheus file_sd文件
//
// 参数:
//   - ctx: 上下文
//
// 使用示例:
//
//	etcd.AutoFetchServices(context.Background())
func AutoFetchServices(ctx context.Context) {
	ticker := time.NewTicker(time.Second * 10)
	prefix := "/service/"
	outputFile := "/etc/prometheus/file_sd_configs/targets.json"
	go func() {
		for range ticker.C {
			targets, err := fetchServicesFromEtcd(ctx, prefix)
			if err != nil {
				zap.L().Error("从ETCD中拉去服务列表失败", zap.Error(err))
				continue
			}
			if err := generatePrometheusFileSd(targets, outputFile); err != nil {
				zap.L().Error("更新Prometheus file_sd文件失败", zap.Error(err))
			}
		}
	}()
}
