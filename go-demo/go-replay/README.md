# Goreplay使用介绍

## 1 介绍
- Goreplay 是用 Golang 写的一个 HTTP 实时流量复制工具。功能更强大，支持流量的放大、缩小，频率限制，还支持把请求记录到文件，方便回放和分析，也支持和 ElasticSearch 集成，将流量存入 ES 进行实时分析。
- GoReplay 不是代理，而是监听网络接口上的流量，不需要更改生产基础架构，而是在与服务相同的计算机上运行 GoReplay 守护程序。

Github地址：https://github.com/buger/goreplay

## 2 安装配置
进入官网下载对应操作系统版本，下载后解压即可使用
> 这里选择最新的v1.3.3为例：https://github.com/buger/goreplay/releases/download/1.3.3/gor_1.3.3_mac.tar.gz


## 3 使用
### 前置准备
本地编写http服务并启动，用于模拟线上服务器
> 勾选goland启动多实例选项，分别启动两个服务：localhost:8888、localhost:9999
> 或者将port端口暴露为参数，将main.go编译为可执行文件，然后分别运行两次

### 用法一：本地录制流量保存到文件并回放到指定服务器

1. 启动goreplay(gor)，执行命令，录制线上流量
```bash
#将端口 8888 流量保存到本地的文件
sudo ./gor --input-raw :8888 --output-file=requests.gor
```

2. 通过curl或postman等接口测试工具，访问服务
3. 暂停gor，观察录制的流量信息
```bash
1 f83d22b80000000169dd18f1 1736750112736464000 0
GET /getInfo HTTP/1.1
User-Agent: Apifox/1.0.0 (https://apifox.com)
Accept: */*
Host: localhost:8888
Accept-Encoding: gzip, deflate, br
Connection: keep-alive
```

4. 回放流量，将流量回放到另一个服务器上
```bash
./gor --input-file requests_0.gor --output-http="http://localhost:9999"
```

### 用法二：将实时流量输出到控制台
1. 执行gor命令录制流量
```bash
# 将服务一的实时流量输出到控制台
sudo ./gor --input-raw :8888 --output-stdout
```
2. 通过curl或其他工具访问服务一，并观察控制台打印信息


### 用法三：将实时流量转发到服务二
1. 执行gor命令转发流量
```bash
# 将服务一的实时流量转发到服务二
sudo ./gor --input-raw "localhost:8888" --output-http="http://localhost:9999"
```
2. 通过curl或其他工具访问服务一，并观察服务二的请求信息

### 用法四：压测（流量放大或缩小）
> goreplay支持将捕获到的生产实际请求流量减少或者放大重播以用于测试环境的压力测试.压力测试一般针对 Input 流量减少或者放大。

1. 执行gor命令进行流量缩放
```bash
# 将流量放大为200%，按照两倍速率去回放,比如：之前两个请求间的间隔为2s，我放大两倍后，请求间隔就变为了1s，相当于qps翻倍
# 如果“input-flie”是多个文件，可以用正则去匹配，如“request*.gor|200%”
sudo ./gor --input-file "requests*.gor|1%" --output-http="http://localhost:9999"

# 除了上面的百分比方式，goreplay还支持绝对值限速,代表该服务每秒只接受1个请求
sudo ./gor --input-file "requests*.gor|1" --output-http="http://localhost:9999"

# input和output两端都支持限速，有两种限速算法，百分比或者绝对值
## 百分比：input端支持缩小或者放大请求流量，基于指定的策略随机丢弃请求流量
## 绝对值：如果单位时间（秒）内达到临界值，则丢弃剩余的请求流量，下一秒临界值还原
## 示例[绝对值限速]：localhost:9999服务，每秒最多接受1个请求，超过则丢弃
sudo ./gor --input-file "requests*.gor" --output-http="http://localhost:9999|1"
sudo ./gor --input-file "requests*.gor|1" --output-http="http://localhost:9999"
## 示例[百分比限速]：
sudo ./gor --input-file "requests*.gor" --output-http="http://localhost:9999|50%"
sudo ./gor --input-file "requests*.gor|1%" --output-http="http://localhost:9999"
```
2. 通过curl或其他工具访问服务一，并观察服务二的请求信息，观察请求是否有被放大或缩小


## 其他命令行参数
–input-raw ：用来捕捉http流量，需要指定ip地址和端口
–input-file ：接收流量
–output-file：保存流量的文件
–input-tcp：将多个Goreplay实例获取的流量聚集到一个Goreplay实例
–output-stdout：终端输出
–output-tcp:将获取的流量转移至另外的Goreplay实例
–output-http:流量释放的对象server，需要指定IP地址和端口
–output-file：录制流量时指定的存储文件
–http-disallow-url :不允许正则匹配的URL
–http-allow-header :允许的Header头
–http-disallow-header:不允许的Header头
–http-allow-method:允许的请求方法，传入值为GET，POST，OPTIONS等
–input-file-loop:无限循环，而不是读完这个文件就停止了
–output-http-workers:并发请求数
–stats --out-http-stats 每5秒输出一次TPS数据（查看统计信息）
–split-output true: 按照轮训方式分割流量
–output-http-timeout 30s：http超时30秒时间设置,默认是5秒


```bash
# 查看帮助文档
./gor --help
```