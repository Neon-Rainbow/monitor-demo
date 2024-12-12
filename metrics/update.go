package metrics

import (
	"monitor/config"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"
	"go.uber.org/zap"
)

// AutoUpdateMetrics 定时更新指标
// 该函数会在一个新的goroutine中运行, 不会阻塞主程序
// 该函数会间隔一定时间更新指标, 间隔时间由配置文件中的Ticker字段决定, 该函数会更新以下指标:
//
// 1. 更新CPU使用率
// 2. 更新内存使用率
// 3. 更新goroutine数量
// 4. 更新进程数量
// 5. 更新自定义计数器
//
// 使用示例:
//
//	metrics.NewMetrics().AutoUpdateMetrics()
func (m *metrics) AutoUpdateMetrics() {
	go func() {
		ticker := time.NewTicker(time.Duration(config.Get().Server.Ticker) * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				m.update()
			}
		}
	}()
}

// update 更新指标
//
// 1. 更新CPU使用率
// 2. 更新内存使用率
// 3. 更新goroutine数量
// 4. 更新进程数量
// 5. 更新自定义计数器
func (m *metrics) update() {

	// 更新CPU使用率
	cpuPercent, err := cpu.Percent(0, false)
	if err != nil {
		zap.L().Error("get cpu percent failed", zap.Error(err))
		m.updateCpuUsage("local", 0)
	} else {
		m.updateCpuUsage("local", cpuPercent[0])
	}

	// 更新内存使用率
	memoryPercent, err := mem.VirtualMemory()
	if err != nil {
		zap.L().Error("get memory percent failed", zap.Error(err))
		m.updateMemoryUsage("local", 0)
	} else {
		m.updateMemoryUsage("local", memoryPercent.UsedPercent)
	}

	// 更新goroutine数量
	m.updateGoroutineNum("local", float64(runtime.NumGoroutine()))

	// 更新进程数量
	processes, err := process.Processes()
	if err != nil {
		zap.L().Error("get process num failed", zap.Error(err))
		m.updateProcessNum("local", 0)
	} else {
		m.updateProcessNum("local", float64(len(processes)))
	}

	// 更新自定义计数器
	m.addCustomCounter("local", 1)
}
