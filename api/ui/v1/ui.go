package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

type UiReq struct {
	g.Meta `path:"/ui" tags:"ui" method:"get" summary:"ui console index"`
}
type UiRes struct {
	g.Meta  `mime:"text/html" example:"string"`
	Content string
}
