package model

import (
	"fmt"
	"math/rand"
	"time"
)

func TestAdds(num int) {
	for i := 1; i <= num; i++ {
		time.Sleep(1 * time.Second)
		TestAdd(i)
	}
}

func TestAdd(num int) {
	fmt.Println("开始填充用户", num)
	t1 := time.Now().UnixNano()
	//获取随机数
	ra := rand.New(rand.NewSource(time.Now().UnixNano()))
	attrs := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "A", "B", "C", "D"}

	j := (num - 1) * 1000000
	p := num*1000000 - 1
	//2000万用户
	for i := j; i <= p; i++ {
		Set("gender", i, ra.Intn(2))
		Set("age", i, ra.Intn(100))
		for j := 0; j < 16; j++ {
			Set(attrs[ra.Intn(30)], i, 1)
		}
	}
	fmt.Println(num, "完成：", (time.Now().UnixNano()-t1)/1000000000, "秒")
}
