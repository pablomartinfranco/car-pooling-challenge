package router

import "net/http"

type Route struct {
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

type Router struct {
	routes Routes
}

func New() *Router {
	return &Router{}
}

func (r *Router) Register(pattern string, handler http.HandlerFunc) {
	var route = Route{
		Method:      http.MethodGet,
		Pattern:     pattern,
		HandlerFunc: handler,
	}
	r.routes = append(r.routes, route)
}

func (r *Router) GetRoutes() Routes {
	return r.routes
}
