package service

import (
	"context"
	"github.com/easy-cloud-Knet/KWS_Control/request"
	"github.com/easy-cloud-Knet/KWS_Control/structure"
)

func GetGuacamoleToken(uuid structure.UUID, ctx *structure.ControlContext) (string, error) {
	core := ctx.FindCoreByVmUUID(uuid)
	if core == nil {
		return "", structure.ErrCoreNotFound(uuid)
	}

	if vm, exists := core.VMInfoIdx[uuid]; exists {
		client := request.NewGuacamoleClient(&ctx.Config)

		err := client.Authenticate(context.Background(), string(uuid), vm.GuacPassword)
		if err != nil {
			return "", err
		}

		return client.AuthToken(), nil
	} else {
		return "", structure.ErrVmNotFound(uuid)
	}
}
