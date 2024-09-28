# phantomJs: 无页面浏览器，支持页面元素查找、js执行；由于不支持css和gui渲染，运行效率比真实浏览器快很多
# Selenium： 可以根据我们的指令，让浏览器自动加载页面，获取需要的数据，甚至页面截屏，或者判断网站上某些动作是否发生。
# Selenium 自己不带浏览器，不支持浏览器的功能，它需要与第三方浏览器结合在一起才能使用。
# 但是我们有时候需要让它内嵌在代码中运行，所以我们可以用一个叫 Phantomjs 的工具代替真实的浏览器。
from selenium import webdriver
from selenium.webdriver.common.by import By
browser = webdriver.PhantomJS("/Users/ziyi/GolandProjects/ziyifast-code_instruction/python-demo/spider-demo/selenium/phantomjs")
browser.get("http://www.baidu.com")
browser.save_screenshot("baidu.png")