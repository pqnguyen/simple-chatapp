package main

import "github.com/pqnguyen/simple-chatapp/backend/server"

func main() {
	srv := server.Server{}
	srv.Start(":8080")
}
