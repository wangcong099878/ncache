package model

import (
	"fmt"
	. "hash/consistenthash"
	"strconv"
	//"math/rand"
	. "mcache/storage"
	"runtime"
	"time"
)

//所有查询条件md5存储
type WMD struct {
	Endtime  int    //过期时间戳
	Filename string //where条件的md5值
}

type WMDS map[string]WMD

//map[md5值]WMD    查询缓存    如果找到条件  并且超时就删除掉
var wmds WMDS

//var unioy MiiDB

type User struct {
	Gender int
	Age    int
	Attr   map[string]int
}

var unioy map[string]map[int]User
var hashmap *Map

//读写锁
var rlock int
var wlock int

func init() {
	//初始化db
	unioy = make(map[string]map[int]User)
	hashmap = New(1, nil)
	//定义节点列表
	for i := 1; i <= runtime.NumCPU(); i++ {
		is := strconv.Itoa(i)
		hashmap.Add(is)
		unioy[is] = make(map[int]User)
	}
}

//设置一列值
func Set(Aname string, uid int, val int) {
	key := hashmap.Get(strconv.Itoa(uid))
	var user User
	if suser, ok := unioy[key][uid]; ok {
		user = suser
	} else {
		user.Attr = make(map[string]int)
	}
	switch Aname {
	case "gender":
		user.Gender = val
		break
	case "age":
		user.Age = val
		break
	default:
		user.Attr[Aname] = val
	}
	unioy[key][uid] = user
}

//设置所有属性
func SetAttrs(uid int, attrs map[string]int) {
	key := hashmap.Get(strconv.Itoa(uid))
	if suser, ok := unioy[key][uid]; ok {
		suser.Attr = attrs
		unioy[key][uid] = suser
	}
}

//清空用户属性
func DelAttrs(uid int) {
	key := hashmap.Get(strconv.Itoa(uid))
	if suser, ok := unioy[key][uid]; ok {
		suser.Attr = make(map[string]int)
		unioy[key][uid] = suser
	}
}

//删除一个用户属性
func DelAttr(uid int, attr string) {
	key := hashmap.Get(strconv.Itoa(uid))
	if suser, ok := unioy[key][uid]; ok {
		delete(suser.Attr, attr)
		unioy[key][uid] = suser
	}
}

//设置一个用户的信息
func addUser(uid int, user User) {
	key := hashmap.Get(strconv.Itoa(uid))
	unioy[key][uid] = user
}

//根据多个uid获取多个用户
func Get(uids []int) []MongoUser {
	var mongoUsers []MongoUser
	for _, v := range uids {
		user := MongoUser{v, GetUser(v)}
		mongoUsers = append(mongoUsers, user)
	}
	return mongoUsers
}

//删除一列值
func Del(tabname string, key int) {
	fmt.Println(tabname, key)
	/*if mii, ok := unioy[tabname]; ok {
		delete(mii, key)
		unioy[tabname] = mii
	}*/
}

//返回总用户数
func Count() int {
	i := 0
	for _, v := range unioy {
		i += len(v)
	}
	return i
}

func ResToFile(gender int, age []int, attr []string, filter []string, md5 string, Uids []map[int]int, count int) {
	var tr TaskResult

	t := time.Now().Unix()
	time := int(t)

	us := marges(Uids)

	tr.Gender = gender
	tr.Age = age
	tr.Attr = attr
	tr.Md5 = md5
	tr.Uids = us
	tr.Count = count
	tr.Endtime = time + 86400
	tr.Time = time
	tr.Filter = filter
	tr.Aid = 0
	tr.Num = 0
	//将tr存入文件
	WF(Gob_encode(tr), md5, "qs")

}

func GetCid(key string, num int, cidkey string) {
	tr := GetQs(key)
	tr.Num = num
	uids := make(map[int]int)
	j := 0
	for k, v := range tr.Uids {
		if j >= num {
			break
		}
		uids[k] = v
		j++
	}
	tr.Uids = uids
	//将tr存入文件
	WF(Gob_encode(tr), cidkey, "qs")
}

func Search(gender int, age []int, attr []string, filter []string, md5 string) int {
	var filters []map[int]int
	if len(filter) > 0 {
		//获取需要过滤的用户
		for _, v := range filter {
			us := GetQsUsers(v)
			filters = append(filters, us)
		}
	}

	mc := make(chan map[int]int, runtime.NumCPU())

	for i := 1; i <= runtime.NumCPU(); i++ {
		is := strconv.Itoa(i)
		go SIten(gender, age, attr, filters, mc, is)
	}

	var results []map[int]int

	countnum := 0
	for i := 1; i <= runtime.NumCPU(); i++ {
		res := <-mc
		results = append(results, res)
		countnum += len(res)
	}
	//将结果异步合并写入文件
	go ResToFile(gender, age, attr, filter, md5, results, countnum)
	return countnum
}

//gender 2表示不限   年龄少于两个参数表示不限   attr可以传多个需要的属性      filter
func SIten(gender int, age []int, attr []string, filters []map[int]int, mc chan map[int]int, dbnum string) {
	res := make(map[int]int)
	sign := false
	for k, v := range unioy[dbnum] {
		sign = false
		if gender != 2 {
			if gender != v.Gender {
				continue
			}
		}
		if len(age) > 1 {
			if v.Age > age[1] || v.Age < age[0] {
				continue
			}
		}

		if len(filters) > 0 {
			for _, filter := range filters {
				if _, ok := filter[k]; ok {
					sign = true
					break
				}
			}
			if sign {
				continue
			}
		}
		if len(attr) > 0 {
			for _, vv := range attr {
				if _, ok := v.Attr[vv]; ok {
					res[k] = 0
					break
				}
			}
		} else {
			res[k] = 0
		}

	}
	mc <- res
}

//根据uid获取一个用户
func GetUser(uid int) User {
	var user User
	key := hashmap.Get(strconv.Itoa(uid))
	if suser, ok := unioy[key][uid]; ok {
		return suser
	}
	return user
}

//获取两个map的交集
func Intersection(a map[int]int, b map[int]int) map[int]int {
	res := make(map[int]int)
	for key, _ := range a {
		if _, ok := b[key]; ok {
			res[key] = 0
		}
	}
	return res
}

//获取两个map的差集   过滤a中的元素   获取不存在b中的元素
func Diff(a map[int]int, b map[int]int) map[int]int {
	for key, _ := range b {
		if _, ok := a[key]; ok {
			delete(a, key)
		}
	}
	return a
}

//合并多个map
func marges(args []map[int]int) map[int]int {
	res, args := args[len(args)-1], args[:len(args)-1]
	for _, v := range args {
		for key, _ := range v {
			res[key] = 0
		}
	}
	return res
}
