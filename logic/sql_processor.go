package logic

import (
	"github.com/jingwanglong/cellnet"
	"github.com/davyxu/golog"
	"db-server/proto/dbproto"
	_"db-server/database"
	"db-server/database"
)

var log *golog.Logger = golog.New("SqlProcessor")


func InitMessageRegister(peer cellnet.Peer) {
	cellnet.RegisterMessage(peer, "dbproto.RunSql", func(ev *cellnet.Event) {
		log.Debugf("receive run sql")
		msg := ev.Msg.(*dbproto.RunSql)
		ack := dbproto.RunSqlResponse{
			OpId: msg.OpId,
		}
		ev.Send(&ack)
	})
	
	cellnet.RegisterMessage(peer, "dbproto.QuerySql", func(ev *cellnet.Event) {
		log.Debugf("receive query sql")
		var ack dbproto.QuerySqlResponse
		msg := ev.Msg.(*dbproto.QuerySql)
		sql := msg.Xml
		sqlHandler := database.GetQueryHandler(*sql)
		var paramList []interface{}
		for _, field := range msg.Params {
			log.Debugf("field: %s ", string(field.Value))
			paramList = append(paramList, string(field.Value))
		}
		rows, err := sqlHandler(paramList)
		ack.OpId = msg.OpId
		var ErrorMsg dbproto.ErrorInfo
		if err != nil {
			*ErrorMsg.Code = 1
			*ErrorMsg.Message = err.Error()
			*ack.Error = ErrorMsg
		}else{
			ack.Rows = rows
		}
		ev.Send(&ack)
	})
	
	cellnet.RegisterMessage(peer, "dbproto.BatchSqlList", func(ev *cellnet.Event) {
		log.Debugf("receive batch sql")
		msg := ev.Msg.(*dbproto.QuerySql)
		ack := dbproto.BatchSqlListResponse{
			OpId: msg.OpId,
		}
		ev.Send(&ack)
	})
}



