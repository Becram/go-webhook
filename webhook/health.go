package webhook

import (
	"encoding/json"
	"net/http"
	"os"
)

type Status struct {
	Code    int    `json:"code"`
	Version string `json:"version"`
}

// The function returns the health status of the application with its version number.
func GetHealth(w http.ResponseWriter, req *http.Request) {
	var status = []Status{{Code: http.StatusOK, Version: os.Getenv("VERSION")}}
	json.NewEncoder(w).Encode(status)

}
