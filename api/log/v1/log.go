package v1

import (
	"cicdgf/internal/model/entity"

	"github.com/gogf/gf/v2/frame/g"
)

type GetOneReq struct {
	g.Meta `path:"/log/{id}" method:"get" tags:"api" summary:"GetOne Log"`
	Id     uint64 `v:"required" json:"id" dc:"log id"`
}
type GetOneRes struct {
	*entity.CicdLog `dc:"log"`
}

type CreateReq struct {
	g.Meta     `path:"/log" method:"post" tags:"Log" summary:"Create Log"`
	PipelineId int    `v:"required-without:agentId" json:"pipelineId" dc:"Pipeline ID"`
	AgentId    int    `v:"required-without:pipelineId" dc:"Agent ID"`
	JobType    string `v:"required|in:BUILD,DEPLOY" json:"jobType" dc:"Job type"`
	JobId      int    `v:"required" json:"jobId" dc:"Job ID"`
	TaskStatus string `v:"required|in:pending,running,failed,success" json:"taskStatus" dc:"Task status"`
	Ipaddr     string `v:"required" json:"ipaddr" dc:"IP address"`
	UpdatedAt  int64  `v:"required" json:"updatedAt" dc:"Update timestamp"`
	Output     string `v:"required" json:"output" dc:"Output content"`
}
type CreateRes struct {
	Id int64 `json:"id" dc:"log id"`
}

type UpdateReq struct {
	g.Meta `path:"/log/{id}" method:"put" tags:"Log" summary:"Update Log"`

	Id         uint64 `v:"required" json:"logId" dc:"Log ID"`
	PipelineId int    `v:"required-without:agentId" json:"pipelineId" dc:"Pipeline ID"`
	AgentId    int    `v:"required-without:pipelineId" json:"agentId" dc:"Agent ID"`
	JobType    string `v:"required|in:BUILD,DEPLOY" json:"jobType" dc:"Job type"`
	JobId      int    `v:"required" json:"jobId" dc:"Job ID"`
	TaskStatus string `v:"required|in:pending,running,failed,success" json:"taskStatus" dc:"Task status"`
	Ipaddr     string `v:"required" json:"ipaddr" dc:"IP address"`
	UpdatedAt  int64  `v:"required" json:"updatedAt" dc:"Update timestamp"`
	Output     string `v:"required" json:"output" dc:"Output content"`
}
type UpdateRes struct{}
