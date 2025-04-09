package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	//  MySQL 연결 설정
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(".env 파일 로드 실패:", err)
	}

	// 환경 변수로부터 DB 정보 로드
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	// DB 연결
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("DB 연결 실패:", err)
	}
	defer db.Close()

	//  OpenPhish 피드 요청
	resp, err := http.Get("https://openphish.com/feed.txt")
	if err != nil {
		log.Fatal("OpenPhish 피드 호출 실패:", err)
	}
	defer resp.Body.Close()

	//  한 줄씩 읽으며 URL 저장
	scanner := bufio.NewScanner(resp.Body)
	count := 0
	for scanner.Scan() {
		url := scanner.Text()
		if err := saveURL(db, url); err != nil {
			log.Println("저장 실패:", url, err)
		} else {
			count++
		}
	}
	fmt.Printf("수집 완료! 총 %d개의 피싱 URL 저장\n", count)
}

// URL 저장 함수
func saveURL(db *sql.DB, url string) error {
	_, err := db.Exec(`
        INSERT IGNORE INTO phishing_urls(url, source, country)
        VALUES (?, 'OpenPhish', 'Global')`,
		url,
	)
	return err
}
