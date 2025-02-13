package cmd

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"

	"cicdgf/internal/controller/hello"
	"cicdgf/internal/controller/ui"
	"cicdgf/internal/controller/user"
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
				group.Middleware(middleware.AuthMiddleware)
				group.Bind(
					ui.NewV1(),
				)

				group.GET("/test", func(r *ghttp.Request) {
					r.Response.Write("Home Page")
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

			})

			s.Run()
			return nil
		},
	}
)
