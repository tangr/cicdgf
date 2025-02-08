package service

import (
	"cicdgf/internal/dao"
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

var Group = groupService{}

type groupService struct{}

type ListGroups struct {
	Id         int    `json:"id"`
	Group_name string `json:"groupName"`
}

func (s *groupService) GetListGroups(ctx context.Context) (groups []ListGroups, err error) {
	if err = dao.CicdGroup.Ctx(ctx).
		Fields("id,group_name").
		Scan(&groups); err != nil {
		g.Log().Error(ctx, "Failed to get groups:", err)
		return nil, gerror.Wrap(err, "get groups failed")
	}

	return
}

func (s *groupService) New(groupname string) int64 {
	ctx := context.Background()

	new_group := g.Map{
		"group_name": groupname,
		"parent_id":  0,
	}

	result, err := dao.CicdGroup.Ctx(ctx).Insert(new_group)
	if err != nil {
		g.Log().Error(ctx, err)
	}

	groupid, err := result.LastInsertId()
	if err != nil {
		g.Log().Error(ctx, err)
	}

	return groupid
}

func (s *groupService) Update(groupid string, groupname string) string {
	ctx := context.Background()

	newgroup := g.Map{
		"group_name": groupname,
	}
	g.Log().Debug(ctx, newgroup)

	_, err := dao.CicdGroup.
		Ctx(ctx).
		Where("id=?", groupid).
		Update(newgroup)
	if err != nil {
		g.Log().Error(ctx, err)
	}

	return groupid
}

func (s *groupService) GetGroupName(group_id string) string {
	ctx := context.Background()

	group_name, err := dao.CicdGroup.Ctx(ctx).
		Fields("group_name").
		Where("id=?", group_id).
		Value()

	if err != nil {
		g.Log().Error(ctx, err)
		return ""
	}

	return group_name.String()
}
