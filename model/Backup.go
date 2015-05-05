package model

import (
	"io/ioutil"
	. "mcache/conf"
	. "mcache/storage"
	"time"
)

func DataGc(step int) {
	timer := time.NewTicker(time.Duration(step) * time.Second)
	for {
		select {
		case <-timer.C:
			go Writelist()
		}
	}
}

func Writelist() {
	if rlock == 0 {
		for tabname, _ := range unioy {
			pushdb(tabname)
		}
	}
}

func Readlist() {
	fileinfos, err := ioutil.ReadDir(GlobalConf.DbPath + "/db/")
	if err != nil {
		return
	}
	for _, v := range fileinfos {
		item := readdb(v.Name())
		unioy[v.Name()] = item
	}
}

func pushdb(tabname string) {
	if data, ok := unioy[tabname]; ok {
		databyte := Gob_encode(data)
		WF(databyte, tabname, "db")
	}
}

func readdb(tabname string) map[int]User {
	item := make(map[int]User)
	databyte := GF(tabname, "db")
	Gob_decode(databyte, &item)
	return item
}
