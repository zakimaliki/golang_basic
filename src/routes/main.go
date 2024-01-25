package routes

import (
	"fmt"
	"golang-be-batch1/src/controllers"
	"golang-be-batch1/src/middleware"
	"net/http"
)

func Router() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello World")
	})
	http.HandleFunc("/products", middleware.XssMiddleware(controllers.Data_products))
	http.HandleFunc("/product/", middleware.XssMiddleware(controllers.Data_product))
}
