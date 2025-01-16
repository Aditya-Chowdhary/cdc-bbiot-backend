package main

import (
	"fmt"

	"github.com/GDGVIT/bbiot-backend/api"
)

func main() {
	server := api.NewServer()

	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}

}
