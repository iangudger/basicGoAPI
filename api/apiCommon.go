package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

func checkParam(values url.Values, key string) bool {
	if values[key] != nil || len(values[key]) == 1 {
		return true
	}
	return false
}

func writeJson(res http.ResponseWriter, result interface{}) {
	res.Header().Set("Content-Type", "application/json;charset=utf-8")
	encoded, err := json.Marshal(result)
	if err != nil {
		fmt.Fprintln(res, err)
		return
	}
	log.Println("API result:", string(encoded))
	fmt.Fprintln(res, string(encoded))
}
