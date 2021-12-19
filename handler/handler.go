package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"cmd/service/main.go/pkg/config"
	easysync "cmd/service/main.go/pkg/easy-sync"

	"github.com/google/uuid"
	"github.com/savsgio/atreugo/v11"
)

// SyncEasyToWP triggers the sync between Easyverein and Wordpress
// and provides an url where the status of the sync can be requested
// Method: GET
// Path: /sync
// Responses: 202
func SyncEasyToWP(ctx *atreugo.RequestCtx) error {
	log.Println("Processing:")
	log.Printf("%s http://%s:%d%s",
		string(ctx.Method()),
		config.GetConfig().API.Host,
		config.GetConfig().API.Port,
		string(ctx.RequestURI()),
	)

	requestID := uuid.New()
	requestIDStr := fmt.Sprintf("%v", requestID)
	pollingURL := fmt.Sprintf("http://%s:%d/sync/status/%s",
		config.GetConfig().API.Host,
		config.GetConfig().API.Port,
		requestIDStr,
	)
	log.Println("Status Polling URL:")
	log.Println(pollingURL)

	go easysync.TriggerSync(requestIDStr)

	msg := fmt.Sprintf("sync request accepted. status url: %s", pollingURL)
	return ctx.JSONResponse(msg, http.StatusAccepted)
}

// GetSyncStatus sync status polling route
// if sync still running status will remain 202 'accepted'
// if sync successful status will be 200 'ok'
// Method: GET
// Path: /sync/status/{requestID}
// Response: 200, 202, 404
func SyncStatus(ctx *atreugo.RequestCtx) error {
	requestID := ctx.UserValue("requestID")
	requestIDStr := fmt.Sprintf("%v", requestID)
	syncStatus, err := easysync.GetSyncStatus(requestIDStr)
	if err != nil {
		log.Println("error polling sync status:", err)
		return ctx.JSONResponse(err.Error(), http.StatusBadRequest)
	}

	response, err := json.Marshal(syncStatus)
	if err != nil {
		msg := fmt.Sprintf("error marshalling syncstatus to json:, %v", err.Error())
		ctx.Error(msg, http.StatusInternalServerError)
		return fmt.Errorf(msg)
	}

	// if lasts sync was not, or is not yet successful
	if !syncStatus.Status {
		return ctx.JSONResponse(response, http.StatusAccepted)
	}

	return ctx.JSONResponse(response)
}
