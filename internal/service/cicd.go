package service

import (
	"cicdgf/internal/dao"
	"context"
	"fmt"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
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
