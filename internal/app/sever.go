package app

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"

	// "testsocket/pkg/repository"

	"github.com/eldar34/goWebSocket/config"
	"github.com/eldar34/goWebSocket/internal/store/mysqlStore"
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
	config := config.NewConfig()
	dbConnection, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.MySQL["user"], config.MySQL["password"], config.MySQL["host"], config.MySQL["port"], config.MySQL["db"]))
	// defer store.Close()

	if err != nil {
		panic(err)
	}

	err = dbConnection.Ping()
	if err != nil {
		panic(err)
	}

	store := mysqlStore.New(dbConnection)

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Get Bearer token
		reqToken := req.Header.Get("Authorization")
		splitToken := strings.Split(reqToken, "Bearer ")
		reqToken = splitToken[1]

		row, err := store.User().FindByToken(reqToken)
		// defer dbConnection.Close()

		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			// panic(err)
			fmt.Println(err)
		}
		fmt.Println(row)
		// u := userInfo{}

		// if reqToken != "longkey" {
		// 	http.Error(w, "Unauthorized", http.StatusUnauthorized)
		// }

		next.ServeHTTP(w, req)

	})
}

func Run() {

	flag.Parse()
	hub := NewHub()
	go hub.Run()
	mux := http.NewServeMux()

	mux.HandleFunc("/", serveHome)
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ServeWs(hub, w, r)
	})
	err := http.ListenAndServeTLS(*addr, "./../../cert/ca-cert.pem", "./../../cert/ca-key.pem", Logging(mux))
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
