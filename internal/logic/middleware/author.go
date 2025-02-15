package middleware

import (
	"cicdgf/internal/logic/common"
	"cicdgf/internal/service"
	"fmt"
	"strconv"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func AuthorMiddleware(r *ghttp.Request) {
	// 获取pipeline_id参数
	pipelineId := r.Get("pipeline_id").Int()

	// 从上下文获取用户ID
	useridVar := r.Context().Value(common.UseridKey)
	if useridVar == nil {
		r.Response.WriteJson(g.Map{
			"code": 401,
			"msg":  "Unauthorized: User not logged in",
		})
		return
	}

	// 转换用户ID为int类型
	userId := 0
	switch v := useridVar.(type) {
	case int:
		userId = v
	case string:
		userId, _ = strconv.Atoi(v)
	default:
		r.Response.WriteJson(g.Map{
			"code": 500,
			"msg":  "Invalid user ID format",
		})
		return
	}

	// 获取用户组ID列表
	groupIdsUser, err := service.User.GetListGroupId(userId)
	if err != nil {
		r.Response.WriteJson(g.Map{
			"code": 500,
			"msg":  "Failed to get user group IDs",
		})
		return
	}

	// 获取pipeline的组ID
	groupIdPipeline, err := service.Pipeline.GetOneGroupId(pipelineId)
	if err != nil {
		r.Response.WriteJson(g.Map{
			"code": 500,
			"msg":  "Failed to get pipeline group ID",
		})
		return
	}

	// 检查权限
	if stringInSlice(fmt.Sprint(groupIdPipeline), groupIdsUser) {
		r.Middleware.Next()
		return
	}

	// 没有权限，返回错误
	r.Response.WriteJson(g.Map{
		"code": 403,
		"msg":  "No permission to access this pipeline",
	})
}
