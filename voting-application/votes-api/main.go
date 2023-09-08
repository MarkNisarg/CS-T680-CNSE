package main

import (
	"flag"
	"fmt"

	"votes-api/api"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	hostFlag    string
	portFlag    uint
	voterAPIURL string
	pollAPIURL  string
)

func processCmdLineFlags() {
	flag.StringVar(&hostFlag, "h", "0.0.0.0", "Listen on all interfaces")
	flag.StringVar(&voterAPIURL, "v", "http://host.docker.internal:1080", "Default voter API location")
	flag.StringVar(&pollAPIURL, "papi", "http://host.docker.internal:1081", "Default poll API location")
	flag.UintVar(&portFlag, "p", 1082, "Default Port")

	flag.Parse()
}

func main() {
	processCmdLineFlags()
	r := gin.Default()
	r.Use(cors.Default())

	// Create a new instance of the VoterAPI handler.
	voterHandler := api.NewVotesHandler(pollAPIURL, voterAPIURL)

	// Register the HealthMiddleware, it will be called for every request.
	r.Use(api.HealthMiddleware(voterHandler))

	// Define the API endpoints and map them to the corresponding handler.
	r.GET("/", voterHandler.WelcomeToVotesAPI)
	r.GET("/votes", voterHandler.ListAllVotes)
	r.GET("/votes/:id", voterHandler.GetVote)
	r.POST("/votes/:id", voterHandler.AddVote)
	r.DELETE("/votes/:id", voterHandler.DeleteVote)
	r.GET("/votes/health", voterHandler.HealthCheck)

	// Start the server.
	serverPath := fmt.Sprintf("%s:%d", hostFlag, portFlag)
	r.Run(serverPath)
}
