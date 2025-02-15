package ui

import (
	"cicdgf/internal/logic/common"
	"cicdgf/internal/service"
	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func (c *ControllerV1) AuthLogin(ctx context.Context, req *AuthLoginReq) (response *ghttp.Response, err error) {
	r := g.RequestFromCtx(ctx)

	returnUrl := r.GetQuery("returnUrl").String()
	if returnUrl == "" {
		returnUrl = r.GetReferer()
		if returnUrl == "" {
			returnUrl = "/dashboard"
		}
	}

	// 将 returnUrl 保存到 session
	if err := r.Session.Set("returnUrl", returnUrl); err != nil {
		g.Log().Error(ctx, "Set session returnUrl error:", err)
	}

	redirectURL := fmt.Sprintf("%s%s?service=%s/callback",
		common.Cfg.CasServerURL,
		common.Cfg.LoginURL,
		common.Cfg.ServiceURL,
	)
	r.Response.RedirectTo(redirectURL, 303)

	return nil, err
}

func (c *ControllerV1) AuthCallback(ctx context.Context, req *AuthCallbackReq) (response *ghttp.Response, err error) {
	r := g.RequestFromCtx(ctx)

	ticket := r.GetQuery("ticket").String()
	if ticket == "" {
		r.Response.Write("Invalid CAS ticket")
		return
	}

	callbackURL := fmt.Sprintf("%s/callback", common.Cfg.ServiceURL)
	casResp, err := common.ValidateSSOSession(r.Context(), ticket, callbackURL)
	if err != nil {
		r.Response.Write("CAS validation failed")
		return
	}

	adminUsers := common.AdminCfg.AdminUsers
	username := casResp.Success.User
	arrayUsers := garray.NewStrArrayFrom(adminUsers)
	if arrayUsers.Contains(username) {
		r.Session.Set("IsAdmin", true)
	} else {
		r.Session.Set("IsAdmin", false)
	}

	userid, err := service.User.GetOneUsername(username)
	if err != nil {
		r.Response.Write("CAS validation failed")
		return
	}

	r.Session.Set("user", username)
	r.Session.Set("userid", userid)
	r.Session.Set("ticket", ticket)
	r.Session.Set("last_validate_time", time.Now())

	// 正确处理 Session.Get() 返回值
	returnUrlVar, err := r.Session.Get("returnUrl")
	var returnUrl string
	if err != nil {
		returnUrl = "/dashboard"
	} else {
		returnUrl = returnUrlVar.String()
	}

	// 使用后删除
	r.Session.Remove("returnUrl")

	if returnUrl != "" {
		r.Response.RedirectTo(returnUrl)
	} else {
		r.Response.RedirectTo("/dashboard")
	}

	return nil, err
}

func (c *ControllerV1) AuthLogout(ctx context.Context, req *AuthLogoutReq) (response *ghttp.Response, err error) {
	r := g.RequestFromCtx(ctx)

	r.Session.RemoveAll()

	logoutURL := fmt.Sprintf("%s%s?service=%s",
		common.Cfg.CasServerURL,
		common.Cfg.LogoutURL,
		common.Cfg.ServiceURL,
	)
	r.Response.RedirectTo(logoutURL)

	return nil, err
}
