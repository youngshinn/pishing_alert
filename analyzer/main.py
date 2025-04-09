from fastapi import FastAPI, Query
from fastapi.middleware.cors import CORSMiddleware
from analyz import analyze_domain

app = FastAPI(
    title="WHOIS Analyzer API",
    description="분석 요청시 도메인 등록 정보를 WHOIS로 조회 후 의심 여부 판단",
    version="1.0"
)

# CORS 허용 (Go에서 요청 허용)
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_methods=["*"],
    allow_headers=["*"],
)

@app.get("/analyze")
async def analyze(domain: str = Query(..., description="도메인 이름 입력")):
    return analyze_domain(domain)
