package service

import (
	"database/sql"
	"fmt"

	vms "github.com/easy-cloud-Knet/KWS_Control/structure"
)

func InsertVMInfo(db *sql.DB, vm *vms.VMInfo, coreIP string, corePort uint16) error {
	if db == nil {
		return fmt.Errorf("InsertVMInfo: db is nil")
	}
	if vm == nil {
		return fmt.Errorf("InsertVMInfo: vm is nil")
	}

	const q = `INSERT INTO vm_info (uuid, core_ip, core_port, ip_vm, guac_password, memory_mib, cpu_cores, disk_mib, is_alive)
               VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
               ON DUPLICATE KEY UPDATE
                 core_ip      = VALUES(core_ip),
                 core_port    = VALUES(core_port),
                 ip_vm        = VALUES(ip_vm),
                 guac_password= VALUES(guac_password),
                 memory_mib   = VALUES(memory_mib),
                 cpu_cores    = VALUES(cpu_cores),
                 disk_mib     = VALUES(disk_mib),
                 is_alive     = VALUES(is_alive)`

	_, err := db.Exec(q,
		vm.UUID,
		coreIP,
		corePort,
		vm.IP_VM,
		vm.GuacPassword,
		vm.Memory,
		vm.Cpu,
		vm.Disk,
		1, // is_alive true on create
	)
	if err != nil {
		return fmt.Errorf("InsertVMInfo: %w", err)
	}
	return nil
}

func DeleteVMInfo(db *sql.DB, uuid vms.UUID) error {
	if db == nil {
		return fmt.Errorf("DeleteVMInfo: db is nil")
	}
	const q = `DELETE FROM vm_info WHERE uuid = ?`
	if _, err := db.Exec(q, uuid); err != nil {
		return fmt.Errorf("DeleteVMInfo: %w", err)
	}
	return nil
}

func SetVMAlive(db *sql.DB, uuid vms.UUID, alive bool) error {
	if db == nil {
		return fmt.Errorf("SetVMAlive: db is nil")
	}
	val := 0
	if alive {
		val = 1
	}
	const q = `UPDATE vm_info SET is_alive = ? WHERE uuid = ?`
	if _, err := db.Exec(q, val, uuid); err != nil {
		return fmt.Errorf("SetVMAlive: %w", err)
	}
	return nil
}

func GetVMHardware(db *sql.DB, uuid vms.UUID) (mem uint32, cpu uint32, disk uint32, err error) {
	if db == nil {
		return 0, 0, 0, fmt.Errorf("GetVMHardware: db is nil")
	}
	const q = `SELECT memory_mib, cpu_cores, disk_mib FROM vm_info WHERE uuid = ?`
	if err = db.QueryRow(q, uuid).Scan(&mem, &cpu, &disk); err != nil {
		return 0, 0, 0, err
	}
	return
}

func GetGuacPassword(db *sql.DB, uuid vms.UUID) (string, error) {
	if db == nil {
		return "", fmt.Errorf("GetGuacPassword: db is nil")
	}
	const q = `SELECT guac_password FROM vm_info WHERE uuid = ?`
	var pass string
	if err := db.QueryRow(q, uuid).Scan(&pass); err != nil {
		return "", err
	}
	return pass, nil
}
