package controller

import (
	"github.com/labstack/echo/v4"
	"github.com/semenovem/report/config"
	"github.com/semenovem/report/internal/action"
	"github.com/semenovem/report/internal/zoo/lg"
	"net/http"
	"time"
)

const (
	AccessCodeName = "access_code"
	SessionIDName  = "session_id"
)

type requestSession struct {
	createAt     time.Time
	fileDownload bool
}

type Controller struct {
	config     *config.Main
	logger     *lg.Lg
	sessions   map[string]*requestSession
	dataMining *action.DataMining
}

func New(c *config.Main, l *lg.Lg, a *action.DataMining) *Controller {
	o := &Controller{
		config:     c,
		logger:     l,
		dataMining: a,
		sessions:   make(map[string]*requestSession),
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
