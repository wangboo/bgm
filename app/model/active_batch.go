package model

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
)

var (
	COL_ACTIVE_BATCH = "active_batches"
)

type ActiveBatch struct {
	Id            bson.ObjectId `bson:"_id"`
	Name          string        `bson:"name"`         // 批次名字
	Desc          string        `bson:"desc"`         // 批次描述
	BeginTime     time.Time     `bson:"begin_time"`   // 开始时间
	EndTime       time.Time     `bson:"end_time"`     //结束时间
	IsMuti        bool          `bson:"is_muti"`      //是否能多次使用
	AllPlatform   bool          `bson:"all_platform"` //是否对所有平台有效
	AllServer     bool          `bson:"all_server"`   // 是否能对所有平台下服务器使用
	LimTimes      int           `bson:"lim_times"`    //多次使用时，可使用次数
	RewardId      bson.ObjectId `bson:"reward_id"`
	ActiveTypeId  bson.ObjectId `bson:"active_type"`
	PlatformMasks []string      `bson:"platform_masks"` //所属平台
	ZoneIds       []int         `bson:"zone_ids"`       // 所属服务器
}

func ColActiveBatch(s *mgo.Session) *mgo.Collection {
	return s.DB(DB_NAME).C(COL_ACTIVE_BATCH)
}

func FindAllActiveBatches() []ActiveBatch {
	s := Session()
	defer s.Close()
	c := ColActiveBatch(s)
	result := []ActiveBatch{}
	c.Find(nil).All(&result)
	return result
}

func CreateActiveBatch(ab *ActiveBatch) {
	s := Session()
	defer s.Close()
	c := ColActiveBatch(s)
	ab.Id = bson.NewObjectId()
	c.Insert(ab)
}
