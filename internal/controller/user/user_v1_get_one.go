package user

import (
	"context"

	v1 "cicdgf/api/user/v1"
	"cicdgf/internal/dao"
)

func (c *ControllerV1) GetOne(ctx context.Context, req *v1.GetOneReq) (res *v1.GetOneRes, err error) {
	res = &v1.GetOneRes{}
	err = dao.User.Ctx(ctx).WherePri(req.Id).Scan(&res.User)
	return
}
