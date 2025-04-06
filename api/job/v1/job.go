package v1

import (
	"cicdgf/internal/model/entity"

	"github.com/gogf/gf/v2/frame/g"
)

type GetOneReq struct {
	g.Meta `path:"/job/{job_id}" method:"get" tags:"api" summary:"Get one job"`
	Id     int64 `v:"required" dc:"user id"`
}

type GetOneRes struct {
	*entity.CicdJob `dc:"job"`
}
