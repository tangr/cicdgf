package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"

	"cicdgf/internal/controller/hello"
	"cicdgf/internal/controller/ui"
	"cicdgf/internal/controller/user"
	"cicdgf/internal/logic/common"
	"cicdgf/internal/logic/middleware"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()
			s.Group("/api", func(group *ghttp.RouterGroup) {
				group.Middleware(ghttp.MiddlewareHandlerResponse)
				group.Bind(
					hello.NewV1(),
					user.NewV1(),
				)
			})

			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(ghttp.MiddlewareHandlerResponse)
				group.Middleware(middleware.AuthMiddleware) // 添加认证中间件
				group.Bind(
					ui.NewV1(),
				)

				group.GET("/test", func(r *ghttp.Request) {
					r.Response.Write("Home Page")
				})
				group.GET("/login", func(r *ghttp.Request) {
					// 重定向到CAS服务器
					redirectURL := fmt.Sprintf("%s%s?service=%s/callback",
						common.Cfg.CasServerURL,
						common.Cfg.LoginURL,
						common.Cfg.ServiceURL,
					)
					r.Response.RedirectTo(redirectURL)
				})

				// callback 处理中相应修改
				group.GET("/callback", func(r *ghttp.Request) {
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

					// 设置 Session
					r.Session.Set("user", casResp.Success.User)
					r.Session.Set("ticket", ticket)
					r.Session.Set("last_validate_time", time.Now())

					r.Response.RedirectTo("/dashboard")
				})

				group.GET("/dashboard", func(r *ghttp.Request) {
					sessionData, err := r.Session.Data()
					if err != nil {
						g.Log().Error(context.Background(), "sessionData:", err)
					}
					g.Log().Debug(context.Background(), "All session data:", sessionData)

					user := r.Session.MustGet("user").String()
					r.Response.Writef("Welcome %s! <a href='/logout'>Logout</a>", user)
				})

				// 登出
				group.GET("/logout", func(r *ghttp.Request) {
					// 清除本地Session
					r.Session.RemoveAll()

					// 重定向到CAS全局登出
					logoutURL := fmt.Sprintf("%s%s?service=%s",
						common.Cfg.CasServerURL,
						common.Cfg.LogoutURL,
						common.Cfg.ServiceURL,
					)
					r.Response.RedirectTo(logoutURL)
				})

			})

			s.Run()
			return nil
		},
	}
)
