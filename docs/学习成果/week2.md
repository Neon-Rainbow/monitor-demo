# 学习成果

## prometheus

首先使用Docker部署prometheus

```yaml
services:
  prometheus:
    image: bitnami/prometheus:latest
    container_name: prometheus
    ports:
      - "9090:9090" # 映射主机的 9090 端口到容器
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml # 挂载 Prometheus 的配置文件
    restart: unless-stopped
```

其中prometheus的配置文件如下:
```yaml
global:
  scrape_interval: 15s # 默认抓取间隔

scrape_configs:
  - job_name: 'prometheus'
    static_configs:
      - targets: ['localhost:9090'] # 抓取 Prometheus 自身

  - job_name: 'etcd'
    static_configs:
      - targets: ['etcd:2379'] # 抓取 etcd 服务

  - job_name: 'monitor'
    static_configs:
      - targets: ['host.docker.internal:8080'] # 抓取监控服务
```

由于目前 prometheus 是运行在Docker内,而项目是运行在宿主机上,因此将`targets`设置为`'host.docker.internal:8080'`, 用于容器访问宿主机的网络接口

部署完成后访问[localhost:9090]( http://localhost:9090)

![](./assets/%E6%88%AA%E5%B1%8F2024-12-10%2015.26.48.png)

成功访问到Prometheus的UI界面

编写一个gin的简单的demo:
```go
// 注册路由
r.GET("/metrics", gin.WrapH(promhttp.Handler()))

r.GET("/ping", func(c *gin.Context) {
    c.JSON(200, gin.H{
        "message": "pong",
    })
})
```

终端输入:
```shell
curl http://localhost:8080/metrics
```

运行结果:
```text
 user@userdeMacBook-Air  ~  curl http://localhost:8080/metrics
# HELP go_gc_duration_seconds A summary of the wall-time pause (stop-the-world) duration in garbage collection cycles.
# TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{quantile="0"} 0
go_gc_duration_seconds{quantile="0.25"} 0
go_gc_duration_seconds{quantile="0.5"} 0
go_gc_duration_seconds{quantile="0.75"} 0
go_gc_duration_seconds{quantile="1"} 0
go_gc_duration_seconds_sum 0
go_gc_duration_seconds_count 0
# HELP go_gc_gogc_percent Heap size target percentage configured by the user, otherwise 100. This value is set by the GOGC environment variable, and the runtime/debug.SetGCPercent function. Sourced from /gc/gogc:percent
# TYPE go_gc_gogc_percent gauge
go_gc_gogc_percent 100
# HELP go_gc_gomemlimit_bytes Go runtime memory limit configured by the user, otherwise math.MaxInt64. This value is set by the GOMEMLIMIT environment variable, and the runtime/debug.SetMemoryLimit function. Sourced from /gc/gomemlimit:bytes
# TYPE go_gc_gomemlimit_bytes gauge
go_gc_gomemlimit_bytes 9.223372036854776e+18
# HELP go_goroutines Number of goroutines that currently exist.
# TYPE go_goroutines gauge
go_goroutines 9
# HELP go_info Information about the Go environment.
# TYPE go_info gauge
go_info{version="go1.23.4"} 1
# HELP go_memstats_alloc_bytes Number of bytes allocated in heap and currently in use. Equals to /memory/classes/heap/objects:bytes.
# TYPE go_memstats_alloc_bytes gauge
go_memstats_alloc_bytes 266304
# HELP go_memstats_alloc_bytes_total Total number of bytes allocated in heap until now, even if released already. Equals to /gc/heap/allocs:bytes.
# TYPE go_memstats_alloc_bytes_total counter
go_memstats_alloc_bytes_total 266304
# HELP go_memstats_buck_hash_sys_bytes Number of bytes used by the profiling bucket hash table. Equals to /memory/classes/profiling/buckets:bytes.
# TYPE go_memstats_buck_hash_sys_bytes gauge
go_memstats_buck_hash_sys_bytes 4504
# HELP go_memstats_frees_total Total number of heap objects frees. Equals to /gc/heap/frees:objects + /gc/heap/tiny/allocs:objects.
# TYPE go_memstats_frees_total counter
go_memstats_frees_total 0
# HELP go_memstats_gc_sys_bytes Number of bytes used for garbage collection system metadata. Equals to /memory/classes/metadata/other:bytes.
# TYPE go_memstats_gc_sys_bytes gauge
go_memstats_gc_sys_bytes 1.521256e+06
# HELP go_memstats_heap_alloc_bytes Number of heap bytes allocated and currently in use, same as go_memstats_alloc_bytes. Equals to /memory/classes/heap/objects:bytes.
# TYPE go_memstats_heap_alloc_bytes gauge
go_memstats_heap_alloc_bytes 266304
# HELP go_memstats_heap_idle_bytes Number of heap bytes waiting to be used. Equals to /memory/classes/heap/released:bytes + /memory/classes/heap/free:bytes.
# TYPE go_memstats_heap_idle_bytes gauge
go_memstats_heap_idle_bytes 1.769472e+06
# HELP go_memstats_heap_inuse_bytes Number of heap bytes that are in use. Equals to /memory/classes/heap/objects:bytes + /memory/classes/heap/unused:bytes
# TYPE go_memstats_heap_inuse_bytes gauge
go_memstats_heap_inuse_bytes 1.998848e+06
# HELP go_memstats_heap_objects Number of currently allocated objects. Equals to /gc/heap/objects:objects.
# TYPE go_memstats_heap_objects gauge
go_memstats_heap_objects 1090
# HELP go_memstats_heap_released_bytes Number of heap bytes released to OS. Equals to /memory/classes/heap/released:bytes.
# TYPE go_memstats_heap_released_bytes gauge
go_memstats_heap_released_bytes 1.769472e+06
# HELP go_memstats_heap_sys_bytes Number of heap bytes obtained from system. Equals to /memory/classes/heap/objects:bytes + /memory/classes/heap/unused:bytes + /memory/classes/heap/released:bytes + /memory/classes/heap/free:bytes.
# TYPE go_memstats_heap_sys_bytes gauge
go_memstats_heap_sys_bytes 3.76832e+06
# HELP go_memstats_last_gc_time_seconds Number of seconds since 1970 of last garbage collection.
# TYPE go_memstats_last_gc_time_seconds gauge
go_memstats_last_gc_time_seconds 0
# HELP go_memstats_mallocs_total Total number of heap objects allocated, both live and gc-ed. Semantically a counter version for go_memstats_heap_objects gauge. Equals to /gc/heap/allocs:objects + /gc/heap/tiny/allocs:objects.
# TYPE go_memstats_mallocs_total counter
go_memstats_mallocs_total 1090
# HELP go_memstats_mcache_inuse_bytes Number of bytes in use by mcache structures. Equals to /memory/classes/metadata/mcache/inuse:bytes.
# TYPE go_memstats_mcache_inuse_bytes gauge
go_memstats_mcache_inuse_bytes 9600
# HELP go_memstats_mcache_sys_bytes Number of bytes used for mcache structures obtained from system. Equals to /memory/classes/metadata/mcache/inuse:bytes + /memory/classes/metadata/mcache/free:bytes.
# TYPE go_memstats_mcache_sys_bytes gauge
go_memstats_mcache_sys_bytes 15600
# HELP go_memstats_mspan_inuse_bytes Number of bytes in use by mspan structures. Equals to /memory/classes/metadata/mspan/inuse:bytes.
# TYPE go_memstats_mspan_inuse_bytes gauge
go_memstats_mspan_inuse_bytes 54400
# HELP go_memstats_mspan_sys_bytes Number of bytes used for mspan structures obtained from system. Equals to /memory/classes/metadata/mspan/inuse:bytes + /memory/classes/metadata/mspan/free:bytes.
# TYPE go_memstats_mspan_sys_bytes gauge
go_memstats_mspan_sys_bytes 65280
# HELP go_memstats_next_gc_bytes Number of heap bytes when next garbage collection will take place. Equals to /gc/heap/goal:bytes.
# TYPE go_memstats_next_gc_bytes gauge
go_memstats_next_gc_bytes 4.194304e+06
# HELP go_memstats_other_sys_bytes Number of bytes used for other system allocations. Equals to /memory/classes/other:bytes.
# TYPE go_memstats_other_sys_bytes gauge
go_memstats_other_sys_bytes 1.165344e+06
# HELP go_memstats_stack_inuse_bytes Number of bytes obtained from system for stack allocator in non-CGO environments. Equals to /memory/classes/heap/stacks:bytes.
# TYPE go_memstats_stack_inuse_bytes gauge
go_memstats_stack_inuse_bytes 425984
# HELP go_memstats_stack_sys_bytes Number of bytes obtained from system for stack allocator. Equals to /memory/classes/heap/stacks:bytes + /memory/classes/os-stacks:bytes.
# TYPE go_memstats_stack_sys_bytes gauge
go_memstats_stack_sys_bytes 425984
# HELP go_memstats_sys_bytes Number of bytes obtained from system. Equals to /memory/classes/total:byte.
# TYPE go_memstats_sys_bytes gauge
go_memstats_sys_bytes 6.966288e+06
# HELP go_sched_gomaxprocs_threads The current runtime.GOMAXPROCS setting, or the number of operating system threads that can execute user-level Go code simultaneously. Sourced from /sched/gomaxprocs:threads
# TYPE go_sched_gomaxprocs_threads gauge
go_sched_gomaxprocs_threads 8
# HELP go_threads Number of OS threads created.
# TYPE go_threads gauge
go_threads 6
# HELP promhttp_metric_handler_requests_in_flight Current number of scrapes being served.
# TYPE promhttp_metric_handler_requests_in_flight gauge
promhttp_metric_handler_requests_in_flight 1
# HELP promhttp_metric_handler_requests_total Total number of scrapes by HTTP status code.
# TYPE promhttp_metric_handler_requests_total counter
promhttp_metric_handler_requests_total{code="200"} 0
promhttp_metric_handler_requests_total{code="500"} 0
promhttp_metric_handler_requests_total{code="503"} 0
```

可以看到应用成功暴露了 Prometheus 的 /metrics 接口,并返回了相关的监控指标

查看[http://localhost:9090/targets](http://localhost:9090/targets),

![](./assets/%E6%88%AA%E5%B1%8F2024-12-10%2018.27.06.png)

可以看到项目的状态为`Up`

 

下面在查询框输入指标进行查询:

![](./assets/%E6%88%AA%E5%B1%8F2024-12-10%2018.28.14.png)

可以看到成功查询出了指标

此时查看项目的日志:

![](./assets/%E6%88%AA%E5%B1%8F2024-12-10%2018.37.49.png)

可以看到prometheus每15s都会访问`/metrics`端口,



然后根据[gopsutil](https://github.com/shirou/gopsutil)来获取系统运行的一些数据

通过查询资料,得知CPU使用率,内存使用率,进程数分别可以根据一下三个函数查询得到:

```go 
func Percent(interval time.Duration, percpu bool) ([]float64, error)
func VirtualMemory() (*VirtualMemoryStat, error)
func Processes() ([]*Process, error)
```

goroutine可以通过`runtime`中的该函数得到:

```go
func NumGoroutine() int
```



## ETCD

使用Docker运行ETCD

Docker-compose.yml如下:

```yaml
etcd:
    image: bitnami/etcd:latest
    container_name: etcd
    environment:
      - ETCD_NAME=etcd-1
      - ETCD_INITIAL_ADVERTISE_PEER_URLS=http://etcd:2380
      - ETCD_LISTEN_PEER_URLS=http://0.0.0.0:2380
      - ETCD_ADVERTISE_CLIENT_URLS=http://etcd:2379
      - ETCD_LISTEN_CLIENT_URLS=http://0.0.0.0:2379
      - ETCD_INITIAL_CLUSTER=etcd-1=http://etcd:2380
      - ETCD_INITIAL_CLUSTER_STATE=new
      - ALLOW_NONE_AUTHENTICATION=yes
    ports:
      - "2379:2379" # 映射 etcd 客户端访问端口
      - "2380:2380" # 映射 etcd 节点间通信端口
    restart: unless-stopped
```

然后定义服务注册时的注册信息:

```json
{
    "id":"1",
    "name":"monitor",
    "host":"localhost",
    "port":8080,
    "lease_time":5,
    "routes":[
        {
            "path":"/metrics",
            "method":"GET"
        }
    ]
}
```

其中, `id`, `name`, `host`, `port`, `lease_time`这几个字段的值的配置在[config.toml](../../config.toml)的`service`中, `routes`由`Gin` 框架中的路由配置动态生成服务的路由信息

在项目中定义结构体来注册这样的信息



定义服务注册时的Key为以下结构:

```
/service/{name}/{id}
```



然后进行服务注册

根据zap日志的输出,可以看到服务注册成功

```json
{"level":"info","ts":1733883794.6131,"caller":"etcd/register.go:111","msg":"服务注册成功","key":"/service/service/service-1"}
```

然后我们查看etcd中是否有服务注册的信息:

![](./assets/%E6%88%AA%E5%B1%8F2024-12-11%2010.25.08.png)

可以看到服务注册成功



然后实现使用ETCD进行服务注册,然后项目定期发现ETCD中注册的服务,然后将信息保存到`/etc/prometheus/file_sd_configs/targets.json`,然后prometheus.yml中进行如下配置:

```yaml
global:
  scrape_interval: 15s # 全局抓取间隔

scrape_configs:
  # Prometheus 自身的监控
  - job_name: 'prometheus'
    static_configs:
      - targets: ['prometheus:9090']

  # 监控 ETCD 的状态
  - job_name: 'etcd'
    static_configs:
      - targets: ['etcd:2379']

  # 动态配置的服务发现
  - job_name: 'dynamic_services'
    file_sd_configs:
      - files:
          - '/etc/prometheus/file_sd_configs/targets.json'
        refresh_interval: 5s
```

进行动态的服务发现,每5s读取`/etc/prometheus/file_sd_configs/targets.json`

其中`/etc/prometheus/file_sd_configs/targets.json`内容如下:

![](./assets/%E6%88%AA%E5%B1%8F2024-12-11%2014.40.55.png)

由于项目使用了Docker,因此各个容器之间使用 bridge 网络,这样 Docker 会为同一网络中的容器自动分配内部 DNS 和 IP 地址,因此在monitor这个项目进行服务注册时,可以直接将host设置为`monitor`,具体的ip由docker分配



首先制作monitor项目的镜像,然后使用docker-compose来启动项目



启动后,访问[http://localhost:9090/targets](http://localhost:9090/targets)

![](./assets/%E6%88%AA%E5%B1%8F2024-12-11%2014.37.38.png)

然后查询项目的指标

![](./assets/%E6%88%AA%E5%B1%8F2024-12-11%2014.38.24.png)

![](./assets/%E6%88%AA%E5%B1%8F2024-12-11%2014.38.55.png)

可以看到各指标可以被正确查询