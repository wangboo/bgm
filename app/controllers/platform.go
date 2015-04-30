package controllers

import (
	"fmt"
	"github.com/revel/revel"
	"github.com/wangboo/bgm/app/model"
	"labix.org/v2/mgo/bson"
)

type Platform struct {
	*revel.Controller
}

// 查询所有平台
func (c *Platform) All() revel.Result {
	platforms := model.FindAllPlatform()
	return c.RenderJson(platforms)
}

func (c *Platform) Find(id string) revel.Result {
	// id := c.Params.Get("id")
	platform := model.FindPlatform(id)
	ss := model.FindAllServerState()
	return c.RenderJson(bson.M{"platform": platform, "ss": ss})
}

// 通过游戏服务器id查询平台id
func (c *Platform) Pid(sid string) revel.Result {
	server, ok := model.FindServer(sid)
	if !ok {
		return c.RenderText(`{"pid":"0"}`)
	}
	rst := fmt.Sprintf(`{"pid":"%s"}`, server.PlatformId.Hex())
	return c.RenderText(rst)
}
