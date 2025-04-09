# 베이스 이미지: Python 3.10 사용
FROM python:3.10-slim

# 작업 디렉토리 설정
WORKDIR /app

# requirements 복사 및 설치
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# 소스 복사
COPY . .

# 서버 실행 (포트 5001)
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "5001"]
