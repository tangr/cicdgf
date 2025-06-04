package ui

import (
	"cicdgf/internal/service"
	"context"
	"fmt"

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

func (c *ControllerV1) PipelineNew(ctx context.Context, req *PipelineNewReq) (response *ghttp.Response, err error) {
	r := g.RequestFromCtx(ctx)

	agents, err := service.Agent.GetListAgents(ctx)
	if err != nil {
		return nil, err
	}

	groups, err := service.Group.GetListGroups(ctx)
	if err != nil {
		return nil, err
	}

	err = r.Response.WriteTpl("pipelines/new.html", g.Map{
		"url":            "/pipelines/",
		"groups":         groups,
		"agents":         agents,
		"newPipelineUrl": "/pipelines/",
	})
	return nil, err
}

func (c *ControllerV1) PipelineCreate(ctx context.Context, req *PipelineCreateReq) (response *ghttp.Response, err error) {
	r := g.RequestFromCtx(ctx)

	var pipeline_name string = r.Get("pipeline_name").String()
	var group_id int = r.Get("group_id").Int()
	var agent_id int = r.Get("agent_id").Int()
	var concurrency int = r.Get("concurrency").Int()
	var pipeline_body string = r.Get("pipeline_body").String()

	pipeline_id := service.Pipeline.New(pipeline_name, group_id, agent_id, concurrency, pipeline_body)

	r.Response.RedirectTo("/pipelines/"+fmt.Sprint(pipeline_id), 303)

	return nil, err
}

func (c *ControllerV1) PipelineGetOne(ctx context.Context, req *PipelineGetOneReq) (response *ghttp.Response, err error) {
	r := g.RequestFromCtx(ctx)

	pipeline_id := r.Get("id").Int()

	pipeline, err := service.Pipeline.GetOne(pipeline_id)
	if err != nil {
		return nil, err
	}

	all_agents, err := service.Agent.GetListAgents(ctx)
	if err != nil {
		return nil, err
	}

	all_groups, err := service.Group.GetListGroups(ctx)
	if err != nil {
		return nil, err
	}

	g.Log().Debug(ctx, "pipeline: ", pipeline)
	g.Log().Debug(ctx, "pipeline_body: ", pipeline.Body)

	err = r.Response.WriteTpl("pipelines/edit.html", g.Map{
		"url":            "/pipelines/",
		"apiurl":         "/pipelines/" + fmt.Sprint(pipeline_id) + "/put",
		"pipeline_name":  pipeline.Pipeline_name,
		"pipeline_id":    pipeline_id,
		"all_agents":     all_agents,
		"all_groups":     all_groups,
		"pipeline_group": pipeline.Group_id,
		"pipeline_agent": pipeline.Agent_id,
		"pipeline_body":  pipeline.Body,
	})
	return nil, err
}

func (c *ControllerV1) PipelineUpdate(ctx context.Context, req *PipelineUpdateReq) (response *ghttp.Response, err error) {
	r := g.RequestFromCtx(ctx)

	var pipeline_id string = r.Get("id").String()
	var group_id string = r.Get("group_id").String()
	var agent_id string = r.Get("agent_id").String()
	var pipeline_body string = r.Get("pipeline_body").String()

	g.Log().Debugf(ctx, "PipelineUpdate pipeline_id: %s, group_id: %s, agent_id: %s, pipeline_body: %s",
		pipeline_id, group_id, agent_id, pipeline_body)

	_ = service.Pipeline.Update(pipeline_id, group_id, agent_id, pipeline_body)
	r.Response.RedirectTo("/pipelines/"+fmt.Sprint(pipeline_id), 303)

	return nil, err
}
