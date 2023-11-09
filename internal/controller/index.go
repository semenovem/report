package controller

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/semenovem/report/config"
	"os"
	"strings"
	"time"
)

func (ct *Controller) Index(c echo.Context) error {
	var (
		ctx = c.Request().Context()
		ll  = ct.logger.Func(ctx, "Index")
	)

	dat, err := os.ReadFile(ct.config.HTMLDir + config.IndexHTMLFileName)
	if err != nil {
		ll.Named("os.ReadFile").Error(err.Error())
		return ct.errResponse(c, err)
	}

	var (
		body      = string(dat)
		sessionID = uuid.NewString()
	)
	ct.sessions[sessionID] = &requestSession{createAt: time.Now()}

	body = strings.ReplaceAll(body, "{{access_code}}", ct.config.AccessCode)
	body = strings.ReplaceAll(body, "{{session_id}}", sessionID)

	ll.Debug("index data send")

	return c.HTML(200, body)
}
