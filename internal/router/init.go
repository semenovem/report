package router

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/semenovem/report/config"
	"github.com/semenovem/report/internal/action"
	"github.com/semenovem/report/internal/controller"
	"github.com/semenovem/report/internal/lg"
	"net/http"
)

func New(
	ctx context.Context,
	logger *lg.Lg,
	config *config.Main,
	bl *action.BL,
) (*Router, error) {
	var (
		ll = logger.Named("router")
		e  = echo.New()
	)

	echo.NotFoundHandler = func(c echo.Context) error {
		return c.JSON(http.StatusNotFound, panicMessage{
			Code:    http.StatusNotFound,
			Message: "method didn't exists",
		})
	}

	corsConfig := middleware.CORSConfig{
		Skipper:          middleware.DefaultSkipper,
		AllowOrigins:     []string{"*"},
		AllowCredentials: true,
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
			echo.HeaderXFrameOptions,
			echo.HeaderXContentTypeOptions,
			echo.HeaderContentSecurityPolicy,
		},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodHead,
			http.MethodPut,
			http.MethodPatch,
			http.MethodPost,
			http.MethodDelete,
			http.MethodOptions,
		},
		MaxAge: 60,
	}

	e.Use(
		middleware.Logger(),
		panicRecover(ll, config.Base.CliMode),
		middleware.CORSWithConfig(corsConfig),
	)

	cnt := controller.New(config, logger, bl)

	r := &Router{
		ctx:    ctx,
		logger: logger.Named("router"),
		server: e,
		config: config,
		cnt:    cnt,
	}

	r.unauth = e.Group("")
	r.auth = r.unauth.Group("", tokenMiddleware(logger))

	r.addRoutes()

	return r, nil
}
