package main

import (
	"fmt"

	"github.com/bartleyg/useragentvaluestore"
)

func main() {
	fmt.Println("starting")
	api := useragentvaluestore.NewApi()

	err := api.Run(":80") // listen and serve on 0.0.0.0:80
	fmt.Println(err)
}
