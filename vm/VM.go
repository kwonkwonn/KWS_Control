package vms

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"sync"
)

func (c *Computer) GetVMList(VMList map[UUID]*VM) (err error) {
	baseUrl := c.IP
	if !strings.HasPrefix(baseUrl, "http://") && !strings.HasPrefix(baseUrl, "https://") {
		baseUrl = "http://" + baseUrl
	}

	resp, err := http.Get(baseUrl + ":8080" + "/getStatus")
	if err != nil {
		return
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var info []VMInfo
	err = json.Unmarshal(body, &info)
	if err != nil {
		return
	}

	for _, v := range info {
		vm := &VM{
			VMInfo:      v,
			IsAlive:     true, // TODO: check alive, allocated
			IsAllocated: false,
			IsLocatedAt: *c,
		}
		VMList[v.UUID] = vm
	}

	return nil
}

func (i *InfraContext) UpdateList(ctx context.Context) {
	var wg sync.WaitGroup
	wg.Add(len(i.Computers))
	for _, c := range i.Computers {
		computer := c
		go func() {
			defer wg.Done()

			select {
			case <-ctx.Done():
				return
			default:
				_ = computer.GetVMList(i.VMPool)
				// TODO: reschedule if failed to fetch VM list
			}
		}()
	}

	wg.Wait()
}
