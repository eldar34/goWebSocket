package app

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"testsocket/internal/entity"
	"testsocket/pkg/repository"

	// "testsocket/config"

	_ "github.com/go-sql-driver/mysql"
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
	http.ServeFile(w, r, "home.html")
}

type userInfo struct {
	id           int
	user_id      int
	access_token string
}

func Logging(next http.Handler) http.Handler {
	db, errors := repository.NewDb()
	if errors != nil {
		panic(errors)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Get Bearer token
		reqToken := req.Header.Get("Authorization")
		splitToken := strings.Split(reqToken, "Bearer ")
		reqToken = splitToken[1]

		row := db.QueryRow("SELECT * from for_go.oauth_access_tokens WHERE access_token = ?", reqToken)
		// defer db.Close()

		u := userInfo{}
		err := row.Scan(&u.id, &u.user_id, &u.access_token)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			// panic(err)
			fmt.Println(err)
		}

		// if reqToken != "longkey" {
		// 	http.Error(w, "Unauthorized", http.StatusUnauthorized)
		// }

		next.ServeHTTP(w, req)

	})
}

func Run() {

	flag.Parse()
	hub := entity.NewHub()
	go hub.Run()
	mux := http.NewServeMux()

	mux.HandleFunc("/", serveHome)
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		entity.ServeWs(hub, w, r)
	})
	err := http.ListenAndServeTLS(*addr, "./../../cert/ca-cert.pem", "./../../cert/ca-key.pem", Logging(mux))
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
