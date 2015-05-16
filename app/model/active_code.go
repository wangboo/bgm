package model

import (
	"fmt"
	"github.com/revel/revel"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

var (
	COL_ACTIVE_CODE = "active_codes"
)

type ActiveCode struct {
	Id            bson.ObjectId `bson:"_id"`
	Code          string        `bson:"code"`
	UseFlag       bool          `bson:"use_flag"`
	Times         int           `bson:"times"` // 批量使用的剩余次数
	ServerUserId  bson.ObjectId `bson:"server_user_id,omitempty"`
	ActiveBatchId bson.ObjectId `bson:"active_batch_id"`
}

func (c *ActiveCode) String() string {
	return fmt.Sprintf("Id=%v,Code=%s", c.Id, c.Code)
}

func ColActiveCode(s *mgo.Session) *mgo.Collection {
	return s.DB(DB_NAME).C(COL_ACTIVE_CODE)
}

// 找到改平台的所有code
func FindAllActiveCodes(id string) []string {
	s := Session()
	defer s.Close()
	c := ColActiveCode(s)
	hex := bson.ObjectIdHex(id)
	all := []ActiveCode{}
	c.Find(bson.M{"active_batch_id": hex}).All(&all)
	codes := []string{}
	for _, c := range all {
		codes = append(codes, c.Code)
	}
	return codes
}

// 批量创建激活码
func CreateActiveCodes(acs []*ActiveCode) {
	revel.INFO.Println("create code size : ", len(acs))
	s := Session()
	defer s.Close()
	c := ColActiveCode(s)
	// batch := []*ActiveCode{}
	batch := []interface{}{}
	for _, ac := range acs {
		batch = append(batch, ac)
		if len(batch) >= 100 {
			err := c.Insert(batch...)
			if err != nil {
				revel.INFO.Printf("CreateActiveCodes err %s \n", err.Error())
			} else {
				revel.INFO.Printf("insert %d codes \n", len(batch))
			}
			// batch = []*ActiveCode{}
			batch = []interface{}{}
		}
	}
	if len(batch) > 0 {
		err := c.Insert(batch...)
		if err != nil {
			revel.ERROR.Println("save batch error ", err.Error())
		}
		revel.INFO.Printf("insert %d codes \n", len(batch))
		revel.INFO.Println("batch = ", batch)
	}
}
