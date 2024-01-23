package helper

import (
	"golang-be-batch1/src/config"
	"golang-be-batch1/src/models"
)

func Migration() {
	config.DB.AutoMigrate(&models.Product{})
}
