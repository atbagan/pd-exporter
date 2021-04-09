package collector

import (
	"github.com/PagerDuty/go-pagerduty"
	"os"
)

var (
	client    = pagerduty.NewClient(authToken)
	authToken = os.Getenv("AUTH_TOKEN")
)
