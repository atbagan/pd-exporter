package collector

import (
	"fmt"
	"github.com/PagerDuty/go-pagerduty"
	"github.com/prometheus/client_golang/prometheus"
	"regexp"
)

type ServiceInfo struct {
	Compliant   int
	ServiceName string
	ServiceTeam string
}

type MyCollector struct {
	totalGaugeDesc            *prometheus.Desc
	businessServicesGaugeDesc *prometheus.Desc
	complianceGaugeDesc       *prometheus.Desc
}

func NewServiceCollector() *MyCollector {
	return &MyCollector{
		totalGaugeDesc:            prometheus.NewDesc("pagerduty_total_services_metric", "The number of total services in AIpagerduty", nil, nil),
		businessServicesGaugeDesc: prometheus.NewDesc("pagerduty_total_business_services_metric", "Shows the total number of business services", nil, nil),
		complianceGaugeDesc:       prometheus.NewDesc("pagerduty_total_services_compliant_metric", "Shows the number of compliant services names", nil, nil),
	}
}

func (c *MyCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.totalGaugeDesc
	ch <- c.businessServicesGaugeDesc
	ch <- c.complianceGaugeDesc
}

func (c *MyCollector) Collect(ch chan<- prometheus.Metric) {
	var pagerdutyServices = pdServices()
	total := totalServices(pagerdutyServices)
	compliance := compliantServices(pagerdutyServices)
	//serviceIds := getCompliantServiceIds(pagerdutyServices)
	serviceInfoSlice := serviceComplianceInfoWithNameAndTeam(pagerdutyServices)
	businessServices := callPagerdutyApiBusinessServices()

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
