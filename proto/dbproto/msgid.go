// Generated by db-server/protoc-gen-msg
// DO NOT EDIT!
// Source: db.server.proto

package dbproto

import (
	"github.com/jingwanglong/cellnet"
	"reflect"
	_ "github.com/jingwanglong/cellnet/codec/pb"
)

func init() {

	// db.server.proto

	cellnet.RegisterMessageMeta("pb", "dbproto.RunSql", reflect.TypeOf((*RunSql)(nil)).Elem(), 11)

	cellnet.RegisterMessageMeta("pb", "dbproto.QuerySql", reflect.TypeOf((*QuerySql)(nil)).Elem(), 22)

	cellnet.RegisterMessageMeta("pb", "dbproto.BatchSqlList", reflect.TypeOf((*BatchSqlList)(nil)).Elem(), 33)

	cellnet.RegisterResponseMsgMeta("pb", "dbproto.RunSqlResponse", reflect.TypeOf((*RunSqlResponse)(nil)).Elem(), 11)

	cellnet.RegisterResponseMsgMeta("pb", "dbproto.QuerySqlResponse", reflect.TypeOf((*QuerySqlResponse)(nil)).Elem(), 22)

	cellnet.RegisterResponseMsgMeta("pb", "dbproto.BatchSqlListResponse", reflect.TypeOf((*BatchSqlListResponse)(nil)).Elem(), 33)

}
