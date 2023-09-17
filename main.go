package main

import (
	"bytes"
	"fmt"
	"log"

	"github.com/ra-khalish/deploy-webhook/services"

	"github.com/ra-khalish/deploy-webhook/routes"
)

func main() {
	result := services.All("kubisa-id")
	//for i, c := range result {
	//	fmt.Printf("%d: %c\n", i, c)
	//}

	split := bytes.Split(result, []byte("\n"))

	for i := 0; i < len(split)-1; i++ {
		fmt.Printf("%d: %s\n", i, split[i])
	}

	log.Println("start server on port 8090")
	routes.Init()

}
