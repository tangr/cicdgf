package ui

import (
	"cicdgf/internal/service"
	"context"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func (c *ControllerV1) UserGetList(ctx context.Context, req *UserGetListReq) (response *ghttp.Response, err error) {
	r := g.RequestFromCtx(ctx)

	users, err := service.User.GetListUsers(ctx)
	if err != nil {
		return nil, err
	}

	err = r.Response.WriteTpl("users/list.html", g.Map{
		"url":         "/users/",
		"users":       users,
		"newUsersUrl": "/usernew",
	})
	return nil, err
}

func (c *ControllerV1) UserNew(ctx context.Context, req *UserNewReq) (response *ghttp.Response, err error) {
	r := g.RequestFromCtx(ctx)
	groups, err := service.Group.GetListGroups(ctx)
	if err != nil {
		return nil, err
	}

	err = r.Response.WriteTpl("users/new.html", g.Map{
		"url":        "/users/",
		"groups":     groups,
		"newUserUrl": "/users",
	})
	return nil, err
}

func (c *ControllerV1) UserCreate(ctx context.Context, req *UserCreateReq) (response *ghttp.Response, err error) {
	r := g.RequestFromCtx(ctx)

	var username string = r.Get("username").String()
	userid := service.User.New(username)
	r.Response.RedirectTo("/users/" + fmt.Sprint(userid))

	return nil, err
}

func (c *ControllerV1) UserGetOne(ctx context.Context, req *UserGetOneReq) (response *ghttp.Response, err error) {
	r := g.RequestFromCtx(ctx)

	userid := r.Get("id").String()
	username := service.User.GetUserName(userid)
	groups, err := service.Group.GetListGroups(ctx)
	if err != nil {
		return nil, err
	}

	g.Log().Debugf(ctx, "userid: %s, username: %s", userid, username)

	err = r.Response.WriteTpl("users/show.html", g.Map{
		"url":      "/users/",
		"apiurl":   "/users/" + userid + "/put",
		"username": username,
		"groups":   groups,
	})
	return nil, err
}
