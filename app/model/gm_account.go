package model

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
)

type GmAccount struct {
	Id       bson.ObjectId `bson:"_id"`
	Username string        `bson:"username"`
	Password string        `bson:"password"`
	LoginAt  time.Time     `bson:"login_at"`
}

const (
	COL_GM_ACCOUNT = "gm_accounts"
)

func ColGmAccount(s *mgo.Session) *mgo.Collection {
	return s.DB(DB_NAME).C(COL_GM_ACCOUNT)
}

// 更新登陆时间
func (g *GmAccount) UpdateLoginAt() {
	s := Session()
	defer s.Close()
	c := ColGmAccount(s)
	c.UpdateId(g.Id, bson.M{"login_at": time.Now()})
}

// 查询账号 返回 账号，是否找到
func FindGmAccountByUsername(name string) (*GmAccount, bool) {
	s := Session()
	defer s.Close()
	c := ColGmAccount(s)
	result := &GmAccount{}
	if err := c.Find(bson.M{"username": name}).One(result); err == mgo.ErrNotFound {
		return nil, false
	}
	return result, true
}
