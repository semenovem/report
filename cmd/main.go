package main

import (
	"context"
	"fmt"
	"github.com/semenovem/report/config"
	"github.com/semenovem/report/internal/lg"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var (
		ctx, cancel = context.WithCancel(context.Background())
		sig         = make(chan os.Signal)
		ll, setter  = lg.New()
	)

	defer func() {
		cancel()
		ll.Info("exiting")
		fmt.Println("exiting")
		time.Sleep(time.Millisecond * 500)
	}()

	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sig
		cancel()
	}()

	cfg, err := config.ParseAPI()
	if err != nil {
		ll.Named("env.Parse").Errorf("can't parse env: ", err)
		cancel()
		return
	}

	setter.SetCli(true)
	setter.SetShowTime(true)
	setter.SetLevel(cfg.Base.LogLevel)

	if err = newApp(ctx, ll, cfg); err != nil {
		_ = ll.NestedWith(err, "can't start app")
		//cancel()
	}

	<-ctx.Done()
}
