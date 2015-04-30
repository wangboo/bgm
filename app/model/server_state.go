package model

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

var (
	COL_SERVER_STATE = "server_states"
)

type ServerState struct {
	Id   bson.ObjectId `bson:"_id"`
	Name string        `bson:"name"`
	Desc string        `bson:"desc"`
	Show string        `bson:"show"`
}

func ColServerState(s *mgo.Session) *mgo.Collection {
	return s.DB(DB_NAME).C(COL_SERVER_STATE)
}

// 查询所有的服务器状态
func FindAllServerState() []ServerState {
	s := Session()
	defer s.Close()
	colSS := ColServerState(s)
	ss := []ServerState{}
	colSS.Find(nil).All(&ss)
	return ss
}
