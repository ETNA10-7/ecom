package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/ETNA10-7/ecom/services/cart"
	"github.com/ETNA10-7/ecom/services/order"
	"github.com/ETNA10-7/ecom/services/product"
	"github.com/ETNA10-7/ecom/services/user"
	"github.com/gorilla/mux"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

// Creates new instances of APIServer and also storing in struct APIServer
func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

// Initializes all the router
func (s *APIServer) Run() error {

	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	//Can Create a different file for Routing
	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRouter(subrouter)

	productStore := product.NewStore(s.db)
	productHandler := product.NewHandler(productStore, userStore)
	productHandler.RegisterRouter(subrouter)

	orderStore := order.NewStore(s.db)

	cartHandler := cart.NewHandler(productStore, orderStore, userStore)
	cartHandler.RegisterRouter(subrouter)
	log.Println("Listening on ", s.addr)
	return http.ListenAndServe(s.addr, router)
}
