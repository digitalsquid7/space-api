package main

import (
	"space-api/pkg/routes"
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.

// @contact.name   digitalsquid7

// @host      localhost:8080
// @BasePath  /main/v1

// @securityDefinitions.basic  BasicAuth
func main() {
	engine := routes.CreateEngine()
	if err := engine.Run(); err != nil {
		panic(err)
	}
}
