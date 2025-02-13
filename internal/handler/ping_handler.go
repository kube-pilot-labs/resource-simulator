package handler

import (
	"net/http"

	"github.com/kube-pilot-labs/resource-simulator/internal/util"
)

func PingHandler(w http.ResponseWriter, r *http.Request) {
	util.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"result": "true",
	})
}
