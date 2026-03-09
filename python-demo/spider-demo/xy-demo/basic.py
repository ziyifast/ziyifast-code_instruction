import json
import requests
import time
import hashlib
import re

# ==================== 主流程 ====================

# 1. 准备请求头（补充缺失的关键字段）
headers = {
    "accept": "application/json",
    "accept-language": "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6",
    "content-type": "application/x-www-form-urlencoded",
    "origin": "https://www.goofish.com",
    "referer": "https://www.goofish.com/",
    "sec-ch-ua": '"Not(A:Brand";v="8", "Chromium";v="144", "Microsoft Edge";v="144"',
    "sec-ch-ua-mobile": "?0",
    "sec-ch-ua-platform": '"macOS"',
    "sec-fetch-dest": "empty",
    "sec-fetch-mode": "cors",
    "sec-fetch-site": "same-site",
    "user-agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36 Edg/144.0.0.0",
    # cookie
    "cookie": "tracknick=xy522420987376; cna=SaRdHxGMUFMCAbaMma9w538Q; t=db398972a201001b1d0522a167ea49b0; unb=2218578999066; _m_h5_tk=f953bd8f64510f47d55f1a09bb6fadca_1773052681378; _m_h5_tk_enc=7514b5a784c65702588621e4733ca6b7"
}
# for循环翻页请求
for i in range(1, 3):
    print("====================>第", i, "页<===================")
    # 2. 准备请求参数（保持 JSON 字符串格式）
    param_json = '{"pageNumber":'+str(i)+',"keyword":"手机","fromFilter":false,"rowsPerPage":30,"sortValue":"","sortField":"","customDistance":"","gps":"","propValueStr":{},"customGps":"","searchReqFromPage":"pcSearch","extraFilterValue":"{}","userPositionJson":"{}"}'
    print(param_json)
    data = {
        'data': param_json
    }

    # 3. 生成签名所需参数
    def get_sign(param_json):
        token = "f953bd8f64510f47d55f1a09bb6fadca"  # 从 _m_h5_tk cookie 中提取（不含时间戳部分）
        app_key = "34839810"
        t = str(int(time.time() * 1000))
        sign_str = token + "&" + t + "&" + app_key + "&" + param_json
        sign = hashlib.md5(sign_str.encode('utf-8')).hexdigest()
        return app_key, t, sign

    # 4. 构造签名
    app_key, t, sign = get_sign(param_json)
    print(f"签名字符串：{sign}")

    # 5. 构造完整 URL
    url = "https://h5api.m.goofish.com/h5/mtop.taobao.idlemtopsearch.pc.search/1.0/"

    # 6. 构建查询参数
    params = {
        "jsv": "2.7.2",
        "appKey": app_key,
        "t": t,
        "sign": sign,
        "v": "1.0",
        "type": "originaljson",
        "accountSite": "xianyu",
        "dataType": "json",
        "timeout": "20000",
        "api": "mtop.taobao.idlemtopsearch.pc.search",
        "sessionOption": "AutoLoginOnly",
        "spm_cnt": "a21ybx.search.0.0",
        "spm_pre": "a21ybx.home.searchInput.0",
    }

    # 7. 发送请求
    response = requests.post(url=url, headers=headers, params=params, data=data)

    # 调试：打印实际请求信息
    print(f"请求 URL: {response.url}")
    print(f"请求体：{response.request.body}")

    # 8. 处理响应
    result_data = response.json()
    print(f"响应状态：{result_data.get('ret', 'Unknown')}")

    # 检查是否成功
    if result_data.get('ret') == ['SUCCESS::调用成功']:
        # 9. 解析结果数据
        result_list = result_data.get("data", {}).get("resultList", [])
        for result in result_list:
            item = result.get("data", {}).get("item", {}).get("main", {})
            if not item:
                continue

            ex_content = item.get("exContent", {})
            title = ex_content.get("title", "")

            # 价格从 price 数组拼接
            price_parts = ex_content.get("price", [])
            price = "".join(part.get("text", "") for part in price_parts if part.get("text"))
            target_url = item.get("targetUrl", "")
            url = "https://www.goofish.com/item?id=" + re.search(r"item\?id=(\d+)", target_url).group(1)

            print(f"标题：{title}")
            print(f"价格：{price}")
            print(f"链接：{url}")
            print("-" * 50)
    else:
        print(f"请求失败：{result_data.get('ret')}")