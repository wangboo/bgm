package model

import (
	// "github.com/revel/revel"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
)

var (
	COL_SERVER = "servers"
)

// 服务器
type Server struct {
	Id            bson.ObjectId `bson:"_id"`
	WorkState     int           `bson:"work_state"`
	Name          string        `bson:"name"`
	Desc          string        `bson:"desc"`
	Ip            string        `bson:"ip"`
	Port          int           `bson:"port"`
	ZoneId        int           `bson:"zone_id"`
	PlatformId    bson.ObjectId `bson:"platform_id"`
	ServerStateId bson.ObjectId `bson:"server_state_id"`
	UpdatedAt     time.Time     `bson:"updated_at"`
	CreatedAt     time.Time     `bson:"created_at"`

	platform Platform
}

func ColServer(s *mgo.Session) *mgo.Collection {
	return s.DB(DB_NAME).C(COL_SERVER)
}

// 查找游戏服务器
func FindServer(id string) (*Server, bool) {
	s := Session()
	defer s.Close()
	colServer := ColServer(s)
	server := &Server{}
	err := colServer.FindId(bson.ObjectIdHex(id)).One(server)
	if err != nil {
		return server, false
	}
	return server, true
}

// 通过服务器id查询服务器
func FindServersByIds(ids []string) []Server {
	s := Session()
	defer s.Close()
	colServer := ColServer(s)
	hexIds := []bson.ObjectId{}
	for _, id := range ids {
		hexIds = append(hexIds, bson.ObjectIdHex(id))
	}
	servers := []Server{}
	colServer.Find(bson.M{"_id": bson.M{"$in": hexIds}}).All(&servers)
	return servers
}

// 通过平台id集合查询服务器
func FindServerByPlatformIds(ids []string) []Server {
	s := Session()
	defer s.Close()
	colServer := ColServer(s)
	hexIds := []bson.ObjectId{}
	for _, id := range ids {
		hexIds = append(hexIds, bson.ObjectIdHex(id))
	}
	servers := []Server{}
	colServer.Find(bson.M{"platform_id": bson.M{"$in": hexIds}}).All(&servers)
	return servers
}

// 通过平台id查询maskId 集合
func FindServerZoneIdByIds(ids []string) []int {
	s := Session()
	defer s.Close()
	c := ColServer(s)
	hexIds := []bson.ObjectId{}
	for _, id := range ids {
		hexIds = append(hexIds, bson.ObjectIdHex(id))
	}
	all := &[]Server{}
	c.Find(bson.M{"_id": bson.M{"$in": hexIds}}).Select(bson.M{"zone_id": 1}).All(all)
	zoneIds := []int{}
	for _, s := range *all {
		zoneIds = append(zoneIds, s.ZoneId)
	}
	return zoneIds
}
