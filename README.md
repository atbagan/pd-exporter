[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Docker Pulls](https://img.shields.io/docker/pulls/agtbagan/pagerduty-exporter.svg?maxAge=604800)](https://hub.docker.com/r/agtbagan/pagerduty-exporter/)
[![Go Report Card](https://goreportcard.com/badge/github.com/atbagan/pd-exporter)](https://goreportcard.com/report/github.com/atbagan/pd-exporter)
# Welcome to the Pagerduty Exporter 

------------------------------------------------------------------------------------------------------------------------
### Table of Contents
<!-- TOC -->
- [Welcome](#welcome-to-the-pagerduty-exporter)
  - [Contributing](#contributing)
    - [Developer Workflow](#developer-workflow)
  - [Getting Started](#getting-started)
  - [Configuration](#configuration)
------------------------------------------------------------------------------------------------------------------------
## Contributing

I welcome any contributions. Please fork the project on GitHub and open
Pull Requests for any proposed changes.

Please note that I will not merge any changes that encourage insecure
behaviour. If in doubt please open an Issue first to discuss your proposal.

### Developer Workflow

```
In order to properly use this repository you will need to use a standard feature branch workflow.
This can be explained further in the link at the bottom of this section. 

```
[`basic branch workflow`](https://gist.github.com/Chaser324/ce0505fbed06b947d962)

------------------------------------------------------------------------------------------------------------------------

## Getting Started 
[Dockerhub](https://hub.docker.com/r/agtbagan/pagerduty-exporter)
`docker pull agtbagan/pagerduty-exporter:0.1.2`

Example `docker-compose.yml`:

`web.telemetry-path` is an example of how to use the command arguments

`PD_ANALYTICS_SETTINGS` example of how to use env vars
```yaml

pagerduty_exporter:
  image: agtbagan/pagerduty-exporter:0.1.2
  command:
    - '--web.telemetry-path=/example'
  environment:
    - AUTH_TOKEN=your_api_key
    - PD_ANALYTICS_SETTINGS=false
  restart: always
  ports:
    - "127.0.0.1:9696:9696"
```
Or, you can build and run it yourself like the following:

`$ docker build -t pd-exporter .`

`$ docker run -e AUTH_TOKEN=your-api-key-here -dp 9696:9696 pd-exporter`

after a few seconds check: `http://localhost:9696/metrics`

I am currently running this in `ECS`

### Configuration

| Argument                | Environment Variable  |Introduced in Version | Description | Default     |
| --------                | --------------------- | -----------          | ----------- | ----------- | 
| web.listen-address      |  WEB_LISTEN_ADDRESS   |   0.1                |  Address to listen on for web server                   | 9696 |
| web.telemetry-path      |  WEB_TELEMETRY_PATH   |   0.1                |  Path where to expose metrics                          | /metrics |
| pd.analytics_settings   |  PD_ANALYTICS_SETTINGS|   0.1                |  Pagerduty Analytics Metrics Settings on/off (boolean) | true |
| pd.services_settings    |  PD_SERVICES_SETTINGS |   0.1                |  Pagerduty Services Metrics Settings on/off (boolean)  | true |
| pd.teams_settings       |  PD_TEAMS_SETTINGS    |   0.1                |  Pagerduty Teams Metrics Settings on/off (boolean)     | true |
| pd.users_settings       |  PD_USERS_SETTINGS    |   0.1                |  Pagerduty Users Metrics Settings on/off (boolean)     | true |
| n/a                     |  AUTH_TOKEN           |   0.1                |  Pagerduty API Key  (required)                          | ""      |

**NOTE** The list of metrics is currently small and will continue to grow over time as I make time to work on this.

**NOTE** Metrics marked with an asterisk(*) are company specific and will likely hold no value unless we work for the same org.
In the future those metrics will be flaggable and will be off by default.

### Metrics
| Name                | Type                  |Help |
| --------            | --------------------- | ----------- |         
| pagerduty_total_users_metric          | Gauge                 | Total number of pagerduty users in your account |
| pagerduty_total_business_services_metric | Gauge             | Total number of business services in your account |
| pagerduty_total_teams_metric            | Gauge              | Total number of teams in your account |
| pagerduty_total_services_metric        | Gauge              | Total number of services in your account |
| pagerduty_mtta_analytics_metric*                  | Gauge              | MTTA for services with compliant naming convention* |
| pagerduty_mttr_analytics_metric*                 | Gauge              | MTTR for services with compliant naming convention* |
| pagerduty_uptime_percentage_analytics_metric* | Gauge | Uptime Percentage for services* |
| pagerduty_service_names_metric*  | Gauge              | Metric to check compliancy of service names. 0 for non 1 for compliant* |
| pagerduty_total_services_compliant_metric* | Gauge | total services compliant with standard naming convention* |
