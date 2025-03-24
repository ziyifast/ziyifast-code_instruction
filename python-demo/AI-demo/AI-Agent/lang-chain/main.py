import requests
from langchain_openai import ChatOpenAI
from langchain.agents import AgentExecutor, create_openai_tools_agent
from langchain_core.prompts import ChatPromptTemplate, MessagesPlaceholder
from langchain_core.tools import tool
from dotenv import load_dotenv, find_dotenv

GD_KEY="xxxx"
# 加载环境变量
_ = load_dotenv(find_dotenv())


# 定义打招呼工具
@tool
def say_hello(name: str) -> str:
    """当需要向某人打招呼时使用此工具。例如：向Tom打招呼。

    Args:
        name: 需要打招呼的对象名称
    Returns:
        个性化问候语
    """
    return f"👋 你好 {name}！今天过得怎么样？"


# 定义获取天气工具
@tool
def get_weather(cityname: str = "成都") -> dict:
    """获取指定城市的天气信息。例如：查询北京天气。

    Args:
        cityname: 需要查询天气的城市名称
    Returns:
        包含天气信息的字典
    """
    abcode, city_name = get_abcode(cityname)
    url = "https://restapi.amap.com/v3/weather/weatherInfo?"
    params = {"key": GD_KEY, "city": abcode, "extensions": "base"}
    try:
        response = requests.get(url=url, params=params)
        response.raise_for_status()
        if "1" == response.json()["status"]:
            return response.json()["lives"][0]
        else:
            return None
    except requests.exceptions.RequestException as e:
        print(f"Error during API request: {e}")
        return f"Error during API request: {e}"


# 获取城市编码
def get_abcode(city: str) -> tuple:
    url = "https://restapi.amap.com/v3/config/district?"
    params = {"key": GD_KEY, "keywords": city, "subdistrict": 0}
    try:
        response = requests.get(url=url, params=params)
        response.raise_for_status()
        if "1" == response.json()["status"]:
            abcode = response.json()["districts"][0]["adcode"]
            city_name = response.json()["districts"][0]["name"]
            return (abcode, city_name)
        else:
            return None
    except requests.exceptions.RequestException as e:
        print(f"Error during API request: {e}")
        return f"Error during API request: {e}"


# 初始化模型和工具
llm = ChatOpenAI(temperature=0.5)
tools = [say_hello, get_weather]

# 构建提示模板
prompt = ChatPromptTemplate.from_messages([
    ("system", """你是一个智能助手，可以使用以下工具：

    {tools}

    使用规则：
    1. 当用户说“向[名字]打招呼”时，必须使用say_hello工具
    2. 当用户说“查询[城市]天气”时，必须使用get_weather工具
    3. 保持自然友好的语气
    4. 始终使用中文
    """),
    MessagesPlaceholder("chat_history", optional=True),
    ("human", "{input}"),
    MessagesPlaceholder("agent_scratchpad")
]).partial(
    tools="\n".join([
        f"工具名称：{t.name}\n功能描述：{t.description}\n参数格式：{t.args}"
        for t in tools
    ]),
    tool_names=", ".join([t.name for t in tools])
)

# 创建Agent工作流
agent = create_openai_tools_agent(llm=llm, tools=tools, prompt=prompt)

# 创建执行器
agent_executor = AgentExecutor(
    agent=agent,
    tools=tools,
    verbose=True,  # 显示详细执行过程
    handle_parsing_errors=True,  # 自动处理解析错误
    return_intermediate_steps=True  # 返回中间步骤
)


# 格式化函数
def format_response(data, tool_name):
    """根据工具类型增强格式化"""
    if tool_name == "get_weather":
        return (
            f"🌆 {data['city']}天气\n"
            f"🌤 现象：{data['weather']}\n"
            f"🌡 温度：{data['temperature']}℃\n"
            f"💨 风力：{data['windpower']}\n"
            f"💧 湿度：{data['humidity']}%\n"
            f"🕒 更新：{data['reporttime']}"
        )
    elif tool_name == "say_hello":
        return f"✨ {data}"
    else:
        return data


# 主循环
while True:
    user_input = input("请输入你的需求（输入退出结束对话）：")
    if user_input == "退出":
        break
    try:
        result = agent_executor.invoke({"input": user_input})

        # 提取工具调用结果
        if "intermediate_steps" in result:
            for step in result["intermediate_steps"]:
                tool_call, tool_result = step
                fn_name = tool_call.tool
                print(format_response(tool_result, fn_name))
    except Exception as e:
        print(f"\n❌ 错误: {str(e)}\n")
