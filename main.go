package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

const baseUploadDir = "./bucket"

// 图片存储目录
const baseImageDir = baseUploadDir + "/i"

// 视频存储目录
const baseVideoDir = baseUploadDir + "/v"

// 音频存储目录
const baseAudioDir = baseUploadDir + "/a"

// 文件存储目录
const baseFileDir = baseUploadDir + "/f"

func main() {
	// 创建存储目录
	if _, err := os.Stat(baseImageDir); os.IsNotExist(err) {
		os.Mkdir(baseImageDir, os.ModePerm)
	}

	router := gin.Default()

	// 上传图片接口
	router.POST("/upload", uploadImage)
	// 获取图片预览接口
	router.GET("/:date/:filename", getImage)
	// 获取图片列表接口
	router.GET("/images", listImages)

	router.Run(":8801")
}

// 上传图片处理函数
func uploadImage(c *gin.Context) {
	file, fileHeader, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file is uploaded"})
		return
	}

	// 验证文件类型（确保是图像）
	if !isValidImage(fileHeader) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type"})
		return
	}

	// 生成当天的日期文件夹
	dateFolder := time.Now().Format("20060102")
	dateDir := filepath.Join(baseImageDir, dateFolder)

	// 如果日期文件夹不存在则创建
	if _, err := os.Stat(dateDir); os.IsNotExist(err) {
		if err := os.Mkdir(dateDir, os.ModePerm); err != nil {
			log.Println("Failed to create date directory:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create directory"})
			return
		}
	}

	// 生成短文件名
	fileName := generateShortFileName(filepath.Ext(fileHeader.Filename))
	uploadPath := filepath.Join(dateDir, fileName)

	// 创建文件
	out, err := os.Create(uploadPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save the file"})
		return
	}
	defer out.Close()
	// 保存文件
	_, err = file.Seek(0, 0) // 回到文件开头
	_, err = out.ReadFrom(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving file"})
		return
	}

	// 返回图片预览的 URL
	imageUrl := fmt.Sprintf("http://%s/%s/%s", c.Request.Host, dateFolder, fileName)
	c.JSON(http.StatusOK, gin.H{"url": imageUrl})
}

// 获取图片预览处理函数
func getImage(c *gin.Context) {
	date := c.Param("date")
	filename := c.Param("filename")
	filePath := filepath.Join(baseImageDir, date, filename)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
		return
	}

	c.File(filePath)
}

// 列出所有已上传图片的函数
func listImages(c *gin.Context) {
	var imageList []string

	err := filepath.Walk(baseImageDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 只添加文件，忽略目录
		if !info.IsDir() {
			relativePath, _ := filepath.Rel(baseImageDir, path)
			imageUrl := fmt.Sprintf("http://%s/%s", c.Request.Host, relativePath)
			imageList = append(imageList, imageUrl)
		}
		return nil
	})

	if err != nil {
		log.Println("Failed to list images:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list images"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"images": imageList})
}

func isValidImage(fileHeader *multipart.FileHeader) bool {
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
	}
	fileType := fileHeader.Header.Get("Content-Type")
	return allowedTypes[fileType]
}

func generateShortFileName(extension string) string {
	// 生成一个 4 字节的随机数（8 个字符）
	randomBytes := make([]byte, 4)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(randomBytes) + extension
}
