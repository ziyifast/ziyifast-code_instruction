# selenium + phantomjs 爬虫教程
> 注意📢：本教程只用于教学，使用爬虫过程中需要遵守相关法律法规，否则后果自负！！！

## 1 selenium：web应用程序测试工具，模拟用户操作浏览器

### 介绍
> Selenium是一个用于Web应用程序测试的工具。Selenium 测试直接运行在浏览器中，就像真正的用户在操作一样。

### 安装环境
1. 安装Google驱动
> 安装Google驱动（打开谷歌浏览器，设置-关于-查看对应Google版本），然后进入下面网址下载
https://googlechromelabs.github.io/chrome-for-testing/#stable

2. 安装selenium
```pycon
pip install selenium
```




### 实战
> 代码参考：selenium/01_selenium_demo.py

> 如果运行项目出现告警：NotOpenSSLWarning: urllib3 v2 only supports OpenSSL 1.1.1+, currently the 'ssl' module is compiled with 'LibreSSL 2.8.3'
> 解决：pip install urllib3==1.26.15


## 2 selenium + phantomjs
### 介绍

### 环境安装
> 下载地址：https://phantomjs.org/download.html
```pycon
# 需要注意最新版的selenium不支持phantomjs
# 如果要使用phantomjs，需要安装之前版本2.48.0
pip uninstall selenium 
pip install selenium==2.48.0
```

 
 
### 实战
> 代码参考：selenium/02_phantomjs_demo.py


## 3 chrome headless 模式：用于替代selenium+phantomjs无页面爬虫
### 概念
>随着Chrome59版本推出Headless模式（无界面模式）以来，越来越多人采用Selenium+Headless Chrome模式
> selenium+headless VS selenium+phantomjs
> - Headless Chrome加载速度比PhantomJS快55% 
> - Headless Chrome消耗内存比PhantomJS少38%
> 数据来源：https://hackernoon.com/benchmark-headless-chrome-vs-phantomjs-e7f44c6956c

### 环境配置
> Chrome
- Unix\Linux 系统需要 chrome >= 59 
- Windows 系统需要 chrome >= 60 Python3.6
             Selenium==3.4.*
             ChromeDriver==2.31


### 实战