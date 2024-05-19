package main

import (
	"encoding/json"
	"fmt"
	"io"
	"module/config"
	"module/model/data"
	"module/model/service"
	"net/http"
	_ "net/http/pprof"
	"time"

	"go.uber.org/zap"

	"github.com/redis/go-redis/v9"
)

type Handler struct {
	CacheClient    *data.RedisClient
	DatabaseClient *data.PostgresClient
	Logger         *zap.Logger
}

func main() {
	// The below command performs service profiling
	// Start the profiling server on a separate goroutine
	go func() {
		http.ListenAndServe("localhost:6060", nil)
	}()

	zapConfig := zap.NewProductionConfig()
	zapConfig.OutputPaths = []string{
		"app.log",
		"stdout",
	}
	logger, err := zapConfig.Build()

	if err != nil {
		panic(err)
	}

	// Instantiate Data stores
	// Instantiate cache client
	// Cache
	cache := &data.RedisClient{
		RedisConfig: data.RedisConfig{
			Addr:     config.RedisConfigAddr,
			Password: config.RedisConfigPassword,
			DB:       config.RedisConfigDB,
		},
		Logger: logger,
	}

	cache.StartRedis()

	// Database
	db := &data.PostgresClient{
		PostgresConfig: data.PostgresConfig{
			Addr:     config.PostgresConfigAddr,
			User:     config.PostgresConfigUser,
			Password: config.PostgresConfigPassword,
			Database: config.PostgreConfigDB,
		},
		Logger: logger,
	}
	db.StartPostgres()
	db.CreateSchema()

	logger.Info("Starting the Handler")

	handler := &Handler{
		CacheClient:    cache,
		DatabaseClient: db,
		Logger:         logger,
	}

	// HTTP server
	logger.Info("Starting the HTTP server")
	http.Handle("/account", handler)
	err = http.ListenAndServe(fmt.Sprintf(":%d", config.HTTPServerPort), nil)
	if err != nil {
		logger.Fatal(err.Error())
	}

	// Sleep to keep the program running
	time.Sleep(10 * time.Minute)

}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	fmt.Println(r.Method)
	switch r.Method {
	case "GET":
		h.HandleGet(w, r)
	case "POST":
		h.HandlePost(w, r)
	default:
		h.Logger.Error("No suitable handler found for this method")
	}
}

func (h *Handler) HandleGet(w http.ResponseWriter, r *http.Request) {
	userQuery := r.URL.Query().Get("user")
	account := &service.Account{}
	res := h.CacheClient.Get(r.Context(), userQuery)
	if res.Err() != nil {
		if res.Err() == redis.Nil {
			h.Logger.Warn("Key not found in cache")
		} else {
			h.Logger.Error("Unexpected Redis exception")
		}
		h.Logger.Info("Reading from main database")
		h.DatabaseClient.GetUser(account, userQuery)
		h.Logger.Debug(account.Name)
		h.Logger.Debug(fmt.Sprintf("%f", account.Balance))

	} else {
		h.Logger.Info("Key found in cache")
		h.Logger.Debug(res.Val())
		err := json.Unmarshal([]byte(res.Val()), &account)
		if err != nil {
			fmt.Println("Error reading from cache")
		}
	}

	b, _ := json.Marshal(account)
	w.Write(b)
}

func (h *Handler) HandlePost(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("POST")
	b, _ := io.ReadAll(r.Body)
	account := &service.Account{}
	json.Unmarshal(b, account)
	h.Logger.Debug(fmt.Sprintf("Setting cache value for %s", account.Name))

	sendChan := make(chan bool)
	sendChan2 := make(chan bool)

	go h.CacheClient.PutInCache(r.Context(), account.Name, b, config.RedisTTL, sendChan)
	recChan := <-sendChan

	go h.DatabaseClient.PutInStore(account, sendChan2)
	recChan2 := <-sendChan2

	if recChan && recChan2 {
		h.Logger.Info("Written successfully to cache and store")
	}

}
