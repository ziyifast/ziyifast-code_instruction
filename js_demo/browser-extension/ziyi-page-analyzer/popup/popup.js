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