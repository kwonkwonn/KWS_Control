package guacamole

import (
	"database/sql"
	"encoding/hex"
	"fmt"

	"github.com/easy-cloud-Knet/KWS_Control/pkg/crypto"
	"github.com/easy-cloud-Knet/KWS_Control/util"
)

func Configure(username, uuid, ip, privateKey string, db *sql.DB) string {
	log := util.GetLogger()

	if db == nil {
		log.Error("guacamole: db connection is nil")
		return ""
	}

	// 1. 무작위 비밀번호 생성
	userPass, err := crypto.GenerateRandomPassword(12)
	if err != nil {
		log.Error("guacamole: failed to generate random password:", err, true)
		return ""
	}

	// 2. 32바이트 Salt 생성
	salt, err := crypto.GenerateRandomSalt(32)
	if err != nil {
		log.Error("guacamole: failed to create random salt:", err, true)
		return ""
	}

	// 3. 해시 계산: SHA256(salt + password)
	passwordHash := fmt.Sprintf("%x", crypto.HashPasswordWithSalt(userPass, salt))
	saltHex := hex.EncodeToString(salt)

	tx, err := db.Begin()
	if err != nil {
		log.Error("guacamole: failed to start transaction: %v", err, true)
		return ""
	}

	defer func() {
		if err != nil {
			log.DebugWarn("guacamole: rolling back transaction due to error")
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Error("guacamole: failed to rollback transaction:", rbErr, true)
			}
		}
	}()

	// 4. Entity 생성
	// 중복 확인 이후
	var entityID int64
	var res sql.Result
	err = tx.QueryRow(`SELECT entity_id FROM guacamole_entity WHERE name = ? AND type = 'USER'`, uuid).Scan(&entityID)
	if err == sql.ErrNoRows {
		res, err = tx.Exec(`INSERT INTO guacamole_entity (name, type) VALUES (?, 'USER')`, uuid)
		if err != nil {
			log.Error("guacamole: failed to create an entity:", err, true)
			return ""
		}
		entityID, err = res.LastInsertId()
		if err != nil {
			log.Error("guacamole: failed to retrieve entity id:", err, true)
			return ""
		}
		log.DebugInfo("guacamole: created new entity with ID: %d", entityID)
	} else if err != nil {
		log.Error("guacamole: failed to check existing entity:", err, true)
		return ""
	} else {
		log.DebugInfo("guacamole: using existing entity with ID: %d", entityID)
	}

	// 5. User 생성
	_, err = tx.Exec(`
		INSERT INTO guacamole_user (entity_id, password_hash, password_salt, password_date)
		VALUES (?, UNHEX(?), UNHEX(?), NOW())
		ON DUPLICATE KEY UPDATE
		password_hash = VALUES(password_hash),
		password_salt = VALUES(password_salt),
		password_date = VALUES(password_date)
	`, entityID, passwordHash, saltHex)
	if err != nil {
		log.Error("guacamole: failed to create user:", err, true)
		return ""
	}

	// 6. Connection 생성
	connectionName := fmt.Sprintf("%s-ssh", uuid)
	res, err = tx.Exec(`
		INSERT INTO guacamole_connection (connection_name, protocol)
		VALUES (?, 'ssh')
		ON DUPLICATE KEY UPDATE connection_name = VALUES(connection_name)
	`, connectionName)
	if err != nil {
		log.Error("guacamole: failed to create connection:", err, true)
		return ""
	}

	connectionID, err := res.LastInsertId()
	if err != nil {
		log.Error("guacamole: failed to retrieve connection id:", err, true)
		return ""
	}

	// 7. Connection parameters 설정
	parameters := map[string]string{
		"hostname":    ip,
		"port":        "22",
		"username":    username,
		"private-key": privateKey,
	}

	for name, value := range parameters {
		_, err = tx.Exec(`
			INSERT INTO guacamole_connection_parameter (connection_id, parameter_name, parameter_value)
			VALUES (?, ?, ?)
			ON DUPLICATE KEY UPDATE parameter_value = VALUES(parameter_value)
		`, connectionID, name, value)
		if err != nil {
			log.Error("guacamole: failed to set parameter %s: %v", name, err, true)
			return ""
		}
	}

	// 8. User에게 Connection 권한 부여
	_, err = tx.Exec(`
		INSERT INTO guacamole_connection_permission (entity_id, connection_id, permission)
		VALUES (?, ?, 'READ')
		ON DUPLICATE KEY UPDATE permission = VALUES(permission)
	`, entityID, connectionID)
	if err != nil {
		log.Error("guacamole: failed to grant permission:", err, true)
		return ""
	}

	// 9. 트랜잭션 커밋
	if err = tx.Commit(); err != nil {
		log.Error("guacamole: failed to commit transaction:", err, true)
		return ""
	}

	return userPass
}

func Cleanup(uuid string, db *sql.DB) error {
	log := util.GetLogger()

	if db == nil {
		log.Error("failed to establish database connection: db is nil")
		return fmt.Errorf("db is nil")
	}

	if err := db.Ping(); err != nil {
		log.Error("database connection test failed: %v", err, true)
		return fmt.Errorf("database connection test failed: %w", err)
	}

	tx, err := db.Begin()
	if err != nil {
		log.Error("failed to start transaction: %v", err, true)
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	defer func() {
		if err != nil {
			log.DebugWarn("rolling back transaction due to error")
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Error("failed to rollback transaction:", rbErr, true)
			}
		}
	}()

	var entityID int64
	err = tx.QueryRow(`SELECT entity_id FROM guacamole_entity WHERE name = ? AND type = 'USER'`, uuid).Scan(&entityID)
	if err == sql.ErrNoRows {
		log.DebugInfo("no entity found for UUID %s, nothing to clean up", uuid)
		return nil
	} else if err != nil {
		log.Error("failed to find entity for UUID %s: %v", uuid, err, true)
		return fmt.Errorf("failed to find entity for UUID %s: %w", uuid, err)
	}

	_, err = tx.Exec(`DELETE FROM guacamole_entity WHERE entity_id = ?`, entityID)
	if err != nil {
		log.Error("failed to delete entity for UUID %s: %v", uuid, err, true)
		return fmt.Errorf("failed to delete entity for UUID %s: %w", uuid, err)
	}

	connectionName := fmt.Sprintf("%s-ssh", uuid)
	_, err = tx.Exec(`
		DELETE FROM guacamole_connection
		WHERE connection_name = ? AND connection_id NOT IN (
			SELECT DISTINCT connection_id FROM guacamole_connection_permission
		)
	`, connectionName)
	if err != nil {
		log.Error("failed to delete orphaned connections: %v", err, true)
		return fmt.Errorf("failed to delete orphaned connections: %w", err)
	}

	if err = tx.Commit(); err != nil {
		log.Error("failed to commit transaction: %v", err, true)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Info("successfully cleaned up configuration for UUID %s", uuid, true)
	return nil
}
