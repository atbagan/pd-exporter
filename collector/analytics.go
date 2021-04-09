package collector

import (
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"strings"
	"time"
)

// AnalyticsCollector is Analytics API Collector
type AnalyticsCollector struct {
	totalAnalytics *prometheus.Desc
}

// AnalyticsData is API Data structure
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

// AnalyticsResponse is API Response
type AnalyticsResponse struct {
	Data            []AnalyticsData  `json:"data,omitempty"`
	AnalyticsFilter *AnalyticsFilter `json:"filters,omitempty"`
	AggregateUnit   string           `json:"aggregate_unit,omitempty"`
	TimeZone        string           `json:"time_zone,omitempty"`
}

// AnalyticsFilter is Analytics API Filter
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

// new collector registered in main
func NewAnalyticsCollector() *AnalyticsCollector {
	return &AnalyticsCollector{
		totalAnalytics: prometheus.NewDesc("pagerduty_total_analytics_services_metric", "The number of total analytics in AIpagerduty", nil, nil),
	}
}

// describe channel for analytics
func (c *AnalyticsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.totalAnalytics
}

// collector for the analytics api
func (c *AnalyticsCollector) Collect(ch chan<- prometheus.Metric) {
	var pagerdutyServices = pdServices()
	serviceIds := getCompliantServiceIds(pagerdutyServices)

	mapOfMetrics := callPagerDutyApiAnalytics(serviceIds)

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
