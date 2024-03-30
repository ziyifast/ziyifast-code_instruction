# Mac golang ebitengine库开发小游戏

## 1 安装依赖库
> Go version >= 1.18
> 官网地址：https://ebitengine.org/en/documents/install.html?os=darwin
```
go get -u github.com/hajimehoshi/ebiten/v2
```
> 如果发现报错：build constraints exclude all Go files in xxx
> 在命令行启用执行：CGO_ENABLED=1,启用CGO
> 执行：go run github.com/hajimehoshi/ebiten/v2/examples/rotate@latest
> 如果出现GUI页面表明环境初始成功


