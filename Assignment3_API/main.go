// main.go
package main

import (
	"fmt"
	"workspace/server"
)

func main() {
	fmt.Println("Server started at http://localhost:8080")
	server.StartServer()
}
