// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package ui

import (
	"context"

	"cicdgf/api/ui/v1"
)

type IUiV1 interface {
	Ui(ctx context.Context, req *v1.UiReq) (res *v1.UiRes, err error)
}
