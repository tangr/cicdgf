package ui

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"cicdgf/api/ui/v1"
)

func (c *ControllerV1) Ui(ctx context.Context, req *v1.UiReq) (res *v1.UiRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
