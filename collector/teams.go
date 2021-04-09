package collector

import (
	"github.com/PagerDuty/go-pagerduty"
	"github.com/prometheus/client_golang/prometheus"
)

// TeamsCollector is struct for teams desc
type TeamsCollector struct {
	totalTeamsGaugeDesc *prometheus.Desc
}

// NewTeamsCollector is new teams collector registered in main
func NewTeamsCollector() *TeamsCollector {
	return &TeamsCollector{
		totalTeamsGaugeDesc: prometheus.NewDesc("pagerduty_total_teams_metric", "The number of total teams in AIpagerduty", nil, nil),
	}
}

// Describe is channel for teams
func (c *TeamsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.totalTeamsGaugeDesc
}

// Collect is teams api collector
func (c *TeamsCollector) Collect(ch chan<- prometheus.Metric) {

	teams := getTotalTeams()

	ch <- prometheus.MustNewConstMetric(
		c.totalTeamsGaugeDesc,
		prometheus.GaugeValue,
		float64(teams),
	)
}

func getTotalTeams() int {
	var opts pagerduty.ListUsersOptions
	var APIList pagerduty.APIListObject
	var Users []pagerduty.User

	for {
		eps, err := client.ListUsers(opts)

		if err != nil {
			panic(err)
		}

		Users = append(Users, eps.Users...)
		APIList.Offset += 25
		opts = pagerduty.ListUsersOptions{APIListObject: APIList}

		if eps.More != true {
			break
		}
	}
	totalUsers := len(Users)

	return totalUsers
}
