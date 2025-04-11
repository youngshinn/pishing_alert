package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var db *sql.DB

func main() {
	// 개발 환경에서는 .env 파일 사용
	if os.Getenv("ENV") != "prod" {
		_ = godotenv.Load(".env") // 실패해도 무시
	}

	// 환경변수 로드
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// DSN (Data Source Name) 구성
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	// DB 연결
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(" DB 연결 실패:", err)
	}
	defer db.Close()

	http.HandleFunc("/api/check-url", checkURLHandler)

	log.Println("API 서버 시작: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func checkURLHandler(w http.ResponseWriter, r *http.Request) {
	inputURL := r.URL.Query().Get("url")
	if inputURL == "" {
		http.Error(w, "url 파라미터가 필요합니다", http.StatusBadRequest)
		return
	}

	// URL 파싱 → 도메인 추출
	parsed, err := url.Parse(inputURL)
	if err != nil {
		http.Error(w, "URL 파싱 실패: "+err.Error(), http.StatusBadRequest)
		return
	}
	domain := parsed.Hostname()

	// 화이트리스트 확인
	var safe bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM safe_urls WHERE url = ?)", domain).Scan(&safe)
	if err != nil {
		http.Error(w, "화이트리스트 확인 실패: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if safe {
		log.Println("[화이트리스트 허용]:", domain)
		resp := map[string]interface{}{
			"isPhishing": false,
			"whois":      "화이트리스트 도메인입니다",
		}
		writeJSON(w, resp)
		logRequest(inputURL, domain, false, getIP(r))
		return
	}

	// 피싱 URL DB 조회
	var isPhishing bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM phishing_urls WHERE url = ?)", inputURL).Scan(&isPhishing)
	if err != nil {
		http.Error(w, "DB 조회 실패: "+err.Error(), http.StatusInternalServerError)
		return
	}

	//  WHOIS 분석기 호출
	whoisResult, err := AnalyzeDomain(domain)
	if err != nil {
		log.Println("WHOIS 분석 실패:", err)
		whoisResult = map[string]interface{}{"error": err.Error()}
	} else {
		//  WHOIS 로그 저장
		_, err = db.Exec(`
			INSERT INTO whois_logs (domain, registrar, creation_date, expiration_date, is_suspicious)
			VALUES (?, ?, ?, ?, ?)`,
			whoisResult["domain"],
			whoisResult["registrar"],
			whoisResult["creation_date"],
			whoisResult["expiration_date"],
			whoisResult["is_suspicious"],
		)
		if err != nil {
			log.Println("WHOIS 로그 저장 실패:", err)
		}
	}

	// 5 요청 로그 저장
	_, err = db.Exec(
		"INSERT INTO request_logs (url, domain, is_phishing, user_ip) VALUES (?, ?, ?, ?)",
		inputURL, domain, isPhishing, r.RemoteAddr,
	)
	if err != nil {
		log.Println("로그 저장 실패:", err)
	}

	//  응답 전송
	resp := map[string]interface{}{
		"isPhishing": isPhishing,
		"whois":      whoisResult,
	}
	writeJSON(w, resp)
}

// 공통: JSON 응답 처리
func writeJSON(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// 공통: 요청 로그 저장
func logRequest(url, domain string, isPhishing bool, ip string) {
	_, err := db.Exec(`
		INSERT INTO request_logs (url, domain, is_phishing, user_ip)
		VALUES (?, ?, ?, ?)`, url, domain, isPhishing, ip)
	if err != nil {
		log.Println("요청 로그 저장 실패:", err)
	}
}

// 공통: 사용자 IP 추출
func getIP(r *http.Request) string {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr // fallback
	}
	return ip
}
