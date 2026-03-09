import json
import requests
import time
import hashlib
import re
import csv
from datetime import datetime

# ==================== 配置区域 ====================
CONFIG = {
    "token": "f953bd8f64510f47d55f1a09bb6fadca",
    "app_key": "34839810",
    "base_url": "https://h5api.m.goofish.com/h5/mtop.taobao.idlemtopsearch.pc.search/1.0/",
    "cookie": "tracknick=xy522420987376; cna=SaRdHxGMUFMCAbaMma9w538Q; t=db398972a201001b1d0522a167ea49b0; unb=2218578999066; _m_h5_tk=f953bd8f64510f47d55f1a09bb6fadca_1773052681378; _m_h5_tk_enc=7514b5a784c65702588621e4733ca6b7",
    "keyword": "手机",
    "rows_per_page": 30,
    "max_pages": 3,
    "delay_seconds": 2,
    "output_file": "xianyu_items.csv",
}

# 请求头
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
    "cookie": CONFIG["cookie"]
}

# ==================== 核心函数 ====================

def get_sign(param_json):
    """生成请求签名"""
    token = CONFIG["token"]
    app_key = CONFIG["app_key"]
    t = str(int(time.time() * 1000))
    sign_str = f"{token}&{t}&{app_key}&{param_json}"
    sign = hashlib.md5(sign_str.encode('utf-8')).hexdigest()
    return app_key, t, sign

def extract_item_id(target_url):
    """从 URL 中提取商品 ID"""
    if not target_url:
        return None
    match = re.search(r"item\?id=(\d+)", target_url)
    return match.group(1) if match else None

def search_page(page):
    """搜索单页数据"""
    param_json = json.dumps({
        "pageNumber": page,
        "keyword": CONFIG["keyword"],
        "fromFilter": False,
        "rowsPerPage": CONFIG["rows_per_page"],
        "sortValue": "",
        "sortField": "",
        "customDistance": "",
        "gps": "",
        "propValueStr": {},
        "customGps": "",
        "searchReqFromPage": "pcSearch",
        "extraFilterValue": "{}",
        "userPositionJson": "{}"
    }, separators=(',', ':'))

    app_key, t, sign = get_sign(param_json)

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

    data = {"data": param_json}

    try:
        response = requests.post(
            url=CONFIG["base_url"],
            headers=headers,
            params=params,
            data=data,
            timeout=20
        )
        return response.json()
    except Exception as e:
        print(f"请求异常：{e}")
        return None

def parse_items(result_data):
    """解析商品列表"""
    items = []
    result_list = result_data.get("data", {}).get("resultList", [])

    for result in result_list:
        item = result.get("data", {}).get("item", {}).get("main", {})
        if not item:
            continue

        ex_content = item.get("exContent", {})
        price_parts = ex_content.get("price", [])
        price = "".join(part.get("text", "") for part in price_parts if part.get("text"))
        target_url = item.get("targetUrl", "")

        # 安全提取商品 ID
        item_id = extract_item_id(target_url)
        full_url = f"https://www.goofish.com/item?id={item_id}" if item_id else target_url

        items.append({
            "title": ex_content.get("title", ""),
            "price": price,
            "url": full_url,
            "item_id": item_id or ""
        })

    return items

def save_to_csv(items, filename):
    """保存数据到 CSV 文件"""
    if not items:
        print("⚠️  没有数据可保存")
        return

    # CSV 字段
    fieldnames = ["序号", "标题", "价格", "商品 ID", "链接", "爬取时间"]

    with open(filename, 'w', newline='', encoding='utf-8-sig') as f:
        writer = csv.DictWriter(f, fieldnames=fieldnames)
        writer.writeheader()

        for i, item in enumerate(items, 1):
            writer.writerow({
                "序号": i,
                "标题": item["title"],
                "价格": item["price"],
                "商品 ID": item["item_id"],
                "链接": item["url"],
                "爬取时间": datetime.now().strftime("%Y-%m-%d %H:%M:%S")
            })

    print(f"✅ 数据已保存到：{filename}")

# ==================== 主流程 ====================

if __name__ == "__main__":
    all_items = []

    print(f"🔍 开始爬取关键词：{CONFIG['keyword']}")
    print(f"📄 最大页数：{CONFIG['max_pages']}")
    print(f"📦 每页数量：{CONFIG['rows_per_page']}")
    print(f"💾 保存文件：{CONFIG['output_file']}")
    print("=" * 60)

    # 翻页爬取
    for page in range(1, CONFIG["max_pages"] + 1):
        print(f"\n【第 {page} 页】")

        # 发送请求
        result_data = search_page(page)

        if not result_data:
            print(f"❌ 第 {page} 页请求失败，停止爬取")
            break

        # 检查响应状态
        ret = result_data.get("ret", [])
        if ret != ["SUCCESS::调用成功"]:
            print(f"❌ 第 {page} 页返回错误：{ret}")
            break

        # 解析数据
        items = parse_items(result_data)
        all_items.extend(items)
        print(f"✅ 获取到 {len(items)} 条商品")

        # 翻页延时
        if page < CONFIG["max_pages"]:
            time.sleep(CONFIG["delay_seconds"])

    # 输出汇总
    print("\n" + "=" * 60)
    print(f"📊 爬取完成！共获取 {len(all_items)} 条商品")
    print("=" * 60)

    # 保存到 CSV
    save_to_csv(all_items, CONFIG["output_file"])

    # 打印前 5 条预览
    if all_items:
        print("\n📋 数据预览（前 5 条）：")
        for i, item in enumerate(all_items[:5], 1):
            print(f"{i}. {item['title'][:30]}... | 价格：{item['price']}")