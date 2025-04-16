# 피싱 탐지 및 경고 시스템 (Phishing Alert System)

##  프로젝트 개요

### 플랫폼 구성
- **AWS EKS (Elastic Kubernetes Service)**: 컨테이너 오케스트레이션
- **MySQL (RDS)**: 피싱 URL 및 요청 로그 저장
- **Grafana + Prometheus**: 실시간 시각화 및 모니터링
- **Slack Webhook**: 피싱 탐지 알림 전송

### 핵심 기능
1. **Chrome 확장 프로그램**
   - 사용자가 접속한 URL을 자동으로 탐지하여 API 서버로 전달

2. **API 서버 (Go)**
   - URL 파싱 및 정규화
   - 화이트리스트 확인
   - 피싱 URL 여부 판단
   - WHOIS 분석기 호출
   - DB 저장 (request_logs, whois_logs)

3. **WHOIS Analyzer (Python / FastAPI)**
   - 도메인 등록일, 만료일, 등록자 정보를 분석
   - 신규 등록이거나 프라이버시 보호 도메인은 의심 판단

4. **MySQL (RDS)**
   - request_logs: 요청된 URL, 도메인, 접속 IP, 탐지 시간 등 저장
   - phishing_urls: 사전 수집된 피싱 URL 저장
   - safe_urls: 화이트리스트 관리

5. **Grafana**
   - 실시간 탐지 요청 그래프 표시
   - 피싱 URL 요청 건수, 도메인별 통계 시각화
   - Safe vs Phishing 비율

6. **Slack 연동 알림**
   - Grafana 알림 조건 설정 (예: 피싱 URL 5건 이상 발생 시)
   - Slack Webhook을 통해 실시간 경고 메시지 전송

---

##  데이터베이스 테이블 구조

### request_logs
- `url`: 접속한 전체 URL
- `domain`: 도메인
- `is_phishing`: 탐지 여부 (1: 피싱, 0: 정상)
- `user_ip`: 요청자 IP
- `requested_at`: 탐지 시각

### whois_logs
- `domain`: 도메인
- `registrar`: 등록 기관
- `creation_date`: 등록일
- `expiration_date`: 만료일
- `is_suspicious`: 의심 여부

### safe_urls
- 화이트리스트에 등록된 안전한 도메인 목록

---

##  DevOps 구성 요소

- **Docker**: api-server, analyzer 등을 컨테이너화
- **GitHub Actions**: CI (도커 이미지 빌드 및 푸시)
- **Argo CD**: GitOps 방식으로 Kubernetes 배포 자동화
- **Helm**: Prometheus, Grafana 배포

---

##  프로젝트 흐름도
1. 사용자가 웹사이트 접속
2. 크롬 확장 프로그램이 URL 수집 → API 서버 전송
3. API 서버가 URL 분석 및 WHOIS 검사 수행
4. MySQL에 기록 저장
5. Grafana에서 실시간 대시보드 시각화
6. 특정 조건 충족 시 Slack으로 경고 전송

---

##  Grafana 대시보드 예시
-  시간대별 피싱 탐지 건수 (Line chart)
-  피싱 도메인 Top 10 (Table chart)
-  Safe vs Phishing 비율 (Pie chart)
-  요청 IP별 요청 수 (Bar chart)

---

##  프로젝트의 가치
- 실시간 피싱 탐지 및 시각화 시스템 구축
- Chrome 확장 → API → 분석기 → 저장소 → 모니터링까지 전체 파이프라인 구현
- DevSecOps 및 클라우드 인프라 실습 기반 완성형 포트폴리오
- Grafana를 활용한 사용자 행동 분석 및 이상 탐지 적용 가능

---

## 📎 추가 가능 기능
- 피싱 탐지 AI 분석기 연결 (머신러닝)
- 관리자용 웹 대시보드 구축
- 다국어 대응 UI (확장 프로그램)
- API 인증 및 접근 제어 기능

