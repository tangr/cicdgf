package middleware

import (
	"fmt"
	"time"

	"cicdgf/internal/logic/common"

	"github.com/gogf/gf/v2/net/ghttp"
)

func AuthMiddleware(r *ghttp.Request) {
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

	user := r.Session.MustGet("user").String()
	ticket := r.Session.MustGet("ticket").String()

	if user == "" || ticket == "" {
		r.Response.RedirectTo("/login")
		return
	}

	lastValidateTime := r.Session.MustGet("last_validate_time").Time()
	if time.Since(lastValidateTime) > 500*time.Minute {
		callbackURL := fmt.Sprintf("%s/callback", common.Cfg.ServiceURL)
		_, err := common.ValidateSSOSession(r.Context(), ticket, callbackURL)
		if err != nil {
			r.Session.RemoveAll()
			r.Response.RedirectTo("/login")
			return
		}
		r.Session.Set("last_validate_time", time.Now())
	}

	r.Middleware.Next()
}
