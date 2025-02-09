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

func (c *ControllerV1) AgentGetOne(ctx context.Context, req *AgentGetOneReq) (response *ghttp.Response, err error) {
	r := g.RequestFromCtx(ctx)

	agent_id := r.Get("id").String()
	agent, err := service.Agent.GetAgent(agent_id)
	if err != nil {
		return nil, err
	}

	g.Log().Debug(ctx, "agent: ", agent)
	g.Log().Debugf(ctx, "AgentGetOne agent_name: %s, agent_ipaddr: %s", agent.Agent_name, agent.Ipaddr)

	err = r.Response.WriteTpl("agents/edit.html", g.Map{
		"url":          "/agents/",
		"apiurl":       "/agents/" + agent_id + "/put",
		"agent_name":   agent.Agent_name,
		"agent_ipaddr": agent.Ipaddr,
	})
	return nil, err
}

func (c *ControllerV1) AgentUpdate(ctx context.Context, req *AgentUpdateReq) (response *ghttp.Response, err error) {
	r := g.RequestFromCtx(ctx)

	var agent_id string = r.Get("id").String()
	var agent_name string = r.Get("agent_name").String()
	var agent_ipaddr string = r.Get("agent_ipaddr").String()

	g.Log().Debugf(ctx, "AgentUpdate agent_id: %s, agent_name: %s, agent_ipaddr: %s", agent_id, agent_name, agent_ipaddr)

	_ = service.Agent.Update(agent_id, agent_name, agent_ipaddr)
	r.Response.RedirectTo("/agents/"+fmt.Sprint(agent_id), 303)

	return nil, err
}
