package main

import (
	"flag"
	"fmt"

	"poll-api/api"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	hostFlag string
	portFlag uint
)

func processCmdLineFlags() {
	flag.StringVar(&hostFlag, "h", "0.0.0.0", "Listen on all interfaces")
	flag.UintVar(&portFlag, "p", 1081, "Default Port")

	flag.Parse()
}

func main() {
	processCmdLineFlags()
	r := gin.Default()
	r.Use(cors.Default())

	// Create a new instance of the PollAPI handler.
	pollHandler := api.NewPollHandler()

	// Register the HealthMiddleware, it will be called for every request.
	r.Use(api.HealthMiddleware(pollHandler))

	// Define the API endpoints and map them to the corresponding handler.
	r.GET("/", pollHandler.WelcomeToPollAPI)
	r.GET("/polls", pollHandler.ListAllVPolls)
	r.GET("/polls/:id", pollHandler.GetPoll)
	r.POST("/polls/:id", pollHandler.AddPoll)
	r.DELETE("/polls", pollHandler.DeleteAllPolls)
	r.DELETE("/polls/:id", pollHandler.DeletePoll)
	r.GET("/polls/:id/options", pollHandler.GetPollOptions)
	r.GET("/polls/:id/options/:optionId", pollHandler.GetPollOption)
	r.POST("/polls/:id/options/:optionId", pollHandler.AddPollOption)
	r.DELETE("/polls/:id/options/:optionId", pollHandler.DeletePollOption)
	r.GET("/polls/health", pollHandler.HealthCheck)

	// Start the server.
	serverPath := fmt.Sprintf("%s:%d", hostFlag, portFlag)
	r.Run(serverPath)
}
