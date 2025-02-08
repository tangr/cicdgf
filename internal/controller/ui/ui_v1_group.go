package ui

import (
	"cicdgf/internal/service"
	"context"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func (c *ControllerV1) GroupGetList(ctx context.Context, req *GroupGetListReq) (response *ghttp.Response, err error) {
	r := g.RequestFromCtx(ctx)

	groups, err := service.Group.GetListGroups(ctx)
	if err != nil {
		return nil, err
	}

	err = r.Response.WriteTpl("groups/list.html", g.Map{
		"url":         "/groups/",
		"groups":      groups,
		"newGroupUrl": "/groupnew",
	})
	return nil, err
}

func (c *ControllerV1) GroupGetOne(ctx context.Context, req *GroupGetOneReq) (response *ghttp.Response, err error) {
	r := g.RequestFromCtx(ctx)

	group_id := r.Get("id").String()
	group_name := service.Group.GetGroupName(group_id)

	g.Log().Debugf(ctx, "group_id: %s, group_name: %s", group_id, group_name)

	err = r.Response.WriteTpl("groups/show.html", g.Map{
		"url":        "/groups/",
		"apiurl":     "/groups/" + group_id + "/put",
		"group_name": group_name,
	})
	return nil, err
}

func (c *ControllerV1) GroupCreate(ctx context.Context, req *GroupCreateReq) (response *ghttp.Response, err error) {
	r := g.RequestFromCtx(ctx)

	var groupname string = r.Get("groupname").String()
	groupid := service.Group.New(groupname)
	r.Response.RedirectTo("/groups/" + fmt.Sprint(groupid))

	return nil, err
}

func (c *ControllerV1) GroupUpdate(ctx context.Context, req *GroupUpdateReq) (response *ghttp.Response, err error) {
	r := g.RequestFromCtx(ctx)

	var groupid string = r.Get("id").String()
	var groupname string = r.Get("groupname").String()
	_ = service.Group.Update(groupid, groupname)
	r.Response.RedirectTo("/groups/"+fmt.Sprint(groupid), 303)

	return nil, err
}

func (c *ControllerV1) GroupNew(ctx context.Context, req *GroupNewReq) (response *ghttp.Response, err error) {
	r := g.RequestFromCtx(ctx)

	err = r.Response.WriteTpl("groups/new.html", g.Map{
		"url": "/groups/",
		// "groups":      groups,
		"newGroupUrl": "/groups",
	})
	return nil, err
}
