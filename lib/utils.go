package lib

import (
	"net/http"

	"github.com/mxschmitt/playwright-go"
)

func cookieToJar(cookies []Cookie) []*http.Cookie {
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
			},
		)
	}
	return cooks
}

func cookieToParam(cookies []Cookie) []playwright.SetNetworkCookieParam {
	var (
		cooks []playwright.SetNetworkCookieParam
	)
	for _, c := range cookies {
		cooks = append(cooks, playwright.SetNetworkCookieParam{
			Name:     c.Name,
			Value:    c.Value,
			Domain:   playwright.String(c.Domain),
			Path:     playwright.String(c.Path),
			Expires:  playwright.Int(c.Expires),
			HttpOnly: playwright.Bool(c.HTTPOnly),
			Secure:   playwright.Bool(c.Secure),
			SameSite: playwright.String(c.SameSite),
		})
	}
	return cooks
}
