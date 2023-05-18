package main

import (
	"fmt"
	"go-webhook/utility"
	"go-webhook/webhook"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

var envs = []string{"SG_API_KEY", "SG_FROM", "SG_FROM_NAME", "SG_TO_LIST", "SG_EMAIL_TMPL_FILE", "PR_PREFIX", "GH_SECRET"}

func init() {
	for _, v := range envs {
		if !utility.DoesEnvExist(v) {
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
	log.Fatal(http.ListenAndServe(":8080", r))
}
