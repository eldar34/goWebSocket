package app

import (
	"flag"
	"log"
	"net/http"

	"testsocket/internal/entity"
)

var addr = flag.String("addr", ":8080", "http service address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "testsocket/internal/usecase/html/home.html")
}

func Run() {
	flag.Parse()
	hub := entity.NewHub()
	go hub.Run()
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		entity.ServeWs(hub, w, r)
	})
	err := http.ListenAndServeTLS(*addr, "./../../cert/ca-cert.pem", "./../../cert/ca-key.pem", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
