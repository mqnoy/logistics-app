package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func DefaultQuery(r *http.Request, key string, defaultValue string) string {
	if value := r.URL.Query().Get(key); value != "" {
		return value
	}
	return defaultValue
}

func GetQuery(r *http.Request, key string) (string, bool) {
	if value := r.URL.Query().Get(key); value != "" {
		return value, true
	}
	return "", false
}

func ParseQueryToBool(query string) *bool {
	result, err := strconv.ParseBool(query)
	if err != nil {
		return nil
	}

	return &result
}

func ParseQueryToInt64Array(query string) (result []int64) {
	if err := json.Unmarshal([]byte(query), &result); err != nil {
		log.Println(err)
		return result
	}
	return result
}
