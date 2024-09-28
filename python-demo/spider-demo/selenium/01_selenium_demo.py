from selenium import webdriver
from selenium.webdriver.common.by import By


# 设置驱动路径、创建浏览器对象
browser = webdriver.Chrome('/Users/ziyi/GolandProjects/ziyifast-code_instruction/python-demo/spider-demo/selenium/chromedriver')

browser.get("https://www.baidu.com")

import time
time.sleep(2)
# 获取文本框对象
input = browser.find_element(By.ID,"kw")
# 输入curry
input.send_keys('curry')
# 休眠2s观察效果
time.sleep(2)
# 获取百度一下按钮
search_button = browser.find_element(By.ID,'su')
# 点击按钮
search_button.click()
time.sleep(2)
# 滑动浏览器页面
js_bottom = 'document.documentElement.scrollTop=100000'
browser.execute_script(js_bottom)
time.sleep(2)
# 获取下一页按钮
next = browser.find_element(By.XPATH,'//a[@class="n"]')
# 点击下一页
next.click()
# 回到上一页面
browser.back()
time.sleep(1)
# 前进一个页面
browser.forward()
time.sleep(2)
# 退出浏览器
browser.quit()