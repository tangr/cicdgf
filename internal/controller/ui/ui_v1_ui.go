package ui

import (
	v1 "cicdgf/api/ui/v1"
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func (c *ControllerV1) Ui(ctx context.Context, req *v1.UiReq) (response *ghttp.Response, err error) {
	r := g.RequestFromCtx(ctx)

	err = r.Response.WriteTpl("ui.html", g.Map{
		"Title":          "111欢迎页面",
		"Content":        "aaaa:bbb:ccc",
		"url":            "/",
		"pipelines":      "pipelines",
		"newPipelineUrl": "/pipelines/new",
	})
	return nil, err
}
