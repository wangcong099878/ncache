package main

/*import (
	"gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"
	//"github.com/astaxie/beego"
	"fmt"
)



//连接到mongo服务器
func Dao(CollectionName string) *mgo.Collection {
	session, err := mgo.Dial("root:plk789@115.29.161.232:28010") //连接数据库
	//session, err := mgo.Dial(beego.AppConfig.String("MDB_HOST"))  //连接数据库
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	return session.DB("unioy").C(CollectionName)
}

//同步数据
func main(){
	getGender()
	getAge()
	getNmas()
}


//读取性别表
func getGender(){
	gender := Dao("gender")
    //实例化一个博文对象
    //实例化主控服务器对象
    var gs []map[string]int
    //有则更新  无则插入
    gender.Find(nil).All(&gs)

    fmt.Println("读取")
    fmt.Println(gs)
}
//读取性别表
func getAge(){
	age := Dao("age")
    //实例化一个博文对象
    //实例化主控服务器对象
    var gs []map[string]int
    //有则更新  无则插入
    age.Find(nil).All(&gs)

    fmt.Println("读取")
    fmt.Println(gs)
}

//读取用户属性
func getNmas(){
	nmas := Dao("nmas")

	type attr struct{
		Uid int `uid`
		Attr string `attr`
	}
    //实例化一个博文对象
    //实例化主控服务器对象
    var gs []attr
    //有则更新  无则插入
    nmas.Find(nil).All(&gs)

    fmt.Println("读取")
    fmt.Println(gs)
}*/
