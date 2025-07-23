package main

import (
	"subscription_service/cmd/app"
	_ "subscription_service/docs"
)

// @title           Subscription Service
// @version         0.0.1
// @description service for subscriptions users

// @host      localhost:8080
// @BasePath  /
func main() {
	app.Run()
}
