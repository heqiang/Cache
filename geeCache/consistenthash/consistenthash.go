package consistenthash

import (
	"fmt"
	"hash/crc32"
	"sort"
	"strconv"
	"strings"
)

type Hash func(data []byte) uint32

type Map struct {
	hash     Hash
	replicas int
	keys     []int
	hashMap  map[int]string
}

func New(replicas int, fn Hash) *Map {
	a := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if a.hash == nil {
		a.hash = crc32.ChecksumIEEE
	}
	return a
}

// Add 换上添加真实节点
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			// 每个真实的节点 都创建m.replicas个虚拟节点
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			//环上添加对应都虚拟节点
			m.keys = append(m.keys, hash)
			//虚拟节点和真实节点都映射关系
			m.hashMap[hash] = key
		}
	}
	sort.Ints(m.keys)
	fmt.Println(m.keys)
}

func (m *Map) Get(key string) string {
	if strings.TrimSpace(key) == "" {
		return ""
	}

	hash := int(m.hash([]byte(key)))
	// 换上都都放置都是我们标记都hash值
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] > hash
	})

	return m.hashMap[m.keys[idx%len(m.keys)]]
}
