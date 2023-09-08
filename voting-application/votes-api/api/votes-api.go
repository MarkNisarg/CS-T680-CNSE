package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	schema "votes-api/Schema"
	"votes-api/votes"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

// The API handler that handles incoming requests.
type VotesAPI struct {
	votesList        *votes.VotesCache
	pollAPIURL       string
	voterAPIURL      string
	apiClient        *resty.Client
	totalCalls       uint64
	errorCalls       uint64
	bootTime         time.Time
	totalRequestTime time.Duration
}

// Create a new instance of VotesAPI with an initialized votes cache.
func NewVotesHandler(pollAPIURL string, voterAPIURL string) *VotesAPI {
	votesCache, _ := votes.NewVotesCache()
	apiClient := resty.New()

	return &VotesAPI{
		votesList:        votesCache,
		pollAPIURL:       pollAPIURL,
		voterAPIURL:      voterAPIURL,
		apiClient:        apiClient,
		totalCalls:       0,
		errorCalls:       0,
		bootTime:         time.Now(),
		totalRequestTime: 0,
	}
}

// The custom middleware to handle health metadata.
func HealthMiddleware(va *VotesAPI) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Update total API calls count.
		va.totalCalls++

		// Record the start time of the request.
		start := time.Now()

		// Process the request.
		c.Next()

		// Update error API calls count if there's an error.
		if c.Writer.Status() >= 400 {
			va.errorCalls++
		}

		// Calculate the request duration.
		duration := time.Since(start)

		// Update the total request time.
		va.totalRequestTime += duration
	}
}

// The root endpoint that welcomes users to the API.
func (va *VotesAPI) WelcomeToVotesAPI(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome to votes API.",
	})
}

// Implementation of GET /votes.
// Returns all Votes with all votes.
func (va *VotesAPI) ListAllVotes(c *gin.Context) {
	Votes, err := va.votesList.GetAllVotes()
	if err != nil {
		log.Println("Error getting Votes: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if Votes == nil {
		Votes = make([]votes.Vote, 0)
	}

	c.JSON(http.StatusOK, Votes)
}

// Implementation of GET /votes/:id.
// Returns a single vote by :id.
func (va *VotesAPI) GetVote(c *gin.Context) {
	voteID := c.Param("id")
	voteIDUint, err := strconv.ParseUint(voteID, 10, 32)
	if err != nil {
		log.Println("Error converting vote ID to uint: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	vote, err := va.votesList.GetVote(uint(voteIDUint))
	if err != nil {
		log.Println("Error getting vote: ", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, vote)
}

// Implementation of POST /votes/:id.
// Add a new voter with :id.
func (va *VotesAPI) AddVote(c *gin.Context) {
	var vote votes.Vote
	if err := c.ShouldBindJSON(&vote); err != nil {
		log.Println("Error binding JSON: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	vID := vote.VoterID

	var voters = []schema.Voter{}
	votersPath := va.voterAPIURL + "/voters"

	_, err := va.apiClient.R().SetResult(&voters).Get(votersPath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Could not find voter in cache"})
		fmt.Println("Error getting voters:", err)
		return
	}

	var foundVoterID bool = false
	for _, voter := range voters {
		if voter.VoterID == vID {
			foundVoterID = true
			break
		}
	}

	if !foundVoterID {
		c.JSON(http.StatusNotFound, gin.H{"error": "Could not find voter in cache"})
		fmt.Println("Error getting voter")
		return
	}

	pID := vote.PollID
	optID := vote.VoteValue

	var polls = []schema.Poll{}
	pollsPath := va.pollAPIURL + "/polls"

	_, err = va.apiClient.R().SetResult(&polls).Get(pollsPath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Could not find poll in cache"})
		fmt.Println("Error getting poll")
		return
	}

	// Check if poll with ID and poll option with ID exist
	var foundPollID bool = false
	var foundPollOptID bool = false
	for _, poll := range polls {
		if poll.PollID == pID {
			foundPollID = true
			for _, option := range poll.PollOptions {
				if option.PollOptionID == optID {
					foundPollOptID = true
				}
			}
			break
		}
	}

	if !foundPollID {
		c.JSON(http.StatusNotFound, gin.H{"error": "Could not find poll in cache"})
		fmt.Println("Error getting poll: " + strconv.FormatUint(uint64(pID), 32))
		return
	}

	if !foundPollOptID {
		c.JSON(http.StatusNotFound, gin.H{"error": "Could not find poll option in cache"})
		fmt.Println("Error getting poll option: " + strconv.FormatUint(uint64(optID), 32))
		return
	}

	voteID := c.Param("id")
	voteIDUint, err := strconv.ParseUint(voteID, 10, 32)
	if err != nil {
		log.Println("Error converting vote ID to uint: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	vote.VoteID = uint(voteIDUint)

	if err := va.votesList.AddVote(vote); err != nil {
		fmt.Println("Error adding vote")
		log.Println("error adding item: ", err)
		c.AbortWithStatus(http.StatusConflict)
		return
	}

	// After successfully adding the vote, add it to the voter's vote history.
	voterID := vote.VoterID

	// Construct the URL for adding the vote to the voter's vote history.
	voterPollURL := fmt.Sprintf("%s/voters/%s/polls/%s", va.voterAPIURL, strconv.Itoa(int(voterID)), strconv.Itoa(int(vote.PollID)))

	// Create a JSON payload for adding the vote to the voter's vote history.
	// You can modify this payload according to your API's requirements.
	voterPollPayload := map[string]interface{}{
		"voteDate": time.Now(),
	}

	// Make an HTTP POST request to add the vote to the voter's vote history using apiClient.
	resp, err := va.apiClient.R().
		SetHeader("Content-Type", "application/json").
		SetBody(voterPollPayload).
		Post(voterPollURL)
	if err != nil {
		log.Println("Error performing POST request to add vote to voter's vote history: ", err)
		log.Println(resp.StatusCode())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, vote)
}

// Implementation of DELETE /Votes/:id.
// Delete a single vote by :id.
func (va *VotesAPI) DeleteVote(c *gin.Context) {
	voteID := c.Param("id")
	voteIDUint, err := strconv.ParseUint(voteID, 10, 32)
	if err != nil {
		log.Println("Error converting vote ID to uint: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	vote, err := va.votesList.GetVote(uint(voteIDUint))
	if err != nil {
		log.Println("Error getting vote: ", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	// Delete the vote from the voter's vote history using the voter API.
	voterID := vote.VoterID
	voterPollURL := fmt.Sprintf("%s/voters/%s/polls/%s", va.voterAPIURL, strconv.Itoa(int(voterID)), strconv.Itoa(int(vote.PollID)))

	// Make an HTTP DELETE request to remove the vote from the voter's vote history.
	resp, err := va.apiClient.R().
		Delete(voterPollURL)
	if err != nil {
		log.Println("Error performing DELETE request to remove vote from voter's vote history: ", err)
		log.Println(resp.StatusCode())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Now, delete the vote from your own votes cache.
	if err := va.votesList.DeleteVote(uint(voteIDUint)); err != nil {
		log.Println("Error deleting vote from cache: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Vote deleted successfully.",
	})
}

// Implementation of GET Votes/health.
// Get the health status of the voter API.
func (va *VotesAPI) HealthCheck(c *gin.Context) {
	uptime := time.Since(va.bootTime).String()
	averageRequestTime := time.Duration(0)
	if va.totalCalls > 0 {
		averageRequestTime = va.totalRequestTime / time.Duration(va.totalCalls)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":             "ok",
		"uptime":             uptime,
		"totalAPICalls":      va.totalCalls,
		"totalAPICallsError": va.errorCalls,
		"bootTime":           va.bootTime,
		"totalRequestTime":   va.totalRequestTime.String(),
		"averageRequestTime": averageRequestTime.String(),
	})
}
