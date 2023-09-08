package api

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"poll-api/poll"

	"github.com/gin-gonic/gin"
)

// The API handler that handles incoming requests.
type PollAPI struct {
	pollList         *poll.PollCache
	totalCalls       uint64
	errorCalls       uint64
	bootTime         time.Time
	totalRequestTime time.Duration
}

// Create a new instance of VoterAPI with an initialized poll cache.
func NewPollHandler() *PollAPI {
	pollCache, _ := poll.NewPollCache()

	return &PollAPI{
		pollList:         pollCache,
		totalCalls:       0,
		errorCalls:       0,
		bootTime:         time.Now(),
		totalRequestTime: 0,
	}
}

// The custom middleware to handle health metadata.
func HealthMiddleware(pa *PollAPI) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Update total API calls count.
		pa.totalCalls++

		// Record the start time of the request.
		start := time.Now()

		// Process the request.
		c.Next()

		// Update error API calls count if there's an error.
		if c.Writer.Status() >= 400 {
			pa.errorCalls++
		}

		// Calculate the request duration.
		duration := time.Since(start)

		// Update the total request time.
		pa.totalRequestTime += duration
	}
}

// The root endpoint that welcomes users to the API.
func (pa *PollAPI) WelcomeToPollAPI(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome to poll API.",
	})
}

// Implementation of GET /polls.
// Returns all polls with all poll options.
func (pa *PollAPI) ListAllVPolls(c *gin.Context) {
	polls, err := pa.pollList.GetAllPolls()
	if err != nil {
		log.Println("Error getting polls: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if polls == nil {
		polls = make([]poll.Poll, 0)
	}

	c.JSON(http.StatusOK, polls)
}

// Implementation of GET /polls/:id.
// Returns a single poll by :id.
func (pa *PollAPI) GetPoll(c *gin.Context) {
	pollID := c.Param("id")
	pollIDUint, err := strconv.ParseUint(pollID, 10, 32)
	if err != nil {
		log.Println("Error converting poll ID to uint: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	poll, err := pa.pollList.GetPoll(uint(pollIDUint))
	if err != nil {
		log.Println("Error getting poll: ", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, poll)
}

// Implementation of POST /polls/:id.
// Add a new poll with :id.
func (pa *PollAPI) AddPoll(c *gin.Context) {
	pollID := c.Param("id")
	pollIDUint, err := strconv.ParseUint(pollID, 10, 32)
	if err != nil {
		log.Println("Error converting poll ID to uint: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var newPoll poll.Poll
	if err := c.ShouldBindJSON(&newPoll); err != nil {
		log.Println("Error binding JSON: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	newPoll = poll.NewPoll(uint(pollIDUint), newPoll.PollTitle, newPoll.PollQuestion)

	if err := pa.pollList.AddPoll(newPoll); err != nil {
		log.Println("Error adding poll: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, newPoll)
}

// Implementation of DELETE /polls.
// Delete all polls.
func (pa *PollAPI) DeleteAllPolls(c *gin.Context) {
	if err := pa.pollList.DeleteAllPolls(); err != nil {
		log.Println("Error deleting polls: ", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "All polls deleted successfully.",
	})
}

// Implementation of DELETE /polls/:id.
// Delete a single poll by :id.
func (pa *PollAPI) DeletePoll(c *gin.Context) {
	pollID := c.Param("id")
	pollIDUint, err := strconv.ParseUint(pollID, 10, 32)
	if err != nil {
		log.Println("Error converting poll ID to uint: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := pa.pollList.DeletePoll(uint(pollIDUint)); err != nil {
		log.Println("Error deleting poll: ", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Poll deleted successfully.",
	})
}

// Implementation of GET /polls/:id/options.
// Get the poll options of a poll by :id.
func (pa *PollAPI) GetPollOptions(c *gin.Context) {
	pollID := c.Param("id")
	pollIDUint, err := strconv.ParseUint(pollID, 10, 32)
	if err != nil {
		log.Println("Error converting poll ID to uint: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	pollOptions, err := pa.pollList.GetPollOptions(uint(pollIDUint))
	if err != nil {
		log.Println("Error getting poll options: ", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, pollOptions)
}

// Implementation of GET /polls/:id/options/:optionid.
// Get a specific poll option from a poll's options with :id & :optionid.
func (pa *PollAPI) GetPollOption(c *gin.Context) {
	pollID := c.Param("id")
	pollIDUint, err := strconv.ParseUint(pollID, 10, 32)
	if err != nil {
		log.Println("Error converting poll ID to uint: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	pollOptionID := c.Param("optionId")
	pollOptionIDUint, err := strconv.ParseUint(pollOptionID, 10, 32)
	if err != nil {
		log.Println("Error converting poll option ID to uint: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	pollOption, err := pa.pollList.GetPollOption(uint(pollIDUint), uint(pollOptionIDUint))
	if err != nil {
		log.Println("Error getting poll option: ", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, pollOption)
}

// Implementation of POST /polls/:id/polls/:optionid.
// Add a new poll to a poll's voting history with :id & :optionid.
func (pa *PollAPI) AddPollOption(c *gin.Context) {
	pollID := c.Param("id")
	pollIDUint, err := strconv.ParseUint(pollID, 10, 32)
	if err != nil {
		log.Println("Error converting poll ID to uint: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	pollOptionID := c.Param("optionId")
	pollOptionIDUint, err := strconv.ParseUint(pollOptionID, 10, 32)
	if err != nil {
		log.Println("Error converting poll option ID to uint: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var requestBody struct {
		OptionText string `json:"optionText"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		log.Println("Error parsing JSON request body: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	newPollOption, err := pa.pollList.AddPollOption(uint(pollIDUint), uint(pollOptionIDUint), requestBody.OptionText)
	if err != nil {
		log.Println("Error adding poll option: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, newPollOption)
}

// Implementation of DELETE /polls/:id/polls/:pollid.
// Delete a specific poll from a poll's voting history with :id & :pollid.
func (pa *PollAPI) DeletePollOption(c *gin.Context) {
	pollID := c.Param("id")
	pollIDUint, err := strconv.ParseUint(pollID, 10, 32)
	if err != nil {
		log.Println("Error converting poll ID to uint: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	pollOptionID := c.Param("optionId")
	pollOptionIDUint, err := strconv.ParseUint(pollOptionID, 10, 32)
	if err != nil {
		log.Println("Error converting poll option ID to uint: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := pa.pollList.DeletePollOption(uint(pollIDUint), uint(pollOptionIDUint)); err != nil {
		log.Println("Error deleting poll option: ", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Poll option deleted successfully.",
	})
}

// Implementation of GET polls/health.
// Get the health status of the poll API.
func (pa *PollAPI) HealthCheck(c *gin.Context) {
	uptime := time.Since(pa.bootTime).String()
	averageRequestTime := time.Duration(0)
	if pa.totalCalls > 0 {
		averageRequestTime = pa.totalRequestTime / time.Duration(pa.totalCalls)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":             "ok",
		"uptime":             uptime,
		"totalAPICalls":      pa.totalCalls,
		"totalAPICallsError": pa.errorCalls,
		"bootTime":           pa.bootTime,
		"totalRequestTime":   pa.totalRequestTime.String(),
		"averageRequestTime": averageRequestTime.String(),
	})
}
