package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Conf struct {
	Host                string
	DbPath              string
	AllowedHosts        map[string]int
	MongoPath           string
	DefaultDB           string
	Syncdelay           int
	CacheExpirationTime int
}

var GlobalConf Conf

//任务列队路径
var TaskPath string

//数据存储路径
var DataPath string

//解析配置文件
func ParseConf(config string) {
	GOPATH := os.Getenv("GOPATH")
	if config == "" {
		//wd, _ := os.Getwd()
		config = GOPATH + "/src/mcache/conf/mcache.conf"
	}

	b, err := ioutil.ReadFile(config)
	if err != nil {

		fmt.Println("没有找到配置文件:", config)
		os.Exit(0)
	}
	err = json.Unmarshal(b, &GlobalConf)
	if err != nil {
		fmt.Println("配置文件解析错误:", config)
		os.Exit(0)
	}

	fmt.Println("初始化完成:")
	fmt.Printf("%#v\n", GlobalConf)

	InitFolder()

}

/**
 * 文件夹说明  complete已完成推送的任务 db 所有数据备份 qs 检索结果  pending 等待通过的任务 task 任务队列 err执行错误
 */
func InitFolder() {
	os.MkdirAll(GlobalConf.DbPath+"/cknum/", 0777)
	os.MkdirAll(GlobalConf.DbPath+"/complete/", 0777)
	os.MkdirAll(GlobalConf.DbPath+"/db/", 0777)
	os.MkdirAll(GlobalConf.DbPath+"/err/", 0777)
	os.MkdirAll(GlobalConf.DbPath+"/log/", 0777)
	os.MkdirAll(GlobalConf.DbPath+"/pending/", 0777)
	os.MkdirAll(GlobalConf.DbPath+"/qs/", 0777)
	os.MkdirAll(GlobalConf.DbPath+"/task/", 0777)
}
