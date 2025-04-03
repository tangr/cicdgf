package notify

import (
	"context"
	"sync"
	"time"

	_ "github.com/gogf/gf/contrib/nosql/redis/v2"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type NotifyItem struct {
	AgentId   int    `v:"required" json:"agentId"   dc:"agentId"`
	AgentName string `v:"required" json:"agentName" dc:"agentName"`
	JobId     int    `v:"required" json:"jobId"     dc:"jobId"`
	JobStatus string `v:"required" json:"jobStatus" dc:"jobStatus"`
	JobOutput string `v:"required" json:"jobOutput" dc:"jobOutput"`
}

type NotifyReq struct {
	g.Meta     `path:"/v1"  method:"post" tags:"Notify" summary:"Notify long polling"`
	Items      []NotifyItem `v:"required" json:"items" dc:"Notification items"`
	TimeoutSec int          `json:"timeoutSec" dc:"Timeout in seconds, default 30 seconds"`
}

// Global variables for storing notifications and handling long polling
var (
	notificationChannelsMap = make(map[string]chan []NotifyItem)
	mutex                   = sync.RWMutex{}
)

type Notify struct{}

// Adds notification items to the corresponding channel
func AddNotificationItems(clientId string, items []NotifyItem) {
	mutex.RLock()
	ch, exists := notificationChannelsMap[clientId]
	mutex.RUnlock()

	if exists {
		// Non-blocking send to avoid issues when client disconnects but channel is not closed
		select {
		case ch <- items:
			// Sent successfully
		default:
			// Channel is full or closed, ignore
		}
	}
}

func (Notify) NotifyV1(ctx context.Context, req *NotifyReq) (res *ghttp.Response, err error) {
	r := g.RequestFromCtx(ctx)

	if req == nil {
		r.Response.WriteStatusExit(400, "Invalid request: request body is required")
		return
	}

	// Set default timeout to 30 seconds
	timeout := 30
	if req.TimeoutSec > 0 {
		timeout = req.TimeoutSec
	}

	// Get client ID from header, fallback to IP if not present
	clientId := r.GetHeader("X-Client-ID")
	if clientId == "" {
		clientId = r.GetClientIp()
	}
	g.Log().Infof(ctx, "clientId: %s", clientId)

	if len(req.Items) > 0 {
		for i, item := range req.Items {
			g.Log().Infof(ctx, "Processing item #%d: Agent=%s(%d), JobId=%d, JobStatus=%s",
				i, item.AgentName, item.AgentId, item.JobId, item.JobStatus)
		}

		AddNotificationItems(clientId, req.Items)
	}

	// Long polling implementation

	// Create notification channel
	tmpnotificationChan := make(chan []NotifyItem, 3)

	// Register channel in global map
	mutex.Lock()
	notificationChannelsMap[clientId] = tmpnotificationChan
	mutex.Unlock()

	// Ensure cleanup on function exit
	defer func() {
		mutex.Lock()
		delete(notificationChannelsMap, clientId)
		mutex.Unlock()
		close(tmpnotificationChan)
	}()

	// Create timeout context
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	// Wait for data or timeout
	select {
	case items := <-tmpnotificationChan:
		// Received notifications, return data
		response := g.Map{
			"code":    0,
			"message": "New notifications received",
			"data": g.Map{
				"items":          items,
				"processedCount": len(items),
			},
		}
		r.Response.WriteStatus(200)
		r.Response.WriteJson(response)

	case <-timeoutCtx.Done():
		// Timeout occurred, return empty response
		g.Log().Infof(ctx, "status code: %d", 304)
		r.Response.WriteStatus(304)
	}

	return
}
