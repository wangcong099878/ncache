package model

import (
	"gopkg.in/mgo.v2"
	. "mcache/conf"
	//"gopkg.in/mgo.v2/bson"
	"fmt"
	"time"
)

//连接到mongo服务器
func Dao(CollectionName string) *mgo.Collection {
	//root:plk789@115.29.161.232:28010
	session, err := mgo.Dial(GlobalConf.MongoPath) //连接数据库
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	return session.DB(GlobalConf.DefaultDB).C(CollectionName)
}

type MongoUser struct {
	Uid      int  `uid`
	Userinfo User `info`
}

//从mongo还原数据
func MReduction(tabname string) {
	t1 := time.Now().UnixNano()
	fmt.Println("开始恢复数据")
	session, err := mgo.Dial(GlobalConf.MongoPath) //连接数据库
	defer session.Close()
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	bakMcache := session.DB(GlobalConf.DefaultDB).C(tabname)

	var items MongoUser
	iter := bakMcache.Find(nil).Iter()
	for iter.Next(&items) {
		addUser(items.Uid, items.Userinfo)
	}
	fmt.Println("恢复耗时：", (time.Now().UnixNano()-t1)/1000000000, "秒")
}

//pushtomongo
func PushToMongo(tabname string) {
	t1 := time.Now().UnixNano()
	var argsNew []interface{}
	var emptyn []interface{}
	//一次性只能插入100万数据
	for _, v := range unioy {
		for k, vv := range v {
			items := &MongoUser{k, vv}
			argsNew = append(argsNew, items)
			if len(argsNew) >= 100000 {
				if !pushItem(argsNew, tabname) {
					fmt.Println("发生异常1，写入耗时：", (time.Now().UnixNano()-t1)/1000000000, "秒")
					return
				} else {
					argsNew = emptyn
				}
			}
		}
	}
	if len(argsNew) > 0 {
		if !pushItem(argsNew, tabname) {
			fmt.Println("发生异常，写入耗时：", (time.Now().UnixNano()-t1)/1000000000, "秒")
			return
		}
	}
	fmt.Println("写入耗时：", (time.Now().UnixNano()-t1)/1000000000, "秒")
	return
}

func pushItem(argsNew []interface{}, tabname string) bool {
	fmt.Println("插入一次数据")
	session, err := mgo.Dial(GlobalConf.MongoPath) //连接数据库
	defer session.Close()
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	bakMcache := session.DB(GlobalConf.DefaultDB).C(tabname)

	err = bakMcache.Insert(argsNew...)
	time.Sleep(1 * time.Second)

	if err != nil {
		//上线后改为写入错误日志
		fmt.Println(err)
		return false
	} else {
		return true
	}
}

//同步数据
func SyncMongo() {
	getGender()
	getAge()
	getNmas()
}

//读取性别表
func getGender() {
	gender := Dao("gender")
	var gs []map[string]int
	gender.Find(nil).All(&gs)

	for _, v := range gs {
		Set("gender", v["uid"], v["gender"])
	}
}

//读取性别表
func getAge() {
	age := Dao("age")
	var gs []map[string]int
	age.Find(nil).All(&gs)
	for _, v := range gs {
		Set("age", v["uid"], v["age"])
	}
}

//读取用户属性
func getNmas() {
	nmas := Dao("nmas")

	type attrs struct {
		Uid  int    `uid`
		Attr string `attr`
	}
	var gs []attrs
	nmas.Find(nil).All(&gs)
	for _, v := range gs {
		Set(v.Attr, v.Uid, 0)
	}
}
