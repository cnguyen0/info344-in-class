package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"sync"
)

// User struct
type User struct {
	// make sure you do the camel casing in the json file lol
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
}

// GetCurrentUser gets users. does some magic with our session packages
func GetCurrentUser(r *http.Request) *User {
	return &User{
		FirstName: "Test",
		LastName:  "User",
	}
}

// NewServiceProxy handles all the addresses and ppl go to diff servers ya feel haha
func NewServiceProxy(addrs []string) *httputil.ReverseProxy {
	nextIndex := 0
	mx := sync.Mutex{}
	return &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			//modify the request to indicate
			//remote host
			user := GetCurrentUser(r)
			userJSON, err := json.Marshal(user)
			if err != nil {
				log.Printf("error marshaling user: %v", err)
			}
			r.Header.Add("X-User", string(userJSON))

			mx.Lock()
			r.URL.Host = addrs[nextIndex%len(addrs)]
			nextIndex++
			r.URL.Scheme = "http"
			mx.Unlock()
		},
	}
}

//RootHandler handles requests for the root resource
func RootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain")
	fmt.Fprintf(w, "Hello from the gateway! Try requesting /v1/time")
}

func main() {
	addr := os.Getenv("ADDR")
	if len(addr) == 0 {
		addr = ":443"
	}

	timesvcAddrs := os.Getenv("TIMESVC_ADDRS")
	splitTimeSvcAddrs := strings.Split(timesvcAddrs, ",")

	nodeSvcAddrs := os.Getenv("NODESVC_ADDRS")
	splitNodeSvcAddrs := strings.Split(nodeSvcAddrs, ",")

	mux := http.NewServeMux()
	mux.HandleFunc("/", RootHandler)
	mux.Handle("/v1/time", NewServiceProxy(splitTimeSvcAddrs))
	mux.Handle("/v1/users/me/hello", NewServiceProxy(splitNodeSvcAddrs))

	log.Printf("server is listening at https://%s...", addr)
	log.Fatal(http.ListenAndServeTLS(addr, "tls/fullchain.pem", "tls/privkey.pem", mux))
}
