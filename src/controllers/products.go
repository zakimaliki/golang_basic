package controllers

import (
	"encoding/json"
	"fmt"
	"golang-be-batch1/src/middleware"
	"golang-be-batch1/src/models"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
)

func Data_products(w http.ResponseWriter, r *http.Request) {
	// Set header CORS untuk memungkinkan akses dari localhost:3000
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	// Menangani permintaan OPTIONS
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Menangani permintaan GET
	if r.Method == "GET" {
		pageOld := r.URL.Query().Get("page")
		limitOld := r.URL.Query().Get("limit")
		page, _ := strconv.Atoi(pageOld)
		if page == 0 {
			page = 1
		}
		limit, _ := strconv.Atoi(limitOld)
		if limit == 0 {
			limit = 5
		}
		offset := (page - 1) * limit
		sort := r.URL.Query().Get("sort")
		if sort == "" {
			sort = "ASC"
		}
		sortby := r.URL.Query().Get("sortBy")
		if sortby == "" {
			sortby = "name"
		}
		sort = sortby + " " + strings.ToLower(sort)
		respons := models.FindCond(sort, limit, offset)
		totalData := models.CountData()
		totalPage := math.Ceil(float64(totalData) / float64(limit))
		result := map[string]interface{}{
			"status":      "Berhasil",
			"data":        respons.Value,
			"currentPage": page,
			"limit":       limit,
			"totalData":   totalData,
			"totalPage":   totalPage,
		}
		// Mengirim respons JSON
		w.Header().Set("Content-Type", "application/json")
		res, err := json.Marshal(result)
		if err != nil {
			http.Error(w, "Gagal Konversi Json", http.StatusInternalServerError)
			return
		}
		w.Write(res)
		return
	}

	// Menangani permintaan POST
	if r.Method == "POST" {
		var product models.Product
		err := json.NewDecoder(r.Body).Decode(&product)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		item := models.Product{
			Name:  product.Name,
			Price: product.Price,
			Stock: product.Stock,
		}
		models.Post(&item)
		// Mengirim respons status Created
		w.WriteHeader(http.StatusCreated)
		msg := map[string]string{
			"Message": "Product Created",
		}
		// Mengirim respons JSON
		w.Header().Set("Content-Type", "application/json")
		res, err := json.Marshal(msg)
		if err != nil {
			http.Error(w, "Gagal Konversi Json", http.StatusInternalServerError)
			return
		}
		w.Write(res)
		return
	}

	// Menangani permintaan selain GET dan POST
	http.Error(w, "method tidak diizinkan", http.StatusMethodNotAllowed)
}

func Data_product(w http.ResponseWriter, r *http.Request) {
	middleware.GetCleanedInput(r)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS,PUT,DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// Check the request method
	if r.Method == "OPTIONS" {
		// Respond to preflight requests
		w.WriteHeader(http.StatusOK)
		return
	}
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

func Handle_upload(w http.ResponseWriter, r *http.Request) {
	const (
		AllowedExtensions = ".jpg,.jpeg,.pdf,.png"
		MaxFileSize       = 2 << 20 // 2 MB
	)

	// Memeriksa method request, harus POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Mendapatkan file dari form
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	ext := filepath.Ext(handler.Filename)
	ext = strings.ToLower(ext)
	allowedExts := strings.Split(AllowedExtensions, ",")
	validExtension := false
	for _, allowedExt := range allowedExts {
		if ext == allowedExt {
			validExtension = true
			break
		}
	}
	if !validExtension {
		http.Error(w, "Invalid file extension", http.StatusBadRequest)
		return
	}

	// Mengecek ukuran file
	fileSize := handler.Size
	if fileSize > MaxFileSize {
		http.Error(w, "File size exceeds the allowed limit", http.StatusBadRequest)
		return
	}

	// Menggunakan timestamp untuk membuat nama file unik
	timestamp := time.Now().Format("20060102_150405")

	// Menginisialisasi konfigurasi Cloudinary
	cloudinaryURL := os.Getenv("CLOUDINARY_URL") // Ambil URL Cloudinary dari variabel lingkungan
	if cloudinaryURL == "" {
		http.Error(w, "Cloudinary URL not found", http.StatusInternalServerError)
		return
	}
	cld, err := cloudinary.NewFromURL(cloudinaryURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Konfigurasi uploader Cloudinary
	uploadParams := uploader.UploadParams{
		PublicID:  fmt.Sprintf("%s_%s", timestamp, handler.Filename),
		Overwrite: true,
	}

	// Mengunggah file ke Cloudinary
	uploadResult, err := cld.Upload.Upload(r.Context(), file, uploadParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Menyampaikan respons berhasil
	msg := map[string]string{
		"Message":        "File uploaded successfully",
		"PublicID":       uploadResult.PublicID,
		"SecureURL":      uploadResult.SecureURL,
		"OriginalWidth":  fmt.Sprintf("%d", uploadResult.Width),
		"OriginalHeight": fmt.Sprintf("%d", uploadResult.Height),
	}
	res, err := json.Marshal(msg)
	if err != nil {
		http.Error(w, "Failed to convert JSON", http.StatusInternalServerError)
		return
	}
	w.Write(res)
}

func SearchProduct(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("search")
	res, err := json.Marshal(models.FindData(keyword).Value)
	if err != nil {
		http.Error(w, "Gagal Konversi Json", http.StatusInternalServerError)
		return
	}
	w.Write(res)
}
