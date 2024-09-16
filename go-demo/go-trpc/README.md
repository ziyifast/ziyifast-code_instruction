# 概念

# 步骤
## 1 编写.proto文件

## 2 生成trpc.go代码
```go
// -p执行.proto文件位置
// --nogomod 不生成go.mod文件
// --rpconly 只生成rpc代码
// --mock=false 不生成mock代码
trpc create -p login.proto --rpconly --nogomod --mock=false
```

## 3 编写服务端配置文件
trpc_go.yaml
```yaml
server:
  service:
#       服务名称
    - name: trpc.login
#      监听地址
      ip: 127.0.0.1
#      服务监听端口
      port: 8000
```

## 4 编写服务端代码

## 5 编写客户端代码

