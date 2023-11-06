package router

func (r *Router) addRoutes() {
	var (
		auth = r.auth
	)

	auth.GET("/", r.cnt.Index)

	auth.GET("/product/list", r.cnt.ReportProducts)
}
