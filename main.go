package main

import (
	"fmt"
	"golang-be-batch1/src/config"
	"golang-be-batch1/src/helper"
	"golang-be-batch1/src/routes"
	"net/http"

	"github.com/subosito/gotenv"
)

func main() {
	// instal linter
	// go get github.com/golangci/golangci-lint/cmd/golangci-lint
	// run linter
	// go test -run=^$ -v ./...
	gotenv.Load()
	config.InitDB()
	helper.Migration()
	defer config.DB.Close()
	routes.Router()
	fmt.Print("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
