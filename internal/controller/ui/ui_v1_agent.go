package ui

import (
	"cicdgf/internal/service"
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func (c *ControllerV1) AgentGetList(ctx context.Context, req *AgentGetListReq) (response *ghttp.Response, err error) {
	r := g.RequestFromCtx(ctx)

	agents, err := service.Agent.GetListAgents(ctx)
	if err != nil {
		return nil, err
	}

	err = r.Response.WriteTpl("agents/list.html", g.Map{
		"url":         "/agents/",
		"agents":      agents,
		"newAgentUrl": "/agentnew",
	})
	return nil, err
}
