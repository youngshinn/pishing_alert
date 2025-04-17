const API_SERVER = "http://ad393aac0be904e93bbda1d41a903123-251081205.ap-northeast-2.elb.amazonaws.com:8081";

chrome.tabs.onUpdated.addListener((tabId, changeInfo, tab) => {
  if (changeInfo.status === 'complete' && tab.url.startsWith("http")) {
    fetch(`${API_SERVER}/api/check-url?url=${encodeURIComponent(tab.url)}`)
      .then(res => {
        // JSON인지 확인 후 파싱
        const contentType = res.headers.get("content-type") || "";
        if (contentType.includes("application/json")) {
          return res.json();
        } else {
          return res.text().then(text => {
            throw new Error("Unexpected response format: " + text);
          });
        }
      })
      .then(data => {
        console.log("[Phishing Check]", tab.url, data);
        if (data.isPhishing || (data.whois && data.whois.is_suspicious)) {
          chrome.tabs.update(tabId, {
            url: chrome.runtime.getURL("warning.html")
          });
        }
      })
      .catch(err => console.error("API 요청 실패:", err));
  }
});
