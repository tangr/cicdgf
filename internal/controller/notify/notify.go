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
	g.Meta     `path:"/v1"  method:"post" tags:"Notify" summary:"Notify long polling"`
	Items      []NotifyItem `v:"required" json:"items" dc:"Notification items"`
	TimeoutSec int          `json:"timeoutSec" dc:"Timeout in seconds, default 30 seconds"`
}

type Notify struct{}

var (
	// 存储agent通知状态
	agentNotifications = make(map[string]string)
	// 存储等待中的请求通道
	notificationChannelsMap = make(map[string]chan string)
	mutex                   = sync.RWMutex{}
)

func AddNotification(agentId, jobId string) {
	mutex.Lock()
	defer mutex.Unlock()

	// 更新通知状态
	agentNotifications[agentId] = jobId

	// 通知等待的通道
	if ch, exists := notificationChannelsMap[agentId]; exists {
		select {
		case ch <- jobId: // 非阻塞发送
		default:
		}
		// 清空已通知的通道
		delete(notificationChannelsMap, agentId)
	}
}

func (Notify) NotifyV1(ctx context.Context, req *NotifyReq) (res *ghttp.Response, err error) {
	r := g.RequestFromCtx(ctx)

	// 参数校验
	if req == nil || len(req.Items) == 0 {
		r.Response.WriteStatusExit(400, "Invalid request")
		return
	}

	// 设置超时时间
	timeout := 30
	if req.TimeoutSec > 0 {
		timeout = req.TimeoutSec
	}
	timeoutDuration := time.Duration(timeout) * time.Second

	// 收集所有agentId
	agentIds := make([]string, 0, len(req.Items))
	for _, item := range req.Items {
		agentIds = append(agentIds, item.AgentId)
	}

	// 第一步：立即检查是否存在已有通知
	mutex.RLock()
	for _, agentId := range agentIds {
		ciAgentKey := "ciagent:" + agentId
		count, redisErr := g.Redis().Exists(ctx, ciAgentKey)
		if redisErr != nil {
			mutex.RUnlock()
			g.Log().Error(ctx, "Redis error:", redisErr)
			err = redisErr
			r.Response.WriteStatus(500)
			return nil, fmt.Errorf("redis exists operation failed: %w", err)
		}
		if count == 0 {
			mutex.RUnlock()
			r.Response.WriteStatus(404)
			r.Response.WriteJson(g.Map{
				"code":    0,
				"message": "agentId Not Found",
				"data": g.Map{
					"agentId": agentId,
				},
			})
			return
		}
		if jobId, exists := agentNotifications[agentId]; exists {
			mutex.RUnlock()
			r.Response.WriteJson(g.Map{
				"code":    0,
				"message": "Notification found",
				"data": g.Map{
					"agentId": agentId,
					"jobId":   jobId,
				},
			})
			return
		}
	}
	mutex.RUnlock()

	// 第二步：进入长轮询
	ctxTimeout, cancel := context.WithTimeout(ctx, timeoutDuration)
	defer cancel()

	done := make(chan struct {
		agentId string
		jobId   string
	}, 1)

	// 注册监听通道
	mutex.Lock()
	for _, agentId := range agentIds {
		if _, exists := notificationChannelsMap[agentId]; !exists {
			notificationChannelsMap[agentId] = make(chan string, 1)
		}
	}
	mutex.Unlock()

	// 清理函数
	defer func() {
		mutex.Lock()
		defer mutex.Unlock()
		for _, agentId := range agentIds {
			delete(notificationChannelsMap, agentId)
		}
	}()

	// 启动监听goroutine
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

	// 等待结果
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

	return
}
