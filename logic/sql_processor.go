package logic

import (
	"github.com/jingwanglong/golog"
	"db-server/database"
	"db-server/proto/dbproto"
	"github.com/jingwanglong/cellnet"
)

var log *golog.Logger = golog.New("SqlProcessor")

func ErrorMsg(errorInfo string) *dbproto.ErrorInfo {
	var errorMsg dbproto.ErrorInfo
	*errorMsg.Code = 1
	*errorMsg.Message = errorInfo
	return &errorMsg
}

func InitMessageRegister(peer cellnet.Peer) {
	cellnet.RegisterMessage(peer, "dbproto.RunSql", func(ev *cellnet.Event) {
		var ack dbproto.RunSqlResponse
		msg := ev.Msg.(*dbproto.RunSql)
		ack.OpId = msg.OpId
		sqlAlias := *msg.Xml
		err := database.ProcessRunSql(sqlAlias, msg.Params)
		if err != nil {
			ack.Error = ErrorMsg(err.Error())
		}
		ev.Send(&ack)
	})
	
	cellnet.RegisterMessage(peer, "dbproto.QuerySql", func(ev *cellnet.Event) {
		var ack dbproto.QuerySqlResponse
		msg := ev.Msg.(*dbproto.QuerySql)
		ack.OpId = msg.OpId
		sqlAlias := *msg.Xml
		rows, err := database.ProcessQuerySql(sqlAlias, msg.Params)
		if err != nil {
			ack.Error = ErrorMsg(err.Error())
		}else{
			ack.Rows = rows
		}
		ev.Send(&ack)
	})
	
	cellnet.RegisterMessage(peer, "dbproto.BatchSqlList", func(ev *cellnet.Event) {
		var ack dbproto.BatchSqlListResponse
		msg := ev.Msg.(*dbproto.BatchSqlList)
		ack.OpId = msg.OpId
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
			if err != nil {
				sqlResult.Error = ErrorMsg(err.Error())
			}
			ack.Result = append(ack.Result, &sqlResult)
		}
		ev.Send(&ack)
	})
}

