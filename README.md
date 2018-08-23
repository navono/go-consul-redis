# go-consul
一个基于`Golang`的`Consul`的示例

# dep
`Consul`和`Redi`服务均运行在一台IP为`192.168.56.101`的虚拟机中。在虚拟机中运行：
1. 启动 `Consul`
> consul agent -dev -client="0.0.0.0"
2. 启动 `Redis`
修改`redis.conf`，开放IP和增加登录密码
> sudo systemctl start redis

3. 启动本程序
> go run .\main.go -addrs '192.168.56.101:6379' -port 8888

# verify
在`Redis`中增加值：
> set test 1234

获取值：
> curl http://localhost:8888/test

# reference
[blog](https://alex.dzyoba.com/blog/go-prometheus-service/)