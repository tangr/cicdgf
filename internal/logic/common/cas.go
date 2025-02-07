package common

import (
	"context"
	"encoding/xml"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
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
	Cfg = CasConfig{
		ServerURL:   "https://sso-prod.yax.tech/cas",
		Service:     "http://localhost:8000",
		ValidateURL: "/serviceValidate",
		LoginURL:    "/login",
		LogoutURL:   "/logout",
	}
)

func ValidateSSOSession(ctx context.Context, ticket string, service string) (*CasResponse, error) {
	validateURL := fmt.Sprintf("%s%s?ticket=%s&service=%s",
		Cfg.ServerURL,
		Cfg.ValidateURL,
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
