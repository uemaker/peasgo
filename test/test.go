package main

import (
	"fmt"
	"peasgo/cache"
	_ "peasgo/cache/memcache"
	"time"
)

func main() {
	//	s := []byte{1, 2}
	//	//	s := "111"
	//	var v interface{} = s

	//	switch v.(type) {
	//	case []byte:
	//		fmt.Println(1)
	//	case string:
	//		fmt.Println(2)
	//	default:
	//		fmt.Println(3)

	//	}
	//	memcache.Register()
	bm, err := cache.NewCache("memcache", `{"servers":"127.0.0.1:11211"}`)

	//	bm, err := cache.NewCache("memory", `{"services":"127.0.0.1:11211"}`)
	if err != nil {
		fmt.Println("init err")
	}
	timeoutDuration := 10000 * time.Second
	if err = bm.Set("name", `{"servers":"127.0.0.1:11211"}`, timeoutDuration); err != nil {
		fmt.Println(err)
	}
	if err = bm.Set("name2", `{"servers":"127.0.0.1:11211"}`, timeoutDuration); err != nil {
		fmt.Println(err)
	}
	v2 := bm.Mget([]string{"name", "name2"})

	fmt.Println(v2)

	//	time.Sleep(30 * time.Second)

	//	if bm.IsExist("astaxie") {
	//		t.Error("check err")
	//	}

	//	a, b := v.([]byte)
	//	fmt.Println(a)
	//	fmt.Println(b)

	//	v, ok := s.([]byte)

	//	fmt.Println(reflect.TypeOf(t).Kind())
}
