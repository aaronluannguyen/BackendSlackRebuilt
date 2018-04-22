package main

import (
	"os"
	"net/http"
	"log"
	"github.com/challenges-aaronluannguyen/servers/gateway/handlers"
	"github.com/go-sql-driver/mysql"
	"database/sql"
	"github.com/challenges-aaronluannguyen/servers/gateway/sessions"
	"github.com/challenges-aaronluannguyen/servers/gateway/models/users"
	"github.com/go-redis/redis"
	"time"
)

func reqEnv(name string) string {
	val := os.Getenv(name)
	if len(val) == 0 {
		log.Fatalf("please set %s variable", name)
	}
	return val
}

//main is the main entry point for the server
func main() {
	/* TODO: add code to do the following
	- Read the ADDR environment variable to get the address
	  the server should listen on. If empty, default to ":80"
	- Create a new mux for the web server.
	- Tell the mux to call your handlers.SummaryHandler function
	  when the "/v1/summary" URL path is requested.
	- Start a web server listening on the address you read from
	  the environment variable, using the mux you created as
	  the root handler. Use log.Fatal() to report any errors
	  that occur when trying to start the web server.
	*/

	sessionKey := reqEnv("SESSIONKEY")
	redisADDR := reqEnv("REDISADDR")
	mysqlAddr := reqEnv("MYSQL_ADDR")
	mysqlDB := reqEnv("MYSQL_DATABASE")
	mysqlPwd := reqEnv("MYSQL_ROOT_PASSWORD")

	addr := os.Getenv("ADDR")
	if len(addr) == 0 {
		addr = ":443"
	}

	DSN := mysql.Config{
		Addr: mysqlAddr,
		User: "root",
		Passwd: mysqlPwd,
		Net: "tcp",
		DBName: mysqlDB,
	}
	db, err := sql.Open("mysql", DSN.FormatDSN())
	if err != nil {
		log.Fatalf("error opening database: %v", err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: redisADDR,
		Password: "",
		DB: 0,
	})

	sessionsStore := sessions.NewRedisStore(redisClient, time.Hour)
	usersStore := users.NewMySQLStore(db)
	hctx := handlers.Context {
		sessionKey,
		sessionsStore,
		usersStore,
	}

	tlsKeyPath := os.Getenv("TLSKEY")
	tlsCertPath := os.Getenv("TLSCERT")
	if len(tlsKeyPath) == 0 || len(tlsCertPath) == 0 {
		log.Fatal("please set TLSKEY and TLSCERT")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/summary", handlers.SummaryHandler)
	mux.HandleFunc("/v1/users", hctx.UsersHandler)
	mux.HandleFunc("/v1/users/", hctx.SpecificUserHandler)
	mux.HandleFunc("/v1/sessions", hctx.SessionsHandler)
	mux.HandleFunc("/v1/sessions/", hctx.SpecificSessionHandler)

	corsWrappedMux := handlers.WrappedCORSHandler(mux)

	log.Printf("Server is listening at https://%s", addr)
	log.Fatal(http.ListenAndServeTLS(addr, tlsCertPath, tlsKeyPath, corsWrappedMux))
}
