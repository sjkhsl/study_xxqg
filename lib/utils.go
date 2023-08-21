package lib

import (
	"net/http"

	"github.com/playwright-community/playwright-go"
)

func CookieToJar(cookies []Cookie) []*http.Cookie {
	var (
		cooks []*http.Cookie
	)
	for _, c := range cookies {
		cooks = append(

			cooks,
			&http.Cookie{
				Name:     c.Name,
				Value:    c.Value,
				Path:     c.Path,
				Domain:   c.Domain,
				Secure:   c.Secure,
				HttpOnly: c.HTTPOnly,
				SameSite: http.SameSiteDefaultMode,
			},
		)
	}
	return cooks
}

func CookieToParam(cookies []Cookie) []playwright.BrowserContextAddCookiesOptionsCookies {
	var (
		cooks []playwright.BrowserContextAddCookiesOptionsCookies
	)

	for _, c := range cookies {
		domain := ".xuexi.cn"
		if c.Name == "acw_tc" || c.Name == "aliyungf_tc" {
			domain = "iflow-api.xuexi.cn\t"
		}
		cooks = append(cooks, playwright.BrowserContextAddCookiesOptionsCookies{
			Name:     playwright.String(c.Name),
			Value:    playwright.String(c.Value),
			Domain:   playwright.String(domain),
			Path:     playwright.String(c.Path),
			Expires:  playwright.Float(float64(c.Expires)),
			HttpOnly: playwright.Bool(c.HTTPOnly),
			Secure:   playwright.Bool(c.Secure),
			SameSite: playwright.SameSiteAttributeStrict,
		})
	}
	return cooks
}
