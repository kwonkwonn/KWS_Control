package vms

func (c *ControlInfra) SelectCoreForNewVM(memory, cpu int) *Core {
	maxScore := 0
	var selectedCore *Core = nil

	for i := range c.Cores {
		core := c.Cores[i]
		if core.FreeMemory < memory || core.FreeCPU < cpu {
			continue
		}

		coreScore := core.FreeMemory + core.FreeCPU*10

		if maxScore < coreScore {
			selectedCore = &c.Cores[i]
		}
	}

	return selectedCore
}
