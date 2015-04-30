package model

import (
	"time"
	// "labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

// 权限
type BgmAuth struct {
	Id   bson.ObjectId `bson:"_id"`
	Uri  string        `bson:"uri"`
	Name string        `bson:"name"`
}

// 用户
type BgmUser struct {
	Id       bson.ObjectId `bson:"_id"`
	Username string        `bson:"username"`
	Password string        `bson:"password"`
	BgmAuth  []string      `bson:"auth"` // 权限
}

type BgmHistory struct {
	Id        bson.ObjectId `bson:"_id"`
	UserId    bson.ObjectId `bson:"userId"`
	Log       string
	CreatedAt time.Time
}
