package service

import (
	"cicdgf/internal/dao"
	"context"

	"github.com/gogf/gf/v2/encoding/gjson"
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
	Pipeline_name string       `json:"pipeline_name"`
	Group_id      string       `json:"group_id"`
	Agent_id      string       `json:"agent_id"`
	Concurrency   string       `json:"concurrency"`
	Body          PipelineBody `json:"body"`
}

type Script struct {
	Args   string `json:"script_args"`
	Script string `json:"script_name"`
}

type PipelineBody struct {
	StageCI Script `json:"stageCI"`
	StageCD Script `json:"stageCD"`
}

type JobScriptValue struct {
	Body string            `json:"scriptBody"`
	Envs map[string]string `json:"scriptEnvs"`
	Args string            `json:"scriptArgs"`
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

func (s *pipelineService) GetOne(pipeline_id int) (*PipelineDetail, error) {
	ctx := context.Background()

	record, err := dao.CicdPipeline.Ctx(ctx).
		Fields("pipeline_name,group_id,agent_id,concurrency,body").
		Where("id=?", pipeline_id).
		One()

	if err != nil {
		g.Log().Error(ctx, err)
		return nil, err
	}

	var pipelineBody PipelineBody
	if err := gjson.DecodeTo(record["body"].String(), &pipelineBody); err != nil {
		return nil, err
	}

	return &PipelineDetail{
		Pipeline_name: record["pipeline_name"].String(),
		Group_id:      record["group_id"].String(),
		Agent_id:      record["agent_id"].String(),
		Concurrency:   record["concurrency"].String(),
		Body:          pipelineBody,
	}, nil
}

func (s *pipelineService) Update(pipeline_id string, group_id string, agent_id string, pipeline_body string) string {
	ctx := context.Background()

	new_pipeline := g.Map{
		"pipeline_id": pipeline_id,
		"group_id":    group_id,
		"agent_id":    agent_id,
		"body":        pipeline_body,
	}
	g.Log().Debug(ctx, new_pipeline)

	_, err := dao.CicdPipeline.
		Ctx(ctx).
		Where("id=?", pipeline_id).
		Update(new_pipeline)
	if err != nil {
		g.Log().Error(ctx, err)
	}

	return pipeline_id
}

func (s *pipelineService) GetOnebody(pipeline_id int) (string, error) {
	ctx := context.Background()

	record, err := dao.CicdPipeline.Ctx(ctx).
		Fields("body").
		Where("id=?", pipeline_id).
		One()

	if err != nil {
		g.Log().Error(ctx, err)
		return "", err
	}

	return record["body"].String(), nil
}

func (s *pipelineService) GetOneGroupId(pipeline_id int) (string, error) {
	ctx := context.Background()

	record, err := dao.CicdPipeline.Ctx(ctx).
		Fields("group_id").
		Where("id=?", pipeline_id).
		One()

	if err != nil {
		g.Log().Error(ctx, err)
		return "", err
	}

	return record["group_id"].String(), nil
}
