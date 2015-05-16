package controllers

import (
	"encoding/base64"
	"fmt"
	"github.com/revel/revel"
	"github.com/wangboo/bgm/app/model"
	"net/url"
	"strings"
	"time"
)

type Query struct {
	*revel.Controller
}

// 通过名字和服务器ids查询服务器
func (c *Query) FindUserByName(name, ids string) revel.Result {
	serverIds := strings.Split(ids, ",")
	servers := model.FindServersByIds(serverIds)
	revel.INFO.Printf("find servers %v \n", servers)
	name = base64.StdEncoding.EncodeToString([]byte(name))
	name = url.QueryEscape(name)
	uri := fmt.Sprintf("/rest/findUserByName?name=%s", name)
	revel.INFO.Printf("uri : %s\n", uri)
	mr := &model.GameServerMapReduce{Servers: servers, Uri: uri}
	data := mr.DoIt()
	return c.RenderJson(data)
}

// 查询个人信息
func (c *Query) FindInfo(sid, uid string) revel.Result {
	uri := fmt.Sprintf("/rest/info?userId=%s", uid)
	data := FindDataByUri(sid, uri)
	return c.RenderJson(data)
}

// 查询阵容
func (c *Query) FindGroup(sid, uid string) revel.Result {
	uri := fmt.Sprintf("/rest/findGroup?userId=%s", uid)
	data := FindDataByUri(sid, uri)
	return c.RenderJson(data)
}

// 查询玩家道具
func (c *Query) FindUserProp(sid, uid string) revel.Result {
	uri := fmt.Sprintf("/rest/findUserProp?userId=%s", uid)
	data := FindDataByUri(sid, uri)
	return c.RenderJson(data)
}

// 查询个人充值信息
func (c *Query) FindUserCharge(sid, uid string) revel.Result {
	uri := fmt.Sprintf("/rest/findUserCharge?userId=%s", uid)
	data := FindDataByUri(sid, uri)
	return c.RenderJson(data)
}

// 查询玩家邮件和聊天信息
func (c *Query) UserSetting(sid, uid string) revel.Result {
	uri := fmt.Sprintf("/rest/findUserSetting?userId=%s", uid)
	data := FindDataByUri(sid, uri)
	return c.RenderJson(data)
}

// 更新玩家聊天消息
func (c *Query) UpdateChat(id int, sid, msg string) revel.Result {
	revel.INFO.Printf("id = %d, msg = %s\n", id, msg)
	msg = base64.StdEncoding.EncodeToString([]byte(msg))
	msg = url.QueryEscape(msg)
	uri := fmt.Sprintf("/rest/updateChat?id=%d&msg=%s", id, msg)
	FindDataByUri(sid, uri)
	return c.RenderText("ok")
}

// 删除玩家聊天
func (c *Query) DeleteChat(id int, sid string) revel.Result {
	uri := fmt.Sprintf("/rest/deleteChat?id=%d", id)
	data := FindDataByUri(sid, uri)
	return c.RenderJson(data)
}

// 删除玩家邮件
func (c *Query) DeleteMail(id int, sid string) revel.Result {
	uri := fmt.Sprintf("/rest/deleteMail?id=%d", id)
	data := FindDataByUri(sid, uri)
	return c.RenderJson(data)
}

func FindDataByUri(sid, uri string) interface{} {
	server, ok := model.FindServer(sid)
	if !ok {
		revel.ERROR.Printf("没有找到服务器%s\n", sid)
		return fail("找不到服务器")
	}
	data := model.GetDataFromGameServer(server, uri)
	return data
}

// 登陆限制
func (c *Query) ForbiddenLogin(login bool, sid, loginTime string, uid int) revel.Result {
	revel.INFO.Printf("login %v, sid=%s, loginTme %s \n", login, sid, loginTime)
	var sec int64
	if login {
		ltime, err := time.Parse("2006-01-02 03:04", loginTime)
		if err != nil {
			return c.RenderJson(fail("日期格式错误,请严格按照(年年年年-月月-日日 时时:分分)格式！"))
		}
		revel.INFO.Printf("ltime = %v \n", ltime)
		sec = ltime.Unix() * 1000
	}
	url := fmt.Sprintf("/rest/forbiddenLogin?uid=%d&login=%v&loginTime=%d", uid, login, sec)
	revel.INFO.Printf("request url = %s \n", url)
	FindDataByUri(sid, url)
	return c.RenderJson(success())
}

// 登陆聊天
func (c *Query) ForbiddenChat(chat bool, sid, chatTime string, uid int) revel.Result {
	var sec int64
	if chat {
		ltime, err := time.Parse("2006-01-02 03:04", chatTime)
		if err != nil {
			return c.RenderJson(fail("日期格式错误,请严格按照(年年年年-月月-日日 时时:分分)格式！"))
		}
		sec = ltime.Unix() * 1000
	}
	url := fmt.Sprintf("/rest/forbiddenChat?uid=%d&chat=%v&chatTime=%d", uid, chat, sec)
	revel.INFO.Printf("request url = %s \n", url)
	FindDataByUri(sid, url)
	return c.RenderJson(success())
}

func success() map[string]bool {
	return map[string]bool{"ok": true}
}

func fail(reason string) map[string]interface{} {
	return map[string]interface{}{"ok": false, "reason": reason}
}
