package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"wishlist-api/controllers"
	"wishlist-api/middleware"
	"wishlist-api/database"

	httpSwagger "github.com/swaggo/http-swagger"
	_ "wishlist-api/docs"
)

// @title Wishlist API
// @version 1.0
// @description API для создания вишлистов с публичным доступом
// @host localhost:8000
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	database.InitDB()

	router := mux.NewRouter()

	public := router.PathPrefix("/").Subrouter()

	public.HandleFunc("/auth/register", controllers.Register).Methods("POST")
	public.HandleFunc("/auth/login", controllers.Login).Methods("POST")

	public.HandleFunc("/public/{token}", controllers.GetPublicWishlist).Methods("GET")
	public.HandleFunc("/public/{token}/book/{itemId}", controllers.BookItem).Methods("POST")

	public.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	protected := router.PathPrefix("/").Subrouter()
	protected.Use(middleware.JwtAuthentication)

	protected.HandleFunc("/wishlists", controllers.CreateWishlist).Methods("POST")
	protected.HandleFunc("/wishlists", controllers.GetWishlists).Methods("GET")
	protected.HandleFunc("/wishlists/{id}", controllers.GetWishlist).Methods("GET")
	protected.HandleFunc("/wishlists/{id}", controllers.UpdateWishlist).Methods("PUT")
	protected.HandleFunc("/wishlists/{id}", controllers.DeleteWishlist).Methods("DELETE")

	protected.HandleFunc("/wishlists/{wishlistId}/items", controllers.CreateItem).Methods("POST")
	protected.HandleFunc("/wishlists/{wishlistId}/items", controllers.GetItems).Methods("GET")
	protected.HandleFunc("/items/{id}", controllers.UpdateItem).Methods("PUT")
	protected.HandleFunc("/items/{id}", controllers.DeleteItem).Methods("DELETE")

	port := os.Getenv("port")
	if port == "" {
		port = "8000"
	}

	fmt.Printf("Server running on port %s\n", port)
	http.ListenAndServe(":"+port, router)
}