const API_SERVER = "http://a702a918641e04b9dbadcc03123a1fa6-465558769.ap-northeast-2.elb.amazonaws.com";

chrome.tabs.onUpdated.addListener((tabId, changeInfo, tab) => {
  if (changeInfo.status === 'complete' && tab.url.startsWith("http")) {
    fetch(`${API_SERVER}/api/check-url?url=${encodeURIComponent(tab.url)}`)
      .then(res => res.json())
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
