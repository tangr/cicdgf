package notify

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// type NotifyReq struct {
// 	g.Meta    `path:"/v1" method:"get" tags:"Notify" summary:"Notify long polling"`
// 	AgentId   int    `v:"required" json:"agentId"   dc:"agentId"`
// 	AgentName string `v:"required" json:"agentName" dc:"agentName"`
// 	JobId     int    `v:"required" json:"jobId"     dc:"jobId"`
// 	JobStatus string `v:"required" json:"jobStatus" dc:"jobStatus"`
// 	JobOutput string `v:"required" json:"jobOutput" dc:"jobOutput"`
// }

type NotifyItem struct {
	AgentId   int    `v:"required" json:"agentId"   dc:"agentId"`
	AgentName string `v:"required" json:"agentName" dc:"agentName"`
	JobId     int    `v:"required" json:"jobId"     dc:"jobId"`
	JobStatus string `v:"required" json:"jobStatus" dc:"jobStatus"`
	JobOutput string `v:"required" json:"jobOutput" dc:"jobOutput"`
}

type NotifyReq struct {
	g.Meta `path:"/v1" method:"post" tags:"Notify" summary:"Notify long polling"`
	Items  []NotifyItem `v:"required" json:"items" dc:"notification items"`
}

// type NotifyRes struct {
// 	g.Meta    `path:"/" method:"get" tags:"Test" summary:"Notify test case"`
// 	AgentId   int    `json:"agentId"   dc:"agentId"`
// 	AgentName string `json:"agentName" dc:"agentName"`
// 	JobId     int    `json:"jobId"     dc:"jobId"`
// 	JobStatus string `json:"jobStatus" dc:"jobStatus"`

// 	Body   string            `json:"scriptBody"`
// 	Envs   map[string]string `json:"scriptEnvs"`
// 	Args   string            `json:"scriptArgs"`
// 	ErrMsg string            `json:"errmsg"`
// }

type Notify struct{}

func (Notify) NotifyV1(ctx context.Context, req *NotifyReq) (res *ghttp.Response, err error) {
	r := g.RequestFromCtx(ctx)

	if req == nil || len(req.Items) == 0 {
		r.Response.WriteStatusExit(400, "Invalid request: items are required")
		return
	}

	for i, item := range req.Items {
		g.Log().Infof(ctx, "Processing item #%d: Agent=%s(%d), JobId=%d, JobStatus=%s",
			i, item.AgentName, item.AgentId, item.JobId, item.JobStatus)

		// 在这里可以根据需求处理每个通知项
		// 例如：保存到数据库、触发事件等
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
