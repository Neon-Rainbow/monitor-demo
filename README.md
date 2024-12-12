# monitor

## 配置文件
项目的配置文件位于[config.toml](./config.toml)中,可以根据需要修改配置文件内容

## 制作项目镜像
```shell
docker build -t monitor-service:latest .
```

然后使用`docker-compose`启动项目
```shell
docker-compose up -d
```

## 项目使用
项目提供了以下可供查询的数据:
+ monitor_cpu_usage
+ monitor_custom_num
+ monitor_goroutine_num
+ monitor_memory_usage
+ monitor_http_request

其中`monitor_http_request`是访问接口的次数

可以访问[http://localhost:9090/query](http://localhost:9090/query) 来通过 prometheus 查看项目数据

其中每个指标的 `job` 字段分别有`job="service"` 以及 `job="consul"` 两种

service表明prometheus通过项目在etcd中注册的服务进行的服务发现,consul表明prometheus通过项目在consul中注册的服务进行的服务发现

同时可以访问[http://localhost:3000/](http://localhost:3000/) 来通过 Grafana 提供更全面的图形化界面,初始用户名为`admin`,初始密码为`admin`

默认已经添加了 prometheus 的数据源

如果项目修改了配置文件, 那么请根据实际的 prometheus 的地址来添加数据源