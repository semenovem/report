package router

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/semenovem/report/config"
	"github.com/semenovem/report/internal/controller"
	"github.com/semenovem/report/internal/lg"
	"net/http"
	"runtime"
)

type panicMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func panicRecover(logger *lg.Lg, cli bool) echo.MiddlewareFunc {
	ll := logger.Named("panicRecover")

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				if r := recover(); r != nil {
					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}

					var (
						stack  = make([]byte, middleware.DefaultRecoverConfig.StackSize)
						length = runtime.Stack(stack, !middleware.DefaultRecoverConfig.DisableStackAll)
						msg    = fmt.Sprintf("[PANIC RECOVER] %v %s\n", err, stack[:length])
					)

					if cli {
						fmt.Println(msg)
					} else {
						ll.Error(msg)
					}

					_ = c.JSON(http.StatusInternalServerError, panicMessage{
						Code:    http.StatusInternalServerError,
						Message: "Internal Server Error",
					})
				}
			}()

			return next(c)
		}
	}
}

func accessTokenMiddleware(
	logger *lg.Lg,
	configMain *config.Main,
) (echo.MiddlewareFunc, error) {
	var (
		ll = logger.Named("accessTokenMiddleware")
	)

	ll.Infof("start")

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			// Для индексной страницы
			if c.Path() == "/:access_code" {
				if configMain.AccessCode == c.Param(controller.AccessCodeName) {
					return next(c)
				}
			} else {
				// Для всех остальных
				if configMain.AccessCode == c.FormValue(controller.AccessCodeName) {
					return next(c)
				}
			}

			logClient(c, ll).Info("forbidden")
			return c.HTML(403, "forbidden")
		}
	}, nil
}

func logClient(c echo.Context, ll *lg.Lg) *lg.Lg {
	return ll.With("ip", c.RealIP()).
		With("user-agent", c.Request().UserAgent()).
		With("query", c.QueryString()).
		With("path", c.Path()).
		With("method", c.Request().Method)
}
