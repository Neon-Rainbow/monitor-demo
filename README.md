# promethues etcd 

## 配置文件
项目的配置文件位于[./config.toml](./config.toml)中,可以根据需要修改配置文件内容

## 制作项目镜像
```shell
docker build -t monitor-service:latest .
```

然后使用`docker-compose`启动项目
```shell
docker-compose up -d
```