package mango

import (
	"regexp"
	"sort"
)

// Methods required by sort.Interface.
type matcherArray []*regexp.Regexp

func specificity(matcher *regexp.Regexp) int {
	return len(matcher.String())
}
func (this matcherArray) Len() int {
	return len(this)
}
func (this matcherArray) Less(i, j int) bool {
	// The sign is reversed below so we sort the matchers in descending order
	return specificity(this[i]) > specificity(this[j])
}
func (this matcherArray) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

func Routing(routes map[string]App) Middleware {
	matchers := matcherArray{}
	handlers := []App{}

	// Compile the matchers
	for matcher, _ := range routes {
		// compile the matchers
		matchers = append([]*regexp.Regexp(matchers), regexp.MustCompile(matcher))
	}

	// sort 'em by descending length
	sort.Sort(matchers)

	// Attach the handlers to each matcher
	for _, matcher := range matchers {
		// Attach them to their handlers
		handlers = append(handlers, routes[matcher.String()])
	}

	return func(env Env, app App) (Status, Headers, Body) {
		for i, matcher := range matchers {
			matches := matcher.FindStringSubmatch(env.Request().URL.Path)
			if len(matches) != 0 {
				// Matched a route; inject matches and return handler
				env["Routing.matches"] = matches
				return handlers[i](env)
			}
		}

		// didn't match any of the other routes. pass upstream.
		return app(env)
	}
}
