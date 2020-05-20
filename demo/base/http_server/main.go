package main

import (
	"log"
	"net/http"
	"time"
)

var (
	Addr = ":1210"
)

func main() {
	//创建路由器
	mux := http.NewServeMux()

	//设置路由规则
	mux.HandleFunc("/bye", sayBye)
	//创建服务器
	server := &http.Server{
		Addr:         Addr,
		Handler:      mux,
		WriteTimeout: time.Second * 3,
	}

	//监听端口并提供任务
	log.Println("Starting httpserver at " + Addr)
	log.Fatal(server.ListenAndServe())
}

func sayBye(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Second)
	w.Write([]byte("bye bye,this is http server"))
}
