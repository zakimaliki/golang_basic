package helper

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

func Upload(w http.ResponseWriter, file multipart.File, handler *multipart.FileHeader) {
	// Membuat format waktu dengan detik
	timestamp := time.Now().Format("20060102_150405")

	// Membuat nama unik untuk file
	filename := fmt.Sprintf("src/uploads/%s_%s", timestamp, handler.Filename)

	// Membuat file untuk menyimpan gambar
	out, err := os.Create(filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer out.Close()

	// Menyalin isi file yang diupload ke file yang baru dibuat
	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
