package main

import (
	"errors"
	"flag"
	log "github.com/Sirupsen/logrus"
	"github.com/gaia-docker/tugbot-result-service/pool"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
)

const defaultLogLevel = log.DebugLevel

var address = flag.String("address", "localhost:8080", "http service address")

var loglevel = flag.String("loglevel", defaultLogLevel.String(), "log level")

var upgrader = websocket.Upgrader{} // use default options

var homeTemplate = template.Must(template.ParseFiles("home.html"))

var wsAddress = url.URL{Scheme: "ws", Host: *address, Path: "/echo"}

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

	go hub.Run()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", home).Methods("GET")
	router.HandleFunc("/echo", echo).Methods("GET")
	router.HandleFunc("/upload-data", upload).Methods("POST")
	log.Fatal(http.ListenAndServe(*address, router))
}

func echo(w http.ResponseWriter, r *http.Request) {

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorf("Failed upgrading connection to web socket %+v", err)
		return
	}
	log.Infof("New websocket connection established %s", wsAddress.String())
	conn := pool.NewConnection(ws)
	hub.Register(conn)
	go conn.Listen(hub.Unregister)
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
	homeTemplate.Execute(w, wsAddress.String())
}

func upload(writer http.ResponseWriter, request *http.Request) {

	retStatus := http.StatusOK
	body, err := getBody(request)
	if err != nil {
		retStatus = http.StatusBadRequest
		log.Error("Error fetching request body. ", err)
	} else {
		hub.Broadcast(body)
	}
	writer.WriteHeader(retStatus)
}

func getBody(request *http.Request) (*string, error) {

	requestBody := request.Body
	if requestBody == nil {
		return nil, errors.New("Empty request body")
	}
	body, err := ioutil.ReadAll(requestBody)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	ret := string(body)

	return &ret, nil
}
