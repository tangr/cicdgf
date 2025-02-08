package service

import (
	"cicdgf/internal/dao"
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

var User = userService{}

type userService struct{}

type ListUsers struct {
	Id         int    `json:"id"`
	Username   string `json:"username"`
	Groups     string `json:"groups"`
	Updated_at int    `json:"updated_at"`
}

func (s *userService) GetListUsers(ctx context.Context) (users []ListUsers, err error) {
	if err = dao.CicdUser.Ctx(ctx).
		Fields("id,username,updated_at").
		Scan(&users); err != nil {
		return nil, gerror.Wrap(err, "get users failed")
	}

	return
}

func (s *userService) New(username string) int64 {
	ctx := context.Background()

	newuseer := g.Map{
		"user_name": username,
	}

	result, err := dao.CicdUser.Ctx(ctx).Insert(newuseer)
	if err != nil {
		g.Log().Error(ctx, err)
	}

	userid, err := result.LastInsertId()
	if err != nil {
		g.Log().Error(ctx, err)
	}

	return userid
}

func (s *userService) GetUserName(user_id string) string {
	ctx := context.Background()

	user_name, err := dao.CicdUser.Ctx(ctx).
		Fields("user_name").
		Where("id=?", user_id).
		Value()

	if err != nil {
		g.Log().Error(ctx, err)
		return ""
	}

	return user_name.String()
}
