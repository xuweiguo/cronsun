package bootstrap

import (
	"net"
	"time"

	"github.com/cockroachdb/cmux"

	"cronsun"
	"cronsun/conf"
	"cronsun/db/entries"
	"cronsun/event"
	"cronsun/log"
	"cronsun/web"
)

func RunWeb(confFile, network string, level int) error {
	if err := initLogger(level); err != nil {
		return err
	}

	if err := cronsun.Init(confFile, true); err != nil {
		log.Errorf(err.Error())
		return err
	}
	entries.EnsureJobLogIndex()

	listener, err := net.Listen(checkNetworkProtocol(network), conf.Config.Web.BindAddr)
	log.Debugf("listen %s:%s", checkNetworkProtocol(network), conf.Config.Web.BindAddr)
	if err != nil {
		log.Errorf(err.Error())
		return err
	}

	mux := cmux.New(listener)
	httpListener := mux.Match(cmux.HTTP1Fast())
	httpServer, err := web.InitServer()
	if err != nil {
		log.Errorf(err.Error())
		return err
	}

	if conf.Config.Mail.Enable {
		var noticer cronsun.Noticer
		if len(conf.Config.Mail.HttpAPI) > 0 {
			noticer = &cronsun.HttpAPI{}
		} else {
			mailer, err := cronsun.NewMail(30 * time.Second)
			if err != nil {
				log.Errorf(err.Error())
				return err
			}
			noticer = mailer
		}
		go cronsun.StartNoticer(noticer)
	}

	period := int64(conf.Config.Web.LogCleaner.EveryMinute)
	var stopCleaner func(interface{})
	if period > 0 {
		closeChan := web.RunLogCleaner(time.Duration(period)*time.Minute, time.Duration(conf.Config.Web.LogCleaner.ExpirationDays)*time.Hour*24)
		stopCleaner = func(i interface{}) {
			close(closeChan)
		}
	}

	go func() {
		if err := httpServer.Serve(httpListener); err != nil {
			panic(err.Error())
		}
	}()

	go mux.Serve()

	log.Infof("cronsun web server started on %s, Ctrl+C or send kill sign to exit", conf.Config.Web.BindAddr)
	event.On(event.EXIT, conf.Exit, stopCleaner)
	event.Wait()
	event.Emit(event.EXIT, nil)
	log.Infof("exit success")
	return nil
}

func checkNetworkProtocol(p string) string {
	switch p {
	case "ipv4":
		return "tcp4"
	case "ipv6":
		return "tcp6"
	}
	return "tcp"
}
