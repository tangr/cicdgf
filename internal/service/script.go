package service

import (
	"cicdgf/internal/dao"
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

var Script = scriptService{}

type scriptService struct{}

type ListScripts struct {
	Id          int    `json:"script_id"`
	Script_name string `json:"script_name"`
	Author      string `json:"author"`
}

type ScriptDetail struct {
	Script_name string `json:"script_name"`
	Script_Body string `json:"script_body"`
}

func (s *scriptService) GetListScripts(ctx context.Context) (scripts []ListScripts, err error) {
	if err = dao.CicdScript.Ctx(ctx).
		Fields("id,script_name,author").
		Scan(&scripts); err != nil {
		return nil, gerror.Wrap(err, "get scripts failed")
	}

	return
}

func (s *scriptService) New(script_name string, script_body string) int64 {
	ctx := context.Background()

	new_script := g.Map{
		"script_name": script_name,
		"script_body": script_body,
	}

	result, err := dao.CicdScript.Ctx(ctx).Insert(new_script)
	if err != nil {
		g.Log().Error(ctx, err)
	}

	script_id, err := result.LastInsertId()
	if err != nil {
		g.Log().Error(ctx, err)
	}

	return script_id
}

func (s *scriptService) GetOne(script_id int) (*ScriptDetail, error) {
	ctx := context.Background()

	record, err := dao.CicdScript.Ctx(ctx).
		Fields("script_name,script_body").
		Where("id=?", script_id).
		One()
	if err != nil {
		g.Log().Error(ctx, err)
		return nil, err
	}

	return &ScriptDetail{
		Script_name: record["script_name"].String(),
		Script_Body: record["script_body"].String(),
	}, nil
}

func (s *scriptService) Update(script_id string, script_body string) string {
	ctx := context.Background()

	new_script := g.Map{
		"script_id":   script_id,
		"script_body": script_body,
	}
	g.Log().Debug(ctx, new_script)

	_, err := dao.CicdScript.
		Ctx(ctx).
		Where("id=?", script_id).
		Update(new_script)
	if err != nil {
		g.Log().Error(ctx, err)
	}

	return script_id
}
