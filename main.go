package main

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/saracen/walker"
	"gorm.io/gorm"
)

type RawImage struct {
	gorm.Model
	Path string
	Type string
}

func getAllImages(path string) (chan string, chan error) {
	resultChan, errorChan := make(chan string), make(chan error)
	go func() {
		walkFn := func(pathname string, fi os.FileInfo) error {
			if !fi.IsDir() &&
				(strings.HasSuffix(pathname, ".jpeg") ||
					strings.HasSuffix(pathname, ".jpg") ||
					strings.HasSuffix(pathname, ".png")) {
				resultChan <- pathname
			}
			return nil
		}

		errCallback := walker.WithErrorCallback(func(pathname string, err error) error {
			if os.IsPermission(err) {
				return nil
			}
			errorChan <- err
			return nil
		})

		walker.Walk(path, walkFn, errCallback)
		close(resultChan)
		close(errorChan)
	}()
	return resultChan, errorChan
}

func main() {

	var router *gin.Engine = gin.Default()
	router.GET("/", func(context *gin.Context) {
		var resultChan, _ = getAllImages("/home/jarusll/Pictures")
	})
	router.Run(":8080")
}
