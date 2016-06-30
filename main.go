package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"github.com/gaia-docker/tugbot-result-service/dataupload"
	"github.com/gaia-docker/tugbot-result-service/pool"
	"github.com/gaia-docker/tugbot-result-service/websocket"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"net/url"
)

const defaultLogLevel = log.DebugLevel

var address = flag.String("address", "localhost:8080", "http service address")

var loglevel = flag.String("loglevel", defaultLogLevel.String(), "log level")

var homeTemplate = template.Must(template.ParseFiles("home.html"))

var hub = pool.NewHub()

func init() {

	flag.Parse()
	level, err := log.ParseLevel(*loglevel)
	if err != nil {
		level = defaultLogLevel
	}
	log.SetLevel(level)
}

func main() {

	log.Info("Starting tugbot-result-service...")
	go hub.Run()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", home).Methods("GET")
	router.Handle("/echo", websocket.NewEchoHandler(hub)).Methods("GET")
	router.Handle("/results", dataupload.NewUploadHandler(hub)).Methods("POST").
		Headers("Content-Type", "application/gzip")
	log.Fatal(http.ListenAndServe(*address, router))
	log.Infof("Listening on %s", address)
}

func home(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	var wsAddress = url.URL{Scheme: "ws", Host: *address, Path: "/echo"}
	homeTemplate.Execute(w, wsAddress.String())
}
