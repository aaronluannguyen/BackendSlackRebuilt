package main

import (
	"os"
	"net/http"
	"log"
	"github.com/challenges-aaronluannguyen/servers/gateway/handlers"
	"database/sql"
	"github.com/challenges-aaronluannguyen/servers/gateway/sessions"
	"github.com/challenges-aaronluannguyen/servers/gateway/models/users"
	"github.com/go-redis/redis"
	"time"
	_ "github.com/go-sql-driver/mysql"
	"github.com/challenges-aaronluannguyen/servers/gateway/indexes"
	"github.com/gorilla/mux"
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
	dsn := reqEnv("DSN")

	summaryServiceAddrs := reqEnv("SUMMARYADDR")
	messagesServiceAddrs := reqEnv("MESSAGESADDR")

	mqAddr := reqEnv("MQADDR")
	mqName := reqEnv("MQNAME")

	addr := os.Getenv("ADDR")
	if len(addr) == 0 {
		addr = ":443"
	}

	db, err := sql.Open("mysql", dsn)
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
	trie, err := usersStore.LoadExistingUsersToTrie()
	if err != nil {
		trie = indexes.NewTrie()
	}
	notifier := handlers.NewNotifier()

	hctx := handlers.Context {
		SigningKey: sessionKey,
		SessionStore: sessionsStore,
		UsersStore: usersStore,
		Trie: trie,
		Notifier: notifier,
	}

	hctx.StartMQ(mqAddr, mqName)

	tlsKeyPath := os.Getenv("TLSKEY")
	tlsCertPath := os.Getenv("TLSCERT")
	if len(tlsKeyPath) == 0 || len(tlsCertPath) == 0 {
		log.Fatal("please set TLSKEY and TLSCERT")
	}

	r := mux.NewRouter()
	r.Handle("/v1/summary", handlers.NewServiceProxy(summaryServiceAddrs, hctx))
	r.HandleFunc("/v1/users", hctx.UsersHandler)
	r.HandleFunc("/v1/users/{userID}", hctx.SpecificUserHandler)
	r.HandleFunc("/v1/sessions", hctx.SessionsHandler)
	r.HandleFunc("/v1/sessions/{currSession}", hctx.SpecificSessionHandler)
	r.Handle("/v1/channels", handlers.NewServiceProxy(messagesServiceAddrs, hctx))
	r.Handle("/v1/channels/{channelID}", handlers.NewServiceProxy(messagesServiceAddrs, hctx))
	r.Handle("/v1/channels/{channelID}/members", handlers.NewServiceProxy(messagesServiceAddrs, hctx))
	r.Handle("/v1/messages/{messageID}", handlers.NewServiceProxy(messagesServiceAddrs, hctx))
	r.Handle("/v1/messages/{messageID}/reactions", handlers.NewServiceProxy(messagesServiceAddrs, hctx))
	r.Handle("/v1/me/starred/messages", handlers.NewServiceProxy(messagesServiceAddrs, hctx))
	r.Handle("/v1/me/starred/messages/{messageID}", handlers.NewServiceProxy(messagesServiceAddrs, hctx))
	r.Handle("/v1/ws", handlers.NewWebSocketsHandler(hctx))

	corsWrappedMux := handlers.WrappedCORSHandler(r)

	log.Printf("Server is listening at https://%s", addr)
	log.Fatal(http.ListenAndServeTLS(addr, tlsCertPath, tlsKeyPath, corsWrappedMux))
}
