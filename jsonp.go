package mango

import (
	"fmt"
	"regexp"
	"strings"
)

var jsonp_valid_callback_matcher *regexp.Regexp = regexp.MustCompile("^[a-zA-Z_$][a-zA-Z_0-9$]*([.]?[a-zA-Z_$][a-zA-Z_0-9$]*)*$")

func JSONP(env Env, app App) (Status, Headers, Body) {
	callback := env.Request().FormValue("callback")

	if callback != "" {
		if !jsonp_valid_callback_matcher.MatchString(callback) {
			return 400, Headers{"Content-Type": []string{"text/plain"}, "Content-Length": []string{"11"}}, "Bad Request"
		}
	}

	status, headers, body := app(env)

	if callback != "" && strings.Contains(headers.Get("Content-Type"), "application/json") {
		headers.Set("Content-Type", strings.Replace(headers.Get("Content-Type"), "json", "javascript", -1))
		body = Body(fmt.Sprintf("%s(%s)", callback, body))
		if headers.Get("Content-Length") != "" {
			headers.Set("Content-Length", fmt.Sprintf("%d", len(body)))
		}
	}

	return status, headers, body
}
