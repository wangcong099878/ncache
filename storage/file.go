package storage

import (
	"bytes"
	"crypto/md5"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	. "mcache/conf"
	"mcache/uuid"
	"os"
	"time"
)

type Task struct {
	Tabname string
	Key     string
	Data    string
	Where   string
}

//任务列队格式
type TaskResult struct {
	Gender     int      //性别
	Age      []int      //年龄
	Attr   []string      //查询条件中间的关系
	Filter  []string    //是否有过滤
	Md5     string      //查询条件的md5编码值
	Uids    map[int]int //所有用户的uid
	Count   int         //uids的长度
	Endtime int         //过期时间  24小时
	Time    int         //创建时间
	Aid     int         //活动id
	Num     int         //厂商选择的人员数量
}

//根据文件名获取一个key值并返回
func GetQs(key string) TaskResult {
	var taskResult TaskResult
	err := Gob_decode(GF(key, "qs"), &taskResult)
	if err != nil {
		FileLog(err.Error())
	}
	return taskResult
}

//根据文件名获取一个key值并返回检索的人员
func GetQsUsers(key string) map[int]int {
	var taskResult TaskResult
	err := Gob_decode(GF(key, "qs"), &taskResult)
	if err != nil {
		FileLog(err.Error())
	}
	return taskResult.Uids
}

func Md5(s string) string {
	h := md5.New()
	h.Write([]byte(s))                    // 需要加密的字符串为 123456
	return hex.EncodeToString(h.Sum(nil)) // 输出加密结果

}

func Res(args ...interface{}) string {
	res := make(map[string]interface{})

	res["code"] = args[0]
	res["msg"] = args[1]
	if len(args) > 2 {
		res["data"] = args[2]
	}
	if len(args) > 3 {
		res["param1"] = args[3]
	}
	if len(args) > 4 {
		res["param2"] = args[4]
	}
	return Json_encode(res)
}

func Json_encodebyte(obj interface{}) []byte {
	i, err := json.Marshal(obj)
	if err != nil {
		return i
	}
	return i
}

func Json_encode(obj interface{}) string {
	i, err := json.Marshal(obj)
	if err != nil {
		return ""
	}
	return string(i)
}

func Json_decode(i string, obj interface{}) bool {
	err := json.Unmarshal([]byte(i), &obj)
	if err != nil {
		return false
	}
	return true
}

//Gob编码
func Gob_encode(data interface{}) []byte {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(data)
	if err != nil {
		var i []byte
		return i
	}
	return buf.Bytes()
}

//Gob解码
func Gob_decode(data []byte, to interface{}) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	return dec.Decode(to)
}

//生成唯一文件名
func Uniqid() string {
	id := uuid.New()
	return id.String()
}

//写入一个随机命名文件到日志文件夹
func F(data []byte) bool {
	filename := Uniqid()
	//写入一个文件
	err := ioutil.WriteFile(GlobalConf.DbPath+"/"+filename, data, 0644)
	if err != nil {
		return false
	}
	return true
}

//写入到指定文件名和文件夹
func WF(data []byte, filename string, folder string) bool {
	//写入一个文件
	err := ioutil.WriteFile(GlobalConf.DbPath+"/"+folder+"/"+filename, data, 0644)
	if err != nil {
		return false
	}
	return true
}

//从指定文件夹读取整个文件
func GF(filename string, folder string) []byte {
	b, err := ioutil.ReadFile(GlobalConf.DbPath + "/" + folder + "/" + filename)
	if err != nil {
		return b
	}
	return b
}

//删除指定文件夹的文件
func DF(filename string, folder string) {
	os.Remove(GlobalConf.DbPath + "/" + folder + "/" + filename)
}

//拷贝一个文件  到指定文件夹   任务最终完成移动到complete
func Copy(oldpath string, filename string, path string) bool {
	b, err := ioutil.ReadFile(GlobalConf.DbPath + "/" + oldpath + "/" + filename)
	if err != nil {
		return false
	}
	err = ioutil.WriteFile(GlobalConf.DbPath+"/"+path+"/"+filename, b, 0644)
	if err != nil {
		return false
	}
	return true
}

func Exist(filename string,path string) bool {
	_, err := os.Stat(GlobalConf.DbPath+"/"+path+"/"+filename)
	return err == nil || os.IsExist(err)
}

func FileLog(str string) {
	day := time.Now().Format("2006-01-02")
	timestr := "[" + time.Now().Format("2006-01-02 15:04:05") + "] "
	filepath := GlobalConf.DbPath+"/log/"+day+".log"
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0777)
	defer file.Close()
	if err != nil {
		return
	}
	file.WriteString(timestr + str + "\r\n")
}
