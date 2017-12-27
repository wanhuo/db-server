package logic

import (
	"github.com/jingwanglong/cellnet"
	"github.com/davyxu/golog"
	"db-server/proto/dbproto"
	"db-server/database"
)

var log *golog.Logger = golog.New("SqlProcessor")

func InitMessageRegister(peer cellnet.Peer) {
	cellnet.RegisterMessage(peer, "dbproto.RunSql", func(ev *cellnet.Event) {
		var ack dbproto.RunSqlResponse
		msg := ev.Msg.(*dbproto.RunSql)
		sqlAlias := *msg.Xml
		err := database.ProcessRunSql(sqlAlias, msg.Params)
		var ErrorMsg dbproto.ErrorInfo
		if err != nil {
			*ErrorMsg.Code = 1
			*ErrorMsg.Message = err.Error()
			ack.Error = &ErrorMsg
		}
		ack.OpId = msg.OpId
		ev.Send(&ack)
	})
	
	cellnet.RegisterMessage(peer, "dbproto.QuerySql", func(ev *cellnet.Event) {
		var ack dbproto.QuerySqlResponse
		msg := ev.Msg.(*dbproto.QuerySql)
		sqlAlias := *msg.Xml
		rows, err := database.ProcessQuerySql(sqlAlias, msg.Params)
		var ErrorMsg dbproto.ErrorInfo
		if err != nil {
			*ErrorMsg.Code = 1
			*ErrorMsg.Message = err.Error()
			ack.Error = &ErrorMsg
		}else{
			ack.Rows = rows
		}
		ack.OpId = msg.OpId
		ev.Send(&ack)
	})
	
	cellnet.RegisterMessage(peer, "dbproto.BatchSqlList", func(ev *cellnet.Event) {
		var ack dbproto.BatchSqlListResponse
		msg := ev.Msg.(*dbproto.BatchSqlList)
		for _, sql := range msg.Sql {
			var sqlResult dbproto.BatchSqlListResponse_OneSqlResult
			var err error
			sqlAlias := *sql.Xml
			sqlResult.SqlId = sql.SqlId
			if *sql.IsQuery {
				var rows []*dbproto.OneRow
				rows, err = database.ProcessQuerySql(sqlAlias, sql.Params)
				if err == nil {
					sqlResult.Rows = rows
				}
			}else{
				err = database.ProcessRunSql(sqlAlias, sql.Params)
			}
			var ErrorMsg dbproto.ErrorInfo
			if err != nil {
				*ErrorMsg.Code = 1
				*ErrorMsg.Message = err.Error()
				sqlResult.Error = &ErrorMsg
			}
			ack.Result = append(ack.Result, &sqlResult)
		}
		ev.Send(&ack)
	})
}

