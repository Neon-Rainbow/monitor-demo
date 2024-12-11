# Week3

## 使用 Docker 搭建服务
首先为项目编写 Dockerfile

```dockerfile
FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN  go build -o monitor .

CMD ["./monitor"]
```
然后制作项目的Docker镜像
然后编写docker-compose.yml, 使用`docker-compose up -d`来运行项目

项目成功运行后查看[http://127.0.0.1:8500/](http://127.0.0.1:8500/)

![](./assets/%E6%88%AA%E5%B1%8F2024-12-11%2018.05.35.png)

可以看到服务被成功注册

然后根据consul中的注册信息来提供给prometheus

![](./assets/%E6%88%AA%E5%B1%8F2024-12-11%2018.06.32.png)

prometheus可以解析出可抓取的目标实例

然后使用grafana

访问[http://localhost:3000/](http://localhost:3000/),用户名与密码分别为`admin`, `admin`

登录后添加Data Source, 在Url中填写 `http://prometheus:9090`(会由Docker自动分配ip)

然后在Dashboard中添加数据

![](./assets/%E6%88%AA%E5%B1%8F2024-12-11%2018.11.48.png)

可以成功看到指标, 并且`monitor_custom_num`随时间增加而增加,符合代码逻辑

并且可以在添加数据时指定`job`字段的值来制定通过etcd还是consul来进行服务发现
