# 安装依赖：pip install lxml
import urllib.request
from lxml import etree

url = "http://www.baidu.com"
response = urllib.request.urlopen(url)
content = response.read().decode('utf-8')
web_html = etree.HTML(content)

# 使用XPath选择器找到id为"su"的<input>元素，并获取其"value"属性值
result = web_html.xpath('//input[@id="su"]/@value')

print(result)
