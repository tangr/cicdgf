package cmd

import (
	"context"
	"encoding/xml"
	"fmt"
	"time"

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

func validateSSOSession(ctx context.Context, ticket string, service string) (*CasResponse, error) {
	validateURL := fmt.Sprintf("%s%s?ticket=%s&service=%s",
		cfg.ServerURL,
		cfg.ValidateURL,
		ticket,
		service,
	)

	res, err := g.Client().Get(ctx, validateURL)
	if err != nil {
		g.Log().Error(ctx, "Validate SSO session error:", err)
		return nil, err
	}
	defer res.Close()

	// 解析CAS响应
	var casResp CasResponse
	if err := xml.Unmarshal(res.ReadAll(), &casResp); err != nil {
		g.Log().Error(ctx, "Parse SSO validation response error:", err)
		return nil, err
	}

	if casResp.Success.User == "" || casResp.Failure != "" {
		return &casResp, fmt.Errorf("invalid SSO session")
	}

	return &casResp, nil
}

// 添加认证中间件
func authMiddleware(r *ghttp.Request) {
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
		callbackURL := fmt.Sprintf("%s/callback", cfg.Service)
		_, err := validateSSOSession(r.Context(), ticket, callbackURL)
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
				group.Middleware(authMiddleware) // 添加认证中间件
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

				// callback 处理中相应修改
				group.GET("/callback", func(r *ghttp.Request) {
					ticket := r.GetQuery("ticket").String()
					if ticket == "" {
						r.Response.Write("Invalid CAS ticket")
						return
					}

					callbackURL := fmt.Sprintf("%s/callback", cfg.Service)
					casResp, err := validateSSOSession(r.Context(), ticket, callbackURL)
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
