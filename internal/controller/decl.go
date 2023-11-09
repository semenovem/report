package controller

import (
	"github.com/labstack/echo/v4"
	"github.com/semenovem/report/config"
	"github.com/semenovem/report/internal/lg"
	"github.com/semenovem/report/internal/provider"
	"net/http"
	"time"
)

const (
	stockFBS         = "fbs"
	stockFBO         = "fbo"
	stockCrossBorder = "crossborder"

	AccessCodeName = "access_code"
	SessionIDName  = "session_id"
)

type requestSession struct {
	createAt     time.Time
	fileDownload bool
}

type Controller struct {
	config   *config.Main
	logger   *lg.Lg
	provider *provider.Provider
	sessions map[string]*requestSession
}

func New(c *config.Main, l *lg.Lg, p *provider.Provider) *Controller {
	o := &Controller{
		config:   c,
		logger:   l,
		provider: p,
		sessions: make(map[string]*requestSession),
	}

	go o.clear()

	return o
}

func (ct *Controller) errResponse(c echo.Context, err error) error {
	return c.JSON(http.StatusInternalServerError, err.Error())
}

func (ct *Controller) clear() {
	const period = time.Minute * 60

	for {
		<-time.After(period)

		now := time.Now().Add(-period)
		for k, s := range ct.sessions {
			if s.createAt.Before(now) {
				delete(ct.sessions, k)
			}
		}
	}
}
