package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"github.com/xyproto/simpleredis/v2"
)

var (
	masterPool *simpleredis.ConnectionPool
)

func ListRangeHandler(rw http.ResponseWriter, req *http.Request) {
	key := mux.Vars(req)["key"]
	list := simpleredis.NewList(masterPool, key) // Use masterPool only
	members := HandleError(list.GetAll()).([]string)
	membersJSON := HandleError(json.MarshalIndent(members, "", "  ")).([]byte)
	rw.Write(membersJSON)
}

func ListPushHandler(rw http.ResponseWriter, req *http.Request) {
	key := mux.Vars(req)["key"]
	value := mux.Vars(req)["value"]
	list := simpleredis.NewList(masterPool, key) // Use masterPool only
	HandleError(nil, list.Add(value))
	ListRangeHandler(rw, req)
}

func InfoHandler(rw http.ResponseWriter, req *http.Request) {
	info := HandleError(masterPool.Get(0).Do("INFO")).([]byte)
	rw.Write(info)
}

func EnvHandler(rw http.ResponseWriter, req *http.Request) {
	environment := make(map[string]string)
	for _, item := range os.Environ() {
		splits := strings.Split(item, "=")
		key := splits[0]
		val := strings.Join(splits[1:], "=")
		environment[key] = val
	}

	envJSON := HandleError(json.MarshalIndent(environment, "", "  ")).([]byte)
	rw.Write(envJSON)
}

func HandleError(result interface{}, err error) (r interface{}) {
	if err != nil {
		panic(err)
	}
	return result
}

func main() {
	// Use only the master Redis instance
	masterPool = simpleredis.NewConnectionPoolHost("redis-master:6379")
	//masterPool = simpleredis.NewConnectionPoolHost("localhost:6379")
	defer masterPool.Close() // Only close masterPool

	r := mux.NewRouter()
	r.Path("/lrange/{key}").Methods("GET").HandlerFunc(ListRangeHandler)
	r.Path("/rpush/{key}/{value}").Methods("GET").HandlerFunc(ListPushHandler)
	r.Path("/info").Methods("GET").HandlerFunc(InfoHandler)
	r.Path("/env").Methods("GET").HandlerFunc(EnvHandler)

	n := negroni.Classic()
	n.UseHandler(r)
	n.Run("0.0.0.0:3000")
	//n.Run("localhost:3000")
}
