package router

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

func tokenMiddleware(
	logger *lg.Lg,
) echo.MiddlewareFunc {
	ll := logger.Named("tokenMiddleware")

	ll.Infof("start")

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			//ll.Infof(c.Request().Header.Get("Authorization"))
			return next(c)
		}
	}
}
