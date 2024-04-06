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

## 2 思路分析
> 1. 程序入口：ebiten.RunGame(model.NewGame())
     > 	- 定义model里的Game类，需要实现ebiten中Game这个interface
            > 		- Update() error：程序会每隔一定时间进行刷新，里面定义刷新逻辑，包括怪物的移动（调整每个怪物的x、y轴坐标），gopher（玩家）的移动，子弹的移动等
> 		- Draw(screen *Image)：通过Update调整好坐标以后，再将怪物、子弹、以及玩家调用Draw方法画到屏幕上
> 		- Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int)：设置页面的布局及大小
> 2. 定义model.Game结构体、同时实现Update、Draw等方法
     > 	- input            *Input：用于监听用户按键，比如：按下空格表示游戏开始
> - ship             *Ship：玩家角色
> - config           *config.Config：配置文件（定义玩家移动速度、怪物移动速度、游戏标题等）
> - bullets          map[*Bullet]struct{}：存储游戏中的子弹
> - monsters         map[*Monster]struct{}：存储游戏中的怪物
> - mode             Mode：标识当前游戏是待开始、已开始、已结束
> - failedCountLimit int：最多能漏掉多少个怪物
> - failedCount      int：当前已经漏掉的怪物个数
> - func1：添加init方法，包括初始化怪物的个数、玩家的位置等
> - func2：实现Update方法：用于更新怪物、玩家、子弹的位置
> - func3：实现Draw方法，重新渲染页面，实现页面动态效果
> 3. 定义GameObj结构体（子弹、怪物、用户角色都需要用到宽高、以及x、y坐标，所以可以抽取出一个Obj）
> - width  int
> - height int
> - x      int
> - y      int
> - func1：Width() int
> - func2：Height() int
> - func3：X() int
> - func4：Y() int
> 4. 定义model.Bullet
> - GameObj 包含x、y坐标（方便后续移动子弹）
> - image：子弹的样式
> - speedFactor：子弹的移动速度
> - fun1：NewBullet
> - func2：实现自己的Draw方法
> - func3：outOfScreen，判断子弹是否移出了屏幕。当子弹超出屏幕时，应当删除，不再维护。
> 5. 定义model.Monster（类比model.Bullet，此处包含怪物的样式、移动速度同时通过GameObj维护怪物坐标x、y）
> - GameObj
> - img         *ebiten.Image
> - speedFactor int
> - fun1：NewMonster
> - func2：Draw
> - func3：OutOfScreen
> 6. 定义model.Ship（类比model.Bullet，此处包含用户角色样式、移动速度）
> - GameObj
> - img         *ebiten.Image
> - speedFactor int
> - fun1：NewShip
> - func2：Draw

> `tips:`
> - 游戏胜负判定规则：
> 1. 胜利win：
     > 		- 遗漏的怪物数<=N（配置文件配置）
> 2. 失败lose：
     > 		- 飞船（用户角色碰到怪物）
     > 		- 遗漏掉太多怪物

## 3 文档地址
> https://blog.csdn.net/weixin_45565886/article/details/137175064