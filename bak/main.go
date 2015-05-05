package main

import (
	//"Middleware/uuid"
	"fmt"
	"io/ioutil"
	"os"
	"time"
	"bytes"
	"encoding/gob"
	"encoding/json"
	"runtime"
	_"unsafe"
)


var rootpath string
var datapath string

//json处理
func Json_encode(obj interface{}) []byte {
	var i []byte
	i, err := json.Marshal(obj)
	if err != nil {
		return i
	}
	return i
}

func Json_decode(i []byte, obj interface{}) bool {
	err := json.Unmarshal(i, &obj)
	if err != nil {
		return false
	}
	return true
}
//Gob处理
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

func Gob_decode(data []byte, to interface{}) error {
    buf := bytes.NewBuffer(data)
    dec := gob.NewDecoder(buf)
    return dec.Decode(to)
}

//写入一个文件
func W(filename string,data []byte) {
	//写入一个文件
	err := ioutil.WriteFile(datapath+filename, data, 0644)
	if err != nil {
		return
	}
}

//读取一个文件
func R(filename string) []byte{
	var b []byte
	//读取整个文件
	b, err := ioutil.ReadFile(datapath+filename)
	if err != nil {
		return b
	}
	return b
}


var mass [][]string


//json占用系统空间更大   100万 json 11秒   Gob 4秒    1000万 json 125秒   Gob 48秒    
func main(){
	runtime.GOMAXPROCS(runtime.NumCPU())

	t1 := time.Now().UnixNano()

	rootpath, _ := os.Getwd()
	datapath = rootpath + "/data/"
	
	t2 := time.Now().UnixNano()
	//组装数据
	makedate()
	fmt.Println("组装数据耗时：",(time.Now().UnixNano()-t2)/1000000000, "秒")
	
	
	
	//fmt.Println(unsafe.Sizeof(mass))
	
	t3 := time.Now().UnixNano()
	//编码数据并写入文件
	data:=Gob_encode(mass)
	fmt.Println("编码数据耗时：",(time.Now().UnixNano()-t3)/1000000000, "秒")
	
	t4 := time.Now().UnixNano()
	W("data",data)
	fmt.Println("写入文件耗时：",(time.Now().UnixNano()-t4)/1000000000, "秒")
	
	
	//读取文件并解码数据
	t5 := time.Now().UnixNano()
	newdata:=R("data")
	fmt.Println("读取文件耗时：",(time.Now().UnixNano()-t5)/1000000000, "秒")
	
	
	t6 := time.Now().UnixNano()
	var tmas [][]string
	Gob_decode(newdata,&tmas)
	fmt.Println("解码数据耗时：",(time.Now().UnixNano()-t6)/1000000000, "秒")
	
	//fmt.Println(tmas)
	fmt.Println("总耗时：",(time.Now().UnixNano()-t1)/1000000000, "秒")
}

func makedate(){
	//组装速度最快
	mas:=[]string{"a","b","c","d","e","f","g","h","i","j","k","l","m","n","o","p","q","r","s","t"}
	//数组形式组装数据比map快
	for i:=1;i<=10000000;i++{
		mass=append(mass,mas)
	}
	fmt.Println(len(mass))
}