package vms

import (
	"net"
)

func (c *ControlInfra) AssignInternalAddress() (net.IP, error) {
	for _, subnet := range c.Config.VmInternalSubnets {
		ip, err := c.findUnusedInternalAddress(subnet)
		if err != nil {
			return nil, err
		}

		if ip != nil {
			return ip, nil
		}
	}

	return nil, nil
}

func (c *ControlInfra) findUnusedInternalAddress(cidr string) (net.IP, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); {
		if ip[3] != 0 && ip[3] != 255 {
			if !c.isAddressUsed(ip) {
				return ip, nil
			}
		}

		for j := len(ip) - 1; j >= 0; j-- {
			ip[j]++
			if ip[j] > 0 {
				break
			}
		}
	}

	return nil, nil
}

func (c *ControlInfra) isAddressUsed(addr net.IP) bool {
	for _, vm := range c.AliveVM {
		if vm.IP_VM == addr.String() {
			return true
		}
	}
	return false
}
