package model

import (
	// "labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

var (
	COL_ACTIVE_TYPE   = "active_types"
	COL_ACTIVE_REWARD = "rewards"
)

type ActiveType struct {
	Id   bson.ObjectId `bson:"_id"`
	Name string        `bson:"name"`
	Desc string        `bson:"desc"`
	Mask string        `bson:"mask"`
}

type ActiveReward struct {
	Id     bson.ObjectId `bson:"_id"`
	Name   string        `bson:"name"`
	Desc   string        `bson:"desc"`
	Reward string        `bson:"reward"`
}

// 查询所有的
func FindAllActiveTypes() []ActiveType {
	s := Session()
	defer s.Close()
	c := s.DB(DB_NAME).C(COL_ACTIVE_TYPE)
	result := []ActiveType{}
	c.Find(nil).All(&result)
	return result
}

// 查询所有的奖品
func FindAllActiveRewards() []ActiveReward {
	s := Session()
	defer s.Close()
	c := s.DB(DB_NAME).C(COL_ACTIVE_REWARD)
	result := []ActiveReward{}
	c.Find(nil).All(&result)
	return result
}
