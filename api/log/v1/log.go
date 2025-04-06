package v1

import (
	"cicdgf/internal/model/entity"

	"github.com/gogf/gf/v2/frame/g"
)

type GetOneReq struct {
	g.Meta `path:"/log/{log_id}" method:"get" tags:"api" summary:"Get one log"`
	Id     int64 `v:"required" dc:"log id"`
}

type GetOneRes struct {
	*entity.CicdLog `dc:"log"`
}
