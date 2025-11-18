package bootstrap

import (
	"cronsun"
	"cronsun/conf"
	"cronsun/event"
	"cronsun/log"
	"cronsun/node"
)

func RunNode(confFile string, level int) error {
	if err := initLogger(level); err != nil {
		return err
	}

	if err := cronsun.Init(confFile, true); err != nil {
		log.Errorf(err.Error())
		return err
	}

	n, err := node.NewNode(conf.Config)
	if err != nil {
		log.Errorf(err.Error())
		return err
	}

	if err = n.Register(); err != nil {
		log.Errorf(err.Error())
		return err
	}

	if err = cronsun.StartProc(); err != nil {
		log.Warnf("[process key will not timeout]proc lease id set err: %s", err.Error())
	}

	if err = n.Run(); err != nil {
		log.Errorf(err.Error())
		return err
	}

	log.Infof("cronsun %s service started, Ctrl+C or send kill sign to exit", n.String())
	event.On(event.EXIT, n.Stop, conf.Exit, cronsun.Exit)
	event.On(event.WAIT, cronsun.Reload)
	event.Wait()
	event.Emit(event.EXIT, nil)
	log.Infof("exit success")
	return nil
}
