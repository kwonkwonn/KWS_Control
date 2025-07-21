package structure

import (
	"database/sql"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/easy-cloud-Knet/KWS_Control/util"
)

type ControlContext struct {
	Config Config
	DB     *sql.DB
	GuacDB *sql.DB
	Cores  []Core // 모든 코어를 관리
}

func (c *ControlContext) FindCoreByVmUUID(uuid UUID) *Core {
	log := util.GetLogger()

	if c.DB == nil {
		log.Error("FindCoreByVmUUID: DB connection is nil", true)
		return nil
	}

	var (
		ip      string
		portInt int
	)
	err := c.DB.QueryRow("SELECT core_ip, core_port FROM vm_info WHERE uuid = ?", uuid).Scan(&ip, &portInt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Error("VM UUID %s not found in database", uuid, true)
		} else {
			log.Error("FindCoreByVmUUID query failed for UUID %s: %v", uuid, err, true)
		}
		return nil
	}

	port := uint16(portInt)
	for i := range c.Cores {
		core := &c.Cores[i]
		if core.IP == ip && core.Port == port {
			if core.IsAlive {
				log.DebugInfo("Found core for VM UUID %s: %s:%d", uuid, core.IP, core.Port)
				return core
			}
			log.Error("Core %s:%d is not alive for VM UUID %s", core.IP, core.Port, uuid, true)
			return nil
		}
	}

	log.Error("Core %s:%d referenced by VM UUID %s not present in local state", ip, port, uuid, true)
	return nil
}

func (c *ControlContext) AssignInternalAddress() (string, error) {
	usedIPs := make(map[string]bool)

	if c.DB != nil {
		rows, err := c.DB.Query("SELECT ip_vm FROM vm_info")
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var ip string
				if err := rows.Scan(&ip); err == nil {
					usedIPs[ip] = true
				}
			}
		}
	}

	// 2. 서브넷을 순회하며 IP를 생성
	for _, cidr := range c.Config.VmInternalSubnets {
		_, ipnet, err := net.ParseCIDR(cidr)
		if err != nil {
			continue
		}

		for ip := ipnet.IP.Mask(ipnet.Mask); ipnet.Contains(ip); incrementIP(ip) {
			ipStr := ip.String()

			if ipStr == ipnet.IP.String() {
				continue
			}

			if strings.HasPrefix(ipStr, "10.5.15.") {
				lastOctet := ip[3]
				if lastOctet <= 10 {
					continue
				}
			}
			if !usedIPs[ipStr] && ipStr != ipnet.IP.String() {
				return ipStr, nil
			}
		}
	}

	return "", fmt.Errorf("No IP available for allocation")
}

func incrementIP(ip net.IP) {
	for i := len(ip) - 1; i >= 0; i-- {
		ip[i]++
		if ip[i] != 0 {
			break
		}
	}
}
