package service

import (
	"cicdgf/internal/dao"
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
)

var Pipeline = pipelineService{}

type pipelineService struct{}

type ListPipelines struct {
	Id            int    `json:"pipeline_id"`
	Pipeline_name string `json:"pipeline_name"`
}

func (s *pipelineService) GetListPipelines(ctx context.Context) (pipelines []ListPipelines, err error) {
	if err = dao.CicdPipeline.Ctx(ctx).
		Fields("id,pipeline_name").
		Scan(&pipelines); err != nil {
		return nil, gerror.Wrap(err, "get pipelines failed")
	}

	return
}
