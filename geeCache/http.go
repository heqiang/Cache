package geeCache

import (
	"Goproject/geeCache/consistenthash"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const (
	defaultBasePath = "/_geecache/"
	defaultReplicas = 50
)

type HttpPool struct {
	self        string
	basePath    string
	mu          sync.Mutex
	peers       *consistenthash.Map
	httpGetters map[string]*httpGetter
}
type httpGetter struct {
	baseURL string
}

func NewHttpPool(self string) *HttpPool {
	return &HttpPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

func (p *HttpPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s\n", p.self, fmt.Sprintf(format, v))
}

func (p *HttpPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.URL.Path, "favicon") {
		return
	}
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTPPool serving unexpected path: " + r.URL.Path)
	}

	p.Log(r.Method, r.URL.Path)
	// /<basepath>/<groupname>/<key> required
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	fmt.Println(parts)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	groupName := parts[0]
	key := parts[1]

	group := GetGroup(groupName)

	if group == nil {
		http.Error(w, "no such group: "+groupName, http.StatusNotFound)
		return
	}

	view, err := group.Get(key)
	if err != nil {
		http.Error(w, "no such group: "+groupName, http.StatusNotFound)
		return
	}
	//w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(view.ByteSlice())
}

func (p *HttpPool) Set(peers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.peers = consistenthash.NewMap(defaultReplicas, nil)
	p.peers.Add(peers...)
	p.httpGetters = make(map[string]*httpGetter)
	for _, peer := range peers {
		p.httpGetters[peer] = &httpGetter{
			baseURL: peer + p.basePath,
		}
	}
}
func (p *HttpPool) PickPeer(key string) (PeerGetter, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if peer := p.peers.Get(key); peer != "" && peer != p.self {
		p.Log("Pick peer %s", peer)
		return p.httpGetters[peer], true
	}
	return nil, false
}

func (h *httpGetter) Get(group, key string) ([]byte, error) {
	u := fmt.Sprintf("%v%v/%v", h.baseURL, url.QueryEscape(group), url.QueryEscape(key))

	res, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned:%v", res.StatusCode)
	}

	byters, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body:%v", err)
	}

	return byters, nil
}

var _ PeerGetter = (*httpGetter)(nil)
