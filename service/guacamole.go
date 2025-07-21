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

	pass, err := GetGuacPassword(ctx.DB, uuid)
	if err != nil {
		return "", err
	}
	client := request.NewGuacamoleClient(&ctx.Config)
	if err := client.Authenticate(context.Background(), string(uuid), pass); err != nil {
		return "", err
	}
	return client.AuthToken(), nil
}
