package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gaia-docker/tugbot-result-service/dataupload"
	"github.com/gaia-docker/tugbot-result-service/pool"
	"github.com/gaia-docker/tugbot-result-service/websocket"
	"github.com/gorilla/mux"
	"github.com/urfave/cli"
	"html/template"
	"net/http"
	"net/url"
	"os"
)

const defaultLogLevel = log.DebugLevel

var address string

var loglevel string

var homeTemplate = template.Must(template.ParseFiles("views/home.html"))

var hub = pool.NewHub()

func main() {

	app := cli.NewApp()
	app.Version = "1.0.0"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "address, a",
			Value:       "localhost:8080",
			Usage:       "http service address",
			Destination: &address,
		},
		cli.StringFlag{
			Name:        "loglevel, l",
			Value:       defaultLogLevel.String(),
			Usage:       "log level",
			Destination: &loglevel,
		},
	}

	app.Name = "tugbot-result-service"
	app.Usage = "Implements Result Service API and exposes websocket which present live stream of test results."
	app.Action = start

	if err := app.Run(os.Args); err != nil {
		log.Error("exiting from main: ", err)
	}
}

func start(c *cli.Context) error {

	log.Info("Starting tugbot-result-service...")
	go hub.Run()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", home).Methods("GET")
	router.Handle("/echo", websocket.NewEchoHandler(hub)).Methods("GET")
	router.Handle("/results", dataupload.NewUploadHandler(hub)).Methods("POST").
		Headers("Content-Type", "application/gzip")
	log.Infof("Listening on %s", address)

	return http.ListenAndServe(address, router)
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
	var wsAddress = url.URL{Scheme: "ws", Host: address, Path: "/echo"}
	homeTemplate.Execute(w, wsAddress.String())
}
