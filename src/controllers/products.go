package controllers

import (
	"encoding/json"
	"golang-be-batch1/src/models"
	"net/http"
)

func Data_products(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		res, err := json.Marshal(models.SelectAll().Value)
		if err != nil {
			http.Error(w, "Gagal Konversi Json", http.StatusInternalServerError)
			return
		}
		w.Write(res)
		w.Header().Set("Content-Type", "application/json")
		return
	} else if r.Method == "POST" {
		var product models.Product
		err := json.NewDecoder(r.Body).Decode(&product)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// if product.Id <= 0 || product.Name == "" || product.Price <= 0 || product.Stock <= 0 {
		// 	w.WriteHeader(http.StatusBadRequest)
		// 	fmt.Fprintf(w, "Invalid product data")
		// 	return
		// }
		item := models.Product{
			Name:  product.Name,
			Price: product.Price,
			Stock: product.Stock,
		}
		models.Post(&item)
		w.WriteHeader(http.StatusCreated)
		msg := map[string]string{
			"Message": "Product Created",
		}
		res, err := json.Marshal(msg)
		if err != nil {
			http.Error(w, "Gagal Konversi Json", http.StatusInternalServerError)
			return
		}
		w.Write(res)
	} else {
		http.Error(w, "method tidak diizinkan", http.StatusMethodNotAllowed)
	}
}

func Data_product(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/product/"):]

	if r.Method == "GET" {
		res, err := json.Marshal(models.Select(id).Value)
		if err != nil {
			http.Error(w, "Gagal Konversi Json", http.StatusInternalServerError)
		}
		w.Write(res)
		w.Header().Set("Content-Type", "application/json")
		return
	} else if r.Method == "PUT" {
		var updateProduct models.Product
		err := json.NewDecoder(r.Body).Decode(&updateProduct)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// if updateProduct.Id <= 0 || updateProduct.Name == "" || updateProduct.Price <= 0 || updateProduct.Stock <= 0 {
		// 	w.WriteHeader(http.StatusBadRequest)
		// 	fmt.Fprintf(w, "Invalid product data")
		// 	return
		// }
		newProduct := models.Product{
			Name:  updateProduct.Name,
			Price: updateProduct.Price,
			Stock: updateProduct.Stock,
		}
		models.Updates(id, &newProduct)
		msg := map[string]string{
			"Message": "Product Updated",
		}
		res, err := json.Marshal(msg)
		if err != nil {
			http.Error(w, "Gagal Konversi Json", http.StatusInternalServerError)
			return
		}
		w.Write(res)
	} else if r.Method == "DELETE" {
		models.Deletes(id)
		msg := map[string]string{
			"Message": "Product Deleted",
		}
		res, err := json.Marshal(msg)
		if err != nil {
			http.Error(w, "Gagal Konversi Json", http.StatusInternalServerError)
			return
		}
		w.Write(res)
	} else {
		http.Error(w, "method tidak diizinkan", http.StatusMethodNotAllowed)
	}
}
