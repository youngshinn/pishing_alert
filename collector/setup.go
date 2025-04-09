package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	// .env 파일 로드
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(" .env 파일 로드 실패:", err)
	}

	// 환경 변수 로딩
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// 1. DB 생성 (없는 경우에만)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/", dbUser, dbPassword, dbHost, dbPort)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("MySQL 접속 실패:", err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + dbName)
	if err != nil {
		log.Fatal("DB 생성 실패:", err)
	}
	fmt.Println("DATABASE '" + dbName + "' 생성 또는 이미 존재함")

	// 2. 테이블 생성 (schema.sql 로드)
	dsnWithDB := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	dbWithDB, err := sql.Open("mysql", dsnWithDB)
	if err != nil {
		log.Fatal("DB 선택 실패:", err)
	}
	defer dbWithDB.Close()

	content, err := os.ReadFile("database/schema.sql")
	if err != nil {
		log.Fatal("schema.sql 파일 읽기 실패:", err)
	}

	stmts := strings.Split(string(content), ";")
	for _, stmt := range stmts {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		_, err := dbWithDB.Exec(stmt)
		if err != nil {
			log.Fatalf("쿼리 실행 실패: %s\n에러: %v", stmt, err)
		}
		tableName := extractTableName(stmt)
		fmt.Printf("테이블 '%s' 생성 또는 이미 존재함\n", tableName)
	}

	fmt.Println("모든 테이블 확인 완료")
}

func extractTableName(stmt string) string {
	stmt = strings.ToLower(stmt)
	stmt = strings.ReplaceAll(stmt, "\n", " ")
	stmt = strings.ReplaceAll(stmt, "\t", " ")
	stmt = strings.Join(strings.Fields(stmt), " ") // 여러 공백을 하나로

	// CREATE TABLE IF NOT EXISTS table_name (
	if strings.Contains(stmt, "create table") {
		parts := strings.Split(stmt, " ")
		for i := 0; i < len(parts)-1; i++ {
			if parts[i] == "exists" {
				return strings.Trim(parts[i+1], "`(")
			} else if parts[i] == "table" && parts[i+1] != "if" {
				return strings.Trim(parts[i+1], "`(")
			}
		}
	}
	return "알 수 없음"
}
