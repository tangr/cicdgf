package middleware

import (
	"fmt"
	"time"

	"cicdgf/internal/logic/common"

	"github.com/gogf/gf/v2/net/ghttp"
)

func AuthMiddleware(r *ghttp.Request) {
	// 不需要验证的路径直接放行
	skipPaths := map[string]bool{
		"/login":    true,
		"/callback": true,
		"/logout":   true,
		"/test":     true,
	}

	if skipPaths[r.URL.Path] {
		r.Middleware.Next()
		return
	}

	// 检查本地 session
	user := r.Session.MustGet("user").String()
	ticket := r.Session.MustGet("ticket").String()

	if user == "" || ticket == "" {
		r.Response.RedirectTo("/login")
		return
	}

	// 定期验证 SSO session
	lastValidateTime := r.Session.MustGet("last_validate_time").Time()
	if time.Since(lastValidateTime) > 5*time.Minute {
		callbackURL := fmt.Sprintf("%s/callback", common.Cfg.Service)
		_, err := common.ValidateSSOSession(r.Context(), ticket, callbackURL)
		if err != nil {
			// SSO session 已失效，清除本地 session
			r.Session.RemoveAll()
			r.Response.RedirectTo("/login")
			return
		}
		// 更新最后验证时间
		r.Session.Set("last_validate_time", time.Now())
	}

	r.Middleware.Next()
}
