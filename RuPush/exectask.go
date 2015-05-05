package main

//处理任务队列
import (
	"flag"
	"gopkg.in/mgo.v2"
	"io/ioutil"
	. "mcache/conf"
	. "mcache/storage"
	//"os"
	"strconv"
	"time"
)

//定义ru表数据类型
type Ru struct {
	Rid    int `rid`
	Time   int `time`
	Status int `status`
	Uid    int `uid`
}

var sem = make(chan string)

//otsmid -c "F:\mygo\src\Middleware\1.conf"
var c *string = flag.String("c", "", "is configfile path")

func main() {
	flag.Parse()
	configpath := *c
	ParseConf(configpath)
	runtask()
}

//执行推送服务
func runtask() {
	FileLog("开始处理.")
	//去掉for循环就是每天定期执行一次
	for {
		go exectask()
		curr := <-sem
		FileLog("处理完成："+curr)
	}
}

func addToMongo(argsNew []interface{}) bool {
	session, err := mgo.Dial(GlobalConf.MongoPath) //连接数据库
	defer session.Close()
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	ruModel := session.DB(GlobalConf.DefaultDB).C("ru")
	err = ruModel.Insert(argsNew...)
	time.Sleep(1 * time.Second)
	if err != nil {
		//上线后改为写入错误日志
		FileLog(err.Error())
		return false
	} else {
		return true
	}
}

//向数据库推送   rid uids
func Push(aid int, uids map[int]int) bool {
	if aid == 0 {
		return false
	}
	var argsNew []interface{}
	var emptyn []interface{}
	
	//一次性只能插入100万数据
	for uid, _ := range uids {
		var ruitem Ru
		ruitem.Uid = uid
		ruitem.Rid = aid
		ruitem.Time = int(time.Now().Unix())
		ruitem.Status = 0
		argsNew = append(argsNew, ruitem)
		delete(uids, uid)
		if len(argsNew) >= 1000000 {
			if !addToMongo(argsNew){
				return false
			}else{
				argsNew = emptyn
			}
		}
	}
	
	if len(argsNew) > 0{
		if !addToMongo(argsNew){
			return false
		}
	}
	return true
}

func exectask() {
	for {
		currentTask := readtask()
		if currentTask != "" {
			time.Sleep(1 * time.Second)
			FileLog("开始处理:" + currentTask)
			b := GF(currentTask, "task")
			var tr TaskResult
			Gob_decode(b, &tr)
			sign := Push(tr.Aid, tr.Uids)
			if sign {
				Log("处理完成", currentTask, tr.Aid, tr.Num)
				Copy("task", currentTask, "complete")
				DF(currentTask, "task")
			} else {
				FileLog("异常:" + currentTask)
				//移动到异常文件夹
				Copy("task", currentTask, "err")
				DF(currentTask, "task")
			}
			sem <- currentTask
			return
		} else {
			//FileLog("全部处理完成.")
			//os.Exit(0)
		}
	}
}

//不间断处理
func readtask() string {
	fileinfos, err := ioutil.ReadDir(GlobalConf.DbPath + "/task")

	if err != nil {
		return ""
	}
	for _, v := range fileinfos {
		return v.Name()
	}
	return ""
}

//组装日志格式  任务名  活动id   推送人数  处理总耗时  日志命名:Y-m-d H
func Log(sign string, taskname string, aid int, num int) {
	rid := strconv.Itoa(aid)
	snum := strconv.Itoa(num)
	str := taskname + " 当前活动：" + rid + ";当前推送人数：" + snum
	FileLog(sign + ":" + str)
}
