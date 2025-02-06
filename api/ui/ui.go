// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package ui

import (
	"context"

	"cicdgf/api/ui/v1"
	"github.com/gogf/gf/v2/net/ghttp"
)

type IUiV1 interface {
	Ui(ctx context.Context, req *v1.UiReq) (res *ghttp.Response, err error)
}
