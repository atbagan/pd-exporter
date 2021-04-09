package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/version"
	"net/http"
	"pd-exporter/collector"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	//This section will start the HTTP server and expose
	//any metrics on the /metrics endpoint.

	var (
		Name          = "pd-exporter"
		listenAddress = kingpin.Flag("web.listen-address",
			"Address to listen on for web server").
			Default(":9696").Envar("WEB_LISTEN_ADDRESS").String()
		metricsPath = kingpin.Flag("web.telemetry-path",
			"Path where to expose metrics").
			Default("/metrics").Envar("WEB_TELEMETRY_PATH").String()
		pdAnalyticsSettings = kingpin.Flag("pd.analytics_settings",
			"Pagerduty Analytics Settings on/off.").
			Default("true").Envar("PD_ANALYTICS_SETTINGS").Bool()
		pdServicesSettings = kingpin.Flag("pd.service_settings",
			"Pagerduty Service Settings on/off.").
			Default("true").Envar("PD_SERVICES_SETTINGS").Bool()
		pdTeamsSettings = kingpin.Flag("pd.teams_settings",
			"Pagerduty Teams Settings on/off.").
			Default("true").Envar("PD_TEAMS_SETTINGS").Bool()
		pdUsersSettings = kingpin.Flag("pd.users_settings",
			"Pagerduty Users Settings on/off.").
			Default("true").Envar("PD_USERS_SETTINGS").Bool()
	)
	kingpin.Version(version.Print(Name))
	kingpin.CommandLine.HelpFlag.Short('h')
	kingpin.Parse()

	if *pdAnalyticsSettings {
		prometheus.MustRegister(collector.NewAnalyticsCollector())
	}

	if *pdServicesSettings {
		prometheus.MustRegister(collector.NewServiceCollector())
	}

	if *pdUsersSettings {
		prometheus.MustRegister(collector.NewUsersCollector())
	}

	if *pdTeamsSettings {
		prometheus.MustRegister(collector.NewTeamsCollector())
	}

	http.Handle(*metricsPath, promhttp.Handler())
	log.Info("Beginning to serve on port " + *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
