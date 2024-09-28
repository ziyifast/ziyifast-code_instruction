# 随着Chrome59版本推出Headless模式（无界面模式）以来，越来越多人采用Selenium+Headless Chrome模式，实现自动化测试+爬虫。
from selenium import webdriver
from selenium.webdriver.chrome.options import Options
chrome_options = Options()
chrome_options.add_argument("--headless")
chrome_options.add_argument("--disable-gpu")
path = r'C:\Program Files (x86)\Google\Chrome\Application\chrome.exe'
chrome_options.binary_location = path
browser = webdriver.Chrome(chrome_options=chrome_options)
browser.get('http://www.baidu.com')