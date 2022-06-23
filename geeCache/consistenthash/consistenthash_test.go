package consistenthash

import (
	"strconv"
	"testing"
)

func Test_consis(t *testing.T) {
	hash := NewMap(3, func(data []byte) uint32 {
		i, _ := strconv.Atoi(string(data))
		return uint32(i)
	})

	// Given the above hash function, this will give replicas with "hashes":
	// 2, 4, 6, 12, 14, 16, 22, 24, 26
	hash.Add("6", "4", "2")

	testCase := map[string]string{
		"2":  "2",
		"11": "2",
		"23": "4",
		"27": "2",
	}
	for k, v := range testCase {
		res := hash.Get(k)
		if res != v {
			t.Errorf("Asking for %s, should have yielded %s,but get is %s", k, v, res)
		}
	}
	testCase["27"] = "8"
	hash.Add("8")

	for k, v := range testCase {
		res := hash.Get(k)
		if res != v {
			t.Errorf("Asking for %s, should have yielded %s,but get is %s", k, v, res)
		}
	}

}
