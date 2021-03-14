package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"

)

func main() {
	//This section will start the HTTP server and expose
	//any metrics on the /metrics endpoint.

	http.Handle("/metrics", promhttp.Handler())
	log.Info("Beginning to serve on port :9798")
	log.Fatal(http.ListenAndServe(":9798", nil))
}