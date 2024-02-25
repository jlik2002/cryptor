package main

import (
	"crypton/crypto"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	fileUploadController := crypto.NewCryptoFileController()

	router.POST("/upload", fileUploadController.EncryptFile)
	router.POST("/decrypt", fileUploadController.DecryptFile)
	router.GET("/download/:name", fileUploadController.DownloadFile)

	router.Run(":8070")
}
