package controllers

import (
	"github.com/revel/revel"
	"github.com/wangboo/bgm/app/model"
	"labix.org/v2/mgo/bson"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type Ac struct {
	*revel.Controller
}

// 创建
func (c *Ac) Create(all, allServer, isMuti bool, name, rType, reward, beginAt, endAt, desc, prefix string,
	size, mutiTimes int, pids, sids string) revel.Result {
	revel.INFO.Println("prefix ", prefix)
	ab := model.ActiveBatch{Name: name, Desc: desc}
	beginAtTime, err := time.Parse("2006-01-02", beginAt)
	if err != nil {
		return c.RenderJson(bson.M{"ok": false, "msg": "开始时间错误，格式必须为：年年年年-月月-日日"})
	}
	endAtTime, err := time.Parse("2006-01-02", endAt)
	if err != nil {
		return c.RenderJson(bson.M{"ok": false, "msg": "结束时间错误，格式必须为：年年年年-月月-日日"})
	}
	ab.BeginTime = beginAtTime
	ab.EndTime = endAtTime
	// 勾选平台判断
	ab.AllPlatform = all
	if !all {
		// 勾选平台
		pidArr := strings.Split(pids, ",")
		// revel.INFO.Println("ab.pidArr = ", pidArr)
		ab.PlatformMasks = model.FindPlatformMaskByIds(pidArr)
		// revel.INFO.Println("ab.PlatformMasks = ", ab.PlatformMasks)
	}
	// 勾选服务器判断
	ab.AllServer = allServer
	if !allServer {
		revel.INFO.Printf("sids = %s\n", sids)
		sidArr := strings.Split(sids, ",")
		ab.ZoneIds = model.FindServerZoneIdByIds(sidArr)
	}
	ab.ActiveTypeId = bson.ObjectIdHex(rType)
	ab.RewardId = bson.ObjectIdHex(reward)
	// 可以次使用判定
	ab.IsMuti = isMuti
	if isMuti {
		ab.LimTimes = mutiTimes
	}
	model.CreateActiveBatch(&ab)
	// 创建 ActiveCode
	codes := createCode(prefix, size)
	acs := []*model.ActiveCode{}
	for _, code := range codes {
		ac := &model.ActiveCode{
			Id:            bson.NewObjectId(),
			ActiveBatchId: ab.Id,
			UseFlag:       false,
			Times:         mutiTimes,
			Code:          code,
		}
		acs = append(acs, ac)
	}
	model.CreateActiveCodes(acs)
	return c.RenderJson(bson.M{"ok": true})
}

func (c *Ac) FindAllActiveBatches() revel.Result {
	all := model.FindAllActiveBatches()
	return c.RenderJson(all)
}

// 查询平台下所有激活码
func (c *Ac) FindAllActiveCodes(id string) revel.Result {
	codes := model.FindAllActiveCodes(id)
	return c.RenderJson(codes)
}

// 查询所有的类型
func (c *Ac) FindAllActiveTypes() revel.Result {
	types := model.FindAllActiveTypes()
	return c.RenderJson(types)
}

func (c *Ac) FindAllActiveRewards() revel.Result {
	rewards := model.FindAllActiveRewards()
	return c.RenderJson(rewards)
}

// 批量创建激活码
func createCode(prefix string, size int) []string {
	rand.Seed(time.Now().UnixNano())
	arr := []string{}
	rst := map[string]bool{}
	if size == 0 {
		return arr
	}
	for {
		dst := "" + prefix
		for j := 0; j < 8; j++ {
			n := rand.Intn(10)
			dst = dst + strconv.Itoa(n)
		}
		if _, ok := rst[dst]; ok {
			continue
		} else {
			rst[dst] = true
			if len(rst) > size {
				break
			}
			arr = append(arr, dst)
		}
	}
	return arr
}
