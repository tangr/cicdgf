package log

import (
	"context"

	v1 "cicdgf/api/log/v1"
	"cicdgf/internal/dao"
	"cicdgf/internal/model/do"
)

func (c *ControllerV1) Create(ctx context.Context, req *v1.CreateReq) (res *v1.CreateRes, err error) {
	insertId, err := dao.CicdLog.Ctx(ctx).Data(do.CicdLog{
		PipelineId: req.PipelineId,
		AgentId:    req.AgentId,
		JobType:    req.JobType,
		JobId:      req.JobId,
		TaskStatus: req.TaskStatus,
		Ipaddr:     req.Ipaddr,
		UpdatedAt:  req.UpdatedAt,
		Output:     req.Output,
	}).InsertAndGetId()
	if err != nil {
		return nil, err
	}
	res = &v1.CreateRes{
		Id: insertId,
	}
	return
}
