package network

import (
	"fmt"
	"strconv"
	"strings"
)

func FindSubnet(lastSubnet string) string {
	value := make([]int, 3)
	j := 0
	for i := 0; i < 3; i++ {
		var temp string
		for lastSubnet[j] != '.' {
			temp = temp + string(lastSubnet[j])
			j++
		}
		value[i], _ = strconv.Atoi(temp)
		j++
	}

	if value[2] >= 255 {
		if value[1] >= 255 {
			if value[0] >= 255 {
				return "err"
			}
			value[0]++
			value[1] = 0
			value[2] = 0
		} else {
			value[1]++
			value[2] = 0
		}
	} else {
		value[2]++
	}

	return fmt.Sprintf("%s.%s.%s.", strconv.Itoa(value[0]), strconv.Itoa(value[1]), strconv.Itoa(value[2]))
}

func GetSubnetFromIP(ip string) (string, error) {
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		return "", fmt.Errorf("invalid IP format: %s", ip)
	}
	return strings.Join(parts[:3], ".") + ".", nil
}
