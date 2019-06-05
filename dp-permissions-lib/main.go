package main

import (
	"html/template"
	"net/http"

	"github.com/ONSdigital/dp-permissions/auth"
	"github.com/ONSdigital/dp-permissions/permissions"
	"github.com/ONSdigital/go-ns/rchttp"
	"github.com/ONSdigital/log.go/log"
	"github.com/gorilla/mux"
)

func main() {
	log.Namespace = "example-api"
	rc := rchttp.NewClient()
	authenticator := permissions.New("http://localhost:8082/permissions", rc)

	auth.Configure("dataset_id", mux.Vars, authenticator)

	temp, _ := template.ParseFiles("ralph.html")

	apiEndpoint := func(w http.ResponseWriter, r *http.Request) {
		// API handler code lives here
		if err := temp.Execute(w, nil); err != nil {
			w.WriteHeader(500)
		}
	}

	r := mux.NewRouter()
	r.HandleFunc("/datasets/{dataset_id}", auth.Require(permissions.CRUD{Read: true}, apiEndpoint))

	log.Event(nil, "starting example service")
	if err := http.ListenAndServe(":8090", r); err != nil {
		log.Event(nil, "boom", log.Error(err))
	}
}
