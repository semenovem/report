package router

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/semenovem/report/config"
	"github.com/semenovem/report/internal/controller"
	"github.com/semenovem/report/internal/lg"
	"github.com/semenovem/report/internal/provider"
	"io"
	"net/http"
	"strings"
)

func New(
	ctx context.Context,
	logger *lg.Lg,
	config *config.Main,
	provider *provider.Provider,
) (*Router, error) {
	var (
		ll = logger.Named("router")
		e  = echo.New()
	)

	echo.NotFoundHandler = func(c echo.Context) error {

		m := map[string]string{
			"code":        "404",
			"message":     "method didn't exists",
			"path":        c.Path(),
			"query":       c.QueryString(),
			"ParamNames":  strings.Join(c.ParamNames(), ", "),
			"ParamValues": strings.Join(c.ParamValues(), ", "),
			"headers":     fmt.Sprintf("%+v", c.Request().Header),
		}
		bodyResp, err := io.ReadAll(c.Request().Body)
		if err != nil {
			ll.Named("io.ReadAll").Error(err.Error())
		}

		m["body"] = string(bodyResp)

		return c.JSON(http.StatusNotFound, m)
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
		//middleware.Logger(),
		panicRecover(ll, config.Base.CliMode),
		middleware.CORSWithConfig(corsConfig),
	)

	cnt := controller.New(config, logger, provider)

	r := &Router{
		ctx:    ctx,
		logger: logger.Named("router"),
		server: e,
		config: config,
		cnt:    cnt,
	}

	accessMiddleware, err := accessTokenMiddleware(logger, config)
	if err != nil {
		ll.Named("accessTokenMiddleware").Debug(err.Error())
		return nil, err
	}

	r.auth = e.Group("", accessMiddleware)

	r.auth.GET("/:access_code", r.cnt.Index)
	r.auth.POST("/product/list", r.cnt.ReportProducts)

	return r, nil
}
