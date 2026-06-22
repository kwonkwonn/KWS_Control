package api

import (
	"net/http"

	"github.com/easy-cloud-Knet/KWS_Control/service"
	"github.com/easy-cloud-Knet/KWS_Control/structure"
	"github.com/easy-cloud-Knet/KWS_Control/util"
)

type ApiVmConnectRequest struct {
	UUID structure.UUID `json:"uuid"`
}

type ApiVmConnectResponse struct {
	AuthToken string `json:"authToken"`
}

func (c *handlerContext) vmConnect(w http.ResponseWriter, r *http.Request) {
	log := util.GetLogger()
	defer r.Body.Close()

	uuidStr := r.URL.Query().Get("uuid")
	if uuidStr == "" {
		util.RespondError(w, http.StatusBadRequest, "Missing 'uuid' query parameter")
		log.Error("vmConnect: missing 'uuid' query parameter", nil, true)
		return
	}

	uuid := structure.UUID(uuidStr)
	authToken, err := service.GetGuacamoleToken(uuid, c.context)
	if err != nil {
		log.Error("vmConnect: failed to get Guacamole token: %v", err, true)
		util.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.RespondJSON(w, http.StatusOK, ApiVmConnectResponse{AuthToken: authToken})
}
