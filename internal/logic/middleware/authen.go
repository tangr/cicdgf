package middleware

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"cicdgf/internal/logic/common"

	"github.com/gogf/gf/v2/net/ghttp"
)

func AuthenMiddleware(r *ghttp.Request) {
	ctx := r.Context()

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

	usernameVar, _ := r.Session.Get("user")
	useridVar, _ := r.Session.Get("userid")
	ticketVar, _ := r.Session.Get("ticket")

	ctx = context.WithValue(ctx, common.UsernameKey, usernameVar)
	ctx = context.WithValue(ctx, common.UseridKey, useridVar)
	r.SetCtx(ctx)

	if usernameVar == nil || ticketVar == nil {
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
