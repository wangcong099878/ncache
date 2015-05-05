package action

import (
	"encoding/json"
	//"fmt"
	"io"
	. "mcache/model"
	. "mcache/storage"
	. "mcache/conf"
	"net/http"
	"strconv"
	"time"
)

func stoi(s string) int {
	a, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return a
}

//用户设置所有属性

//写入一行数据
func ApiAdd(w http.ResponseWriter, r *http.Request) {
	tabname := r.FormValue("tabname")
	key := stoi(r.FormValue("key"))
	val := stoi(r.FormValue("val"))
	Set(tabname, key, val)

	io.WriteString(w, Res(1, "设置成功"))
	return
}

//删除一行属性
func ApiGet(w http.ResponseWriter, r *http.Request) {
	uidstr := r.FormValue("uids")

	var uids []int
	err := json.Unmarshal([]byte(uidstr), &uids)

	if err != nil {
		io.WriteString(w, Res(2, "参数错误", uidstr))
		return
	}

	io.WriteString(w, Res(1, "获取成功", Get(uids)))
	return
}

//删除一行属性
func ApiDel(w http.ResponseWriter, r *http.Request) {
	tabname := r.FormValue("tabname")
	key := stoi(r.FormValue("key"))
	Del(tabname, key)
	io.WriteString(w, Res(1, "删除完成"))
	return
}

//查询数据   返回人员数量    并将结果集写入文件
func ApiSearch(w http.ResponseWriter, r *http.Request) {
	genderstr := r.FormValue("gender")
	agestr := r.FormValue("age")
	attrstr := r.FormValue("attr")
	filterstr := r.FormValue("filter") //过滤条件
	cache := r.FormValue("cache")      //过滤条件  "true"  "false"

	md5 := Md5(genderstr + agestr + attrstr + filterstr)

	if cache == "true" {
		tr := GetQs(md5)
		if tr.Md5 == md5 {
			t := time.Now().Unix()
			time := int(t)
			etime := time - tr.Time
			if etime < GlobalConf.CacheExpirationTime {
				io.WriteString(w, Res(1, "查询成功", md5, tr.Count, "is_cache"))
				return
			}
		}
	}
	//gender int, age []int, attr []string, filters []map[int]int
	//解析数据
	gender := stoi(genderstr)
	var age []int
	var attr []string
	var filter []string

	json.Unmarshal([]byte(agestr), &age)
	json.Unmarshal([]byte(attrstr), &attr)
	json.Unmarshal([]byte(filterstr), &filter)

	//查找用户
	count := Search(gender, age, attr, filter, md5)

	io.WriteString(w, Res(1, "查询成功", md5, count))
	return
}

//选择想要的人数   并且过滤想要过滤的人员
func ApiGetCid(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")
	num := stoi(r.FormValue("num"))
	//生成唯一key
	cidkey := Uniqid()
	go GetCid(key, num, cidkey)
	io.WriteString(w, Res(1, "ok", cidkey))
	return
}

//保存一个活动 ApiSaveAct
func ApiSaveAct(w http.ResponseWriter, r *http.Request) {
	cidkey := r.FormValue("key")
	aid := stoi(r.FormValue("aid"))

	tr := GetQs(cidkey)
	tr.Aid = aid

	//将tr存入文件
	WF(Gob_encode(tr), cidkey, "qs")

	io.WriteString(w, Res(1, "保存活动结果成功"))
	return
}

//后台审核通过
func ApiPassAct(w http.ResponseWriter, r *http.Request) {
	cidkey := r.FormValue("key")

	//将任务文件移动到推送区
	if Copy("qs", cidkey, "task") {
		io.WriteString(w, Res(1, "审核通过成功"))
	} else {
		io.WriteString(w, Res(2, "通过错误"))
	}
	return
}

//获取一个任务结果集
func ApiGetQs(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")
	tr := GetQs(key)
	res := make(map[string]interface{})
	res["Gender"] = tr.Gender
	res["Age"] = tr.Age
	res["Attr"] = tr.Attr
	res["Md5"] = tr.Md5
	//res["Uids"] = tr.Uids
	res["Count"] = len(tr.Uids)
	res["Endtime"] = tr.Endtime
	res["Time"] = tr.Time
	res["Filter"] = tr.Filter
	res["Aid"] = tr.Aid
	res["Num"] = tr.Num

	io.WriteString(w, Res(1, "ok", res))
	return
}

//查询并返回条目数量
func ApiCount(w http.ResponseWriter, r *http.Request) {
	i := Count()
	io.WriteString(w, Res(1, "获取用户信息成功", i))
	return
}

//强制push所有数据到磁盘
func ApiPush(w http.ResponseWriter, r *http.Request) {
	/*tabname := r.FormValue("tabname")
	key := r.FormValue("key")*/
	go Writelist()
	io.WriteString(w, Res(1, "发送成功"))
	return
}

//同步mongo数据库数据
func ApiSyncMongo(w http.ResponseWriter, r *http.Request) {
	SyncMongo()
	/*tabname := r.FormValue("tabname")
	key := r.FormValue("key")*/
	io.WriteString(w, Res(1, "导入成功"))
	return
}

func ApiGetUser(w http.ResponseWriter, r *http.Request) {
	uid := stoi(r.FormValue("uid"))
	io.WriteString(w, Res(1, "ok", GetUser(uid)))
	return
}

func ApiPushToMongo(w http.ResponseWriter, r *http.Request) {
	tabname := r.FormValue("tabname")
	if tabname == "" {
		tabname = "bakcache"
	}

	go PushToMongo(tabname)
	io.WriteString(w, Res(1, "请求执行中,数据将备份到："+tabname))
	return
}

func ApiMReduction(w http.ResponseWriter, r *http.Request) {
	tabname := r.FormValue("tabname")
	if tabname == "" {
		tabname = "bakcache"
	}
	go MReduction(tabname)
	io.WriteString(w, Res(1, "请求执行中,将从"+tabname+"恢复数据"))
	return
}

func ApiTestAdds(w http.ResponseWriter, r *http.Request) {
	num := stoi(r.FormValue("num"))
	go TestAdds(num)
	io.WriteString(w, Res(1, "请求成功正在创建数据"))
	return
}

//SetAttrs
//清空用户属性
func ApiSetAttrs(w http.ResponseWriter, r *http.Request) {
	uid := stoi(r.FormValue("uid"))
	attrstr := r.FormValue("attrs")
	attrs := make(map[string]int)

	err := json.Unmarshal([]byte(attrstr), &attrs)
	if err == nil {
		SetAttrs(uid, attrs)
	}
	io.WriteString(w, Res(1, "设置成功"))
	return
}

//清空用户属性
func ApiDelAttrs(w http.ResponseWriter, r *http.Request) {
	uid := stoi(r.FormValue("uid"))
	DelAttrs(uid)
	io.WriteString(w, Res(1, "清除成功"))
	return
}

//删除一个用户属性
func ApiDelAttr(w http.ResponseWriter, r *http.Request) {
	uid := stoi(r.FormValue("uid"))
	attr := r.FormValue("attr")
	DelAttr(uid, attr)
	io.WriteString(w, Res(1, "删除成功"))
	return
}
