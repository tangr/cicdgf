# GoFrame Template For SingleRepo

Quick Start:

- https://goframe.org/quick

```bash
time curl -X POST \
  -H "Content-Type: application/json" \
  "http://localhost:8000/api/notifys/v1" \
  -d '{
    "items": [{
      "agentId": 1000,
      "agentName": "agent-name-example",
      "jobId": 456,
      "jobStatus": "completed",
      "jobOutput": "job execution output here"
    }],
    "timeoutSec": 5
  }'

time curl -vX POST \
  -H "Content-Type: application/json" \
  "http://localhost:8000/api/notifys/v1" \
  -d '{
    "Items": [{
      "AgentId": 1000,
      "AgentName": "agent-name-example"
    }],
    "timeoutSec": 50
  }'

curl -vX POST \
  -H "Content-Type: " \
  "http://127.0.0.1:8000/api/notifys/v1" \
  -d '{
    "items": [{
      "agentId": 1000,
      "agentName": "agent-test-1"
    }],
    "timeoutSec": 3
  }'

curl -X POST \
  -H "Content-Type: application/json" \
  "http://localhost:8000/api/log?PipelineId=1000&AgentId=1000&JobType=BUILD&JobId=10001&TaskStatus=pending&Ipaddr=127.0.0.1&UpdatedAt=`date +%s`&Output=aaaa"

```

## Pipeline config sample

```json
{
  "stageCD": {
    "args": "",
    "envs": {
      "env1": ["env1"],
      "env2": ["env2"],
      "PKGRDM": ""
    },
    "script": "nginx-conf-cd"
  },
  "stageCI": {
    "args": "1111args",
    "envs": {
      "aaa": ["a", "b", "c"],
      "bbb": ["a", "b", "c"],
      "REPO_nginx_conf": "git|gitlab|144",
      "REPO_nginx_conf_ssl": "git|gitlab|144"
    },
    "script": "nginx-conf-ci"
  }
}
```
