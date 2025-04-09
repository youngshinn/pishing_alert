--  피싱 URL 저장 테이블
CREATE TABLE IF NOT EXISTS phishing_urls (
    id INT AUTO_INCREMENT PRIMARY KEY,
    url VARCHAR(512) UNIQUE,
    domain VARCHAR(255),
    source VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

--  화이트리스트 URL 테이블
CREATE TABLE IF NOT EXISTS safe_urls (
    id INT AUTO_INCREMENT PRIMARY KEY,
    url VARCHAR(255) UNIQUE,
    added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

--  WHOIS 분석 이력 로그
CREATE TABLE IF NOT EXISTS whois_logs (
    id INT AUTO_INCREMENT PRIMARY KEY,
    domain VARCHAR(255),
    registrar VARCHAR(255),
    creation_date DATETIME,
    expiration_date DATETIME,
    is_suspicious BOOLEAN,
    requested_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

--  URL 검사 API 호출 로그 (선택)
CREATE TABLE IF NOT EXISTS request_logs (
    id INT AUTO_INCREMENT PRIMARY KEY,
    url TEXT,
    domain VARCHAR(255),
    is_phishing BOOLEAN,
    user_ip VARCHAR(45),
    requested_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
