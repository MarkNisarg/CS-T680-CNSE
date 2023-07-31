package main

import (
	"flag"
	"fmt"

	"voter-api/api"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	hostFlag string
	portFlag uint
)

func processCmdLineFlags() {
	flag.StringVar(&hostFlag, "h", "0.0.0.0", "Listen on all interfaces")
	flag.UintVar(&portFlag, "p", 1080, "Default Port")

	flag.Parse()
}

func main() {
	processCmdLineFlags()
	r := gin.Default()
	r.Use(cors.Default())

	// Create a new instance of the VoterAPI handler.
	voterHandler := api.NewVoterHandler()

	// Register the HealthMiddleware, it will be called for every request.
	r.Use(api.HealthMiddleware(voterHandler))

	// Define the API endpoints and map them to the corresponding handler.
	r.GET("/", voterHandler.WelcomeToVoterAPI)
	r.GET("/voters", voterHandler.ListAllVoters)
	r.GET("/voters/:id", voterHandler.GetVoter)
	r.POST("/voters/:id", voterHandler.AddVoter)
	r.PUT("/voters/:id", voterHandler.UpdateVoter)
	r.DELETE("/voters", voterHandler.DeleteAllVoters)
	r.DELETE("/voters/:id", voterHandler.DeleteVoter)
	r.GET("/voters/:id/polls", voterHandler.GetVoterHistory)
	r.GET("/voters/:id/polls/:pollid", voterHandler.GetVoterPoll)
	r.POST("/voters/:id/polls/:pollid", voterHandler.AddVoterPoll)
	r.PUT("/voters/:id/polls/:pollid", voterHandler.UpdateVoterPoll)
	r.DELETE("/voters/:id/polls/:pollid", voterHandler.DeleteVoterPoll)
	r.GET("/voters/health", voterHandler.HealthCheck)

	// Start the server.
	serverPath := fmt.Sprintf("%s:%d", hostFlag, portFlag)
	r.Run(serverPath)
}
