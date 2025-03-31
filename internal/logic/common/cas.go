package common

import (
	"context"
	"encoding/xml"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

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

type CasConfig struct {
	CasServerURL string `json:"CasServerURL"`
	ServiceURL   string `json:"serviceURL"`
	ValidateURL  string `json:"validateURL"`
	LoginURL     string `json:"loginURL"`
	LogoutURL    string `json:"logoutURL"`
}

type AdminConfig struct {
	AuthorEnable bool     `json:"author_enable"`
	AdminUsers   []string `json:"admin_users"`
}

type ContextKey string

const (
	UsernameKey ContextKey = "username"
	UseridKey   ContextKey = "userid"
)

var (
	Cfg      *CasConfig
	AdminCfg *AdminConfig
)

func init() {
	ctx := gctx.New()
	if err := g.Cfg().MustGet(ctx, "cas-sso").Scan(&Cfg); err != nil {
		panic(err)
	}
	if err := g.Cfg().MustGet(ctx, "admin").Scan(&AdminCfg); err != nil {
		panic(err)
	}
}

func ValidateSSOSession(ctx context.Context, ticket string, service string) (*CasResponse, error) {
	validateURL := fmt.Sprintf("%s%s?ticket=%s&service=%s",
		Cfg.CasServerURL,
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

	var casResp CasResponse
	if err := xml.Unmarshal(res.ReadAll(), &casResp); err != nil {
		g.Log().Error(ctx, "Parse SSO validation response error:", err)
		return nil, err
	}

	g.Log().Debug(ctx, "cas username:", casResp.Success.User)

	if casResp.Success.User == "" || casResp.Failure != "" {
		return &casResp, fmt.Errorf("invalid SSO session")
	}

	return &casResp, nil
}
