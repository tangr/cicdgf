package ui

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func (c *ControllerV1) GroupGetList(ctx context.Context, req *GroupGetListReq) (response *ghttp.Response, err error) {
	r := g.RequestFromCtx(ctx)

	err = r.Response.WriteTpl("groups/list.html", g.Map{
		"url": "/groups/",
		// "groups":      groups,
		"newGroupUrl": "/groups/new",
	})
	return nil, err
}

func (c *ControllerV1) GroupGetOne(ctx context.Context, req *GroupGetOneReq) (response *ghttp.Response, err error) {
	r := g.RequestFromCtx(ctx)

	err = r.Response.WriteTpl("groups/list.html", g.Map{
		"url": "/groups/",
		// "groups":      groups,
		"newGroupUrl": "/groups/new",
	})
	return nil, err
}

func (c *ControllerV1) GroupCreate(ctx context.Context, req *GroupCreateReq) (response *ghttp.Response, err error) {
	r := g.RequestFromCtx(ctx)

	err = r.Response.WriteTpl("groups/list.html", g.Map{
		"url": "/groups/",
		// "groups":      groups,
		"newGroupUrl": "/groups/new",
	})
	return nil, err
}
