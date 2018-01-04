package main

import (
	"fmt"
	"log"

	"github.com/kardianos/service"
	"github.com/jingwanglong/cellnet"
	"github.com/jingwanglong/cellnet/socket"
	"db-server/config"
	"db-server/logic"
	_ "db-server/database"
)

var logger service.Logger

type program struct{}

func RunServer(){
	queue := cellnet.NewEventQueue()

	host := fmt.Sprintf("%s:%d", config.Host["ip"], int(config.Host["port"].(float64)))
	peer := socket.NewAcceptor(queue).Start(host)
	peer.SetName("server")
	logic.InitMessageRegister(peer)
	queue.StartLoop()
	queue.Wait()
	peer.Stop()
}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}
func (p *program) run() {
	// Do work here
	RunServer()
}
func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	return nil
}

func main() {
	svcConfig := &service.Config{
		Name:        "web-db-server.exe",
		DisplayName: "YoMailWeb-DBServer(Go)",
		Description: "为web服务器提供数据支持",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	logger, err = s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}
	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}
