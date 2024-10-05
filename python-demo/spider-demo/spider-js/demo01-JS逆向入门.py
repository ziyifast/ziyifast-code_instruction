import requests
import json
# 执行命令安装：pip install pyExecJs
import execjs


headers = {
    "accept": "application/json",
    "accept-language": "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6",
    "content-type": "application/json",
    "origin": "https://www.xiniudata.com",
    "priority": "u=1, i",
    "referer": "https://www.xiniudata.com/industry/newest?from=data",
    "sec-ch-ua": "\"Chromium\";v=\"128\", \"Not;A=Brand\";v=\"24\", \"Microsoft Edge\";v=\"128\"",
    "sec-ch-ua-mobile": "?0",
    "sec-ch-ua-platform": "\"macOS\"",
    "sec-fetch-dest": "empty",
    "sec-fetch-mode": "cors",
    "sec-fetch-site": "same-origin",
    "user-agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36 Edg/128.0.0.0"
}
cookies = {
    "btoken": "ZOH2JO8CQ3CUE1YH9PC3K6TG1Z3HCEAE",
    "hy_data_2020_id": "192383e583c10e1-06189fe412f809-7e433c49-1484784-192383e583d26f6",
    "hy_data_2020_js_sdk": "%7B%22distinct_id%22%3A%22192383e583c10e1-06189fe412f809-7e433c49-1484784-192383e583d26f6%22%2C%22site_id%22%3A211%2C%22user_company%22%3A105%2C%22props%22%3A%7B%7D%2C%22device_id%22%3A%22192383e583c10e1-06189fe412f809-7e433c49-1484784-192383e583d26f6%22%7D",
    "Hm_lvt_42317524c1662a500d12d3784dbea0f8": "1727520463,1727571559",
    "HMACCOUNT": "CA33EC3F1CCD3E5F",
    "utoken": "Z1MNZUNI27IG3HYGMABLAIRRP67D1DF4",
    "username": "Clay",
    "Hm_lpvt_42317524c1662a500d12d3784dbea0f8": "1728100324"
}
url = "https://www.xiniudata.com/api2/service/x_service/person_industry_list/list_industries_by_sort"
data = {
    "payload": "LBc3V0I6ZGB5bXsxTCQnPRBuDgQVcDhbICcmb2x3AjI=",
    "sig": "45B0ECB73CAE7AEA531F5A39B29023A0",
    "v": 1
}
data = json.dumps(data, separators=(',', ':'))
response = requests.post(url, headers=headers, cookies=cookies, data=data)

# 解析json响应，获取加密后数据
# print(response.json()["d"])

# 调用我们编写好的逆向JS函数，解密数据
decode_data = execjs.compile(open('./demo01.js', 'r', encoding='utf-8').read()).call('decode_data', response.json()["d"])
print(decode_data)