package structure

import (
	"context"
	"database/sql"
	"time"

	"github.com/easy-cloud-Knet/KWS_Control/util"
)

// MySQL VMRepository 구현체
type MySQLVMRepository struct {
	DB *sql.DB
}

func NewMySQLVMRepository(db *sql.DB) *MySQLVMRepository {
	return &MySQLVMRepository{DB: db}
}

func (r *MySQLVMRepository) AddInstance(instanceInfo *VMInfo, coreIdx int) error {
	log := util.GetLogger()
	tx, err := r.DB.Begin()
	if err != nil {
		log.Error("Failed to start transaction %v", err)
		return err
	}
	defer tx.Rollback()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = tx.ExecContext(ctx, "INSERT INTO inst_info (uuid, inst_ip, guac_pass, inst_mem, inst_vcpu, inst_disk) VALUES (?, ?, ?, ?, ?, ?)",
		string(instanceInfo.UUID),
		instanceInfo.IP_VM,
		instanceInfo.GuacPassword,
		instanceInfo.Memory,
		instanceInfo.Cpu,
		instanceInfo.Disk)
	if err != nil {
		log.Error("Failed to insert instance info: %v", err)
		return err
	}
	_, err = tx.ExecContext(ctx, "INSERT INTO inst_loc (uuid, core) VALUES (?, ?)",
		string(instanceInfo.UUID),
		coreIdx)
	if err != nil {
		log.Error("Failed to insert instance core relation: %v", err)
		return err
	}
	return tx.Commit()
}

func (r *MySQLVMRepository) UpdateInstance(instanceInfo *VMInfo) error {
	log := util.GetLogger()
	tx, err := r.DB.Begin()
	if err != nil {
		log.Error("Failed to start transaction %v", err)
		return err
	}
	defer tx.Rollback()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err = tx.ExecContext(ctx, "UPDATE inst_info SET inst_ip = ?, inst_mem = ?, inst_vcpu = ?, inst_disk = ? WHERE uuid = ?",
		instanceInfo.IP_VM,
		instanceInfo.Memory,
		instanceInfo.Cpu,
		instanceInfo.Disk,
		string(instanceInfo.UUID))
	if err != nil {
		log.Error("Failed to update instance info: %v", err)
		return err
	}
	return tx.Commit()
}

func (r *MySQLVMRepository) DeleteInstance(uuid UUID) error {
	log := util.GetLogger()
	tx, err := r.DB.Begin()
	if err != nil {
		log.Error("Failed to start transaction: %v", err)
		return err
	}
	defer tx.Rollback()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = tx.ExecContext(ctx, "DELETE FROM inst_info WHERE uuid = ?", uuid)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, "DELETE FROM inst_loc WHERE uuid = ?", uuid)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (r *MySQLVMRepository) GetInstance(uuid UUID) (*VMInfo, error) {
	log := util.GetLogger()
	tx, err := r.DB.Begin()
	if err != nil {
		log.Error("Failed to start transaction: %v", err)
		return nil, err
	}
	defer tx.Rollback()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var instance VMInfo
	err = tx.QueryRowContext(ctx, "SELECT uuid, inst_ip, guac_pass, inst_mem, inst_vcpu, inst_disk FROM inst_info WHERE uuid = ?", uuid).Scan(
		&instance.UUID,
		&instance.IP_VM,
		&instance.GuacPassword,
		&instance.Memory,
		&instance.Cpu,
		&instance.Disk)
	if err != nil {
		log.Error("Failed to get instance info: %v", err)
		return nil, err
	}
	return &instance, tx.Commit()
}

func (r *MySQLVMRepository) GetInstanceLocation(uuid UUID) (int, error) {
	log := util.GetLogger()
	tx, err := r.DB.Begin()
	if err != nil {
		log.Error("Failed to start transaction: %v", err)
		return 0, err
	}
	defer tx.Rollback()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var coreIdx int
	err = tx.QueryRowContext(ctx, "SELECT core FROM inst_loc WHERE uuid = ?", uuid).Scan(&coreIdx)
	if err != nil {
		log.Error("Failed to get instance location: %v", err)
		return 0, err
	}
	return coreIdx, tx.Commit()
}

func (r *MySQLVMRepository) GetAllInstanceInfo() ([]VMInfo, []int, error) {
	log := util.GetLogger()
	tx, err := r.DB.Begin()
	if err != nil {
		log.Error("Failed to start transaction: %v", err)
		return nil, nil, err
	}
	defer tx.Rollback()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rows *sql.Rows
	rows, err = tx.QueryContext(ctx, "SELECT info.uuid, loc.core, info.inst_ip, info.guac_pass, info.inst_vcpu, info.inst_mem, info.inst_disk FROM inst_loc loc JOIN inst_info info ON loc.uuid = info.uuid")
	if err != nil {
		log.Error("Failed to get joined instance info: %v", err)
		return nil, nil, err
	}

	var coreIdxList []int
	var VMInfoList []VMInfo

	for rows.Next() {
		var coreIdx int
		var info VMInfo

		if err := rows.Scan(&info.UUID, &coreIdx, &info.IP_VM, &info.GuacPassword, &info.Cpu, &info.Memory, &info.Disk); err != nil {
			log.Error("Failed to scan instance location: %v", err)
			return nil, nil, err
		}
		log.DebugInfo("Found instance: %s on core %d", info.UUID, coreIdx)
		VMInfoList = append(VMInfoList, info)
		coreIdxList = append(coreIdxList, coreIdx)
	}
	return VMInfoList, coreIdxList, tx.Commit()
}
