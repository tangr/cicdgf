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
	notificationChannelsMap = make(map[string][]chan string)
	mutex                   = sync.RWMutex{}
)

func AddNotification(agentId, jobId string) {
	mutex.Lock()
	defer mutex.Unlock()

	// 更新通知状态
	agentNotifications[agentId] = jobId

	// 通知所有等待的通道
	if channels, exists := notificationChannelsMap[agentId]; exists {
		for _, ch := range channels {
			select {
			case ch <- jobId: // 非阻塞发送
			default:
			}
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

	done := make(chan string, 1)
	var channels []chan string

	// 注册监听通道
	mutex.Lock()
	for _, agentId := range agentIds {
		ch := make(chan string, 1)
		channels = append(channels, ch)
		notificationChannelsMap[agentId] = append(notificationChannelsMap[agentId], ch)
	}
	mutex.Unlock()

	// 清理函数
	defer func() {
		mutex.Lock()
		defer mutex.Unlock()
		for i, agentId := range agentIds {
			// 从通知通道列表中移除
			remaining := make([]chan string, 0)
			for _, c := range notificationChannelsMap[agentId] {
				if c != channels[i] {
					remaining = append(remaining, c)
				}
			}
			if len(remaining) > 0 {
				notificationChannelsMap[agentId] = remaining
			} else {
				delete(notificationChannelsMap, agentId)
			}
			close(channels[i])
		}
	}()

	// 启动监听goroutine
	go func() {
		for _, ch := range channels {
			go func(c <-chan string) {
				select {
				case jobId := <-c:
					select {
					case done <- jobId:
					default:
					}
				case <-ctxTimeout.Done():
				}
			}(ch)
		}
	}()

	// 等待结果
	select {
	case jobId := <-done:
		r.Response.WriteJson(g.Map{
			"code":    0,
			"message": "Notification received",
			"data": g.Map{
				// "agentId": agentId,
				"jobId": jobId,
			},
		})
	case <-ctxTimeout.Done():
		r.Response.WriteStatus(304)
	}

	return
}
