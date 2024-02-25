package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type CryptoFileController struct{}

func NewCryptoFileController() *CryptoFileController {
	return &CryptoFileController{}
}

func readFileContent(file *multipart.FileHeader) ([]byte, error) {
	openedFile, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer openedFile.Close()

	return io.ReadAll(openedFile)
}

func encryptContent(content []byte, passphrase string) ([]byte, error) {
	block, err := aes.NewCipher([]byte(passphrase))
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, content, nil), nil
}

func decryptContent(path string, endFilePath string, passphrase string) error {
	cipherText, err := os.ReadFile(path)
	if err != nil {
		log.Print(err)
		return err
	}
	block, err := aes.NewCipher([]byte(passphrase))
	if err != nil {
		return err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}
	nonce := cipherText[:gcm.NonceSize()]
	cipherText = cipherText[gcm.NonceSize():]
	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		log.Print("decrypt file err: %v", err.Error())
		return err
	}
	err = os.WriteFile(endFilePath, plainText, 0777)
	if err != nil {
		log.Print("write file err: %v", err.Error())
		return err
	}
	return nil
}

func hashPassPhrase(passPhrase string) (string, error) {
	if passPhrase == "" {
		return "", errors.New("passPhrase is required")
	}
	md5Hash := md5.Sum([]byte(passPhrase))
	return hex.EncodeToString(md5Hash[:]), nil
}

func (ctrl *CryptoFileController) EncryptFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	content, err := readFileContent(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to read uploaded file"})
		return
	}

	passPhrase, err := hashPassPhrase(c.PostForm("passPhrase"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "passPhrase is required"})
		return
	}

	cipherText, err := encryptContent(content, passPhrase)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt file"})
		return
	}

	path := "./uploads/" + file.Filename + "_" + passPhrase
	if err := os.WriteFile(path, cipherText, 0777); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save encrypted file"})
		return
	}

	c.File(path)
}

func (ctrl *CryptoFileController) DecryptFile(c *gin.Context) {
	var req DecryptRequest

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	passPhrase, err := hashPassPhrase(req.PassPhrase)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "passPhrase is required"})
		return
	}
	path := "./uploads/" + req.FileName + "_" + passPhrase
	endFilePath := "./uncrypted/" + req.FileName

	if err := decryptContent(path, endFilePath, passPhrase); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decrypt file"})
		return
	}
	c.File(endFilePath)
	os.Remove(endFilePath)
}

func (ctrl *CryptoFileController) DownloadFile(c *gin.Context) {
	name := c.Param("name")
	path := "./uploads/" + name + "_crypto"
	c.Header("Content-Disposition", "attachment; filename=\""+name+"_crypto\"")
	c.File(path)
}
