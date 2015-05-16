package controllers

import (
	"github.com/revel/revel"
	"github.com/wangboo/bgm/app/model"
	"labix.org/v2/mgo/bson"
	"net/http"
	"time"
)

type Auth struct {
	*revel.Controller
}

const (
	AUTH_COOKIE_NAME = "bgm-auth"
)

// 登陆
func (c *Auth) Login(username, password string) revel.Result {
	revel.INFO.Printf("username = %s, password = %s \n", username, password)
	acc, ok := model.FindGmAccountByUsername(username)
	if !ok {
		revel.WARN.Println("找不到账号")
		return c.RenderJson(bson.M{"ok": false, "msg": "账号不存在"})
	}
	if acc.Password != password {
		return c.RenderJson(bson.M{"ok": false, "msg": "账号密码错误"})
	}
	acc.UpdateLoginAt()
	expiredTime := time.Hour * 8
	expiredAt := time.Now().Add(expiredTime)
	cookie := &http.Cookie{Name: AUTH_COOKIE_NAME, Expires: expiredAt, MaxAge: int(expiredTime.Seconds())}
	c.SetCookie(cookie)
	return c.RenderJson(bson.M{"ok": true})
}

// 验证cookie
func (c *Auth) Validate() revel.Result {
	now := time.Now().Second()
	for _, cookie := range c.Request.Cookies() {
		if cookie.Name == AUTH_COOKIE_NAME && cookie.Expires.Second() < now {
			return c.RenderJson(bson.M{"ok": true})
		}
	}
	return c.RenderJson(bson.M{"ok": false})
}

// var AuthFiler = func(c *revel.Controller, fc []revel.Filter) {
// 	ok, now := false, time.Now().Second()
// 	for _, cookie := range c.Request.Cookies() {
// 		if cookie.Name == AUTH_COOKIE_NAME && cookie.Expires.Seconds() < now {
// 			ok = true
// 		}
// 	}
// 	if ok {
// 		fc[0](c, fc[1:])
// 		return
// 	}
// 	// 权限验证失败
// 	c.Request.RequestURI = "/tmpl/login.html"
// }
