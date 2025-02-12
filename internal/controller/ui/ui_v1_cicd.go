package ui

import (
	"cicdgf/internal/service"
	"context"
	"fmt"
	"strconv"

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
		"body_url":      "/" + fmt.Sprint(pipeline_id) + "/body",
		"newJobUrl":     "/" + fmt.Sprint(pipeline_id) + "/newjob",
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

func (c *ControllerV1) CicdBodyGetOne(ctx context.Context, req *CicdBodyGetOneReq) (response *ghttp.Response, err error) {
	r := g.RequestFromCtx(ctx)

	var pipeline_id int = r.Get("id").Int()
	// if !service.CheckAuthor(r.Context(), pipeline_id) {
	// 	r.Response.RedirectTo(UrlPrefix + "/forbidden")
	// }

	pipeline_body, err := service.Pipeline.GetOnebody(pipeline_id)
	if err != nil {
		return nil, err
	}

	r.Response.WriteExit(pipeline_body)
	return nil, nil

}

func (c *ControllerV1) CicdJobCreate(ctx context.Context, req *CicdJobCreateReq) (response *ghttp.Response, err error) {
	r := g.RequestFromCtx(ctx)

	var pipeline_id int = r.Get("pipeline_id").Int()
	// if !service.CheckAuthor(r.Context(), pipeline_id) {
	// 	r.Response.RedirectTo(UrlPrefix + "/forbidden")
	// }
	// var username string = service.Session.GetUser(r.Context()).Email
	var username string = "tangshoubin"
	envs := r.GetFormMap()
	job_id, err := service.Cicd.CreateJob(ctx, pipeline_id, envs, username)
	g.Log().Debug(ctx, "CicdJobCreate job_id: ", job_id)
	if err != nil {
		g.Log().Debug(ctx, "CicdJobCreate err: ", err)
		return nil, err
	}

	r.Response.RedirectTo("/" + fmt.Sprint(pipeline_id) + "/" + strconv.FormatInt(job_id, 10))

	// var agent_name string = r.Get("agent_name").String()
	// var agent_ipaddr string = r.Get("agent_ipaddr").String()
	// agent_id := service.Agent.New(agent_name, agent_ipaddr)

	// r.Response.RedirectTo("/agents/"+fmt.Sprint(agent_id), 303)

	return nil, err
}

func (c *ControllerV1) CicdJobGetOne(ctx context.Context, req *CicdJobGetOneReq) (response *ghttp.Response, err error) {
	r := g.RequestFromCtx(ctx)

	// var pipeline_id int = r.Get("pipeline_id").Int()
	// var job_id int = r.Get("job_id").Int()
	// // if !service.CheckAuthor(r.Context(), pipeline_id) {
	// // 	r.Response.RedirectTo(UrlPrefix + "/forbidden")
	// // }

	var pipeline_id int = r.Get("pipeline_id").Int()
	// if !service.CheckAuthor(r.Context(), pipeline_id) {
	// 	r.Response.RedirectTo("/forbidden")
	// }
	var job_id int = r.Get("job_id").Int()

	tasks := service.Cicd.GetListJobTasks(ctx, pipeline_id, job_id)
	pipeline, err := service.Pipeline.GetOne(pipeline_id)
	if err != nil {
		g.Log().Debug(ctx, "CicdJobCreate err: ", err)
		return nil, err
	}

	pipeline_name := pipeline.Pipeline_name
	job, err := service.Cicd.GetOneJob(job_id)
	if err != nil {
		g.Log().Debug(ctx, "CicdJobCreate err: ", err)
		return nil, err
	}
	concurrency, job_type, job_status := job.Concurrency, job.JobType, job.JobStatus
	params := g.Map{
		"url":           "/" + fmt.Sprint(pipeline_id) + "/",
		"apiurl":        "/v1/" + fmt.Sprint(pipeline_id, "/", job_id),
		"pipeline_name": pipeline_name,
		"pipeline_id":   pipeline_id,
		"job_id":        job_id,
		"concurrency":   concurrency,
		"job_type":      job_type,
		"job_status":    job_status,
		"tasks":         tasks,
		"taskurl":       "/" + fmt.Sprint(pipeline_id) + "/",
	}
	if job_type == "BUILD" {
		// r.Response.WriteExit(params)
		r.Response.WriteTpl("cicd/job_build.html", params)
	} else {
		params["progressurl"] = "/" + fmt.Sprint(pipeline_id, "/", job_id) + "/progress"
		r.Response.WriteTpl("cicd/job_deploy.html", params)
	}

	pipeline_body, err := service.Pipeline.GetOnebody(pipeline_id)
	if err != nil {
		return nil, err
	}

	r.Response.WriteExit(pipeline_body)
	return nil, nil

}
