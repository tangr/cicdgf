package ui

import (
	"cicdgf/internal/service"
	"context"
	"fmt"

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

func (c *ControllerV1) AgentNew(ctx context.Context, req *AgentNewReq) (response *ghttp.Response, err error) {
	r := g.RequestFromCtx(ctx)

	err = r.Response.WriteTpl("agents/new.html", g.Map{
		"url":         "/agents",
		"newAgentUrl": "/agents",
	})
	return nil, err
}

func (c *ControllerV1) AgentCreate(ctx context.Context, req *AgentCreateReq) (response *ghttp.Response, err error) {
	r := g.RequestFromCtx(ctx)

	var agent_name string = r.Get("agent_name").String()
	var agent_ipaddr string = r.Get("agent_ipaddr").String()
	agent_id := service.Agent.New(agent_name, agent_ipaddr)

	r.Response.RedirectTo("/agents/"+fmt.Sprint(agent_id), 303)

	return nil, err
}
