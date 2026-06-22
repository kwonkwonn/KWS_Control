package api

import (
	"encoding/json"
	"net/http"

	"github.com/easy-cloud-Knet/KWS_Control/service"
	"github.com/easy-cloud-Knet/KWS_Control/structure"
	"github.com/easy-cloud-Knet/KWS_Control/util"
)

type ApiTakeSnapshotRequest struct {
	UUID     structure.UUID `json:"uuid"`
	SnapName string         `json:"snapName"`
}

type ApiDeleteSnapshotRequest struct {
	UUID    structure.UUID `json:"uuid"`
	SnapKey string         `json:"snapKey"`
}

func (c *handlerContext) takeSnapshot(w http.ResponseWriter, r *http.Request) {
	log := util.GetLogger()
	defer r.Body.Close()

	var req ApiTakeSnapshotRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.UUID == "" || req.SnapName == "" {
		util.RespondError(w, http.StatusBadRequest, "uuid and snapName are required")
		return
	}

	if err := service.TakeSnapshot(req.UUID, req.SnapName, c.context); err != nil {
		log.Error("takeSnapshot: %v", err, true)
		util.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.RespondJSON(w, http.StatusOK, map[string]string{"message": "snapshot taken"})
}

func (c *handlerContext) listSnapshots(w http.ResponseWriter, r *http.Request) {
	log := util.GetLogger()

	uuid := structure.UUID(r.URL.Query().Get("uuid"))
	if uuid == "" {
		util.RespondError(w, http.StatusBadRequest, "uuid query parameter is required")
		return
	}

	keys, err := service.ListSnapshots(uuid)
	if err != nil {
		log.Error("listSnapshots: %v", err, true)
		util.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.RespondJSON(w, http.StatusOK, map[string][]string{"snapshots": keys})
}

func (c *handlerContext) deleteSnapshot(w http.ResponseWriter, r *http.Request) {
	log := util.GetLogger()
	defer r.Body.Close()

	var req ApiDeleteSnapshotRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.UUID == "" || req.SnapKey == "" {
		util.RespondError(w, http.StatusBadRequest, "uuid and snapKey are required")
		return
	}

	if err := service.DeleteSnapshot(req.UUID, req.SnapKey); err != nil {
		log.Error("deleteSnapshot: %v", err, true)
		util.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}
