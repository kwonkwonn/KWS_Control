package service

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/easy-cloud-Knet/KWS_Control/structure"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func getDBConnection() string {
	log := logrus.New()
	log.SetReportCaller(true)

	if err := godotenv.Load(); err != nil {
		log.Warn(".env file not found,,, using default values")
	}

	// .env에서 가져옴 or 기본값 사용
	dbUser := getEnvOrDefault("DB_USER", "root")
	dbPassword := getEnvOrDefault("DB_PASSWORD", "password")
	dbHost := getEnvOrDefault("DB_HOST", "localhost")
	dbPort := getEnvOrDefault("DB_PORT", "3306")
	dbName := getEnvOrDefault("DB_NAME", "guacamole_db")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	fmt.Println(dsn)

	return dsn
}

func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func GuacamoleConfig(Username string, UUID string, Ip string, PrivateKey string, config structure.Config) string {
	log := logrus.New()
	log.SetReportCaller(true)

	db, err := sql.Open("mysql", getDBConnection())
	if err != nil {
		log.Error("guacamole: failed to establish the database connection:", err)
		return ""
	}
	defer db.Close()

	// 1. 무작위 비밀번호 생성
	userPass, err := generateRandomPassword(12)
	if err != nil {
		log.Error("guacamole: failed to generate random password:", err)
		return ""
	}
	fmt.Println("생성된 비밀번호:", userPass)

	// 2. 32바이트 Salt 생성
	salt, err := generateRandomSalt(32)
	if err != nil {
		log.Error("guacamole: failed to create random salt:", err)
		return ""
	}
	fmt.Println("생성된 salt:", hex.EncodeToString(salt))

	// 3. 해시 계산: SHA256(salt + password)
	passwordHash := fmt.Sprintf("%x", hashPasswordWithSalt(userPass, salt))
	saltHex := fmt.Sprintf("%x", salt)

	tx, err := db.Begin()
	if err != nil {
		log.Error("guacamole: failed to start transaction:", err)
		return ""
	}

	defer func() {
		if err != nil {
			log.Warn("guacamole: rolling back transaction due to error")
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Error("guacamole: failed to rollback transaction:", rbErr)
			}
		}
	}()

	// 4. Entity 생성
	res, err := tx.Exec(`INSERT INTO guacamole_entity (name, type) VALUES (?, 'USER')`, UUID)
	if err != nil {
		log.Error("guacamole: failed to create an entity:", err)
		return ""
	}
	entityID, err := res.LastInsertId()
	if err != nil {
		log.Error("guacamole: failed to retrieve entity id:", err)
		return ""
	}

	// 5. User 생성 (salt 포함)
	_, err = tx.Exec(`
		INSERT INTO guacamole_user (entity_id, full_name, password_hash, password_salt, password_date)
		VALUES (?, ?, UNHEX(?), UNHEX(?), ?)
	`, entityID, Username, passwordHash, saltHex, time.Now())
	if err != nil {
		log.Error("guacamole: failed to create user:", err)
		return ""
	}
	fmt.Println("User 생성 완료")

	// 6. SSH Connection 생성
	res, err = tx.Exec(`
		INSERT INTO guacamole_connection (connection_name, protocol)
		VALUES (?, 'ssh')
	`, "My SSH Connection")
	if err != nil {
		log.Error("guacamole: failed to create connection:", err)
		return ""
	}
	connID, err := res.LastInsertId()
	if err != nil {
		log.Error("guacamole: failed to retrieve connection id:", err)
		return ""
	}
	fmt.Println("Connection 생성 완료, ID:", connID)

	// 7. Connection Parameter 설정
	params := map[string]string{
		"hostname":    Ip,
		"port":        "22",
		"username":    Username,
		"private-key": PrivateKey,
	}
	for name, value := range params {
		_, err = tx.Exec(`
			INSERT INTO guacamole_connection_parameter (connection_id, parameter_name, parameter_value)
			VALUES (?, ?, ?)
		`, connID, name, value)
		if err != nil {
			log.Errorf("guacamole: failed to set parameter %s: %v", name, err)
			return ""
		}
	}
	fmt.Println("Connection Parameter 설정 완료")

	// 8. 권한 부여
	_, err = tx.Exec(`
		INSERT INTO guacamole_connection_permission (entity_id, connection_id, permission)
		VALUES (?, ?, 'READ')
	`, entityID, connID)
	if err != nil {
		log.Error("guacamole: failed to grant permission:", err)
		return ""
	}

	if err = tx.Commit(); err != nil {
		log.Error("guacamole: failed to commit transaction:", err)
		return ""
	}

	fmt.Println("Permission 부여 완료")

	return userPass
}

// SHA256 해시 함수 (salt 포함)
func hashPasswordWithSalt(password string, salt []byte) []byte {
	hash := sha256.New()

	var temp = hex.EncodeToString(salt)

	temp = strings.ToUpper(temp)
	hash.Write([]byte(password))
	hash.Write([]byte(temp))
	return hash.Sum(nil)
}

// 랜덤 salt 생성 (32바이트)
func generateRandomSalt(length int) ([]byte, error) {
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

// 안전한 랜덤 비밀번호 생성 함수
func generateRandomPassword(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes)[:length], nil
}
