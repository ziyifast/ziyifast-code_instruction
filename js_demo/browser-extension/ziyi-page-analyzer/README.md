# 从0到1开发并上线浏览器插件 — —  Edge插件开发
# 一、实战开发
> 我们将从0到1创建一个实用的"页面分析助手"插件，它可以显示当前页面的字数统计、阅读时间和主要关键词。
> - 官方插件文档链接：https://learn.microsoft.com/zh-cn/microsoft-edge/extensions-chromium/landing/

整体项目结构：
```json
ziyi-page-analyzer/
├── manifest.json          # 插件配置文件
├── popup/
│   ├── popup.html         # 弹出窗口HTML,用于展示阅读时长所需时间等
│   ├── popup.css          # 弹出窗口样式
│   └── popup.js           # 弹出窗口逻辑，插件核心逻辑，读取html页面，根据字符数预估阅读时长
├── content.js             # 页面内容脚本
└── icons/
    ├── icon16.png         # 16x16图标
    ├── icon48.png         # 48x48图标
    └── icon128.png        # 128x128图标
```
![在这里插入图片描述](https://i-blog.csdnimg.cn/direct/3cae0ceb6c2744808d1b10439d528a8c.png)

## 1.1 环境准备

1. 安装Microsoft Edge浏览器
2. 安装代码编辑器（推荐VS Code）
3. 在Edge中启用开发者模式：
- 打开Edge浏览器
- 访问 edge://extensions/
- 打开右上角的"开发人员模式"开关

## 1.2 manifest.json插件配置文件
> 主要实现对插件icon、描述、版本的配置等
```json
{
    "manifest_version": 3,
    "name": "ziyi-页面分析助手",
    "version": "1.0",
    "description": "分析当前页面的内容，提供字数统计和阅读时间估算",
    "icons": {
      "16": "icons/icon16.png",
      "48": "icons/icon48.png",
      "128": "icons/icon128.png"
    },
    "action": {
      "default_popup": "popup/popup.html",
      "default_icon": {
        "16": "icons/icon16.png",
        "48": "icons/icon48.png"
      }
    },
    "permissions": ["activeTab", "scripting"],
    "content_scripts": [
      {
        "matches": ["<all_urls>"],
        "js": ["content.js"],
        "run_at": "document_idle"
      }
    ]
  }
```
## 1.3 popup插件弹窗逻辑

### ①popup.html插件弹窗内容
> 主要实现点击插件后的弹窗页面。
```html
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <title>ziyi-页面分析助手</title>
  <link rel="stylesheet" href="popup.css">
</head>
<body>
  <div class="container">
    <h1>ziyi-页面分析助手</h1>
    
    <div class="stats-container">
      <div class="stat-card">
        <h2>字数统计</h2>
        <p id="word-count">0</p>
      </div>
      
      <div class="stat-card">
        <h2>阅读时间</h2>
        <p id="reading-time">0分钟</p>
      </div>
    </div>
    
    <div class="keywords-container">
      <h2>关键词</h2>
      <div id="keywords-list"></div>
    </div>
    
    <button id="analyze-btn">分析当前页面</button>
  </div>
  
  <script src="popup.js"></script>
</body>
</html>
```

### ②popup.css插件弹窗样式
> 主要实现插件弹窗页面的样式

```css
body {
    width: 300px;
    padding: 15px;
    font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
    background-color: #f8f9fa;
    color: #212529;
  }
  
  .container {
    display: flex;
    flex-direction: column;
    gap: 15px;
  }
  
  h1 {
    font-size: 1.2rem;
    margin: 0 0 10px 0;
    color: #0d6efd;
    text-align: center;
  }
  
  .stats-container {
    display: flex;
    justify-content: space-between;
    gap: 10px;
  }
  
  .stat-card {
    flex: 1;
    background: white;
    border-radius: 8px;
    padding: 12px;
    box-shadow: 0 2px 5px rgba(0,0,0,0.05);
    text-align: center;
  }
  
  .stat-card h2 {
    font-size: 0.9rem;
    margin: 0 0 5px 0;
    color: #6c757d;
  }
  
  .stat-card p {
    font-size: 1.5rem;
    font-weight: bold;
    margin: 0;
    color: #0d6efd;
  }
  
  .keywords-container {
    background: white;
    border-radius: 8px;
    padding: 12px;
    box-shadow: 0 2px 5px rgba(0,0,0,0.05);
  }
  
  .keywords-container h2 {
    font-size: 0.9rem;
    margin: 0 0 10px 0;
    color: #6c757d;
  }
  
  #keywords-list {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
  }
  
  .keyword {
    background-color: #e7f1ff;
    color: #0d6efd;
    padding: 4px 8px;
    border-radius: 12px;
    font-size: 0.8rem;
  }
  
  #analyze-btn {
    background-color: #0d6efd;
    color: white;
    border: none;
    border-radius: 6px;
    padding: 10px;
    font-weight: bold;
    cursor: pointer;
    transition: background-color 0.2s;
  }
  
  #analyze-btn:hover {
    background-color: #0b5ed7;
  }
```
### ③popup.js弹窗后逻辑
> 实现插件点击后，触发弹窗，执行的动作。比如统计当前页面字符数等
> - 这里是我们实现本文章插件的核心逻辑。用户点击插件，触发弹窗后，会执行js脚本，统计当前页面所有元素字符，然后预估读完的时间。

```js
document.addEventListener('DOMContentLoaded', () => {
    const analyzeBtn = document.getElementById('analyze-btn');
    const wordCountEl = document.getElementById('word-count');
    const readingTimeEl = document.getElementById('reading-time');
    const keywordsListEl = document.getElementById('keywords-list');
    
    // 从当前标签页获取分析数据
    async function getPageData() {
      const [tab] = await chrome.tabs.query({active: true, currentWindow: true});
      
      // 使用chrome.scripting.executeScript执行内容脚本
      const results = await chrome.scripting.executeScript({
        target: {tabId: tab.id},
        func: analyzePageContent
      });
      
      return results[0].result;
    }
    
    // 分析页面内容的函数（将在内容脚本中执行）
    function analyzePageContent() {
      // 获取页面主要内容
      const bodyText = document.body.innerText || "";
      
      // 计算字数
      const words = bodyText.trim().split(/\s+/).filter(word => word.length > 0);
      const wordCount = words.length;
      
      // 估算阅读时间（按200字/分钟）
      const readingTimeMinutes = Math.ceil(wordCount / 200);
      
      // 提取关键词（简化版）
      const wordFrequency = {};
      words.forEach(word => {
        const cleanWord = word.toLowerCase().replace(/[^\w]/g, '');
        if (cleanWord.length > 3) {
          wordFrequency[cleanWord] = (wordFrequency[cleanWord] || 0) + 1;
        }
      });
      
      // 获取前5个高频词
      const keywords = Object.entries(wordFrequency)
        .sort((a, b) => b[1] - a[1])
        .slice(0, 5)
        .map(entry => entry[0]);
      
      return { wordCount, readingTimeMinutes, keywords };
    }
    
    // 更新UI显示分析结果
    function updateUI(data) {
      wordCountEl.textContent = data.wordCount.toLocaleString();
      readingTimeEl.textContent = `${data.readingTimeMinutes}分钟`;
      
      keywordsListEl.innerHTML = '';
      data.keywords.forEach(keyword => {
        const keywordEl = document.createElement('div');
        keywordEl.classList.add('keyword');
        keywordEl.textContent = keyword;
        keywordsListEl.appendChild(keywordEl);
      });
    }
    
    // 点击分析按钮
    analyzeBtn.addEventListener('click', async () => {
      try {
        const pageData = await getPageData();
        updateUI(pageData);
      } catch (error) {
        console.error('分析失败:', error);
        alert('无法分析当前页面，请确保页面已完全加载');
      }
    });
    
    // 页面加载时自动分析当前页面
    getPageData().then(updateUI).catch(console.error);
  });
```

## 1.4 content.js 页面加载逻辑
> 页面加载完成时，执行的逻辑
```js
// 这个脚本会在页面加载时自动执行
console.log("ziyi-页面分析助手内容脚本已加载");

// 这里可以添加在页面上下文中直接执行的代码
// 例如：监听页面变化并通知后台脚本
```

## 1.5 icons 插件图标
> 这里推荐一个在线平台，生成icon：https://uutool.cn/chrome-icon/
> - 上传图片之后，即可生成不同大小的icon，同时可免费下载，然后将其放到icons目录下。
    ![在这里插入图片描述](https://i-blog.csdnimg.cn/direct/d7b9decf13294307b45c270c58db4268.png)

## 1.6 安装
1. 在Edge浏览器中访问 edge://extensions/
2. 打开"开发人员模式"开关
   ![在这里插入图片描述](https://i-blog.csdnimg.cn/direct/7a9d82a49a80417fb726a447fdd34e01.png)

3. 点击"加载解压缩的扩展"
4. 选择你的插件目录（包含manifest.json的目录）
5. 插件图标将出现在浏览器工具栏
   ![在这里插入图片描述](https://i-blog.csdnimg.cn/direct/6220b546f46d4424a98ab9f446409b4c.png)
## 1.7 使用
1. 访问任何网页（如新闻文章）
2. 点击工具栏中的插件图标
3. 插件将显示当前页面的：总字数、预估阅读时间、关键词分析

效果：
![在这里插入图片描述](https://i-blog.csdnimg.cn/direct/9945902190574e70862534468f19a2c2.png)




# 二、打包上线
>主要流程：
> 1. 完成开发后，访问Microsoft Edge 外接程序开发人员中心
> 2. 创建新提交
> 3. 上传ZIP格式的插件包
> 4. 填写详细信息并提交审核
## 2.1 打包插件
将插件目录打包为zip压缩包。
![在这里插入图片描述](https://i-blog.csdnimg.cn/direct/3410e761bb294eb682d131a23b5658f3.png)


## 2.2 申请开发者身份
在Edge浏览器上线插件需要Micro开发者身份。
1. 申请账号，开通开发者：https://partner.microsoft.com/zh-cn/dashboard/microsoftedge/overview
   ![在这里插入图片描述](https://i-blog.csdnimg.cn/direct/b339b27429784d1f9e68d60b67cd308b.png)

2. 在Program模块选择Microsoft Edge项目，点击Get started
   ![在这里插入图片描述](https://i-blog.csdnimg.cn/direct/5ba4f27594804fbdbef9be2faab5a9ff.png)
3. 填写申请基本信息
> 注意：申请信息的城市和身份需要填写缩写。比如城市为四川成都，省份就填写SC，城市就填写CD。
> 问题：报错Error code 2931
解决：https://learn.microsoft.com/zh-cn/answers/questions/1923071/edge-code-2931-correlation-id-7348301b-da4f-4898-b

![在这里插入图片描述](https://i-blog.csdnimg.cn/direct/67c904e0b801404ebc7bc8c8a7ce55a5.png)

4. 查看是否注册成功
   ![在这里插入图片描述](https://i-blog.csdnimg.cn/direct/5482313093b64f2688d809f0606b12e9.png)

## 2.3 填写申请上线插件
进入已加入的Edge项目，上传插件：
![在这里插入图片描述](https://i-blog.csdnimg.cn/direct/e971b96c4a614e648fce5ae350fbfa7a.png)

然后按照左侧步骤，一步步填写内容，比如：插件描述、插件logo等。
填写完成后，就可以进行发布：
![在这里插入图片描述](https://i-blog.csdnimg.cn/direct/37538a0464744b07b418a71db74f9a7c.png)

新增描述，方便审核人员测试使用，最后点击发布：
![在这里插入图片描述](https://i-blog.csdnimg.cn/direct/6d9a3c6ded264933a4fde5f57095c37f.png)
最后效果：
![在这里插入图片描述](https://i-blog.csdnimg.cn/direct/491a4bdccb494fcfbca38e40ab00351f.png)
到这里，就大功告成了，等待审核人员审核通过就可以在Edge插件市场搜索到我们的插件了。


参考文章：
https://learn.microsoft.com/zh-cn/microsoft-edge/extensions-chromium/landing/