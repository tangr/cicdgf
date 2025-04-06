package log

import (
	"context"

	v1 "cicdgf/api/log/v1"
	"cicdgf/internal/dao"
	"cicdgf/internal/model/do"
)

func (c *ControllerV1) Update(ctx context.Context, req *v1.UpdateReq) (res *v1.UpdateRes, err error) {
	_, err = dao.CicdLog.Ctx(ctx).Data(do.CicdLog{
		PipelineId: req.PipelineId,
		AgentId:    req.AgentId,
		JobType:    req.JobType,
		JobId:      req.JobId,
		TaskStatus: req.TaskStatus,
		Ipaddr:     req.Ipaddr,
		UpdatedAt:  req.UpdatedAt,
		Output:     req.Output,
	}).WherePri(req.Id).Update()
	return
}
