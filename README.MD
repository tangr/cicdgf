# GoFrame Template For SingleRepo

Quick Start:

- https://goframe.org/quick

```bash
curl -X POST \
  -H "Content-Type: application/json" \
  "http://localhost:8000/notifys/v1" \
  -d 'items=[{
    "agentId": 123,
    "agentName": "agent-name-example",
    "jobId": 456,
    "jobStatus": "completed",
    "jobOutput": "job execution output here"
  }]'

```
