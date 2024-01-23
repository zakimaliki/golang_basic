package routes

import (
	"fmt"
	"golang-be-batch1/src/controllers"
	"net/http"
)

func Router() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello World")
	})
	http.HandleFunc("/products", controllers.Data_products)
	http.HandleFunc("/product/", controllers.Data_product)
}
