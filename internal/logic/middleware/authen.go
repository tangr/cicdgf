package middleware

import (
	"fmt"
	"net/url"
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

	userVar, _ := r.Session.Get("user")
	ticketVar, _ := r.Session.Get("ticket")

	if userVar == nil || ticketVar == nil {
		currentUrl := r.URL.String()
		loginUrl := fmt.Sprintf("/login?returnUrl=%s", url.QueryEscape(currentUrl))
		r.Response.RedirectTo(loginUrl)
		return
	}

	lastValidateTimeVar, _ := r.Session.Get("last_validate_time")
	if lastValidateTimeVar != nil {
		lastValidateTime := lastValidateTimeVar.Time()
		if time.Since(lastValidateTime) > 500*time.Minute {
			callbackURL := fmt.Sprintf("%s/callback", common.Cfg.ServiceURL)
			_, err := common.ValidateSSOSession(r.Context(), ticketVar.String(), callbackURL)
			if err != nil {
				r.Session.RemoveAll()
				currentUrl := r.URL.String()
				loginUrl := fmt.Sprintf("/login?returnUrl=%s", url.QueryEscape(currentUrl))
				r.Response.RedirectTo(loginUrl)
				return
			}
			r.Session.Set("last_validate_time", time.Now())
		}
	}

	r.Middleware.Next()
}
