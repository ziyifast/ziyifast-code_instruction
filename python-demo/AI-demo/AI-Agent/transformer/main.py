# 这里的GD_KEY替换为自己的高德API Key即可
from tool import GD_KEY
from transformers import AutoModelForCausalLM, AutoTokenizer
import datetime
import re
import json
import requests
# 基于transformers实现简易AI Agent：查询天气、打招呼等...

# 选择大模型，这里以阿里的千问为例
model_name_or_path = "Qwen/Qwen2.5-1.5B-Instruct"
tokenizer = AutoTokenizer.from_pretrained(model_name_or_path)
# 加载配置模型参数
model = AutoModelForCausalLM.from_pretrained(
    model_name_or_path,
    torch_dtype="auto",  # 自动选择精度
    device_map="auto",  # 自动分配GPU/CPU
)


# 调用高德API：根据城市名称获取区划代码
def get_abcode(city):
    url = "https://restapi.amap.com/v3/config/district?"
    params = {
        "key": GD_KEY,
        "keywords": city,
        "subdistrict": 0,
    }
    try:
        response = requests.get(url=url, params=params)
        response.raise_for_status()  # 检查请求是否成功
        if "1" == response.json()["status"]:
            abcode = response.json()["districts"][0]["adcode"]
            city_name = response.json()["districts"][0]["name"]
            return (abcode, city_name)
        else:
            return None
    except requests.exceptions.RequestException as e:
        # 处理请求异常
        print(f"Error during API request: {e}")
        return f"Error during API request: {e}"
    pass


def hello(name: str):
    """Say Hi.

        Args:
            name: 对谁打招呼
        Returns:
            text: 打招呼的内容
        """
    return {"text": f"{name} hello"}


#  根据城市名获取天气，
def get_weather(cityname: str = "成都"):
    # 这里需要添加函数注解，否则transformers会报：Cannot generate JSON schema
    # 同时这里的函数注解会喂给大模型，大模型会根据你的需求调用不同的tools工具来完成你的需求
    """Get current weather at a location.

    Args:
        cityname:获取天气的城市, in the format "City".
    Returns:
        province: 省份名称,
        city: 市级城市名称,
        adcode: 城市的abcode,
        weather: 对于天气现象的描述,
        temperature: 实时气温，单位：摄氏度,
        winddirection: 风向描述,
        windpower:风力级别，单位：级,
        humidity: 空气湿度,
        reporttime: 数据发布的时间,
        temperature_float: 实时气温，单位：摄氏度 的float格式的字符串,
        humidity_float: 空气湿度 的float格式的字符串,
    """
    abcode, city_name = get_abcode(cityname)
    url = "https://restapi.amap.com/v3/weather/weatherInfo?"
    params = {"key": GD_KEY, "city": abcode, "extensions": "base"}
    try:
        # 发送请求
        response = requests.get(url=url, params=params)
        response.raise_for_status()  # 检查请求是否成功
        if "1" == response.json()["status"]:
            return response.json()["lives"][0]
        else:
            return None
    except requests.exceptions.RequestException as e:
        print(f"Error during API request: {e}")
        return f"Error during API request: {e}"


def get_function_by_name(name):
    if name == "get_weather":
        return get_weather
    elif name == "hello":
        return hello


# AI 定义工具库（如：获取天气、打招呼...）
tools = [get_weather, hello]


def try_parse_tool_calls(content: str):
    """Try parse the tool calls."""
    tool_calls = []
    offset = 0
    for i, m in enumerate(re.finditer(r"<tool_call>\n(.+)?\n</tool_call>", content)):
        if i == 0:
            offset = m.start()
        try:
            func = json.loads(m.group(1))
            tool_calls.append({"type": "function", "function": func})
            if isinstance(func["arguments"], str):
                func["arguments"] = json.loads(func["arguments"])
        except json.JSONDecodeError as e:
            print(f"Failed to parse tool calls: the content is {m.group(1)} and {e}")
            pass
    if tool_calls:
        if offset > 0 and content[:offset].strip():
            c = content[:offset]
        else:
            c = ""
        return {"role": "assistant", "content": c, "tool_calls": tool_calls}
    return {"role": "assistant", "content": re.sub(r"<\|im_end\|>$", "", content)}


def get_current_data():
    dt = datetime.datetime.now()
    return dt.strftime("%Y-%m-%d")


def format_response(data, fn_name):
    if fn_name == "get_weather":
        return (
            f"🌆 {data['city']}天气\n"
            f"🌤 天气现象：{data['weather']}\n"
            f"🌡 实时气温：{data['temperature']}℃\n"
            f"💨 风力等级：{data['windpower']}\n"
            f"💧 空气湿度：{data['humidity']}%\n"
            f"🕒 更新时间：{data['reporttime']}"
        )
    elif fn_name == "hello":
        return data


while True:
    input_1 = input("请输入需要查询天气的城市名或者输入“结束”来结束程序\n")
    if "结束" == input_1:
        break
    dt = datetime.datetime.now()
    formatted_date = dt.strftime("%Y-%m-%d")
    MESSAGES = [
        {
            # 系统参数设置默认大模型角色
            "role": "system",
            "content": f"You are Qwen, created by Alibaba Cloud. You are a helpful assistant.\n\nCurrent Date: {formatted_date}",
        },
        {
            # 接受用户输入的参数
            "role": "user",
            "content": f"{input_1}",
        },
    ]
    messages = MESSAGES[:]
    # 配置聊天模版以及对应的工具tools
    text = tokenizer.apply_chat_template(
        messages, tools=tools, add_generation_prompt=True, tokenize=False
    )
    # 将用户输入的字符转换为模型可识别的数据
    inputs = tokenizer(text, return_tensors="pt").to(model.device)
    # 限制最大处理token数为512，同时根据用户输入调用对应的自定义工具tools
    outputs = model.generate(**inputs, max_new_tokens=512)
    # print(f"🤖️ AI outputs: {outputs}\n")
    # 将AI模型返回的数据转换为字符
    output_text = tokenizer.batch_decode(outputs)[0][len(text):]
    print(f"🤖️ AI output_text: {output_text}\n")
    # 根据AI返回的output_text获取需要调用的工具以及对应参数
    # 例：<tool_call>
    # {"name": "get_weather", "arguments": {"cityname": "北京"}}
    # </tool_call><|im_end|>
    response = try_parse_tool_calls(output_text)
    print(f"🤖️ AI try_parse_tool_calls response: {response}\n")
    # 判断是否调用了自定义工具
    if not response.get("tool_calls", None):
        output = response.get("content", "")
        print(f"🤖️ AI: {output}\n")
        continue
    try:
        for tool_call in response.get("tool_calls", None):
            if fn_call := tool_call.get("function"):
                fn_name: str = fn_call["name"]
                fn_args: dict = fn_call["arguments"]
                # 调用API获取数据
                fn_res: str = json.dumps(
                    get_function_by_name(fn_name)(**fn_args), ensure_ascii=False
                )
                print(f"✅ 处理成功")
                print(format_response(json.loads(fn_res), fn_name))

    except Exception as e:
        print(f"❌ 查询失败：\n{e}")
        pass
