package api

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"voter-api/voter"

	"github.com/gin-gonic/gin"
)

// The API handler that handles incoming requests.
type VoterAPI struct {
	voterList        *voter.VoterList
	totalCalls       uint64
	errorCalls       uint64
	bootTime         time.Time
	totalRequestTime time.Duration
}

// Create a new instance of VoterAPI with an initialized voter list.
func NewVoterHandler() *VoterAPI {
	return &VoterAPI{
		voterList:        voter.NewVoterList(),
		totalCalls:       0,
		errorCalls:       0,
		bootTime:         time.Now(),
		totalRequestTime: 0,
	}
}

// The custom middleware to handle health metadata.
func HealthMiddleware(va *VoterAPI) gin.HandlerFunc {
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
func (va *VoterAPI) WelcomeToVoterAPI(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome to voter API.",
	})
}

// Implementation of GET /voters.
// Returns all voters with all voter history.
func (va *VoterAPI) ListAllVoters(c *gin.Context) {
	voters := va.voterList.GetAllVoters()

	if voters == nil {
		voters = make([]voter.Voter, 0)
	}

	c.JSON(http.StatusOK, voters)
}

// Implementation of GET /voters/:id.
// Returns a single voter by :id.
func (va *VoterAPI) GetVoter(c *gin.Context) {
	voterID := c.Param("id")
	voterIDUint, err := strconv.ParseUint(voterID, 10, 32)
	if err != nil {
		log.Println("Error converting voter ID to uint: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	voter, err := va.voterList.GetVoter(uint(voterIDUint))
	if err != nil {
		log.Println("Error getting voter: ", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, voter)
}

// Implementation of POST /voters/:id.
// Add a new voter with :id.
func (va *VoterAPI) AddVoter(c *gin.Context) {
	voterID := c.Param("id")
	voterIDUint, err := strconv.ParseUint(voterID, 10, 32)
	if err != nil {
		log.Println("Error converting voter ID to uint: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var newVoter voter.Voter
	if err := c.ShouldBindJSON(&newVoter); err != nil {
		log.Println("Error binding JSON: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	newVoter = voter.NewVoter(uint(voterIDUint), newVoter.FirstName, newVoter.LastName)

	if err := va.voterList.AddVoter(newVoter); err != nil {
		log.Println("Error adding voter: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, newVoter)
}

// Implementation of PUT /voters/:id.
// Update an existing voter with :id.
func (va *VoterAPI) UpdateVoter(c *gin.Context) {
	voterID := c.Param("id")
	voterIDUint, err := strconv.ParseUint(voterID, 10, 32)
	if err != nil {
		log.Println("Error converting voter ID to uint: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var voter voter.Voter
	if err := c.ShouldBindJSON(&voter); err != nil {
		log.Println("Error binding JSON: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	voter.VoterID = uint(voterIDUint)
	if err := va.voterList.UpdateVoter(voter); err != nil {
		log.Println("Error updating voter: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, va.voterList.Voters[voter.VoterID])
}

// Implementation of DELETE /voters.
// Delete all voters.
func (va *VoterAPI) DeleteAllVoters(c *gin.Context) {
	if err := va.voterList.DeleteAllVoters(); err != nil {
		log.Println("Error deleting voters: ", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "All voters deleted successfully.",
	})
}

// Implementation of DELETE /voters/:id.
// Delete a single voter by :id.
func (va *VoterAPI) DeleteVoter(c *gin.Context) {
	voterID := c.Param("id")
	voterIDUint, err := strconv.ParseUint(voterID, 10, 32)
	if err != nil {
		log.Println("Error converting voter ID to uint: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := va.voterList.DeleteVoter(uint(voterIDUint)); err != nil {
		log.Println("Error deleting voter: ", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Voter deleted successfully.",
	})
}

// Implementation of GET /voters/:id/polls.
// Get the voting history of a voter by :id.
func (va *VoterAPI) GetVoterHistory(c *gin.Context) {
	voterID := c.Param("id")
	voterIDUint, err := strconv.ParseUint(voterID, 10, 32)
	if err != nil {
		log.Println("Error converting voter ID to uint: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	voterHistory, err := va.voterList.GetVoterHistory(uint(voterIDUint))
	if err != nil {
		log.Println("Error getting voter history: ", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, voterHistory)
}

// Implementation of GET /voters/:id/polls/:pollid.
// Get a specific poll from a voter's voting history with :id & :pollid.
func (va *VoterAPI) GetVoterPoll(c *gin.Context) {
	voterID := c.Param("id")
	voterIDUint, err := strconv.ParseUint(voterID, 10, 32)
	if err != nil {
		log.Println("Error converting voter ID to uint: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	pollID := c.Param("pollid")
	pollIDUint, err := strconv.ParseUint(pollID, 10, 32)
	if err != nil {
		log.Println("Error converting poll ID to uint: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	voterPoll, err := va.voterList.GetVoterPoll(uint(voterIDUint), uint(pollIDUint))
	if err != nil {
		log.Println("Error getting voter poll: ", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, voterPoll)
}

// Implementation of POST /voters/:id/polls/:pollid.
// Add a new poll to a voter's voting history with :id & :pollid.
func (va *VoterAPI) AddVoterPoll(c *gin.Context) {
	voterID := c.Param("id")
	voterIDUint, err := strconv.ParseUint(voterID, 10, 32)
	if err != nil {
		log.Println("Error converting voter ID to uint: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	pollID := c.Param("pollid")
	pollIDUint, err := strconv.ParseUint(pollID, 10, 32)
	if err != nil {
		log.Println("Error converting poll ID to uint: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	newVoterPoll, err := va.voterList.AddVoterPoll(uint(voterIDUint), uint(pollIDUint))
	if err != nil {
		log.Println("Error adding voter poll: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, newVoterPoll)
}

// Implementation of PUT /voters/:id/polls/:pollid.
// Update an existing poll in a voter's voting history with :id & :pollid.
func (va *VoterAPI) UpdateVoterPoll(c *gin.Context) {
	voterID := c.Param("id")
	voterIDUint, err := strconv.ParseUint(voterID, 10, 32)
	if err != nil {
		log.Println("Error converting voter ID to uint: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	pollID := c.Param("pollid")
	pollIDUint, err := strconv.ParseUint(pollID, 10, 32)
	if err != nil {
		log.Println("Error converting poll ID to uint: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	updatedVoterPoll, err := va.voterList.UpdateVoterPoll(uint(voterIDUint), uint(pollIDUint))
	if err != nil {
		log.Println("Error updating voter poll: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, updatedVoterPoll)
}

// Implementation of DELETE /voters/:id/polls/:pollid.
// Delete a specific poll from a voter's voting history with :id & :pollid.
func (va *VoterAPI) DeleteVoterPoll(c *gin.Context) {
	voterID := c.Param("id")
	voterIDUint, err := strconv.ParseUint(voterID, 10, 32)
	if err != nil {
		log.Println("Error converting voter ID to uint: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	pollID := c.Param("pollid")
	pollIDUint, err := strconv.ParseUint(pollID, 10, 32)
	if err != nil {
		log.Println("Error converting poll ID to uint: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := va.voterList.DeleteVoterPoll(uint(voterIDUint), uint(pollIDUint)); err != nil {
		log.Println("Error deleting voter poll: ", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Voter poll deleted successfully.",
	})
}

// Implementation of GET voters/health.
// Get the health status of the voter API.
func (va *VoterAPI) HealthCheck(c *gin.Context) {
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
