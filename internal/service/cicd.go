package service

import (
	"cicdgf/internal/dao"
	"context"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gpage"
)

var Cicd = cicdService{}

type cicdService struct{}

type ListJobs struct {
	Id          int    `json:"job_id"`
	Pipeline_id int    `json:"pipeline_id"`
	Agent_id    int    `json:"agent_id"`
	Job_type    string `json:"job_type"`
	Job_status  string `json:"job_status"`
	Comment     string `json:"comment"`
	Author      string `json:"author"`
	Created_at  int    `json:"created_at"`
}

type ListTasks struct {
	Id          int    `json:"log_id"`
	Job_id      int    `json:"job_id"`
	Job_type    string `json:"job_type"`
	Agent_id    int    `json:"agent_id"`
	Pipeline_id int    `json:"pipeline_id"`
	Task_status string `json:"task_status"`
	Ipaddr      string `json:"ipaddr"`
	Actived     int    `json:"Actived"`
	Updated_at  int    `json:"updated_at"`
}

type JobDetail struct {
	Concurrency int    `json:"concurrency"`
	JobType     string `json:"job_type"`
	JobStatus   string `json:"job_status"`
}

type LogDetail struct {
	Task_status string `json:"status"`
	Updated_at  int    `json:"updated_at"`
	Output      string `json:"output"`
}

func (a *cicdService) GetListCicd(r *ghttp.Request) {
	// group_ids := service.GetUserGroupIds(r.Context())
	// pipelines := service.Cicd.ListCicd(group_ids)
	// params := g.Map{
	// 	"url":            UrlPrefix + "/",
	// 	"pipelines":      pipelines,
	// 	"newPipelineUrl": UrlPrefix + "/pipelines/new",
	// }
	// r.Response.WriteTpl("cicd/list.html", params)
}

func (s *cicdService) GetListAgents(ctx context.Context) (agents []ListAgents, err error) {
	if err = dao.CicdAgent.Ctx(ctx).
		Fields("id,agent_name,ipaddr,updated_at").
		Scan(&agents); err != nil {
		return nil, gerror.Wrap(err, "get agents failed")
	}

	return
}

func (s *cicdService) GetListJobs(ctx context.Context, pipeline_id int, pageIndex int, pageSize int) (jobs []ListJobs, totalSize int, err error) {
	offset := pageSize * (pageIndex - 1)

	err = dao.CicdJob.Ctx(ctx).
		Fields("id,pipeline_id,agent_id,job_type,job_status,comment,author,created_at").
		Where("pipeline_id", pipeline_id).
		Order("id DESC").
		Limit(pageSize).
		Offset(offset).
		Scan(&jobs)
	if err != nil {
		return nil, 0, gerror.Wrap(err, "get jobs failed")
	}

	totalSize = len(jobs)

	return jobs, totalSize, nil
}

func (s *cicdService) PageContent(page *gpage.Page) string {
	page.NextPageTag = `<i class="angle right icon"></i>`
	page.PrevPageTag = `<i class="angle left icon"></i>`
	pageStr := page.PrevPage()
	pageStr += fmt.Sprint(page.CurrentPage)
	pageStr += page.NextPage()
	return pageStr
}

func (s *cicdService) CreateJob(ctx context.Context, pipeline_id int, envs map[string]interface{}, username string) (int64, error) {
	var script_args string
	var jobtype string
	var job_envs map[string]string = Comm.ParseEnvs(envs)
	var comment string = job_envs["COMMENT"]
	var job_type string = job_envs["JOBTYPE"]

	g.Log().Debug(ctx, "CreateJob pipeline_id: ", pipeline_id)
	pipeline, err := Pipeline.GetOne(pipeline_id)
	g.Log().Debug(ctx, "CreateJob pipeline: ", pipeline)
	g.Log().Debug(ctx, "CreateJob err: ", err)
	if err != nil {
		return 0, nil
	}

	pipeline_body := pipeline.Body
	pipeline_name := pipeline.Pipeline_name
	agent_id := pipeline.Agent_id
	concurrency := pipeline.Concurrency

	// pipeline_name, agent_id, concurrency, pipeline_body := Pipeline.GetOne(pipeline_id)
	if job_type == "BUILD" {
		jobtype = job_type
		// script_name := pipeline_body.StageCI.Script
		script_args = pipeline_body.StageCI.Args
		job_envs["PKGRDM"] = Comm.RandSeq(20)
	} else if job_type == "DEPLOY" {
		type JobStatus struct {
			Id        int64  `json:"job_id"`
			JobStatus string `json:"job_status"`
		}
		var last_job_status JobStatus
		last_job := g.Map{"pipeline_id": pipeline_id, "job_type": "DEPLOY"}
		err := dao.CicdJob.
			Ctx(ctx).
			Fields("id,job_status").
			Where(last_job).
			OrderDesc("id").
			Limit(1).
			Scan(&last_job_status)
		if err != nil {
			g.Log().Error(ctx, err)
		}
		if last_job_status.JobStatus != "success" && last_job_status.JobStatus != "failed" && last_job_status.Id != 0 {
			return last_job_status.Id, nil
		}
		jobtype = job_type
		// script_name := pipeline_body.StageCD.Script
		script_args = pipeline_body.StageCD.Args
	} else {
		g.Log().Errorf(ctx, "unsupported job_type: %s", job_type)
	}
	job_envs["PIPELINEID"] = fmt.Sprint(pipeline_id)
	job_envs["PIPELINENAME"] = strings.Split(pipeline_name, ":")[0]
	job_envs["USERNAME"] = username
	script_body := "Script.GetScriptBody(script_name)"
	new_jobscript := new(JobScriptValue)
	new_jobscript.Envs = job_envs
	new_jobscript.Args = script_args
	new_jobscript.Body = script_body

	new_job := g.Map{
		"pipeline_id": pipeline_id,
		"agent_id":    agent_id,
		"concurrency": concurrency,
		"job_type":    jobtype,
		"job_status":  "pending",
		"script":      new_jobscript,
		"comment":     comment,
		"author":      username,
		"created_at":  gtime.Now().Timestamp(),
	}
	result, err := dao.CicdJob.Ctx(ctx).
		Data(new_job).
		Save()
	if err != nil {
		g.Log().Error(ctx, err)
	}
	job_id, err := result.LastInsertId()
	if err != nil {
		g.Log().Error(ctx, err)
	}
	return job_id, nil
}

func (s *cicdService) CheckJobid(ctx context.Context, pipeline_id int, job_id int) bool {
	num, err := dao.CicdJob.Ctx(ctx).
		Where(g.Map{"pipeline_id": pipeline_id, "id": job_id}).
		Count()
	if err != nil {
		g.Log().Error(ctx, err)
		return false
	}
	if num < 1 {
		return false
	}
	return true
}

func (s *cicdService) GetListJobTasks(ctx context.Context, pipeline_id int, job_id int) []ListTasks {
	var agentStatusMap map[string]int
	tasks := ([]ListTasks)(nil)
	if !s.CheckJobid(ctx, pipeline_id, job_id) {
		return tasks
	}
	err := dao.CicdLog.Ctx(ctx).
		Fields("id,job_id,job_type,agent_id,pipeline_id,task_status,ipaddr,updated_at").
		Order("id desc").
		Where(g.Map{"job_id": job_id}).
		Scan(&tasks)
	if err != nil {
		g.Log().Error(ctx, err)
	}
	// // var agentStatus string
	// status_url := fmt.Sprint(WsServerAPI, pipeline_id, "/", job_id, "/status")
	// r, err := g.Client().Get(status_url)
	// if err != nil {
	// 	g.Log().Error(err)
	// } else {
	// 	defer r.Close()
	// }
	// agentStatus := r.ReadAllString()
	// json.Unmarshal([]byte(agentStatus), &agentStatusMap)
	for idx, v := range tasks {
		if v.Job_type == "BUILD" {
			mapk := fmt.Sprint("CI-", v.Agent_id, "-", v.Ipaddr)
			tasks[idx].Actived = agentStatusMap[mapk]
		} else {
			mapk := fmt.Sprint("CD-", v.Pipeline_id, "-", v.Ipaddr)
			tasks[idx].Actived = agentStatusMap[mapk]
		}
	}
	return tasks
}

func (s *cicdService) GetOneJob(job_id int) (*JobDetail, error) {
	ctx := context.Background()

	record, err := dao.CicdJob.Ctx(ctx).
		Fields("concurrency,job_type,job_status").
		Where("id=?", job_id).
		One()
	if err != nil {
		g.Log().Error(ctx, err)
		return nil, err
	}

	return &JobDetail{
		Concurrency: record["concurrency"].Int(),
		JobType:     record["job_type"].String(),
		JobStatus:   record["job_status"].String(),
	}, nil
}

func (s *cicdService) GetOneLog(ctx context.Context, pipeline_id int, log_id int) (*LogDetail, error) {
	output := (*LogDetail)(nil)
	// if !s.CheckTaskid(pipeline_id, log_id) {
	// 	return output
	// }
	err := dao.CicdLog.Ctx(ctx).
		Fields("task_status,updated_at,output").
		Where(g.Map{"id": log_id}).
		Scan(&output)
	if err != nil {
		g.Log().Error(ctx, err)
	}
	return output, nil
}
