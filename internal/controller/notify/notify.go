package notify

import (
	"context"
	"sync"
	"time"

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
	Items      []NotifyItem `v:"required" json:"items" dc:"notification items"`
	TimeoutSec int          `json:"timeoutSec" dc:"超时时间(秒), 默认30秒"`
}

// 全局变量，用于存储通知和处理long polling
var (
	notificationChannels = make(map[string]chan []NotifyItem)
	mutex                = sync.RWMutex{}
)

type Notify struct{}

// AddNotificationItems 添加通知项到对应的channel
func AddNotificationItems(clientId string, items []NotifyItem) {
	mutex.RLock()
	ch, exists := notificationChannels[clientId]
	mutex.RUnlock()

	if exists {
		// 非阻塞方式发送，避免client已断开连接但channel未关闭的情况
		select {
		case ch <- items:
			// 发送成功
		default:
			// channel已满或已关闭，忽略
		}
	}
}

func (Notify) NotifyV1(ctx context.Context, req *NotifyReq) (res *ghttp.Response, err error) {
	r := g.RequestFromCtx(ctx)

	if req == nil {
		r.Response.WriteStatusExit(400, "Invalid request: request body is required")
		return
	}

	// 设置默认超时时间为30秒
	timeout := 30
	if req.TimeoutSec > 0 {
		timeout = req.TimeoutSec
	}

	// 如果请求中包含通知项，则立即处理
	if len(req.Items) == 0 {
		for i, item := range req.Items {
			g.Log().Infof(ctx, "Processing item #%d: Agent=%s(%d), JobId=%d, JobStatus=%s",
				i, item.AgentName, item.AgentId, item.JobId, item.JobStatus)
		}

		response := g.Map{
			"code":    0,
			"message": "Notifications processed successfully",
			"data": g.Map{
				"processedCount": len(req.Items),
			},
		}

		r.Response.WriteJson(response)
		return
	}

	// 以下是Long Polling实现部分
	// 从请求中获取客户端ID，如果没有则使用IP地址作为ID
	clientId := r.GetHeader("X-Client-ID")
	if clientId == "" {
		clientId = r.GetClientIp()
	}
	g.Log().Infof(ctx, "clientId: %s", clientId)

	// 创建通知通道
	notificationChan := make(chan []NotifyItem, 1)

	// 将通道注册到全局map
	mutex.Lock()
	notificationChannels[clientId] = notificationChan
	mutex.Unlock()

	// 确保在函数退出时清理资源
	defer func() {
		mutex.Lock()
		delete(notificationChannels, clientId)
		mutex.Unlock()
		close(notificationChan)
	}()

	// 设置超时上下文
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	// 等待数据或超时
	select {
	case items := <-notificationChan:
		// 收到通知，返回数据
		response := g.Map{
			"code":    0,
			"message": "New notifications received",
			"data": g.Map{
				"items":          items,
				"processedCount": len(items),
			},
		}
		r.Response.WriteJson(response)

	case <-timeoutCtx.Done():
		// 超时，返回空数据
		response := g.Map{
			"code":    0,
			"message": "No new notifications within timeout period",
			"data": g.Map{
				"items":          []NotifyItem{},
				"processedCount": 0,
			},
		}
		r.Response.WriteJson(response)
	}

	return
}
