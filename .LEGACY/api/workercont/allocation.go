package WorkerCont

import (
	"sync"

	vms "github.com/easy-cloud-Knet/KWS_Control/vm"
)

func updateFreeCPUAndMemory(core *vms.Core) error {
	cpuTask := NewGetCoreMachineFreeCpuInfoTask(core)
	cpuInfo, err := cpuTask.Await()
	if err != nil {
		return err
	}

	memTask := NewGetCoreMachineFreeMemoryInfoTask(core)
	memInfo, err := memTask.Await()
	if err != nil {
		return nil
	}

	core.FreeCPU = cpuInfo.Information.Idle
	core.FreeMemory = memInfo.Information.Available

	return nil
}

func updateAllCoresFreeResources(control *vms.ControlInfra) {
	var wg sync.WaitGroup

	semaphore := make(chan struct{}, 10)

	for i := range control.Cores {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(index int) {
			defer wg.Done()
			defer func() { <-semaphore }()

			core := &control.Cores[index]
			_ = updateFreeCPUAndMemory(core)
		}(i)
	}

	wg.Wait()
}

func UpdateCoreAndSelectCoreForNewVM(control *vms.ControlInfra, memory uint64, cpu float64) *vms.Core {
	updateAllCoresFreeResources(control)

	maxScore := 0.0
	var selectedCore *vms.Core = nil

	for i := range control.Cores {
		core := control.Cores[i]
		if core.FreeMemory < memory || core.FreeCPU < cpu {
			continue
		}

		coreScore := float64(core.FreeMemory) + core.FreeCPU*10

		if maxScore < coreScore {
			selectedCore = &control.Cores[i]
			maxScore = coreScore
		}
	}

	return selectedCore
}
