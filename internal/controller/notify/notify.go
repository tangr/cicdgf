package notify

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type NotifyReq struct {
	g.Meta    `path:"/" method:"get" tags:"Test" summary:"Notify test case"`
	AgentId   int    `v:"required" json:"agentId"   dc:"agentId"`
	AgentName string `v:"required" json:"agentName" dc:"agentName"`
	JobId     int    `v:"required" json:"jobId"     dc:"jobId"`
	JobStatus string `v:"required" json:"jobStatus" dc:"jobStatus"`
	JobOutput string `v:"required" json:"jobOutput" dc:"jobOutput"`
}

type NotifyRes struct {
	g.Meta    `path:"/" method:"get" tags:"Test" summary:"Notify test case"`
	AgentId   int    `json:"agentId"   dc:"agentId"`
	AgentName string `json:"agentName" dc:"agentName"`
	JobId     int    `json:"jobId"     dc:"jobId"`
	JobStatus string `json:"jobStatus" dc:"jobStatus"`

	Body   string            `json:"scriptBody"`
	Envs   map[string]string `json:"scriptEnvs"`
	Args   string            `json:"scriptArgs"`
	ErrMsg string            `json:"errmsg"`
}

type Notify struct{}

func (Notify) NotifyV1(ctx context.Context, req *NotifyReq) (res *ghttp.Response, err error) {
	r := g.RequestFromCtx(ctx)

	r.Response.WriteExit("pipeline_body")

	return
}
