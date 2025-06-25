package startup

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/easy-cloud-Knet/KWS_Control/request"
	"github.com/easy-cloud-Knet/KWS_Control/structure"
	"golang.org/x/sync/errgroup"

	_ "gopkg.in/yaml.v3"
)

func Initialize(dataPath, configPath string) (structure.ControlContext, error) {
	b, err := os.ReadFile(dataPath)
	if err != nil {
		return structure.ControlContext{}, err
	}

	var infra structure.ControlContext
	if err := json.Unmarshal(b, &infra); err != nil {
		return structure.ControlContext{}, fmt.Errorf("failed to parse JSON: %v", err)
	}

	// 모든 Core 정의
	infra.VMLocation = make(map[structure.UUID]*structure.Core)
	for i := range infra.Cores {
		for vmUUID := range infra.Cores[i].VMInfoIdx {
			infra.VMLocation[vmUUID] = &infra.Cores[i]
		}
		infra.Cores[i].IsAlive = false
	}

	// config에 설정된 코어에 대해서 정보 업데이트
	config, err := readConfig(configPath)
	if err != nil {
		return structure.ControlContext{}, fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}

	// 환경변수에서 먼저 설정을 불러오고 값이 없을 때만 config 파일에서 가져옴
	var coreAddresses []string
	coresEnv := os.Getenv("CORES")
	if coresEnv != "" {
		coreAddresses = strings.Split(coresEnv, ",")
	} else {
		coreAddresses = config.Cores
	}

	g, ctx := errgroup.WithContext(context.Background())
	for _, coreAddress := range coreAddresses {
		parts := strings.Split(coreAddress, ":")
		if len(parts) != 2 {
			return structure.ControlContext{}, fmt.Errorf("invalid core address format in '%s': expected ip:port", coreAddress)
		}
		ip := parts[0]
		port, err := strconv.Atoi(parts[1])
		if err != nil {
			return structure.ControlContext{}, fmt.Errorf("error converting port number from %s: %w", coreAddress, err)
		}

		core := findCore(infra.Cores, ip, uint16(port))
		if core == nil {
			newCore := structure.Core{
				IP:      ip,
				Port:    uint16(port),
				IsAlive: true,
			}
			infra.Cores = append(infra.Cores, newCore)
			core = &infra.Cores[len(infra.Cores)-1]
		} else {
			core.IsAlive = true
		}

		currentCore := core
		g.Go(func() error {
			client := request.NewCoreClient(currentCore)

			cpuResp, err := client.GetCoreMachineCpuInfo(ctx)
			if err != nil {
				currentCore.IsAlive = false
				return fmt.Errorf("failed to get CPU info for core %s:%d: %w", currentCore.IP, currentCore.Port, err)
			}
			memResp, err := client.GetCoreMachineMemoryInfo(ctx)
			if err != nil {
				currentCore.IsAlive = false
				return fmt.Errorf("failed to get Memory info for core %s:%d: %w", currentCore.IP, currentCore.Port, err)
			}
			diskResp, err := client.GetCoreMachineDiskInfo(ctx)
			if err != nil {
				currentCore.IsAlive = false
				return fmt.Errorf("failed to get Disk info for core %s:%d: %w", currentCore.IP, currentCore.Port, err)
			}

			currentCore.CoreInfoIdx.Cpu = uint32(cpuResp.Idle)
			currentCore.CoreInfoIdx.Memory = uint32(memResp.Total)
			currentCore.CoreInfoIdx.Disk = uint32(diskResp.Total)

			currentCore.FreeDisk = uint32(diskResp.Free)
			currentCore.FreeMemory = uint32(memResp.Available)
			currentCore.FreeCPU = uint32(cpuResp.Idle)

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return structure.ControlContext{}, fmt.Errorf("failed to get core info: %w", err)
	}

	infra.Config = config
	return infra, nil
}

func findCore(cores []structure.Core, ip string, port uint16) *structure.Core {
	for i := range cores {
		if cores[i].IP == ip && cores[i].Port == port {
			return &cores[i]
		}
	}
	return nil
}
