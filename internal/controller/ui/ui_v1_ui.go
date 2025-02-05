package ui

import (
	v1 "cicdgf/api/ui/v1"
	"context"
)

func (c *ControllerV1) Ui(ctx context.Context, req *v1.UiReq) (res *v1.UiRes, err error) {
	return &v1.UiRes{
		Content: `<html>Your UI Content</html>`,
	}, nil
}
