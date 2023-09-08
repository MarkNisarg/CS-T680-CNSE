package voter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/go-resty/resty/v2"
	"github.com/nitishm/go-rejson/v4"
)

const (
	RedisNilError        = "redis: nil"
	RedisDefaultLocation = "redis:6379"
	RedisKeyPrefix       = "voter:"
)

// voterPoll represents the voting information for a specific poll.
type voterPoll struct {
	PollID   uint      `json:"pollId"`
	VoteDate time.Time `json:"voteDate"`
}

// Voter represents a voter with a unique ID and voting history.
type Voter struct {
	VoterID     uint        `json:"voterId"`
	FirstName   string      `json:"firstName"`
	LastName    string      `json:"lastName"`
	VoteHistory []voterPoll `json:"voteHistory"`
}

// VoterList is a collection of voters.
type VoterList struct {
	Voters map[uint]Voter
}

type cache struct {
	cacheClient *redis.Client
	jsonHelper  *rejson.Handler
	context     context.Context
}

// The reference to a cache object.
type VoterCache struct {
	cache
	apiClient *resty.Client
}

// The constructor function that returns a pointer to a new VoterCache.
// It uses the default Redis URL with the companion constructor newVoterCacheInstance.
func NewVoterCache() (*VoterCache, error) {
	redisUrl := os.Getenv("REDIS_URL")

	if redisUrl == "" {
		redisUrl = RedisDefaultLocation
	}

	return newVoterCacheInstance(redisUrl)
}

// The constructor function that returns a pointer to a new VoterCache.
func newVoterCacheInstance(url string) (*VoterCache, error) {
	apiClient := resty.New()
	client := redis.NewClient(&redis.Options{
		Addr: url,
	})

	ctx := context.Background()

	err := client.Ping(ctx).Err()
	if err != nil {
		log.Println("Error connecting to redis" + err.Error())
		return nil, err
	}

	jsonHelper := rejson.NewReJSONHandler()
	jsonHelper.SetGoRedisClientWithContext(ctx, client)

	return &VoterCache{
		cache: cache{
			cacheClient: client,
			jsonHelper:  jsonHelper,
			context:     ctx,
		},
		apiClient: apiClient,
	}, nil
}

// Get a string that can be used as a key in redis.
func redisKeyFromId(id uint) string {
	return fmt.Sprintf("%s%d", RedisKeyPrefix, id)
}

// Create a new Voter instance with the provided details.
func NewVoter(id uint, firstName string, lastName string) Voter {
	voter := Voter{
		VoterID:     id,
		FirstName:   firstName,
		LastName:    lastName,
		VoteHistory: make([]voterPoll, 0),
	}

	return voter
}

// Helper to return a Voter from redis provided a key.
func (vc *VoterCache) getItemFromRedis(key string, item *Voter) error {
	itemObject, err := vc.jsonHelper.JSONGet(key, ".")
	if err != nil {
		return err
	}

	err = json.Unmarshal(itemObject.([]byte), item)
	if err != nil {
		return err
	}

	return nil
}

// Return a slice of all voters from the VoterCache.
func (vc *VoterCache) GetAllVoters() ([]Voter, error) {
	var voters []Voter

	pattern := fmt.Sprintf("%s*", RedisKeyPrefix)
	keys, err := vc.cacheClient.Keys(vc.context, pattern).Result()

	if err != nil {
		return voters, err
	}

	for _, key := range keys {
		var voter Voter

		err := vc.getItemFromRedis(key, &voter)
		if err != nil {
			return voters, err
		}

		voters = append(voters, voter)
	}

	return voters, nil
}

// Retrieve a single voter from the VoterCache by voterID.
func (vc *VoterCache) GetVoter(voterID uint) (Voter, error) {
	var voter Voter

	redisKey := redisKeyFromId(voterID)
	if err := vc.getItemFromRedis(redisKey, &voter); err != nil {
		return Voter{}, errors.New("voter does not exist")
	}

	return voter, nil
}

// Add a new voter to the VoterCache.
func (vc *VoterCache) AddVoter(voter Voter) error {
	if _, err := vc.GetVoter(voter.VoterID); err == nil {
		return errors.New("voter already exists")
	}

	redisKey := redisKeyFromId(voter.VoterID)
	if _, setErr := vc.jsonHelper.JSONSet(redisKey, ".", voter); setErr != nil {
		return setErr
	}

	return nil
}

// Update an existing voter in the VoterCache.
func (vc *VoterCache) UpdateVoter(voter Voter) (Voter, error) {
	existingVoter, err := vc.GetVoter(voter.VoterID)

	if err != nil {
		return Voter{}, errors.New("voter does not exist")
	}

	existingVoter.FirstName = voter.FirstName
	existingVoter.LastName = voter.LastName

	redisKey := redisKeyFromId(voter.VoterID)
	if _, setErr := vc.jsonHelper.JSONSet(redisKey, ".", existingVoter); setErr != nil {
		return Voter{}, setErr
	}

	return existingVoter, nil
}

// Delete all voters from the VoterCache.
func (vc *VoterCache) DeleteAllVoters() error {
	pattern := fmt.Sprintf("%s*", RedisKeyPrefix)
	keys, err := vc.cacheClient.Keys(vc.context, pattern).Result()
	if err != nil {
		return err
	}

	for _, key := range keys {
		if _, deleteErr := vc.cacheClient.Del(vc.context, key).Result(); deleteErr != nil {
			return deleteErr
		}
	}

	return nil
}

// Delete a single voter from the VoterCache by voterID.
func (vc *VoterCache) DeleteVoter(voterID uint) error {
	if _, err := vc.GetVoter(voterID); err != nil {
		return errors.New("voter does not exist")
	}

	redisKey := redisKeyFromId(voterID)
	if _, deleteErr := vc.jsonHelper.JSONDel(redisKey, "."); deleteErr != nil {
		return deleteErr
	}

	return nil
}

// Retrieve the vote history of a voter by voterID.
func (vc *VoterCache) GetVoterHistory(voterID uint) ([]voterPoll, error) {
	voter, err := vc.GetVoter(voterID)
	if err != nil {
		return nil, errors.New("voter does not exist")
	}

	return voter.VoteHistory, nil
}

// Retrieve a specific voter poll by voterID and pollID.
func (vc *VoterCache) GetVoterPoll(voterID, pollID uint) (voterPoll, error) {
	voter, err := vc.GetVoter(voterID)
	if err != nil {
		return voterPoll{}, errors.New("voter does not exist")
	}

	for _, poll := range voter.VoteHistory {
		if poll.PollID == pollID {
			return poll, nil
		}
	}

	return voterPoll{}, errors.New("voter poll not found")
}

// Add a new voter poll to the vote history of a voter.
func (vc *VoterCache) AddVoterPoll(voterID, pollID uint, voteDate time.Time) (voterPoll, error) {
	voter, err := vc.GetVoter(voterID)
	if err != nil {
		return voterPoll{}, errors.New("voter does not exist")
	}

	for _, poll := range voter.VoteHistory {
		if poll.PollID == pollID {
			return voterPoll{}, errors.New("voter has already voted in this poll")
		}
	}

	newVoterPoll := voterPoll{
		PollID:   pollID,
		VoteDate: voteDate,
	}

	voter.VoteHistory = append(voter.VoteHistory, newVoterPoll)

	redisKey := redisKeyFromId(voter.VoterID)
	if _, setErr := vc.jsonHelper.JSONSet(redisKey, ".", voter); setErr != nil {
		return voterPoll{}, setErr
	}

	return newVoterPoll, nil
}

// Update an existing voter poll in the vote history of a voter.
func (vc *VoterCache) UpdateVoterPoll(voterID, pollID uint, voteDate time.Time) (voterPoll, error) {
	voter, err := vc.GetVoter(voterID)
	if err != nil {
		return voterPoll{}, errors.New("voter does not exist")
	}

	var updatedVoterPoll voterPoll
	for i, poll := range voter.VoteHistory {
		if poll.PollID == pollID {
			updatedVoterPoll = voterPoll{
				PollID:   pollID,
				VoteDate: voteDate,
			}
			voter.VoteHistory[i] = updatedVoterPoll
			break
		}
	}

	if updatedVoterPoll.PollID == 0 {
		return voterPoll{}, errors.New("voter poll not found")
	}

	redisKey := redisKeyFromId(voter.VoterID)
	if _, setErr := vc.jsonHelper.JSONSet(redisKey, ".", voter); setErr != nil {
		return voterPoll{}, setErr
	}

	return updatedVoterPoll, nil
}

// Remove a specific voter poll from the vote history of a voter.
func (vc *VoterCache) DeleteVoterPoll(voterID, pollID uint) error {
	voter, err := vc.GetVoter(voterID)
	if err != nil {
		return errors.New("voter does not exist")
	}

	var updatedVoteHistory []voterPoll
	for _, poll := range voter.VoteHistory {
		if poll.PollID != pollID {
			updatedVoteHistory = append(updatedVoteHistory, poll)
		}
	}

	if len(updatedVoteHistory) == len(voter.VoteHistory) {
		return errors.New("voter poll not found")
	}

	voter.VoteHistory = updatedVoteHistory
	redisKey := redisKeyFromId(voter.VoterID)
	if _, setErr := vc.jsonHelper.JSONSet(redisKey, ".", voter); setErr != nil {
		return setErr
	}

	return nil
}
