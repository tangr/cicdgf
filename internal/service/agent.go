package service

import (
	"cicdgf/internal/dao"
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

var Agent = agentService{}

type agentService struct{}

type ListAgents struct {
	Id         int    `json:"id"`
	Agent_name string `json:"agent_name"`
	Ipaddr     string `json:"ipaddr"`
	Updated_at int    `json:"updated_at"`
}

type AgentDetail struct {
	Agent_name string `json:"agent_name"`
	Ipaddr     string `json:"ipaddr"`
}

func (s *agentService) GetListAgents(ctx context.Context) (agents []ListAgents, err error) {
	if err = dao.CicdAgent.Ctx(ctx).
		Fields("id,agent_name,ipaddr,updated_at").
		Scan(&agents); err != nil {
		return nil, gerror.Wrap(err, "get agents failed")
	}

	return
}

func (s *agentService) New(agent_name string, agent_ipaddr string) int64 {
	ctx := context.Background()

	new_agent := g.Map{
		"agent_name": agent_name,
		"ipaddr":     agent_ipaddr,
	}

	result, err := dao.CicdAgent.Ctx(ctx).Insert(new_agent)
	if err != nil {
		g.Log().Error(ctx, err)
	}

	agent_id, err := result.LastInsertId()
	if err != nil {
		g.Log().Error(ctx, err)
	}

	return agent_id
}

func (s *agentService) GetAgent(agent_id string) (*AgentDetail, error) {
	ctx := context.Background()

	record, err := dao.CicdAgent.Ctx(ctx).
		Fields("agent_name, ipaddr").
		Where("id=?", agent_id).
		One()

	if err != nil {
		g.Log().Error(ctx, err)
		return nil, err
	}

	return &AgentDetail{
		Agent_name: record["agent_name"].String(),
		Ipaddr:     record["ipaddr"].String(),
	}, nil
}

func (s *agentService) Update(agent_id string, agent_name string, agent_ipaddr string) string {
	ctx := context.Background()

	new_agent := g.Map{
		"agent_name": agent_name,
		"ipaddr":     agent_ipaddr,
		"updated_at": gtime.Now().Timestamp(),
	}
	g.Log().Debug(ctx, "Update new_agent:", new_agent)
	g.Log().Debug(ctx, "Update agent_ipaddr:", agent_ipaddr)

	result, err := dao.CicdAgent.
		Ctx(ctx).
		Where("id=?", agent_id).
		Update(new_agent)
	if err != nil {
		g.Log().Error(ctx, err)
		g.Log().Error(ctx, result)
	}

	return agent_id
}
