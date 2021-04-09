package collector

import (
	"github.com/PagerDuty/go-pagerduty"
	"github.com/prometheus/client_golang/prometheus"
)

type UsersCollector struct {
	totalUsersGaugeDesc *prometheus.Desc
}

func NewUsersCollector() *UsersCollector {
	return &UsersCollector{
		totalUsersGaugeDesc: prometheus.NewDesc("pagerduty_total_users_metric", "The number of total users in AIpagerduty", nil, nil),
	}
}

func (c *UsersCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.totalUsersGaugeDesc
}

func (c *UsersCollector) Collect(ch chan<- prometheus.Metric) {
	users := getTotalUsers()
	ch <- prometheus.MustNewConstMetric(
		c.totalUsersGaugeDesc,
		prometheus.GaugeValue,
		float64(users),
	)
}

func getTotalUsers() int {
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
