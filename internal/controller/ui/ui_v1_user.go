package ui

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func (c *ControllerV1) User(ctx context.Context, req *UserGetListReq) (response *ghttp.Response, err error) {
	r := g.RequestFromCtx(ctx)

	err = r.Response.WriteTpl("users/list.html", g.Map{
		"Title":          "111欢迎页面",
		"Content":        "aaaa:bbb:ccc",
		"url":            "/",
		"pipelines":      "pipelines",
		"newPipelineUrl": "/pipelines/new",
	})
	return nil, err
}
