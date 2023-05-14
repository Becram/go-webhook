package webhook

import (
	"encoding/json"
	"net/http"
)

type Status struct {
	Code string `json:"code"`
}

func GetHealth(w http.ResponseWriter, req *http.Request) {

	var status = []Status{{Code: "200"}}
	json.NewEncoder(w).Encode(status)

}
