package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"example_project/internal/config"
	"example_project/internal/models"
	rep "example_project/internal/repository"
	"example_project/internal/service"
)

type server struct {
	service service.Service
}

// responds with a random order
func (s *server) getRandEncoded(w http.ResponseWriter, r *http.Request) {
	randomEmails := []string{
		"dsdsd@yandex.ru",
		"2323232@gmail.com",
		"090909090@yahoo.com",
	}
	ctx := r.Context()

	isInvalid := r.URL.Query().Has("invalid")

	ids, err := s.service.GetRandomItemIDs(ctx)
	if err != nil {
		http.Error(w, "unable to get random item ids", http.StatusInternalServerError)
		return
	}

	// if 'invalid' query-param had been passed then add some invalid id to the order
	if isInvalid {
		ids = append(ids, 99999)
	}

	order := models.Order{
		Email:   randomEmails[rand.Intn(len(randomEmails))],
		ItemIDs: ids,
	}

	json.NewEncoder(w).Encode(&order)
}

// processes order
func (s *server) processOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var order models.Order

	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, "unable to decode request body", http.StatusBadRequest)
		return
	}

	err := s.service.ProcessOrder(ctx, &order)
	if err != nil && err != service.ErrItemNotFound {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if err == service.ErrItemNotFound {
		http.Error(w, "unknown item in order item list", http.StatusBadRequest)
		return
	}

	resp := struct {
		OK bool `json:"ok"`
	}{
		OK: true,
	}
	json.NewEncoder(w).Encode(&resp)
}

func main() {
	// parsing the config
	var c config.Config
	err := envconfig.Process("", &c)
	if err != nil {
		log.Fatalf("unable to parse config: %v", err.Error())
	}

	fmt.Println(c.DatabaseDSN)

	// connecting to db
	db, err := sql.Open("postgres", c.DatabaseDSN)
	if err != nil {
		log.Fatalf("unable to connect to db: %v", err)
	}

	// trying to ping db until it's successful
	for {
		err := db.Ping()
		if err != nil {
			log.Printf("db ping returned: %v\n", err)
			time.Sleep(time.Second)
			continue
		}
		break
	}

	// applying migrations with golang-migrate
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations",
		"postgres", driver)
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("unable to apply migrations on db: %v", err)
	}

	log.Println("connected to DB")

	// initializing application layers
	itemRep, orderRep := rep.NewItemRepository(db), rep.NewOrderRepository(db)
	srv := server{
		service: service.NewService(orderRep, itemRep),
	}

	// initializing the servier
	r := mux.NewRouter()

	// adding recovery and logging middlewares to the server
	r.Use(
		handlers.RecoveryHandler(),
		func(next http.Handler) http.Handler {
			return handlers.LoggingHandler(os.Stdout, next)
		},
	)

	// registering routes
	r.HandleFunc("/api/orders/process", srv.processOrder).Methods(http.MethodPost)
	r.HandleFunc("/api/orders/rand_encoded", srv.getRandEncoded).Methods(http.MethodGet)

	// running pprof and metrics server in a different goroutine
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":8081", nil)
	}()

	// running the application server
	log.Printf("running server on %s\n", c.Addr)
	http.ListenAndServe(c.Addr, r)
}
