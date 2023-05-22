package main

import (
	"fmt"
	"go-webhook/webhook"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

var envs = []string{"SG_API_KEY",
	"SG_FROM",
	"SG_FROM_NAME",
	"SG_TO_LIST",
	"SG_EMAIL_TMPL_FILE",
	"GH_SECRET",
	"ALERT_SERVICE_LIST",
	"POSTGRESQL_PASSWORD",
	"POSTGRESQL_DATABASE",
	"POSTGRESQL_USERNAME",
	"POSTGRESQL_ADDRESS",
}

func init() {
	for _, v := range envs {
		if !webhook.DoesEnvExist(v) {
			fmt.Printf("Environment variable %s Doesn't exit", v)
			os.Exit(0)
		}
		fmt.Printf("Set: %s\n", v)
	}

}

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {

		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

var routes = Routes{
	Route{
		"postWebhook",
		"POST",
		"/webhook",
		webhook.GetWebhookData,
	},
	Route{
		"getHealth",
		"GET",
		"/health",
		webhook.GetHealth,
	},
	Route{
		"getAll",
		"GET",
		"/all-releases",
		webhook.GetAllRelease,
	},
	Route{
		"getAll",
		"GET",
		"/release/{arthur}",
		webhook.GetRelease,
	},
}

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}

func main() {
	// Init()
	r := NewRouter()
	fmt.Print("Serving http request at localhost:8080.....\n")
	fmt.Printf("Alert for PR with labels  \"%s\"  will be sent to %s\n", os.Getenv("ALERT_SERVICE_LIST"), os.Getenv("SG_TO_LIST"))
	log.Fatal(http.ListenAndServe(":8080", r))
}
