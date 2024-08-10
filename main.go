package main

import (
	"fmt"
	"os"

	"github.com/saracen/walker"
)

func main() {
	resultChan, errorChan := make(chan string), make(chan error)
	const path = `/home/jarusll`
	go func() {
		walkFn := func(pathname string, fi os.FileInfo) error {
			if !fi.IsDir() {
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
	for {
		select {
		case result := <-resultChan:
			_ = result
		case err := <-errorChan:
			fmt.Println(err)
			return
		}
	}
}
