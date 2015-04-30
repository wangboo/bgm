package model

import (
	"github.com/revel/revel"
	"labix.org/v2/mgo"
)

var (
	session *mgo.Session
	DB_NAME = "jiyu_development"
)

func Init() {
	sess, err := mgo.Dial("localhost/" + DB_NAME)
	if err != nil {
		revel.INFO.Printf("mongo dial err %s\n", err.Error())
		return
	}
	sess.SetMode(mgo.Monotonic, true)
	session = sess
	// mgo.SetLogger(revel.INFO)
	// mgo.SetDebug(true)
	// revel.INFO.Println("mgo debug on")
}

func Session() *mgo.Session {
	return session.Copy()
}

func DB(s *mgo.Session) *mgo.Database {
	return s.DB(DB_NAME)
}
