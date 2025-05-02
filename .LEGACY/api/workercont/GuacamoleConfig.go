package WorkerCont

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"log"

	//"math/rand"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql" // 도커 올릴때 확인
)

func GuacamoleConfig(Username string, UUID string, Ip string, PrivateKey string) {
	// MySQL 연결 정보 수정 (내부 IP 사용)
	fmt.Println()
	db, err := sql.Open("mysql", "root:3c0GBUDtUBhyqHd@tcp(223.194.20.119:25118)/guacamole_db")
	if err != nil {
		log.Fatal("DB 연결 실패:", err)
	}
	defer db.Close()

	// 사용자 생성
	//userID := UUID
	//println(userID)
	userPass := 11111
	fmt.Println(userPass)
	//salt := "20250307"
	println("4444")
	passwordHash := fmt.Sprintf("%x", sha256Sum(strconv.Itoa(userPass)))
	fmt.Println("passwordHash:", passwordHash)
	println("5555")

	// Entity 생성
	res, err := db.Exec(`INSERT INTO guacamole_entity (name, type) VALUES (?, 'USER')`, UUID)
	if err != nil {
		log.Fatal("Entity 생성 실패:", err)
	}

	// 새로 생성된 entity_id 가져오기
	entityID, err := res.LastInsertId()
	if err != nil {
		log.Fatal("Entity ID 조회 실패:", err)
	}
	_, err = db.Exec(`
      INSERT INTO guacamole_user (entity_id, full_name, password_hash, password_date)
      VALUES (?, ?, UNHEX(?), ?)
   `, entityID, Username, passwordHash, time.Now())
	println("6666")
	if err != nil {
		log.Fatal("User 생성 실패:", err)
	}
	fmt.Println("User 생성 완료")
	println("8888")

	// SSH Connection 생성
	res, err = db.Exec(`
      INSERT INTO guacamole_connection (connection_name, protocol)
      VALUES (?, 'ssh')
   `, "My SSH Connection")
	println("9999")
	if err != nil {
		log.Fatal("Connection 생성 실패:", err)
	}
	connID, _ := res.LastInsertId()

	fmt.Println("Connection 생성 완료, ID:", connID)

	// Connection Parameter 설정
	params := map[string]string{
		"hostname":    Ip, // SSH 접속할 VM의 IP
		"port":        "22",
		"username":    Username,
		"private-key": PrivateKey,
	}

	for name, value := range params {
		_, err = db.Exec(`
         INSERT INTO guacamole_connection_parameter (connection_id, parameter_name, parameter_value)
         VALUES (?, ?, ?)
      `, connID, name, value)
		if err != nil {
			log.Fatalf("Parameter %s 설정 실패: %v", name, err)
		}
	}
	fmt.Println("Connection Parameter 설정 완료")

	// 사용자 권한 부여
	_, err = db.Exec(`
	INSERT INTO guacamole_connection_permission (entity_id, connection_id, permission)
	VALUES (?, ?, 'READ')
 `, entityID, connID)
	if err != nil {
		log.Fatal("Permission 부여 실패:", err)
	}

	fmt.Println("Permission 부여 완료")
}

func sha256Sum(input string) []byte {
	hash := sha256.New()
	hash.Write([]byte(input))
	return hash.Sum(nil)
}
