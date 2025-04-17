

#  피싱 탐지 및 경고 시스템 (Phishing Alert System)

  실시간 URL 모니터링, WHOIS 기반 분석, 시각화 및 자동 알림까지 갖춘 **DevOps 기반 보안 시스템**

---

##  프로젝트 개요

이 프로젝트는 사용자 보호를 위한 **실시간 피싱 탐지 및 경고 시스템**입니다.  
웹에서의 피싱 위협을 빠르게 탐지하고, 시각화 및 Slack을 통한 알림으로 즉각적인 대응을 가능하게 합니다.

###  목표
- Chrome 확장 프로그램을 통해 웹 브라우징 중 URL을 수집
- 백엔드 서버에서 URL을 파싱하고, 화이트리스트/WHOIS 정보를 기반으로 판단
- 실시간 시각화(Grafana) 및 Slack 알림으로 피싱 URL 경고 제공
- GitOps 기반의 자동화된 인프라 운영

---

##  기술 스택

| 분류 | 기술 |
|------|------|
| 언어 | Go, Python 3.10 |
| 프레임워크 | FastAPI |
| 인프라 | AWS EKS, RDS (MySQL) |
| DevOps | Docker, Kubernetes, Helm, GitHub Actions, Argo CD |
| 모니터링 | Prometheus, Grafana |
| 알림 | Slack Webhook |
| 확장 프로그램 | Chrome Extension (Manifest v3, JS) |

---

##  시스템 아키텍처 구성

- **Chrome Extension** → URL 수집
- **API Server (Go)** → 피싱 탐지 및 WHOIS 분석기 호출
- **WHOIS Analyzer (Python)** → 도메인 정보 분석
- **MySQL (RDS)** → 요청, 분석 결과 저장
- **Grafana + Prometheus** → 실시간 통계 시각화
- **Slack Webhook** → 이상 탐지 시 관리자 알림

---

##  핵심 기능

### 1. Chrome 확장 프로그램
- 사용자가 접속한 웹사이트의 URL을 감지하여 API 서버로 전송

![Image](https://github.com/user-attachments/assets/cb4d804c-1555-4564-a3d5-b3c9134052d0)

### 2. API 서버 (Go)
- URL 정규화 및 도메인 추출
- 화이트리스트 체크
- 피싱 여부 판단 + WHOIS 요청
- 결과 DB 저장

![Image](https://github.com/user-attachments/assets/487f7bc5-fc91-47bd-8532-3313b3260c60)

### 3. WHOIS Analyzer (Python / FastAPI)
- 도메인 등록일, 만료일, 등록자 분석
- 신규 등록/프라이버시 보호 도메인을 의심 도메인으로 판단

![Image](https://github.com/user-attachments/assets/c7ebd2b2-2d0c-4493-9940-9641a79be732)

### 4. DB 스키마 (MySQL on RDS)
- `phishing_urls`: 사전 수집 피싱 URL 목록
- `safe_urls`: 화이트리스트 URL 목록
- `whois_logs`: WHOIS 분석 결과 저장
- `request_logs`: 요청 URL, IP, 탐지 결과 저장

### 5. Grafana 대시보드 예시
- 시간대별 탐지 건수 (Line chart)
- 피싱 도메인 TOP 10 (Table)
- Safe vs Phishing 비율 (Pie)
- 요청 IP 별 요청 수 (Bar)

![Image](https://github.com/user-attachments/assets/ae278e03-3bf5-4f1c-b2cf-f59e5d48a9ae)


### 6. Slack 알림
- 피싱 URL 요청이 일정 수 이상 발생 시 실시간 경고

![Image](https://github.com/user-attachments/assets/16fc87da-c3ea-41e5-8e0d-9c7c332cd9e1)

### 7. Argo CD를 활용한 gitops 구현 
- Apps of app 방식을 활용한 gitops 방식 구현 

![Image](https://github.com/user-attachments/assets/eb8111ed-1143-4a1c-be8f-addea3f65494)

---

##  DevOps 구성요소

| 구성 요소 | 설명 |
|-----------|------|
| Docker | 모든 서비스 컨테이너화 |
| GitHub Actions | CI - 이미지 빌드 및 Docker Hub 푸시 |
| Helm | 마이크로서비스별 Helm Chart 작성 |
| Argo CD | GitOps 방식의 CD 구현 (root-app.yaml 사용) |
| Prometheus & Grafana | 모니터링 및 경고 시각화 |

---

##  시스템 흐름도

```text
[User Browser]
    ↓ (URL 수집)
[Chrome Extension]
    ↓
[API Server (Go)]
    ↓
[WHOIS Analyzer (FastAPI)] ←→ [MySQL (RDS)]
    ↓
[Grafana + Slack]
```

---

##  폴더 구조 요약

```
alert/
├── api-server/         # Go 기반 API 서버
├── analyzer/           # Python WHOIS 분석기
├── collector/          # 초기 데이터 수집기 및 DB Job
├── chrome-extension/   # URL 수집 확장 프로그램
├── helm/               # Helm Chart (App of Apps 구성)
├── k8s/                # YAML 배포파일 (개별 테스트용)
└── .github/workflows/  # GitHub Actions CI 설정
```

---

##  프로젝트 가치

- 피싱 위협에 대한 **즉각 탐지 및 시각화** 기능 제공
- GitOps 기반의 **배포 자동화 및 보안 강화**
- **DevSecOps 관점**에서의 통합 시스템 설계
- 실무에서 바로 적용 가능한 **E2E 파이프라인 구축 경험**

---

##  기여도
- 모든 백엔드, 크롤러, 분석기, 확장 프로그램 개발
- Docker/Helm/K8s 기반 인프라 구성 및 배포
- GitHub Actions, Argo CD 활용한 CI/CD 자동화 구성
- Grafana 대시보드 구성 및 Slack 연동 경고 시스템 설정

---

