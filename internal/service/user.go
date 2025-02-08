package service

import (
	"cicdgf/internal/dao"
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

var User = userService{}

type userService struct{}

type ListUsers struct {
	Id         int    `json:"id"`
	Username   string `json:"username"`
	Groups     string `json:"groups"`
	Updated_at int    `json:"updated_at"`
}

type UserDetail struct {
	Username string `json:"username"`
	GroupId  string `json:"groups"`
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

	new_user := g.Map{
		"user_name": username,
	}

	result, err := dao.CicdUser.Ctx(ctx).Insert(new_user)
	if err != nil {
		g.Log().Error(ctx, err)
	}

	userid, err := result.LastInsertId()
	if err != nil {
		g.Log().Error(ctx, err)
	}

	return userid
}

func (s *userService) GetUserDetail(user_id string) (*UserDetail, error) {
	ctx := context.Background()

	record, err := dao.CicdUser.Ctx(ctx).
		Fields("username, group_id").
		Where("id=?", user_id).
		One()

	if err != nil {
		g.Log().Error(ctx, err)
		return nil, err
	}

	return &UserDetail{
		Username: record["username"].String(),
		GroupId:  record["group_id"].String(),
	}, nil
}

func (s *userService) Update(userid string, username string, groups []string) string {
	ctx := context.Background()

	newuser := g.Map{
		"group_id":   groups,
		"updated_at": gtime.Now().Timestamp(),
	}
	g.Log().Debug(ctx, "newuser:", newuser)
	g.Log().Debug(ctx, "userid:", userid)

	_, err := dao.CicdUser.
		Ctx(ctx).
		Where("id=?", userid).
		Update(newuser)
	if err != nil {
		g.Log().Error(ctx, err)
	}

	return userid
}
