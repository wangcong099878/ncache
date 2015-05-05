package model

/*
import (
	"strconv"
	//"time"
)

//所有查询条件md5存储
type WMD struct {
	Endtime  int    //过期时间戳
	Filename string //where条件的md5值
}

type WMDS map[string]WMD

//map[md5值]WMD    查询缓存    如果找到条件  并且超时就删除掉
var wmds WMDS
var unioy MiiDB

//读写锁
var rlock int
var wlock int

func init() {
	//初始化db
	unioy = make(MiiDB)
}

//设置一列值
func Set(tabname string, key int, val int) {
	if _, ok := unioy[tabname]; !ok {
		unioy[tabname] = make(Mii)
	}
	unioy[tabname][key] = val
}

//获取一列值  用户与前端打印
func Get(tabname string) Msi {
	msi := make(Msi)
	if mii, ok := unioy[tabname]; ok {
		for k, v := range mii {
			msi[strconv.Itoa(k)] = v
		}
	}
	return msi
}

//删除一列值
func Del(tabname string, key int) {
	if mii, ok := unioy[tabname]; ok {
		delete(mii, key)
		unioy[tabname] = mii
	}
}

//返回单表总条数
func Count(tabname string) int {
	if mii, ok := unioy[tabname]; ok {
		return len(mii)
	}
	return 0
}

//查询条件
type Where struct {
	Type   string //查询的类型  all返回全部  size按大小条件筛选  eq按等于筛选   index前缀查询    preg正则
	Params []int  //参数
}

type Wheres map[string]Where

//多条件处理  返回检索符合条件的结果集并集
func Search(and Wheres, or Wheres, assoc string) map[int]int {

	if len(and) == 0 && len(or) != 0 {
		return SearchOr(or)
	}
	if len(and) != 0 && len(or) == 0 {
		return SearchAnd(and)
	}

	mc := make(chan map[int]int, 2)

	go func(a Wheres, mc chan map[int]int) {
		mc <- SearchOr(or)
	}(or, mc)

	go func(a Wheres, mc chan map[int]int) {
		mc <- SearchAnd(and)
	}(or, mc)

	result := <-mc
	result2 := <-mc

	if assoc == "and" {
		return Intersection(result, result2)
	} else {
		for key, _ := range result2 {
			result[key] = 0
		}
	}
	return result
}

//处理并关系
func SearchAnd(a Wheres) map[int]int {
	result := make(map[int]int)
	if len(a) == 0 {
		return result
	}

	c := make(chan map[int]int, len(a))

	for tabname, w := range a {
		if w.Type == "size" {
			go Srange(tabname, w.Params[0], w.Params[1], c)
		}
		if w.Type == "eq" {
			go Find(tabname, w.Params[0], c)
		}
		if w.Type == "all" {
			go GetAll(tabname, c)
		}
	}

	for i := 0; i < len(a); i++ {
		res := <-c
		if i == 0 {
			result = res
		} else {
			result = Intersection(result, res)
		}
	}
	return result
}

//处理或关系
func SearchOr(a Wheres) map[int]int {
	result := make(map[int]int)
	if len(a) == 0 {
		return result
	}

	c := make(chan map[int]int, len(a))

	for tabname, w := range a {
		if w.Type == "size" {
			go Srange(tabname, w.Params[0], w.Params[1], c)
		}
		if w.Type == "eq" {
			go Find(tabname, w.Params[0], c)
		}
		if w.Type == "all" {
			go GetAll(tabname, c)
		}
	}
	for i := 0; i < len(a); i++ {
		res := <-c
		for key, _ := range res {
			result[key] = 0
		}
	}
	return result
}

//获取一列值
func GetAll(tabname string, c chan map[int]int) {
	msi := make(map[int]int)
	if mii, ok := unioy[tabname]; ok {
		msi = map[int]int(mii)
	}
	c <- msi
	return
}

//根据值的范围查找  返回所有符合的列
func Srange(tabname string, start int, end int, c chan map[int]int) {
	res := make(map[int]int)
	if mii, ok := unioy[tabname]; ok {
		for key, val := range mii {
			if val <= end && val >= start {
				res[key] = val
			}
		}
	}
	c <- res
	return
}

//根据val查找  所有val完全等于的列     多核心查找方案
func Find(tabname string, value int, c chan map[int]int) {
	res := make(map[int]int)
	if mii, ok := unioy[tabname]; ok {
		for key, val := range mii {
			if val == value {
				res[key] = val
			}
		}
	}
	c <- res
	return
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
func marges(args map[int]map[int]int) map[int]int {
	res := args[1]
	delete(args, 1)
	for _, v := range args {
		for key, _ := range v {
			res[key] = 0
		}
	}
	return res
}
*/
