package ui

import (
	"cicdgf/internal/service"
	"context"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func (c *ControllerV1) ScriptGetList(ctx context.Context, req *ScriptGetListReq) (response *ghttp.Response, err error) {
	r := g.RequestFromCtx(ctx)

	scripts, err := service.Script.GetListScripts(ctx)
	if err != nil {
		return nil, err
	}

	err = r.Response.WriteTpl("scripts/list.html", g.Map{
		"url":          "/scripts/",
		"scripts":      scripts,
		"newScriptUrl": "/scriptnew",
	})
	return nil, err
}

func (c *ControllerV1) ScriptNew(ctx context.Context, req *ScriptNewReq) (response *ghttp.Response, err error) {
	r := g.RequestFromCtx(ctx)

	err = r.Response.WriteTpl("scripts/new.html", g.Map{
		"url":          "/scripts/",
		"newScriptUrl": "/scripts/",
	})
	return nil, err
}

func (c *ControllerV1) ScriptCreate(ctx context.Context, req *ScriptCreateReq) (response *ghttp.Response, err error) {
	r := g.RequestFromCtx(ctx)

	var script_name string = r.Get("script_name").String()
	var script_body string = r.Get("pipeline_body").String()
	script_body = strings.Replace(script_body, "\r\n", "\n", -1)

	script_id := service.Script.New(script_name, script_body)

	r.Response.RedirectTo("/scripts/"+fmt.Sprint(script_id), 303)

	return nil, err
}

func (c *ControllerV1) ScriptGetOne(ctx context.Context, req *ScriptGetOneReq) (response *ghttp.Response, err error) {
	r := g.RequestFromCtx(ctx)

	script_id := r.Get("id").Int()

	script_obj, err := service.Script.GetOne(script_id)
	if err != nil {
		return nil, err
	}

	err = r.Response.WriteTpl("scripts/show.html", g.Map{
		"url":         "/scripts/",
		"apiurl":      "/scripts/" + fmt.Sprint(script_id) + "/put",
		"script_id":   script_id,
		"script_name": script_obj.Script_name,
		"script_body": script_obj.Script_Body,
	})
	return nil, err
}

func (c *ControllerV1) ScriptUpdate(ctx context.Context, req *ScriptUpdateReq) (response *ghttp.Response, err error) {
	r := g.RequestFromCtx(ctx)

	var script_id string = r.Get("id").String()
	var script_body string = r.Get("script_body").String()

	g.Log().Debugf(ctx, "ScriptUpdate script_id: %s, script_body: %s",
		script_id, script_body)

	_ = service.Script.Update(script_id, script_body)
	r.Response.RedirectTo("/scripts/"+fmt.Sprint(script_id), 303)

	return nil, err
}
