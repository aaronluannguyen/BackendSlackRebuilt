package handlers

import (
	"net/http/httputil"
	"strings"
	"net/http"
	"sync"
	"github.com/challenges-aaronluannguyen/servers/gateway/sessions"
	"encoding/json"
)

const headerUser = "X-User"

//NewServiceProxy returns a new ReverseProxy
//for a microservice given a comma-delimited
//list of network addresses
func NewServiceProxy(addrs string, ctx Context) *httputil.ReverseProxy {
	splitAddrs := strings.Split(addrs, ",")
	nextAddr := 0
	mx := sync.Mutex{}

	return &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			r.URL.Scheme = "http"
			mx.Lock()
			r.URL.Host = splitAddrs[nextAddr]
			nextAddr = (nextAddr + 1) % len(splitAddrs)
			mx.Unlock()

			r.Header.Del(headerUser)
			currentState := &SessionState{}
			user, err := sessions.GetState(r, ctx.SigningKey, ctx.SessionStore, currentState)
			if err != nil {
				return
			}
			userJSON, _ := json.Marshal(user)
			r.Header.Set(headerUser, string(userJSON))
		},
	}
}