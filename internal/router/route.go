package router

func (r *Router) addRoutes() {
	var (
		auth = r.auth
	)

	auth.GET("/report1", r.cnt.Report1)
}
