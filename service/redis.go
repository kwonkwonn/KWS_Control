package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/easy-cloud-Knet/KWS_Control/client/model"
	"github.com/easy-cloud-Knet/KWS_Control/structure"
	"github.com/easy-cloud-Knet/KWS_Control/util"
	"github.com/redis/go-redis/v9"
)

// redisVMInfo is the internal structure for storing instance info in Redis.
// For Separation of Concerns, this is not exposed outside the service layer.
// service DTOs are converted to/from this format for Redis operations.
type redisVMInfo struct {
	UUID   structure.UUID `json:"uuid"`
	CPU    uint32         `json:"cpu"`
	Memory uint32         `json:"memory"` // MiB
	Disk   uint32         `json:"disk"`   // MiB
	IP     string         `json:"ip"`
	Status string         `json:"status"`
	Time   int64          `json:"time"`
}

func (r redisVMInfo) toServiceDTO() VMInfo {
	return VMInfo{
		UUID:   r.UUID,
		CPU:    r.CPU,
		Memory: r.Memory,
		Disk:   r.Disk,
		IP:     r.IP,
		Status: r.Status,
	}
}

// StoreVMInfoToRedis will serialize the given VMInfo and store it in Redis under the key of its UUID.
func StoreVMInfoToRedis(ctx context.Context, rdb *redis.Client, vmInfo VMInfo, timestamp int64) error {
	log := util.GetLogger()

	key := string(vmInfo.UUID)
	stored := redisVMInfo{
		UUID:   vmInfo.UUID,
		CPU:    vmInfo.CPU,
		Memory: vmInfo.Memory,
		Disk:   vmInfo.Disk,
		IP:     vmInfo.IP,
		Status: vmInfo.Status,
		Time:   timestamp,
	}

	jsonData, err := json.Marshal(stored)
	if err != nil {
		log.Error("failed to marshal vm for redis %v", err, true)
		return fmt.Errorf("failed to marshal vm for redis: %w", err)
	}

	if err := rdb.Set(ctx, key, string(jsonData), 0).Err(); err != nil {
		log.Error("failed to store vm in redis: UUID=%s, error=%v", key, err, true)
		return fmt.Errorf("failed to store vm in redis: %w", err)
	}

	log.Info("vm info stored in redis: UUID=%s, CPU=%d, Memory=%d, Disk=%d, IP=%s, Status=%s, Time=%d",
		key, stored.CPU, stored.Memory, stored.Disk, stored.IP, stored.Status, stored.Time, true)
	return nil
}

func RemoveVMInfoFromRedis(ctx context.Context, rdb *redis.Client, uuid structure.UUID) error {
	log := util.GetLogger()

	key := string(uuid)
	result, err := rdb.Del(ctx, key).Result()
	if err != nil {
		log.Error("failed to delete vm from redis: UUID=%s, error=%v", key, err, true)
		return fmt.Errorf("failed to delete vm from redis: %w", err)
	}

	if result == 0 {
		log.Warn("vm info not found in redis during deletion: UUID=%s", key, true)
	} else {
		log.Info("vm info removed from redis: UUID=%s", key, true)
	}
	return nil
}

// GetVMInfoFromRedis는 service DTO를 반환
func GetVMInfoFromRedis(ctx context.Context, rdb *redis.Client, uuid structure.UUID) (VMInfo, error) {
	log := util.GetLogger()

	stored, err := loadRedisVMInfo(ctx, rdb, uuid)
	if err != nil {
		return VMInfo{}, err
	}
	log.DebugInfo("vm info retrieved from redis: UUID=%s", string(uuid))
	return stored.toServiceDTO(), nil
}

// loadRedisVMInfo는 Redis에서 raw 저장 포맷을 읽음 (service 계층)
// UpdateVMStatusInRedis에서의 Time 필드 보존을 위해 필요
func loadRedisVMInfo(ctx context.Context, rdb *redis.Client, uuid structure.UUID) (redisVMInfo, error) {
	log := util.GetLogger()

	key := string(uuid)
	var stored redisVMInfo

	jsonData, err := rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			log.Warn("vm info not found in redis: UUID=%s", key, true)
			return stored, fmt.Errorf("vm info not found in redis: %s", key)
		}
		log.Error("failed to get vm from redis: UUID=%s, error=%v", key, err, true)
		return stored, fmt.Errorf("failed to get vm from redis: %w", err)
	}

	if err := json.Unmarshal([]byte(jsonData), &stored); err != nil {
		log.Error("failed to unmarshal vm from redis: UUID=%s, error=%v", key, err, true)
		return stored, fmt.Errorf("failed to unmarshal vm from redis: %w", err)
	}
	return stored, nil
}

func UpdateVMStatusInRedis(ctx context.Context, rdb *redis.Client, uuid structure.UUID, status string, timestamp int64) error {
	log := util.GetLogger()

	key := string(uuid)
	stored, err := loadRedisVMInfo(ctx, rdb, uuid)
	if err != nil {
		log.Error("failed to get existing vm info from redis: UUID=%s, error=%v", key, err, true)
		return fmt.Errorf("failed to get existing vm info from redis: %w", err)
	}

	stored.Status = status
	stored.Time = timestamp

	jsonData, err := json.Marshal(stored)
	if err != nil {
		log.Error("failed to marshal updated vm for redis: UUID=%s, error=%v", key, err, true)
		return fmt.Errorf("failed to marshal updated vm for redis: %w", err)
	}

	if err := rdb.Set(ctx, key, string(jsonData), 0).Err(); err != nil {
		log.Error("failed to update vm status in redis: UUID=%s, error=%v", key, err, true)
		return fmt.Errorf("failed to update vm status in redis: %w", err)
	}

	log.Info("vm status updated in redis: UUID=%s, Status=%s, Time=%d", key, status, timestamp, true)
	return nil
}

// VM 상태 상수는 client/model에서 재노출 (api/update_redis.go도 동일 상수 정의, 추후 통합 예정)
const (
	VMStatusPrepareBegin = model.VMStatusPrepareBegin
	VMStatusStartBegin   = model.VMStatusStartBegin
	VMStatusStarted      = model.VMStatusStarted
	VMStatusStopped      = model.VMStatusStopped
	VMStatusRelease      = model.VMStatusRelease
	VMStatusMigrate      = model.VMStatusMigrate
	VMStatusRestore      = model.VMStatusRestore
	VMStatusUnknown      = model.VMStatusUnknown
)
