package main

import (
	"flag"
	"io"
	"log"
	. "mcache/action"
	. "mcache/conf"
	"net/http"
	"runtime"
	"time"
)

type myHandler struct{}

var mux map[string]func(w http.ResponseWriter, r *http.Request)
var allowedHosts map[string]int

//otsmid -c "F:\mygo\src\Middleware\1.conf"
var c *string = flag.String("c", "", "is configfile path")

func main() {
	//初始化配置文件
	flag.Parse()
	configpath := *c
	ParseConf(configpath)

	runtime.GOMAXPROCS(runtime.NumCPU())

	s := &http.Server{
		Addr:           GlobalConf.Host,
		Handler:        &myHandler{},
		ReadTimeout:    200 * time.Second,
		WriteTimeout:   200 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	mux = make(map[string]func(w http.ResponseWriter, r *http.Request))
	mux["/add"] = ApiAdd
	mux["/get"] = ApiGet
	mux["/del"] = ApiDel
	mux["/search"] = ApiSearch
	mux["/getqs"] = ApiGetQs
	mux["/getcid"] = ApiGetCid
	mux["/saveact"] = ApiSaveAct
	mux["/passact"] = ApiPassAct
	
	mux["/getuserinfo"] = ApiGetUser

	//返回所有用户数量
	mux["/count"] = ApiCount
	//删除属性
	mux["/setattrs"] = ApiSetAttrs
	mux["/delattrs"] = ApiDelAttrs
	mux["/delattr"] = ApiDelAttr

	//备份与测试
	mux["/syncmongo"] = ApiSyncMongo
	//备份数据到硬盘
	mux["/push"] = ApiPush
	//添加测试用户
	mux["/testadds"] = ApiTestAdds
	
	//从mongo备份与恢复
	mux["/pushtomongo"] = ApiPushToMongo
	mux["/mreduction"] = ApiMReduction

	allowedHosts = GlobalConf.AllowedHosts
	allowedHosts["127.0.0.1"] = 1

	log.Fatal(s.ListenAndServe())
}

func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, ok := allowedHosts[r.Host]; ok {
		w.WriteHeader(403)
		io.WriteString(w, "非法请求"+r.URL.String())
		return
	}

	if h, ok := mux[r.URL.String()]; ok {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		h(w, r)
		return
	}

	io.WriteString(w, r.URL.String())
	return
}
