package main

import (
	"Goproject/geeCache"
	"fmt"
	"log"
	"net/http"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

// 此测试仅仅是通过http服务访问的方式测试缓存中存在时会直接冲缓存中拿取
// 缓存中不存咋时会出数据源中通过callback函数拿取并且load到缓存中
func main() {
	geeCache.NewGroup("scores", 2<<10, geeCache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
	addr := "localhost:9999"
	peers := geeCache.NewHttpPool(addr)
	log.Println("geecache is running at", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}

func createGroup() *geeCache.Group {
	return geeCache.NewGroup("scores", 2<<10, geeCache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
}
