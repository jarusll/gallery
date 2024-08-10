package main

import (
	"os"
	"strings"

	"github.com/saracen/walker"
	"gorm.io/gorm"
)

type RawImage struct {
	gorm.Model
	Path string
	Type string
}

func getAllImages(path string) (chan string, chan error, chan bool) {
	resultChan, errorChan, doneChan := make(chan string), make(chan error), make(chan bool)
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
		doneChan <- true
		close(doneChan)
		close(resultChan)
		close(errorChan)
	}()
	return resultChan, errorChan, doneChan
}

func main() {
	// db, err := gorm.Open(sqlite.Open("dev.db"), &gorm.Config{})
	// if err != nil {
	// 	println("Error initializing DB")
	// }
	picturePaths, picturesErr, picDone := getAllImages("/home/jarusll/Pictures")
	for {
		select {
		case path := <-picturePaths:
			println(path)
		case pictureError := <-picturesErr:
			println("Error ", pictureError)
		case <-picDone:
			return
		}
	}
	// var router *gin.Engine = gin.Default()
	// router.GET("/", func(context *gin.Context) {
	// })
	// router.Run(":8080")
}
