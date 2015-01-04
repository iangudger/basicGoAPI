package api

import (
	"log"
	"net/http"

	"github.com/iangudger/basicGoAPI/common"
	"github.com/iangudger/basicGoAPI/database"
)

func Handler(res http.ResponseWriter, req *http.Request, dbconn database.DB) {
	var result interface{} = req.URL.Query()

	// Check for duplicate parameters.
	for _, value := range req.URL.Query() {
		if len(value) != 1 {
			common.InternalServerError(res, duplicateParameters)
			return
		}
	}

	// Check for action
	if !checkParam(req.URL.Query(), "action") {
		common.InternalServerError(res, missingAction)
		return
	}
	actionName := req.URL.Query()["action"][0]
	log.Printf("Action: %s\n", actionName)

	// Get requested action
	action := actions[actionName]

	// Execute requested action
	if action != nil {
		result = action(res, req, dbconn)
	} else {
		result = ApiError{invalidAction}
	}

	// Write result if no error has been written
	if result != nil {
		writeJson(res, result)
	}
}
