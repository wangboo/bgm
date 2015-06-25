package model

import (
	"encoding/json"
	"fmt"
	"github.com/revel/revel"
	"io/ioutil"
	"net/http"
	"strings"
)

// 游戏服务器请求map/reduce任务
type GameServerMapReduce struct {
	// 任务分工的服务器
	Servers []Server
	// 请求地址
	Uri string
}

type MapResponse struct {
	id     string
	result interface{}
}

func (m *MapResponse) String() string {
	return fmt.Sprintf("id=%s,result=%v\n", m.id, m.result)
}

// 从游戏服务器获取数据
func GetDataFromGameServer(s *Server, uri string) interface{} {
	url := fmt.Sprintf("http://%s:%d/jiyu/admin%s", s.Ip, s.Port, uri)
	revel.INFO.Printf("GetDataFromGameServer %s \n", url)
	resp, err := http.Get(url)
	if err != nil {
		revel.ERROR.Printf("请求游戏服务器失败:%s\n", url)
		return fail(err.Error())
	}
	bytes, _ := ioutil.ReadAll(resp.Body)
	revel.INFO.Printf("服务器应答:%s\n", bytes)
	var data interface{}
	str := string(bytes)
	if strings.HasPrefix(str, "[") {
		data = []map[string]interface{}{}
	} else {
		data = map[string]interface{}{}
	}
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return fail(err.Error())
	}
	return data
}

func worker(uri string, s *Server, c chan *MapResponse) {
	revel.INFO.Printf("Worker 启动\n")
	data := GetDataFromGameServer(s, uri)
	resp := &MapResponse{id: s.Id.Hex(), result: data}
	c <- resp
}

// 失败
func fail(reason string) map[string]interface{} {
	return map[string]interface{}{"ok": false, "reason": reason}
}

// 任务执行
func (g *GameServerMapReduce) DoIt() map[string]interface{} {
	channel := make(chan *MapResponse)
	for _, s := range g.Servers {
		go worker(g.Uri, &s, channel)
	}
	reduce := map[string]interface{}{}
	count := 0
	size := len(g.Servers)
	for {
		select {
		case resp := <-channel:
			revel.INFO.Printf("收到工作线程应答%v\n", resp)
			reduce[resp.id] = resp.result
		}
		count += 1
		if count >= size {
			break
		}
	}
	return reduce
}
