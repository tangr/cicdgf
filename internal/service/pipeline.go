package service

import (
	"cicdgf/internal/dao"
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

var Pipeline = pipelineService{}

type pipelineService struct{}

type ListPipelines struct {
	Id            int    `json:"pipeline_id"`
	Pipeline_name string `json:"pipeline_name"`
}

type PipelineDetail struct {
	Pipeline_name string `json:"pipeline_name"`
	Agent_id      string `json:"agent_id"`
	Concurrency   string `json:"concurrency"`
	Body          string `json:"body"`
}

func (s *pipelineService) GetListPipelines(ctx context.Context) (pipelines []ListPipelines, err error) {
	if err = dao.CicdPipeline.Ctx(ctx).
		Fields("id,pipeline_name").
		Scan(&pipelines); err != nil {
		return nil, gerror.Wrap(err, "get pipelines failed")
	}

	return
}

func (s *pipelineService) New(pipeline_name string, group_id int, agent_id int, concurrency int, pipeline_body string) int64 {
	ctx := context.Background()

	new_pipeline := g.Map{
		"pipeline_name": pipeline_name,
		"group_id":      group_id,
		"agent_id":      agent_id,
		"concurrency":   concurrency,
		"pipeline_body": pipeline_body,
	}

	result, err := dao.CicdPipeline.Ctx(ctx).Insert(new_pipeline)
	if err != nil {
		g.Log().Error(ctx, err)
	}

	pipeline_id, err := result.LastInsertId()
	if err != nil {
		g.Log().Error(ctx, err)
	}

	return pipeline_id
}

func (s *pipelineService) GetOne(pipeline_id string) (*PipelineDetail, error) {
	ctx := context.Background()

	record, err := dao.CicdPipeline.Ctx(ctx).
		Fields("pipeline_name,agent_id,concurrency,body").
		Where("id=?", pipeline_id).
		One()

	if err != nil {
		g.Log().Error(ctx, err)
		return nil, err
	}

	return &PipelineDetail{
		Pipeline_name: record["pipeline_name"].String(),
		Agent_id:      record["agent_id"].String(),
		Concurrency:   record["concurrency"].String(),
		Body:          record["body"].String(),
	}, nil
}
