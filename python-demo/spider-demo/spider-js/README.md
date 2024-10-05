# JS逆向分析+Python爬虫结合
> `特别声明📢：本教程只用于教学，大家在使用爬虫过程中需要遵守相关法律法规，否则后果自负！！！`
## 1 概念
有时我们通过Python爬虫爬取数据会发现服务器返回的是密文数据，我们拿到密文之后，需要对它进行解密才能转换为可用数据。
- 前端对服务器返回的密文数据是通过JS来进行解密的，这时我们就需要分析并找到前端解密的代码，应用到我们的爬虫，达到爬虫抓取并解密数据的效果。
- 总结：JS逆向分析，大致就是我们 ①分析前端JS代码 + ②通过JS代码解密服务器返回的密文数据。

## 2 JS逆向分析
> 以烯牛平台数据为例，我们登录平台，访问目前国内最新赛道，服务发现服务器返回的都是加密后数据，这时我们就需要使用XHR断点，一步步进行逆向分析。
> ![在这里插入图片描述](https://i-blog.csdnimg.cn/direct/b22605e6282246bf95b83987b5d8a3e3.png)
> - 平台链接：https://www.xiniudata.com/industry/newest?from=data

### 思路分析 & 定位JS文件
 ![在这里插入图片描述](https://i-blog.csdnimg.cn/direct/eab9990fc7494dea9727a600308b307e.png)
 1. 可以看到服务器返回回来的数据在`d`字段中，这种命名太过常见，因此我们不能通过关键字定位搜索对应的js代码。
 2. 因为返回的是json数据，前端js肯定会进行解码，所以我们可以搜索`与JSON.parse()`有关的方法，找到对应解码位置。
 3. f12打开开发者工具，抓请求包，找到请求发起的js文件，点击进入
 ![在这里插入图片描述](https://i-blog.csdnimg.cn/direct/35e253c912a84e2ca8224c97a27e60cd.png)



### 格式化JS代码 & 打XHR断点
1. 点击进入js代码后，点击代码片段左下角的花括号，格式化js代码，同时根据上面的分析，搜索JSON.parse方法
![在这里插入图片描述](https://i-blog.csdnimg.cn/direct/b4b36fbd91d24e3a9b2b542880a1fe90.png)
2. 可以看到搜索结果有11个，但此处我们应该找入参由他自定义解析出的数据，不要选成了由内置的函数解析的数据，如：l = JSON.stringfy(n)，这种就该放过，不打断点。
![在这里插入图片描述](https://i-blog.csdnimg.cn/direct/2d7135a86aaf4ad79c54ad06984b050a.png)
3. y由前端自定义的方法而来，我们应该在此处打断点，因为这里很可能是解码的js部分。
![在这里插入图片描述](https://i-blog.csdnimg.cn/direct/fb5ed64b81e44555a8b3ab7aa60643eb.png)
> 继续寻找，直到整个js文件搜索完成。最后我们在第30行、第38行为其打上断点。
> ![在这里插入图片描述](https://i-blog.csdnimg.cn/direct/ef264cb517ed42129b64ab9d4f9b0b49.png)



### 刷新页面 & 移除无用断点
> 断点我们已经打好了，下面就需要我们刷新页面或者下滑，继续请求数据，以触发XHR断点。

1. 刷新页面，触发断点
![在这里插入图片描述](https://i-blog.csdnimg.cn/direct/0186e076018441b5bf7658fd3d547a9b.png)
2. 我们发现断点走在了第30行的位置，这样我们就可以把其他断点去掉
![在这里插入图片描述](https://i-blog.csdnimg.cn/direct/371c438f48dc4762aeb4955e5e030faf.png)
### 控制台 & 逆向分析JS
1. 点击下一步，让程序往下执行一步，然后观察返回的数据是否是我们所需要的明文数据：
![在这里插入图片描述](https://i-blog.csdnimg.cn/direct/25fc15c08d68421e8cd3cd5e0ef0cd0b.png)
2. 来到控制台打印m的值，观察是否是我们需要的数据
![在这里插入图片描述](https://i-blog.csdnimg.cn/direct/8698f57f78a142538c64025d82a6e04f.png)
3. 上面m发现是我们的个人信息，并非是我们所需要的数据，因此可直接跳过当前断点
![在这里插入图片描述](https://i-blog.csdnimg.cn/direct/b5d80cd8f00e4ce4917a9e7b309968b6.png)
4. 点击resume跳过断点后，发现程序再次走到JS代码的第30行。此时我们单步往下走一行，让程序计算出m的值
![在这里插入图片描述](https://i-blog.csdnimg.cn/direct/07d626ae378c47828d71349faa92894d.png)
5. 控制台输入m，回车，观察控制台返回的是否是我们所需的数据
![在这里插入图片描述](https://i-blog.csdnimg.cn/direct/e62231a7403c410b8f7563ce5f781c9c.png)
> 可以看到，成功解析出明文数据，表明断点出的这个方法就是前端的JS解密逻辑处。

### 定位前端代码 & 完善JS
1. 新建逆向JS代码，并拷贝这段JS代码

```javascript
var d = Object(u.a)(s)
 , y = Object(u.b)(d)
 , m = JSON.parse(y);
```

![在这里插入图片描述](https://i-blog.csdnimg.cn/direct/29b00b383d2b4eb8b9082f23fef69d60.png)


2. 控制台输入对应函数名，获取对应函数详细信息
> 可以看到上面的JS代码，我们并不知道Object(u.a)是什么意思，那么我们就放行，单步执行断点，待断点执行到Object(u.a)之后，就可以从控制台获取其具体值
> ![在这里插入图片描述](https://i-blog.csdnimg.cn/direct/d522238283574f3db032e926b8e669c0.png)
3. 控制台获取函数具体值：Object(u.a)，点击控制台返回结果，查看函数详情
![在这里插入图片描述](https://i-blog.csdnimg.cn/direct/476b5684a64b48399f7295fdea027036.png)
4. 拷贝函数到我们逆向的JS代码
![在这里插入图片描述](https://i-blog.csdnimg.cn/direct/774d10269ec74c9bb793d4661ba9153e.png)
我们逆向的JS代码就成了：

```javascript
// 将u.a替换为了d1
var d = Object(d1)(s)
 , y = Object(u.b)(d)
 , m = JSON.parse(y);

function d1(e) {
    var t, n, r, o, i, a, u = "", c = 0;
    for (e = e.replace(/[^A-Za-z0-9\+\/\=]/g, ""); c < e.length;)
        t = _keyStr.indexOf(e.charAt(c++)) << 2 | (o = _keyStr.indexOf(e.charAt(c++))) >> 4,
            n = (15 & o) << 4 | (i = _keyStr.indexOf(e.charAt(c++))) >> 2,
            r = (3 & i) << 6 | (a = _keyStr.indexOf(e.charAt(c++))),
            u += String.fromCharCode(t),
        64 != i && (u += String.fromCharCode(n)),
        64 != a && (u += String.fromCharCode(r));
    return u
}
```
> 重复此步骤，不断完善我们自己新建的JS代码，直到能正确解析出服务器返回的密文数据。
> - 服务器返回的密文数据获取：
> ![在这里插入图片描述](https://i-blog.csdnimg.cn/direct/465e8d2bad9948e7a0ecf886015d1adf.png)
> - 拷贝到自己的逆向JS代码，然后执行，观察是否能解析成功
> ![在这里插入图片描述](https://i-blog.csdnimg.cn/direct/b203c9af6dfa4cd3a71661e33de82e3f.png)
> - 执行JS代码，解析成功
> ![在这里插入图片描述](https://i-blog.csdnimg.cn/direct/9f12981a21974559b7c1e60c372f67aa.png)










## 3 JS+Python爬虫结合实战
###  curl自动生成Python爬虫
1. 复制请求为curl

![在这里插入图片描述](https://i-blog.csdnimg.cn/direct/64b6d5192329497cba4bd97ff4ac5099.png)

2. 将curl转换为Python代码
>网站链接：https://spidertools.cn/#/curl2Request
![在这里插入图片描述](https://i-blog.csdnimg.cn/direct/b9dc1edfe0ec479e800712506a5d102a.png)
### 调用execjs执行我们编写的逆向JS
> Python代码和js逆向代码都就绪了，剩下的就是我们通过execjs在Python代码中使用了。