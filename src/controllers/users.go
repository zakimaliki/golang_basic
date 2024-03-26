package controllers

import (
	"encoding/json"
	"fmt"
	"golang-be-batch1/src/helper"
	"golang-be-batch1/src/models"
	"net/http"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS,PUT,DELETE")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

        // Check the request method
        if r.Method == "OPTIONS" {
            // Respond to preflight requests
            w.WriteHeader(http.StatusOK)
            return
        }
	if r.Method == "POST" {
		var input models.User
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "invalid request body")
			return
		}
		hashPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		Password := string(hashPassword)
		newUser := models.User{
			Email:    input.Email,
			Password: Password,
		}
		models.CreateUser(&newUser)
		msg := map[string]string{
			"Message": "Register successfully",
		}
		res, err := json.Marshal(msg)
		if err != nil {
			http.Error(w, "Gagal Konversi Json", http.StatusInternalServerError)
			return
		}
		w.Write(res)
	} else {
		http.Error(w, "", http.StatusBadRequest)
	}
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS,PUT,DELETE")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

        // Check the request method
        if r.Method == "OPTIONS" {
            // Respond to preflight requests
            w.WriteHeader(http.StatusOK)
            return
        }
	if r.Method == "POST" {
		var input models.User
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "invalid request body")
			return
		}
		ValidateEmail := models.FindEmail(&input)
		if len(ValidateEmail) == 0 {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(w, "Email is not Found")
			return
		}
		var PasswordSecond string
		for _, user := range ValidateEmail {
			PasswordSecond = user.Password
		}
		if err := bcrypt.CompareHashAndPassword([]byte(PasswordSecond), []byte(input.Password)); err != nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Password not Found")
			return
		}
		jwtKey := os.Getenv("SECRETKEY")
		token, _ := helper.GenerateToken(jwtKey, input.Email)
		item := map[string]string{
			"Email": input.Email,
			"Token": token,
		}
		res, _ := json.Marshal(item)
		w.Write(res)
	} else {
		http.Error(w, "", http.StatusBadRequest)
	}
}
