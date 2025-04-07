package notify

import (
	"context"
	"fmt"
	"sync"
	"time"

	_ "github.com/gogf/gf/contrib/nosql/redis/v2"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type NotifyItem struct {
	AgentId   string `v:"required" json:"agentId"   dc:"agentId"`
	AgentName string `v:"required" json:"agentName" dc:"agentName"`
	// JobId     int    `v:"required" json:"jobId"     dc:"jobId"`
	// JobStatus string `v:"required" json:"jobStatus" dc:"jobStatus"`
	// JobOutput string `v:"required" json:"jobOutput" dc:"jobOutput"`
}

type NotifyReq struct {
	g.Meta     `path:"/notifys/v1"  method:"post" tags:"Notify" summary:"Notify long polling"`
	Items      []NotifyItem `v:"required" json:"items" dc:"Notification items"`
	TimeoutSec int          `json:"timeoutSec" dc:"Timeout in seconds, default 30 seconds"`
}

type Notify struct{}

var (
	// Store agent notification status
	agentNotifications = make(map[string]string)
	// Store channels for waiting requests
	notificationChannelsMap = make(map[string]chan string)
	mutex                   = sync.RWMutex{}
)

func AddNotification(agentId, jobId string) {
	mutex.Lock()
	defer mutex.Unlock()

	// Update notification status
	agentNotifications[agentId] = jobId

	// Notify waiting channels
	if ch, exists := notificationChannelsMap[agentId]; exists {
		select {
		case ch <- jobId: // Non-blocking send
		default:
		}
		// Clear the notified channel
		delete(notificationChannelsMap, agentId)
	}
}

func (Notify) NotifyV1(ctx context.Context, req *NotifyReq) (res *ghttp.Response, err error) {
	r := g.RequestFromCtx(ctx)

	// Parameter validation
	if err = validateRequest(req); err != nil {
		r.Response.WriteStatusExit(400, err.Error())
		return
	}

	// Set timeout duration
	timeout := getTimeout(req.TimeoutSec)
	timeoutDuration := time.Duration(timeout) * time.Second

	// Collect all agentIds
	agentIds := collectAgentIds(req.Items)

	// Step 1: Check immediately if notifications already exist
	result, found, err := checkExistingNotifications(ctx, r, agentIds)
	if err != nil || found {
		return result, err
	}

	// Step 2: Start long polling
	return performLongPolling(ctx, r, agentIds, timeoutDuration)
}

func validateRequest(req *NotifyReq) error {
	if req == nil || len(req.Items) == 0 {
		return fmt.Errorf("invalid request")
	}
	return nil
}

func getTimeout(requestedTimeout int) int {
	if requestedTimeout > 0 {
		return requestedTimeout
	}
	return 30
}

func collectAgentIds(items []NotifyItem) []string {
	agentIds := make([]string, 0, len(items))
	for _, item := range items {
		agentIds = append(agentIds, item.AgentId)
	}
	return agentIds
}

func checkExistingNotifications(ctx context.Context, r *ghttp.Request, agentIds []string) (*ghttp.Response, bool, error) {
	mutex.RLock()
	defer mutex.RUnlock()

	for _, agentId := range agentIds {
		// Check if there are existing notifications
		if jobId, exists := agentNotifications[agentId]; exists {
			r.Response.WriteStatus(200)
			r.Response.WriteJson(g.Map{
				"code":    0,
				"message": "Notification found",
				"data": g.Map{
					"agentId": agentId,
					"jobId":   jobId,
				},
			})
			delete(agentNotifications, agentId)
			return nil, true, nil
		}

		// Check if agent exists
		ciAgentKey := "ciagent:" + agentId
		count, redisErr := g.Redis().Exists(ctx, ciAgentKey)
		if redisErr != nil {
			g.Log().Error(ctx, "Redis error:", redisErr)
			r.Response.WriteStatus(500)
			return nil, true, fmt.Errorf("redis exists operation failed: %w", redisErr)
		}

		if count == 0 {
			r.Response.WriteStatus(404)
			r.Response.WriteJson(g.Map{
				"code":    0,
				"message": "agentId Not Found",
				"data": g.Map{
					"agentId": agentId,
				},
			})
			return nil, true, nil
		}

	}

	return nil, false, nil
}

func performLongPolling(ctx context.Context, r *ghttp.Request, agentIds []string, timeoutDuration time.Duration) (*ghttp.Response, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, timeoutDuration)
	defer cancel()

	done := make(chan struct {
		agentId string
		jobId   string
	}, 1)

	// Register listening channels
	registerNotificationChannels(agentIds)
	// Cleanup function
	defer cleanupNotificationChannels(agentIds)

	// Start listener goroutines
	startListenerGoroutines(ctxTimeout, agentIds, done)

	// Wait for results
	select {
	case notification := <-done:
		r.Response.WriteJson(g.Map{
			"code":    0,
			"message": "Notification received",
			"data": g.Map{
				"agentId": notification.agentId,
				"jobId":   notification.jobId,
			},
		})
	case <-ctxTimeout.Done():
		r.Response.WriteStatus(304)
	}

	return nil, nil
}

func registerNotificationChannels(agentIds []string) {
	mutex.Lock()
	defer mutex.Unlock()

	for _, agentId := range agentIds {
		if _, exists := notificationChannelsMap[agentId]; !exists {
			notificationChannelsMap[agentId] = make(chan string, 1)
		}
	}
}

func cleanupNotificationChannels(agentIds []string) {
	mutex.Lock()
	defer mutex.Unlock()

	for _, agentId := range agentIds {
		delete(notificationChannelsMap, agentId)
	}
}

func startListenerGoroutines(ctxTimeout context.Context, agentIds []string, done chan<- struct {
	agentId string
	jobId   string
}) {
	for _, agentId := range agentIds {
		go func(id string) {
			select {
			case jobId := <-notificationChannelsMap[id]:
				select {
				case done <- struct {
					agentId string
					jobId   string
				}{agentId: id, jobId: jobId}:
				default:
				}
			case <-ctxTimeout.Done():
			}
		}(agentId)
	}
}
