package cmd

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"

	"cicdgf/internal/controller/hello"
	"cicdgf/internal/controller/ui"
	"cicdgf/internal/controller/user"
)

type CasConfig struct {
	ServerURL   string // CAS服务器地址 例如: https://cas.example.com/cas
	Service     string // 本服务地址 例如: http://localhost:8199
	ValidateURL string // 票据验证地址
	LoginURL    string // 登录地址
	LogoutURL   string // 登出地址
}

var (
	cfg = CasConfig{
		ServerURL:   "https://sso-prod.yax.tech/cas",
		Service:     "http://localhost:8000",
		ValidateURL: "/serviceValidate",
		LoginURL:    "/login",
		LogoutURL:   "/logout",
	}

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
				group.Bind(
					ui.NewV1(),
				)

				group.GET("/test", func(r *ghttp.Request) {
					r.Response.Write("Home Page")
				})
				group.GET("/login", func(r *ghttp.Request) {
					// 重定向到CAS服务器
					redirectURL := fmt.Sprintf("%s%s?service=%s/callback",
						cfg.ServerURL,
						cfg.LoginURL,
						cfg.Service,
					)
					r.Response.RedirectTo(redirectURL)
				})

			})
			s.Run()
			return nil
		},
	}
)
