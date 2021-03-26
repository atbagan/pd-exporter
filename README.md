[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
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
------------------------------------------------------------------------------------------------------------------------
## Contributing

### Developer Workflow

```

In order to properly use this repository you will need to use a standard feature branch workflow.
This can be explained further in the link at the bottom of this section. The general work flow is the following:

clone the project
$ git clone https://gitlab.com/moneng/reusable/prom-exporters/pagerduty.git

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
`AUTH_TOKEN` : Pagerduty API Token

## Reasoning

## Getting Started 
`$ docker build -t pd-exporter .`

`$ docker run -e AUTH_TOKEN=your-api-key-here -dp 9798:9798 pd-exporter`

after a few seconds check: `http://localhost:9798/metrics`

I am currently running this in `ECS`

I am currently working on the `service file` for non containerized deployments.
`Docker` is the only deployment method currently.

