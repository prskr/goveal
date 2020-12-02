package routing

import (
	"net/http"
	"regexp"
)

type regexpRule struct {
	pattern *regexp.Regexp
	handler http.Handler
}

type RegexpRouter struct {
	rules []regexpRule
}

func (r *RegexpRouter) AddRule(pattern string, handler http.Handler) (err error) {
	var exp *regexp.Regexp
	if exp, err = regexp.Compile(pattern); err != nil {
		return
	}
	r.rules = append(r.rules, regexpRule{
		pattern: exp,
		handler: handler,
	})
	return
}

func (r *RegexpRouter) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	for idx := range r.rules {
		rule := r.rules[idx]
		if rule.pattern.MatchString(request.URL.Path) {
			rule.handler.ServeHTTP(writer, request)
			return
		}
	}
	writer.WriteHeader(404)
}
