package log

import (
	"context"

	v1 "cicdgf/api/log/v1"
	"cicdgf/internal/dao"
)

func (c *ControllerV1) GetOne(ctx context.Context, req *v1.GetOneReq) (res *v1.GetOneRes, err error) {
	res = &v1.GetOneRes{}
	err = dao.CicdLog.Ctx(ctx).WherePri(req.Id).Scan(&res.CicdLog)
	return
}
