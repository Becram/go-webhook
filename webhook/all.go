package webhook

import (
	"encoding/json"
	"fmt"
	"go-webhook/psql"
	"net/http"

	"github.com/gorilla/mux"
)

// The function returns the health status of the application with its version number.
func GetAllRelease(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := psql.Connect()
	defer db.Close()
	var notification []psql.Notification
	err := db.Model(&notification).Select()

	psql.CheckNilErr(err)

	json.NewEncoder(w).Encode(notification)
}

func GetRelease(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := psql.Connect()
	defer db.Close()
	param := mux.Vars(req)
	fmt.Printf("Parameter received %s\n", param["arthur"])
	release := &psql.Notification{Arthur: param["arthur"]}
	var notification []psql.Notification
	err := db.Model(&notification).Where("arthur = ?", release.Arthur).Select()

	psql.CheckNilErr(err)

	json.NewEncoder(w).Encode(notification)
}
