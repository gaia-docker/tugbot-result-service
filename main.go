package main

import (
	"errors"
	"flag"
	logger "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

var address = flag.String("address", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options

var homeTemplate = template.Must(template.ParseFiles("home.html"))

var wsAddress = url.URL{Scheme: "ws", Host: *address, Path: "/echo"}

func main() {
	flag.Parse()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", home).Methods("GET")
	router.HandleFunc("/echo", echo).Methods("GET")
	router.HandleFunc("/upload-data", upload).Methods("POST")
	logger.Fatal(http.ListenAndServe(*address, router))
}

func echo(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Errorf("Error during upgrade to web socket protocol %+v", err)
		return
	}
	defer conn.Close()
	mt, message, err := conn.ReadMessage()
	if err != nil {
		logger.Errorf("Error while reading from web socket %+v", err)
	}
	logger.Infof("echo Receive: %s", message)
	err = conn.WriteMessage(mt, message)
	if err != nil {
		logger.Errorf("Error while writing to web socket %+v", err)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, wsAddress.String())
}

func upload(writer http.ResponseWriter, request *http.Request) {

	retStatus := http.StatusOK
	body, err := getBody(request)
	if err != nil {
		retStatus = http.StatusBadRequest
		logger.Error("Error fetching request body. ", err)
	} else {
		logger.Infof("connecting to %s", wsAddress.String())

		conn, _, err := websocket.DefaultDialer.Dial(wsAddress.String(), nil)
		if err != nil {
			logger.Fatal("dial:", err)
		}
		defer conn.Close()
		err = conn.WriteMessage(websocket.TextMessage, []byte(*body))
		if err != nil {
			logger.Errorf("Error in writing to web socket %+v", err)
			return
		}
		var msg = make([]byte, 512)
		_, msg, err = conn.ReadMessage()
		if err != nil {
			log.Fatal(err)
		}
		logger.Infof("upload Receive: %s", msg)
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
		logger.Error(err)
		return nil, err
	}
	ret := string(body)

	return &ret, nil
}
