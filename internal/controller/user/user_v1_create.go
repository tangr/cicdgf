package user

import (
	"context"

	v1 "cicdgf/api/user/v1"
	"cicdgf/internal/dao"
	"cicdgf/internal/model/do"
)

func (c *ControllerV1) Create(ctx context.Context, req *v1.CreateReq) (res *v1.CreateRes, err error) {
	insertId, err := dao.User.Ctx(ctx).Data(do.User{
		Name:   req.Name,
		Status: v1.StatusOK,
		Age:    req.Age,
	}).InsertAndGetId()
	if err != nil {
		return nil, err
	}
	res = &v1.CreateRes{
		Id: insertId,
	}
	return
}
