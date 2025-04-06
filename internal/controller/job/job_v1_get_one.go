package job

import (
	"context"

	v1 "cicdgf/api/job/v1"
	"cicdgf/internal/dao"
)

func (c *ControllerV1) GetOne(ctx context.Context, req *v1.GetOneReq) (res *v1.GetOneRes, err error) {
	res = &v1.GetOneRes{}
	err = dao.CicdJob.Ctx(ctx).WherePri(req.Id).Scan(&res.CicdJob)
	return
}
