package ui

import (
	"cicdgf/internal/service"
	"context"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func (c *ControllerV1) CicdGetList(ctx context.Context, req *CicdGetListReq) (response *ghttp.Response, err error) {
	r := g.RequestFromCtx(ctx)

	pipelines, err := service.Pipeline.GetListPipelines(ctx)
	if err != nil {
		return nil, err
	}

	err = r.Response.WriteTpl("cicd/list.html", g.Map{
		"url":            "/",
		"pipelines":      pipelines,
		"newPipelineUrl": "/pipelinenew",
	})
	return nil, err
}

func (c *ControllerV1) CicdGetOne(ctx context.Context, req *CicdGetOneReq) (response *ghttp.Response, err error) {
	r := g.RequestFromCtx(ctx)

	var pipeline_id int = r.Get("id").Int()
	// if !service.CheckAuthor(r.Context(), pipeline_id) {
	// 	r.Response.RedirectTo(UrlPrefix + "/forbidden")
	// }

	pipeline, err := service.Pipeline.GetOne(pipeline_id)
	if err != nil {
		return nil, err
	}

	pageid := r.Get("page").Int()
	jobs, totalSize, err := service.Cicd.GetListJobs(ctx, pipeline_id, pageid, 10)
	if err != nil {
		return nil, err
	}

	page := r.GetPage(totalSize, 10)

	err = r.Response.WriteTpl("cicd/show.html", g.Map{
		"url":           "/" + fmt.Sprint(pipeline_id),
		"apiurl":        "/v1/" + fmt.Sprint(pipeline_id) + "/body",
		"newJobUrl":     "/v1/" + fmt.Sprint(pipeline_id) + "/newjob",
		"pkgurl":        "/v1/" + fmt.Sprint(pipeline_id) + "/pkgs",
		"pipeline_name": pipeline.Pipeline_name,
		"pipeline_id":   pipeline_id,
		"jobs":          jobs,
		"page":          service.Cicd.PageContent(page),
		"envurl":        "/v1/" + fmt.Sprint(pipeline_id) + "/",

		// "url":           "/pipelines/",
		// "apiurl": "/v1/pipelines/" + fmt.Sprint(pipeline_id),
		// "pipeline_name": pipeline.Pipeline_name,
		// "pipeline_id": pipeline_id,
		// "agents": agents,
		// "groups": groups,
	})
	return nil, err
}
