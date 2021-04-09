[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Docker Pulls](https://img.shields.io/docker/pulls/agtbagan/pagerduty-exporter.svg?maxAge=604800)]
# Welcome to the Pagerduty Exporter 

------------------------------------------------------------------------------------------------------------------------
### Table of Contents
<!-- TOC -->
- [Welcome](#welcome-to-the-pagerduty-exporter)
  - [Contributing](#contributing)
    - [Developer Workflow](#developer-workflow)
  - [Maintenance](#maintenance)
  - [Environment](#environment-variables)
  - [Reasoning](#reasoning)
  - [Getting Started](#getting-started)
  - [Configuration](#configuration)
------------------------------------------------------------------------------------------------------------------------
## Contributing

### Developer Workflow

```

In order to properly use this repository you will need to use a standard feature branch workflow.
This can be explained further in the link at the bottom of this section. The general work flow is the following:

clone the project
$ git clone https://github.com/atbagan/pd-exporter.git

create a branch to work on
$ git checkout -b <your_branch_name>

write your code, make changes, etc..
$ git commit -am "committing my branch"

push your branch to this repo 
$ git push origin <your_branch_name>

Create a merge request (MR) for your recently pushed branch

Wait for review of your MR



```
[`basic branch workflow`](https://docs.gitlab.com/ee/gitlab-basics/feature_branch_workflow.html)

------------------------------------------------------------------------------------------------------------------------
## Maintenance

## Environment Variables
| ENVIRONMENT VARIABLES   | Introduced in Version | Description | Default     |
| --------                | --------------------- | ----------- | ----------- |
| AUTH_TOKEN              | 0.1                   | Api Token   | None        |
| WEB_LISTEN_ADDRESS      | 0.1                   |  Address to listen on for web server | 9696 |
| WEB_TELEMETRY_PATH      | 0.1                   |  Path where to expose metrics        | /metrics |
| PD_ANALYTICS_SETTINGS   | 0.1                   |  Pagerduty Analytics Metrics Settings on/off (boolean)| false |
| PD_SERVICES_SETTINGS    | 0.1                   |  Pagerduty Services Metrics Settings on/off (boolean)| false |
| PD_TEAMS_SETTINGS       | 0.1                   |  Pagerduty Teams Metrics Settings on/off (boolean)| false |
| PD_USERS_SETTINGS       | 0.1                   |  Pagerduty Users Metrics Settings on/off (boolean)| false |

## Reasoning

## Getting Started 
`$ docker build -t pd-exporter .`

`$ docker run -e AUTH_TOKEN=your-api-key-here -dp 9696:9696 pd-exporter`

after a few seconds check: `http://localhost:9696/metrics`

I am currently running this in `ECS`

I am currently working on the `service file` for non containerized deployments.
`Docker` is the only deployment method currently.

### Configuration

| Argument                | Introduced in Version | Description | Default     |
| --------                | --------------------- | ----------- | ----------- |
| web.listen-address      | 0.1                   |  Address to listen on for web server | 9696 |
| web.telemetry-path      | 0.1                   |  Path where to expose metrics        | /metrics |
| pd.analytics_settings   | 0.1                   |  Pagerduty Analytics Metrics Settings on/off (boolean)| false |
| pd.services_settings    | 0.1                   |  Pagerduty Services Metrics Settings on/off (boolean)| false |
| pd.teams_settings       | 0.1                   |  Pagerduty Teams Metrics Settings on/off (boolean)| false |
| pd.users_settings       | 0.1                   |  Pagerduty Users Metrics Settings on/off (boolean)| false |
