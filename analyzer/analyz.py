import whois
from datetime import datetime

#  화이트리스트 도메인 목록
WHITELIST_DOMAINS = {
    "naver.com",
    "google.com",
    "kakao.com",
    "daum.net",
    "youtube.com",
    "amazon.com",
    "webflow.io"
}

def analyze_domain(domain: str) -> dict:
    # 화이트리스트 먼저 검사
    if domain in WHITELIST_DOMAINS:
        return {
            "domain": domain,
            "is_suspicious": False,
            "note": "화이트리스트 도메인입니다."
        }

    try:
        info = whois.whois(domain)

        reg_date = info.creation_date
        exp_date = info.expiration_date
        if isinstance(reg_date, list):
            reg_date = reg_date[0]
        if isinstance(exp_date, list):
            exp_date = exp_date[0]

        is_suspicious = False
        if "privacy" in str(info.get("org", "")).lower():
            is_suspicious = True
        if reg_date and reg_date > datetime.now().replace(year=datetime.now().year - 1):
            is_suspicious = True

        return {
            "domain": domain,
            "registrar": info.registrar,
            "creation_date": str(reg_date),
            "expiration_date": str(exp_date),
            "is_suspicious": is_suspicious
        }

    except Exception as e:
        return {"error": str(e), "is_suspicious": True}
