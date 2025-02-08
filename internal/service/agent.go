package service

import (
	"cicdgf/internal/dao"
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

var Agent = agentService{}

type agentService struct{}

type ListAgents struct {
	Id         int    `json:"id"`
	Username   string `json:"username"`
	Groups     string `json:"groups"`
	Updated_at int    `json:"updated_at"`
}

func (s *agentService) GetListAgents(ctx context.Context) (users []ListAgents, err error) {
	if err = dao.CicdAgent.Ctx(ctx).
		Fields("id,agent_name,updated_at").
		Scan(&users); err != nil {
		return nil, gerror.Wrap(err, "get agents failed")
	}

	return
}

func (s *agentService) New(agent_name string, agent_ipaddr string) int64 {
	ctx := context.Background()

	new_agent := g.Map{
		"agent_name":   agent_name,
		"agent_ipaddr": agent_ipaddr,
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
