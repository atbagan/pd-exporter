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
	"time"
)

var authToken = os.Getenv("AUTH_TOKEN")

//var authToken = "2zBBzbY8Qm5B-b8zLCbb"
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
	var pagerdutyServices = pdServices()
	total := totalServices(pagerdutyServices)
	compliance := compliantServices(pagerdutyServices)
	serviceIds := getCompliantServiceIds(pagerdutyServices)
	serviceInfoSlice := serviceComplianceInfoWithNameAndTeam(pagerdutyServices)

	users := callPagerdutyApiUsers()
	teams := callPagerdutyApiTeams()
	businessServices := callPagerdutyApiBusinessServices()
	mapOfMetrics := callPagerDutyApiAnalytics(serviceIds)

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
	for _, v := range serviceInfoSlice {

		metric, err := prometheus.NewConstMetric(
			prometheus.NewDesc("pagerduty_service_names_metric", fmt.Sprintf("service names "), nil, prometheus.Labels{"allServiceNames": v.ServiceName, "allTeamNames": v.ServiceTeam}),
			prometheus.GaugeValue,
			float64(v.Compliant),
		)
		if err != nil {
			panic(err)
		}
		ch <- metric
	}
	for k, v := range mapOfMetrics {
		metric, err := prometheus.NewConstMetric(
			prometheus.NewDesc("pagerduty_mtta_analytics_metric", fmt.Sprintf("Mean seconds to first ack "), nil, prometheus.Labels{"compliantServiceName": k}),
			prometheus.GaugeValue,
			float64(v[0]),
		)
		if err != nil {
			panic(err)
		}
		ch <- metric
	}
	for k, v := range mapOfMetrics {
		metric, err := prometheus.NewConstMetric(
			prometheus.NewDesc("pagerduty_mttr_analytics_metric", fmt.Sprintf("Mean seconds to resolve"), nil, prometheus.Labels{"compliantServiceName": k}),
			prometheus.GaugeValue,
			float64(v[1]),
		)
		if err != nil {
			panic(err)
		}
		ch <- metric
	}
}

func NewMyCollector() *MyCollector {
	return &MyCollector{
		totalGaugeDesc:            prometheus.NewDesc("pagerduty_total_services_metric", "The number of total services in AIpagerduty", nil, nil),
		complianceGaugeDesc:       prometheus.NewDesc("pagerduty_total_services_compliant_metric", "Shows the number of compliant services names", nil, nil),
		usersGaugeDesc:            prometheus.NewDesc("pagerduty_total_users_metric", "Shows the total number of users", nil, nil),
		teamsGaugeDesc:            prometheus.NewDesc("pagerduty_total_teams_metric", "Shows the total number of teams", nil, nil),
		businessServicesGaugeDesc: prometheus.NewDesc("pagerduty_total_business_services_metric", "Shows the total number of business services", nil, nil),
	}
}

type ServiceInfo struct {
	Compliant   int
	ServiceName string
	ServiceTeam string
}

func pdServices() []pagerduty.Service {
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
	return Services
}

func totalServices(pagerdutyServices []pagerduty.Service) int {
	services := pagerdutyServices
	return len(services)
}

func compliantServices(pagerdutyServices []pagerduty.Service) int {
	services := pagerdutyServices

	complianceCount := 0
	for k, _ := range services {

		re := regexp.MustCompile("_SVC+")

		if re.MatchString(services[k].Name) {
			complianceCount += 1
		}
	}
	return complianceCount
}

func getCompliantServiceIds(pagerdutyServices []pagerduty.Service) []string {
	services := pagerdutyServices
	serviceIds := make([]string, 0)

	for k, _ := range services {

		re := regexp.MustCompile("_SVC+")

		if re.MatchString(services[k].Name) {
			serviceIds = append(serviceIds, services[k].ID)
		}
	}
	return serviceIds
}

func serviceComplianceInfoWithNameAndTeam(pagerdutyServices []pagerduty.Service) []ServiceInfo {

	services := pagerdutyServices
	var serviceInfo ServiceInfo
	serviceInfoSlice := make([]ServiceInfo, 0)

	for k, v := range services {

		tmp := v.Teams

		if len(tmp) != 0 {
			re := regexp.MustCompile("_SVC+")

			if re.MatchString(services[k].Name) {
				serviceInfo = ServiceInfo{
					Compliant:   1,
					ServiceName: services[k].Name,
					ServiceTeam: services[k].Teams[0].APIObject.Summary,
				}

				serviceInfoSlice = append(serviceInfoSlice, serviceInfo)
			} else {
				serviceInfo = ServiceInfo{
					Compliant:   0,
					ServiceName: services[k].Name,
					ServiceTeam: services[k].Teams[0].APIObject.Summary,
				}

				serviceInfoSlice = append(serviceInfoSlice, serviceInfo)
			}
		}
	}

	return serviceInfoSlice
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

//func callPagerdutyApiIncidents() {
//	var opts pagerduty.ListIncidentsOptions
//	var APIList pagerduty.APIListObject
//	var Incidents []pagerduty.Incident
//	for {
//		eps, err := client.ListIncidents(opts)
//		if err != nil {
//			panic(err)
//		}
//
//		Incidents = append(Incidents, eps.Incidents...)
//		APIList.Offset += 25
//		opts = pagerduty.ListIncidentsOptions{APIListObject: APIList}
//
//		if eps.More != true {
//			break
//		}
//	}
//
//	for k := range Incidents {
//		fmt.Println(Incidents[k].Service.ID)
//	}
//}

type AnalyticsData struct {
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

type AnalyticsResponse struct {
	Data            []AnalyticsData  `json:"data,omitempty"`
	AnalyticsFilter *AnalyticsFilter `json:"filters,omitempty"`
	AggregateUnit   string           `json:"aggregate_unit,omitempty"`
	TimeZone        string           `json:"time_zone,omitempty"`
}

type AnalyticsFilter struct {
	CreatedAtStart string   `json:"created_at_start,omitempty"`
	CreatedAtEnd   string   `json:"created_at_end,omitempty"`
	Urgency        string   `json:"urgency,omitempty"`
	Major          bool     `json:"major,omitempty"`
	ServiceIDs     []string `json:"service_ids,omitempty"`
	TeamIDs        []string `json:"team_ids,omitempty"`
	PriorityIDs    []string `json:"priority_ids,omitempty"`
	PriorityName   []string `json:"priority_name,omitempty"`
}

func callPagerDutyApiAnalytics(serviceIds []string) map[string][]int {
	url := "https://api.pagerduty.com/analytics/metrics/incidents/services"
	endTime := time.Now().Format(time.RFC3339)
	startTime := time.Now().AddDate(-1, 0, 0).Format(time.RFC3339)
	chunks := make([][]string, 0)
	chunkSize := 10

	for i := 0; i < len(serviceIds); i += chunkSize {
		end := i + chunkSize
		if end > len(serviceIds) {
			end = len(serviceIds)
		}

		chunks = append(chunks, serviceIds[i:end])
	}

	var analyticsResponse AnalyticsResponse
	arrOfThings := make([]int, 0)
	mapOfThings := make(map[string][]int, 0)

	for _, chunk := range chunks {

		formattedPayloadString := fmt.Sprintf("{\"filters\":{\"created_at_start\":\"%s\",\"created_at_end\":\"%s\",\"service_ids\":[\"%s\"]}}", startTime, endTime, strings.Join(chunk, "\", \""))
		payload := strings.NewReader(formattedPayloadString)
		req, err := http.NewRequest("POST", url, payload)
		if err != nil {
			panic(err)
		}

		req.Header.Add("X-EARLY-ACCESS", "analytics-v2")
		req.Header.Add("Accept", "application/vnd.pagerduty+json;version=2")
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", "Token token="+authToken)

		resp, _ := http.DefaultClient.Do(req)
		err = decodeJSON(resp, &analyticsResponse)
		tmp := analyticsResponse

		for _, v := range tmp.Data {
			arrOfThings = nil
			arrOfThings = append(arrOfThings, v.MeanSecondsToFirstAck, v.MeanSecondsToResolve)
			mapOfThings[v.ServiceName] = arrOfThings
		}

	}
	return mapOfThings
}

func decodeJSON(resp *http.Response, payload interface{}) error {
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(payload)
}
