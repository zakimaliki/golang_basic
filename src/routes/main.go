package routes

import (
	"fmt"
	"golang-be-batch1/src/controllers"
	"golang-be-batch1/src/middleware"
	"net/http"

	"github.com/goddtriffin/helmet"
)

func Router() {
	helmet := helmet.Default()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello World")
	})
	http.Handle("/products", helmet.Secure(middleware.XssMiddleware(http.HandlerFunc(controllers.Data_products))))
	http.Handle("/product/", helmet.Secure(middleware.XssMiddleware(http.HandlerFunc(controllers.Data_product))))
	http.Handle("/upload", http.HandlerFunc(controllers.Handle_upload))

}
