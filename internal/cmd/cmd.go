package cmd

import (
	"context"
	"encoding/xml"
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

type CasResponse struct {
	XMLName xml.Name `xml:"serviceResponse"`
	Success struct {
		User       string `xml:"user"`
		Attributes struct {
			CredentialType                         string `xml:"credentialType"`
			ClientIpAddress                        string `xml:"clientIpAddress"`
			IsFromNewLogin                         bool   `xml:"isFromNewLogin"`
			AuthenticationDate                     string `xml:"authenticationDate"`
			AuthenticationMethod                   string `xml:"authenticationMethod"`
			SuccessfulAuthenticationHandlers       string `xml:"successfulAuthenticationHandlers"`
			ServerIpAddress                        string `xml:"serverIpAddress"`
			UserAgent                              string `xml:"userAgent"`
			LongTermAuthenticationRequestTokenUsed bool   `xml:"longTermAuthenticationRequestTokenUsed"`
		} `xml:"attributes"`
	} `xml:"authenticationSuccess"`
	Failure string `xml:"authenticationFailure"`
}

var (
	cfg = CasConfig{
		ServerURL:   "https://sso-prod.yax.tech/cas",
		Service:     "http://localhost:8000",
		ValidateURL: "/serviceValidate",
		// ValidateURL: "/validate",
		LoginURL:  "/login",
		LogoutURL: "/logout",
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

				// CAS回调处理
				group.GET("/callback", func(r *ghttp.Request) {
					ticket := r.GetQuery("ticket").String()
					if ticket == "" {
						r.Response.Write("Invalid CAS ticket")
						return
					}

					// 验证CAS票据
					validateURL := fmt.Sprintf("%s%s?ticket=%s&service=%s/callback",
						cfg.ServerURL,
						cfg.ValidateURL,
						ticket,
						cfg.Service,
					)

					g.Log().Debug(context.Background(), "validateURL:", validateURL)

					res, err := g.Client().Get(context.Background(), validateURL)
					if err != nil {
						g.Log().Error(context.Background(), "CAS validation error:", err)
						r.Response.Write("CAS validation failed")
						return
					}
					defer res.Close()
					g.Log().Debug(context.Background(), "CAS validation response:", res)

					// 解析CAS响应
					var casResp CasResponse
					g.Log().Debug(context.Background(), "CAS response status:", res.StatusCode)
					body := res.ReadAll()
					g.Log().Debug(context.Background(), "CAS response body:", string(body))
					if err := xml.Unmarshal(body, &casResp); err != nil {
						g.Log().Error(context.Background(), "XML unmarshal error:", err)
						r.Response.Write("Invalid CAS response")
						return
					}

					g.Log().Debug(context.Background(), "Parsed CAS response:", casResp)

					if casResp.Failure != "" {
						r.Response.Write("CAS authentication failed1: " + casResp.Failure)
						return
					}

					if casResp.Success.User == "" {
						g.Log().Error(context.Background(), "Empty user in CAS response")
					}

					// 登录成功，设置Session
					sessionUser := casResp.Success.User
					if sessionUser != "" {
						_ = r.Session.Set("user", sessionUser)
						r.Response.RedirectTo("/dashboard")
						return
					}

					r.Response.Write("CAS authentication failed2")
				})

				// 仪表盘（需要登录）
				group.GET("/dashboard", func(r *ghttp.Request) {
					user := r.Session.MustGet("user").String()
					if user == "" {
						r.Response.RedirectTo("/login")
						return
					}
					r.Response.Writef("Welcome %s! <a href='/logout'>Logout</a>", user)
				})

				// 登出
				group.GET("/logout", func(r *ghttp.Request) {
					// 清除本地Session
					_ = r.Session.RemoveAll()

					// 重定向到CAS全局登出
					logoutURL := fmt.Sprintf("%s%s?service=%s",
						cfg.ServerURL,
						cfg.LogoutURL,
						cfg.Service,
					)
					r.Response.RedirectTo(logoutURL)
				})

			})

			s.Run()
			return nil
		},
	}
)
