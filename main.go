package main

import (
	"fmt"
	"db-server/config"
	"db-server/logic"
	"github.com/jingwanglong/cellnet/socket"
	"github.com/jingwanglong/cellnet"
	"github.com/davyxu/golog"
	_ "db-server/database"
)

var log *golog.Logger = golog.New("main")

func main() {
	queue := cellnet.NewEventQueue()

	host := fmt.Sprintf("%s:%d", config.Host["ip"], int(config.Host["port"].(float64)))
	peer := socket.NewAcceptor(queue).Start(host)
	peer.SetName("server")

	logic.InitMessageRegister(peer)

	queue.StartLoop()

	queue.Wait()

	peer.Stop()
}
