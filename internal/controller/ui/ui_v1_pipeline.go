package ui

import (
	"cicdgf/internal/service"
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func (c *ControllerV1) PipelineGetList(ctx context.Context, req *PipelineGetListReq) (response *ghttp.Response, err error) {
	r := g.RequestFromCtx(ctx)

	pipelines, err := service.Pipeline.GetListPipelines(ctx)
	if err != nil {
		return nil, err
	}

	err = r.Response.WriteTpl("pipelines/list.html", g.Map{
		"url":            "/pipelines/",
		"pipelines":      pipelines,
		"newPipelineUrl": "/pipelinenew",
	})
	return nil, err
}
