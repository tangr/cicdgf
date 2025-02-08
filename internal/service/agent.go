package service

import (
	"cicdgf/internal/dao"
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
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
