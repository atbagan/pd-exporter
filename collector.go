package main

import (
	"encoding/json"
	"fmt"
	"github.com/PagerDuty/go-pagerduty"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var authToken = os.Getenv("AUTH_TOKEN")

var client = pagerduty.NewClient(authToken)

func init() {
	//Register metrics with prometheus
	prometheus.MustRegister(NewMyCollector())
}

type MyCollector struct {
	totalGaugeDesc            *prometheus.Desc
	complianceGaugeDesc       *prometheus.Desc
	usersGaugeDesc            *prometheus.Desc
	teamsGaugeDesc            *prometheus.Desc
	businessServicesGaugeDesc *prometheus.Desc
}

func (c *MyCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.totalGaugeDesc
	ch <- c.complianceGaugeDesc
	ch <- c.usersGaugeDesc
	ch <- c.teamsGaugeDesc
	ch <- c.businessServicesGaugeDesc
}

func (c *MyCollector) Collect(ch chan<- prometheus.Metric) {
	total, compliance := callPagerdutyApi()
	users := callPagerdutyApiUsers()
	teams := callPagerdutyApiTeams()
	businessServices := callPagerdutyApiBusinessServices()
	callPagerDutyApiAnalytics()
	ch <- prometheus.MustNewConstMetric(
		c.totalGaugeDesc,
		prometheus.GaugeValue,
		float64(total),
	)
	ch <- prometheus.MustNewConstMetric(
		c.complianceGaugeDesc,
		prometheus.GaugeValue,
		float64(compliance),
	)
	ch <- prometheus.MustNewConstMetric(
		c.usersGaugeDesc,
		prometheus.GaugeValue,
		float64(users),
	)
	ch <- prometheus.MustNewConstMetric(
		c.teamsGaugeDesc,
		prometheus.GaugeValue,
		float64(teams),
	)
	ch <- prometheus.MustNewConstMetric(
		c.businessServicesGaugeDesc,
		prometheus.GaugeValue,
		float64(businessServices),
	)
}
func NewMyCollector() *MyCollector {
	return &MyCollector{
		totalGaugeDesc:            prometheus.NewDesc("total_pagerduty_services_metric", "The number of total services in AIpagerduty", nil, nil),
		complianceGaugeDesc:       prometheus.NewDesc("compliancy_pagerduty_services_metric", "Shows the number of compliant services names", nil, nil),
		usersGaugeDesc:            prometheus.NewDesc("total_pagerduty_users_metric", "Shows the total number of users", nil, nil),
		teamsGaugeDesc:            prometheus.NewDesc("total_pagerduty_teams_metric", "Shows the total number of teams", nil, nil),
		businessServicesGaugeDesc: prometheus.NewDesc("total_pagerduty_business_services_metric", "Shows the total number of business services", nil, nil),
	}
}

func callPagerdutyApi() (int, int) {

	var opts pagerduty.ListServiceOptions
	var APIList pagerduty.APIListObject
	var Services []pagerduty.Service

	for {
		eps, err := client.ListServices(opts)

		if err != nil {
			panic(err)
		}

		Services = append(Services, eps.Services...)
		APIList.Offset += 25
		opts = pagerduty.ListServiceOptions{APIListObject: APIList}

		if eps.More != true {
			break
		}

	}

	totalServices := len(Services)
	complianceCount := 0
	for k, _ := range Services {
		re := regexp.MustCompile("_SVC+")
		if re.MatchString(Services[k].Name) {
			complianceCount += 1
		}
	}
	return totalServices, complianceCount
}

func callPagerdutyApiUsers() int {
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

func callPagerdutyApiTeams() int {
	var opts pagerduty.ListTeamOptions
	var APIList pagerduty.APIListObject
	var Teams []pagerduty.Team

	for {
		eps, err := client.ListTeams(opts)

		if err != nil {
			panic(err)
		}

		Teams = append(Teams, eps.Teams...)
		APIList.Offset += 25
		opts = pagerduty.ListTeamOptions{APIListObject: APIList}

		if eps.More != true {
			break
		}
	}
	totalTeams := len(Teams)

	return totalTeams
}

func callPagerdutyApiBusinessServices() int {
	var opts pagerduty.ListBusinessServiceOptions
	var APIList pagerduty.APIListObject

	totalBusinessServices := 0
	for {
		eps, err := client.ListBusinessServices(opts)

		if err != nil {
			panic(err)
		}

		APIList.Offset += 25
		opts = pagerduty.ListBusinessServiceOptions{APIListObject: APIList}
		totalBusinessServices = len(eps.BusinessServices)
		if eps.More != true {
			break
		}
	}
	return totalBusinessServices
}

func callPagerdutyApiIncidents(){
	var opts pagerduty.ListIncidentsOptions
	var APIList pagerduty.APIListObject
	var Incidents []pagerduty.Incident
	for {
		eps, err := client.ListIncidents(opts)
		if err != nil{panic(err)}

		Incidents = append(Incidents, eps.Incidents...)
		APIList.Offset += 25
		opts = pagerduty.ListIncidentsOptions{APIListObject: APIList}

		if eps.More != true{
			break
		}
	}

	for k, _ := range Incidents {
		fmt.Println(Incidents[k].Service.ID)
	}
}

type AnalyticsD struct {
	ServiceID                      string `json:"service_id"`
	ServiceName                    string `json:"service_name"`
	TeamID                         string `json:"team_id"`
	TeamName                       string `json:"team_name"`
	MeanSecondsToResolve           int    `json:"mean_seconds_to_resolve"`
	MeanSecondsToFirstAck          int    `json:"mean_seconds_to_first_ack"`
	MeanSecondsToEngage            int    `json:"mean_seconds_to_engage"`
	MeanSecondsToMobilize          int    `json:"mean_seconds_to_mobilize"`
	MeanEngagedSeconds             int    `json:"mean_engaged_seconds"`
	MeanEngagedUserCount           int    `json:"mean_engaged_user_count"`
	TotalEscalationCount           int    `json:"total_escalation_count"`
	MeanAssignmentCount            int    `json:"mean_assignment_count"`
	TotalBusinessHourInterruptions int    `json:"total_business_hour_interruptions"`
	TotalSleepHourInterruptions    int    `json:"total_sleep_hour_interruptions"`
	TotalOffHourInterruptions      int    `json:"total_off_hour_interruptions"`
	TotalSnoozedSeconds            int    `json:"total_snoozed_seconds"`
	TotalEngagedSeconds            int    `json:"total_engaged_seconds"`
	TotalIncidentCount             int    `json:"total_incident_count"`
	UpTimePct                      int    `json:"up_time_pct"`
	UserDefinedEffortSeconds       int    `json:"user_defined_effort_seconds"`
	RangeStart                     string `json:"range_start"`
}

type AnalyticsR struct {
	Data            []AnalyticsD  `json:"data,omitempty"`
	AnalyticsFilter *AnalyticsF `json:"filters,omitempty"`
	AggregateUnit   string           `json:"aggregate_unit,omitempty"`
	TimeZone        string           `json:"time_zone,omitempty"`
}

type AnalyticsF struct {
	CreatedAtStart string   `json:"created_at_start,omitempty"`
	CreatedAtEnd   string   `json:"created_at_end,omitempty"`
	Urgency        string   `json:"urgency,omitempty"`
	Major          bool     `json:"major,omitempty"`
	ServiceIDs     []string `json:"service_ids,omitempty"`
	TeamIDs        []string `json:"team_ids,omitempty"`
	PriorityIDs    []string `json:"priority_ids,omitempty"`
	PriorityName   []string `json:"priority_name,omitempty"`
}

func callPagerDutyApiAnalytics(){
	url := "https://api.pagerduty.com/analytics/metrics/incidents/services"
	payload := strings.NewReader("{\"aggregate_unit\":\"week\"}")

	var analyticsR AnalyticsR

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {panic(err)}

	req.Header.Add("X-EARLY-ACCESS", "analytics-v2")
	req.Header.Add("Accept","application/vnd.pagerduty+json;version=2")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Token token="+authToken)

	resp, _ := http.DefaultClient.Do(req)
	err  = decodeJSON(resp, &analyticsR)
	fmt.Println(resp)
	for _, v := range analyticsR.Data {
		fmt.Println(v.ServiceName, v.MeanSecondsToFirstAck, v.MeanSecondsToResolve)
	}
}

func decodeJSON(resp *http.Response, payload interface{}) error {
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(payload)
}