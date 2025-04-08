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
		return structure.ControlContext{}, err
	}

	g, ctx := errgroup.WithContext(context.Background())
	for _, coreAddress := range config.Cores {
		parts := strings.Split(coreAddress, ":")
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

			core = &newCore
		}

		client := request.NewCoreClient(core)
		g.Go(func() error {
			cpuResp, err := client.GetCoreMachineCpuInfo(ctx)
			if err != nil {
				return err
			}
			memResp, err := client.GetCoreMachineMemoryInfo(ctx)
			if err != nil {
				return err
			}
			diskResp, err := client.GetCoreMachineDiskInfo(ctx)
			if err != nil {
				return err
			}

			core.CoreInfoIdx.Cpu = uint32(cpuResp.Idle)
			core.CoreInfoIdx.Memory = uint32(memResp.Available)
			core.CoreInfoIdx.Disk = uint32(diskResp.Free)

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return structure.ControlContext{}, fmt.Errorf("failed to get core info: %w", err)
	}

	return infra, nil
}

func findCore(cores []structure.Core, ip string, port uint16) *structure.Core {
	for _, c := range cores {
		if c.IP == ip && c.Port == port {
			return &c
		}
	}
	return nil
}
